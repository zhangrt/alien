package alien

import "net/http"

type Context struct {
	pattern string
	pathVariables map[string]string
	w http.ResponseWriter
	r *http.Request
	params map[string][]string
}


