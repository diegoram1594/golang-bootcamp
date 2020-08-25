package main

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"net/http"
)

func main() {
	products := data.ReadProducts()
	users := data.ReadUsers()
	//PrintProducts(products)
	user := findUserById("1",users)
	user.PrintCart(products)
	user = findUserById("2",users)
	data.WriteUsers(users)

	server := Servidor{port: ":8000"}
	server.Listen()
}

func PrintProducts(products []model.Product) string {
	var stringProducts string
	for _, element := range products{
		stringProducts += fmt.Sprintf("Name: %s --- Price $%.2f USD, $%.0f COP  \n",element.GetName(),
			element.GetPriceUSD(),element.GetPriceCOP())
	}
	return stringProducts
}

func findUserById(id string, users []*model.User) *model.User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

type Servidor struct {
	port string
}

func (s Servidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(data.ReadProducts())
}

func (s *Servidor) Listen() error {
	fmt.Println("Server ")
	http.Handle("/", s)
	return  http.ListenAndServe(s.port,nil)
}
