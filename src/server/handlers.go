package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"net/http"
)


func HandleArticles(w http.ResponseWriter, r *http.Request)  {
	json.NewEncoder(w).Encode(data.ReadProducts())
}

func HandleRoot(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w,"root")
}
