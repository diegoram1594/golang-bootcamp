package server

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"net/http"
	"strings"
)

type IDataBase interface {
	NewDB()
	CloseDB()
}
type DBCart interface {
	ReadProductCartUser(idUser,idProduct string) (int,error)
	UpdateProductCartUser(idUser,idProduct string, quantity uint) error
	AddProductCartUser(idUser,idProduct string, quantity uint) error
	RemoveOneProductCartUser(idUser,idProduct string) error
	DeleteProductCartUser(idUser,idProduct string) error
	InsertProductCartUser(cart model.UserCart) error
	DeleteAllProductsCartUser(idUser string) error
}
type DBProducts interface {
	ReadProducts() ([]model.InternetProduct,error)
	ReadProductById(id string) (model.Product,error)
	ReadProductByIdGO(id string, quantity uint,channel chan data.ResultProductUser)
}
type DBUser interface {
	NewUser(user model.User) error
	ReadUserById(id string) (*model.User,error)
}

type ProductHandler struct {
	productDB DBProducts
}
type UserHandler struct {
	userDB DBUser
}
type CartHandler struct {
	cartDB DBCart
	userDB DBUser
	productDB DBProducts
}

func NewProductHandler(db DBProducts) *ProductHandler {
	return &ProductHandler{productDB: db}
}
func NewUserHandler(db DBUser) *UserHandler {
	return &UserHandler{userDB: db}
}
func NewCartHandler(cartDB DBCart,userDB DBUser,productDB DBProducts) *CartHandler {
	return &CartHandler{
		cartDB:    cartDB,
		userDB:    userDB,
		productDB: productDB,
	}
}

func (h ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request)  {
	path := strings.Split(r.URL.Path, "/")
	switch len(path){
	case 2:
		//All products
		products, err := h.productDB.ReadProducts()
		if err!= nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(products)
	case 3:
		//Product by id
		res,err := h.productDB.ReadProductById(path[2])
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

func (h UserHandler)HandleNewUser(w http.ResponseWriter, r *http.Request)  {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	validationMessage := h.validateUser(user)
	if len(validationMessage) > 0{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,validationMessage)
		return
	}
	err = h.userDB.NewUser(user)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h UserHandler)HandleGetUser(w http.ResponseWriter, r *http.Request)  {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user,err := h.userDB.ReadUserById(path[2])
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

func (h CartHandler)HandleRemoveItemsCart(w http.ResponseWriter, r *http.Request)  {
	var userCart model.UserCart
	json.NewDecoder(r.Body).Decode(&userCart)

	if len(userCart.ProductId) > 0{
		//Remove one item
		_, err := h.cartDB.ReadProductCartUser(userCart.UserId,userCart.ProductId)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = h.cartDB.RemoveOneProductCartUser(userCart.UserId,userCart.ProductId)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}else{
		//Remove All items
		err := h.cartDB.DeleteAllProductsCartUser(userCart.UserId)
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

func (h CartHandler) HandleAddItemCart(w http.ResponseWriter, r *http.Request)  {
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
	product,err := h.productDB.ReadProductById(userCart.ProductId)
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
	user,err := h.userDB.ReadUserById(userCart.UserId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if user == nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	quantity,_ := h.cartDB.ReadProductCartUser(userCart.UserId,userCart.ProductId)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	if quantity > 0{
		if userCart.Total{
			err = h.cartDB.UpdateProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
		} else{
			err = h.cartDB.AddProductCartUser(userCart.UserId,userCart.ProductId,userCart.Quantity)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
		}
	}else{
		err = h.cartDB.InsertProductCartUser(userCart)
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


func (h UserHandler)validateUser(user model.User) string{
	if len(user.Name) < 2{
		return "Title should have at least 2 characters"
	}
	if len(user.Id) == 0 {
		return "Id should have at least 1 character"
	}
	if user.Currency == "COP" || user.Currency == "USD"{
		userDB, _ := h.userDB.ReadUserById(user.Id)
		if userDB != nil{
			return "Id duplicated"
		}
	}else{
		return "Currency should be COP or USD"
	}
	return ""
}