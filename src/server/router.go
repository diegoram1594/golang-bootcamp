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
	handler,ok := router.FindHandler(r.URL.Path,r.Method)

	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler(w,r)
}

func (router Router) FindHandler(pathRequest string, method string) (http.HandlerFunc, bool) {
	path := strings.Split(pathRequest, "/")
	if len(path)< 2{
		return nil,false
	}
	handler,exist := router.rules[method]["/"+path[1]]
	return handler,exist
}