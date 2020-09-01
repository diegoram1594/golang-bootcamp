package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"net/http"
	"strings"
)

var db data.IDatabase

func InitDB()  {
	db = data.Db{}
	db.OpenDB()
}

func HandleProducts(w http.ResponseWriter, r *http.Request)  {
	path := strings.Split(r.URL.Path, "/")
	switch len(path){
	case 2:
		//All products
		json.NewEncoder(w).Encode(db.ReadProducts())
	case 3:
		//Product by id
		res :=db.ReadProductById(path[2])
		if res == nil{
			w.WriteHeader(http.StatusNotFound)
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	validationMessage := validateUser(user)
	if len(validationMessage) > 0{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationMessage)
		return
	}
	ok := db.NewUser(user)
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func HandleGetUser(w http.ResponseWriter, r *http.Request)  {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := db.ReadUserById(path[2])
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
		_, ok := db.ReadProductCartUser(userCart.UserId,userCart.ProductId)
		if !ok{
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ok = db.RemoveOneProductCartUser(userCart.UserId,userCart.ProductId)
		if ok{
			w.WriteHeader(http.StatusOK)
			return
		}

	}else{
		//Remove All items
		ok := db.DeleteAllProductsCartUser(userCart.UserId)
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
	if userCart.Total {
		if userCart.Quantity == 0{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if userCart.Quantity == 0{
		userCart.Quantity = 1
	}
	//validate Product
	product := db.ReadProductById(userCart.ProductId)
	if product == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Validate user
	user := db.ReadUserById(userCart.UserId)
	if user == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_,ok := db.ReadProductCartUser(userCart.UserId,userCart.ProductId)
	if ok{
		if userCart.Total{
			db.UpdateProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
		} else{
			db.AddProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
		}
	}else{
		db.InsertProductCartUser(userCart)
	}
	
	w.WriteHeader(http.StatusOK)
}

func HandleRoot(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w,"root")
}


func validateUser(user model.User) string{
	if len(user.Name) < 2{
		return "Title should have at least 2 characters"
	}
	if len(user.Id) == 0 {
		return "Id should have at least 1 character"
	}
	if user.Currency == "COP" || user.Currency == "USD"{
		if db.ReadUserById(user.Id) != nil{
			return "Id duplicated"
		}
	}else{
		return "Currency should be COP or USD"
	}
	return ""
}