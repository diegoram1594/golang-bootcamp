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
	s.handle(http.MethodGet,"/", HandleRoot)
	s.handle(http.MethodGet,"/articles", HandleProducts)
	s.handle(http.MethodPost,"/user", HandleNewUser)
	s.handle(http.MethodGet,"/user", HandleGetUser)
	s.handle(http.MethodPut,"/cart", HandleAddItemCart)
	s.handle(http.MethodDelete,"/cart",HandleRemoveItemsCart)
	return s
}

func (s *Server) handle(method,path string, handlerFunc http.HandlerFunc)  {
	_, exist := s.router.rules[method]
	if !exist{
		s.router.rules[method] = make(map[string]http.HandlerFunc)
	}
	s.router.rules[method][path] = handlerFunc
}

func (s *Server) Listen() error {
	fmt.Println("Server ")
	http.Handle("/", s.router)
	return  http.ListenAndServe(s.Port,nil)
}