package alien

import (
	"net"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)

type Router struct {
	mutex *sync.RWMutex
	getEntries map[string] *routerEntry
	postEntries map[string] *routerEntry
	putEntries map[string] *routerEntry
	deleteEntries map[string] *routerEntry
	hosts bool
}

type routerEntry struct {
	explicit bool
	handler Handler
	regex *regexp.Regexp
	pathVariables map[int]string
	pattern string
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(ctx *Context)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx)
}

func NewRouter() *Router  {
	return &Router{
		mutex: new(sync.RWMutex),
		getEntries: make(map[string] *routerEntry),
		postEntries: make(map[string] *routerEntry),
		putEntries: make(map[string] *routerEntry),
		deleteEntries: make(map[string] *routerEntry),
		hosts: false}
}

func (router *Router) Get(pattern string, handler HandlerFunc)  {
	p, entry := router.generateEntry(pattern, handler)
	router.getEntries[p] = entry
}

func (router *Router) Post(pattern string, handler HandlerFunc)  {
	p, entry := router.generateEntry(pattern, handler)
	router.postEntries[p] = entry
}

func (router *Router) Put(pattern string, handler HandlerFunc)  {
	p, entry := router.generateEntry(pattern, handler)
	router.putEntries[p] = entry
}

func (router *Router) Delete(pattern string, handler HandlerFunc)  {
	p, entry := router.generateEntry(pattern, handler)
	router.deleteEntries[p] = entry
}

func (router *Router) generateEntry(p string, handler HandlerFunc) (pattern string, entry *routerEntry) {
	pattern = cleanPath(p)

	// 如果path中含有'：'说明有带参数的url模板
	if strings.IndexByte(pattern, ':') == -1 {
		entry = &routerEntry{explicit:false, handler: HandlerFunc(handler), pattern: pattern}
	} else {
		parts := strings.Split(pattern, "/")
		pathVariables := make(map[int]string)
		for i, part := range parts {
			if strings.HasPrefix(part, ":") {
				expr := "([^/]+)"

				//a user may choose to override the defult expression
				// similar to expressjs: ‘/user/:id([0-9]+)’

				if index := strings.Index(part, "("); index != -1 {
					expr = part[index:]
					part = part[:index]
				}
				pathVariables[i] = part
				parts[i] = expr
			}
		}

		//recreate the url pattern, with parameters replaced
		//by regular expressions. then compile the regex

		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {

			//TODO add error handling here to avoid panic
			panic(regexErr)
			return
		}

		entry = &routerEntry{explicit:false, handler: HandlerFunc(handler), regex: regex, pathVariables: pathVariables, pattern: pattern}
	}
	return
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.RequestURI == "*" {
		w.Header().Set("Connection", "close")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, pattern, pathVariables := router.Handler(r)
	ctx := new(Context)
	ctx.pattern = pattern
	ctx.pathVariables = pathVariables
	ctx.w = w
	ctx.r = r
	h.ServeHTTP(ctx)
}

func (router *Router) Handler(r *http.Request) (h Handler, pattern string, pathVariables map[string]string) {

	// All other requests have any port stripped and path cleaned
	// before passing to mux.handler.
	host := stripHostPort(r.Host)
	path := cleanPath(r.URL.Path)
	method := r.Method

	return router.handler(host, path, method)
}

func (router *Router) handler(host, path string, method string) (h Handler, pattern string, pathVariables map[string]string) {
	router.mutex.RLock()
	defer router.mutex.RUnlock()

	h, pattern, pathVariables = router.match(path, method)

	if h == nil {
		h, pattern, pathVariables = http.NotFoundHandler(), "", make(map[string]string)
	}
	return
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (router *Router) match(path string, method string) (h Handler, pattern string, pathVariables map[string]string) {
	entries := router.fetchEntries(method)
	pathVariables = make(map[string]string)

	// Check for exact match first.
	v, ok := entries[path]
	if ok {
		return v.handler, v.pattern, pathVariables
	}

	// Check for longest valid match.
	for k, entry := range entries {
		if entry.regex.MatchString(path) {
			parts := strings.Split(path, "/")
			for index, pathVariable := range entry.pathVariables{
				pathVariables[pathVariable] = parts[index]
			}
			h = entry.handler
			pattern = k
			break
		}
	}
	return
}

func (router *Router) fetchEntries(method string) (entries map[string] *routerEntry)  {
	if method == "GET"{
		entries = router.getEntries
	} else if method == "POST" {
		entries = router.postEntries
	} else if method == "PUT" {
		entries = router.putEntries
	} else if method == "DELETE" {
		entries = router.deleteEntries
	}
	return
}

// stripHostPort returns h without any trailing ":<port>".
func stripHostPort(h string) string {
	// If no port on host, return unchanged
	if strings.IndexByte(h, ':') == -1 {
		return h
	}
	host, _, err := net.SplitHostPort(h)
	if err != nil {
		return h // on error, return unchanged
	}
	return host
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)

	if np[len(np)-1] != '/' {
		np += "/"
	}

	return np
}