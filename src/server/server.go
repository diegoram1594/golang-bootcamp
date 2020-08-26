package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port string
	router *Router
}

func NewServer() *Server  {
	s := &Server{
		Port:   ":8000",
		router: NewRouter(),
	}
	s.handle("/", HandleRoot)
	s.handle("/articles", HandleArticles)
	return s
}

func (s *Server) handle(path string, handlerFunc http.HandlerFunc)  {
	s.router.rules[path] = handlerFunc
}

func (s *Server) Listen() error {
	fmt.Println("Server ")
	http.Handle("/", s.router)
	return  http.ListenAndServe(s.Port,nil)
}