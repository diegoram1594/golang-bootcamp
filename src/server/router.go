package server

import (
	"net/http"
)

type Router struct {
	rules  map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		rules:make(map[string]http.HandlerFunc),
			}
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	handler,ok := router.findHandler(r.URL.Path)
	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler(w,r)

}

func (router Router) findHandler(path string) (http.HandlerFunc, bool) {
	handler,exist := router.rules[path]
	return handler,exist
}