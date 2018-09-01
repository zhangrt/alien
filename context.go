package alien

import "net/http"

type Context struct {
	pattern string
	pathVariables map[string]string
	ResponseWriter http.ResponseWriter
	Request *http.Request
	params map[string][]string
}


