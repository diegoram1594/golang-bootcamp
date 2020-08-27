package server

import (
	"net/http"
	"strings"
)

type Router struct {
	rules  map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		rules:make(map[string]map[string]http.HandlerFunc),
			}
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	handler,ok := router.findHandler(r.URL.Path,r.Method)

	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler(w,r)
}

func (router Router) findHandler(path string, method string) (http.HandlerFunc, bool) {
	p := strings.Split(path, "/")
	if len(p)< 2{
		return nil,false
	}
	handler,exist := router.rules[method]["/"+p[1]]
	return handler,exist
}