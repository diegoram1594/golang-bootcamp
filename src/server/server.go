package server

import (
	"fmt"
	"golangbootcamp/src/data"
	"net/http"
)

type Server struct {
	Port string
	router *Router
}

func NewServer(db data.DB) *Server  {
	s := &Server{
		Port:   ":8000",
		router: NewRouter(),
	}
	userHandler := NewUserHandler(db)
	productHandler := NewProductHandler(db)
	cartHandler:= NewCartHandler(db,db,db)
	s.handle(http.MethodGet,"/", HandleRoot)
	s.handle(http.MethodGet,"/articles", productHandler.HandleProducts)
	s.handle(http.MethodPost,"/user", userHandler.HandleNewUser)
	s.handle(http.MethodGet,"/user", userHandler.HandleGetUser)
	s.handle(http.MethodPut,"/cart", cartHandler.HandleAddItemCart)
	s.handle(http.MethodDelete,"/cart",cartHandler.HandleRemoveItemsCart)
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