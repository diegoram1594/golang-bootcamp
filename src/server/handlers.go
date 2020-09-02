package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"net/http"
	"strings"
)

var db data.DB

func InitDB(newDB *data.DB)  {
	db = *newDB
}

func HandleProducts(w http.ResponseWriter, r *http.Request)  {
	path := strings.Split(r.URL.Path, "/")
	switch len(path){
	case 2:
		//All products
		products, err := db.ReadProducts()
		if err!= nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(products)
	case 3:
		//Product by id
		res,err :=db.ReadProductById(path[2])
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,validationMessage)
		return
	}
	err = db.NewUser(user)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
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
	user,err := db.ReadUserById(path[2])
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
		_, err := db.ReadProductCartUser(userCart.UserId,userCart.ProductId)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = db.RemoveOneProductCartUser(userCart.UserId,userCart.ProductId)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}else{
		//Remove All items
		err := db.DeleteAllProductsCartUser(userCart.UserId)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
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
	product,err := db.ReadProductById(userCart.ProductId)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	if product == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Validate user
	user,err := db.ReadUserById(userCart.UserId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if user == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	quantity,_ := db.ReadProductCartUser(userCart.UserId,userCart.ProductId)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	if quantity > 0{
		if userCart.Total{
			err = db.UpdateProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
		} else{
			err = db.AddProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
		}
	}else{
		err = db.InsertProductCartUser(userCart)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
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
		userDB, err := db.ReadUserById(user.Id)
		if err != nil{
			return err.Error()
		}
		if userDB != nil{
			return "Id duplicated"
		}
	}else{
		return "Currency should be COP or USD"
	}
	return ""
}