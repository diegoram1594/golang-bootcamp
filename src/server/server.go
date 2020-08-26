package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"net/http"
)

type Server struct {
	Port string
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(data.ReadProducts())
}

func (s *Server) Listen() error {
	fmt.Println("Server ")
	http.Handle("/", s)
	return  http.ListenAndServe(s.Port,nil)
}