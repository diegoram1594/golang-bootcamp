package data

import (
	"encoding/json"
	"golangbootcamp/src/model"
	"net/http"
)

func (db DB) ReadProducts() ([]model.InternetProduct,error) {
	var products []model.InternetProduct
	res, err := http.Get("https://challenge.getsandbox.com/articles")
	if err != nil{
		return nil,err
	}
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil,err
	}
	return products,nil
}

func (db DB) ReadProductById(id string) (model.Product,error) {

	res, err := http.Get("https://challenge.getsandbox.com/articles/"+id)
	if err != nil{
		return nil,err
	}
	var product model.InternetProduct
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		return nil,err
	}
	return product,nil
}

func (db DB) ReadProductByIdGO(id string, quantity uint,channel chan ResultProductUser) {
	item := ResultProductUser{
	}
	res, err := http.Get("https://challenge.getsandbox.com/articles/"+id)
	if err != nil{
		item.Error = err
		channel <- item
		return
	}
	var product model.InternetProduct
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		item.Error = err
		channel <- item
		return
	}
	item.ProductUser.Product = product
	item.ProductUser.Quantity = quantity
	channel <- item
}



