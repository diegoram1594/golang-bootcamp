package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"net/http"
	"strings"
)

type localError struct {
	description string
}

func HandleArticles(w http.ResponseWriter, r *http.Request)  {
	p := strings.Split(r.URL.Path, "/")
	switch len(p){
	case 2:
		//All products
		json.NewEncoder(w).Encode(data.ReadProducts())
	case 3:
		//Product by id
		res :=data.ReadProductById(p[2])
		if res == nil{
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(localError{description: "article not found"})
			return
		}
		json.NewEncoder(w).Encode(res)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func HandleNewUser(w http.ResponseWriter, r *http.Request)  {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		le := localError{description: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(le)
		return
	}
	le,ok := validateUser(user)
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(le)
		return
	}
	user.Cart = make(map[string]uint)
	users := data.ReadUsers()
	users = append(users, &user)
	data.WriteUsers(users)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func HandleUser(w http.ResponseWriter, r *http.Request)  {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 3{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := data.ReadUserById(p[2])
	if user != nil{
		json.NewEncoder(w).Encode(user)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func HandleRemoveItemsCart(w http.ResponseWriter, r *http.Request)  {
	var userCart model.UserCart
	json.NewDecoder(r.Body).Decode(&userCart)
	if len(userCart.ProductId) > 0{
		//Remove one item
		ok := data.RemoveProductCartUser(userCart.UserId,userCart.ProductId)
		if ok{
			w.WriteHeader(http.StatusOK)
			return
		}

	}else{
		//Remove All items
		ok := data.DeleteCartUser(userCart.UserId)
		if ok{
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func HandleAddItemCart(w http.ResponseWriter, r *http.Request)  {
	var userCart model.UserCart
	err := json.NewDecoder(r.Body).Decode(&userCart)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	if userCart.Quantity == 0 {
		userCart.Quantity = 1
	}
	product := data.ReadProductById(userCart.ProductId)
	if product == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res := data.AddProductCartUser(userCart)
	if !res{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleRoot(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w,"root")
}


func validateUser(user model.User) (localError,bool){
	var le localError
	if len(user.Name) < 2{
		le.description = "Title should have at least 2 characters"
		return le,false
	}
	if len(user.Id) == 0 {
		le.description = "Id should have at least 1 character"
		return le,false
	}
	if user.Currency == "COP" || user.Currency == "USD"{
		if data.ReadUserById(user.Id) != nil{
			le.description = "Id duplicated"
			return le,false
		}
	}else{
		le.description = "Currency should be COP or USD"
		return le,false
	}
	return le,true
}