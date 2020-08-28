package data

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/model"
	"io/ioutil"
	"net/http"
)

func ReadUsers() []*model.User {

	data, err := ioutil.ReadFile("src/data/users.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var slice []*model.User
	err = json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println(err)
	}
	return slice
}

func WriteUsers(users []*model.User)  {
	b, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile( "src/data/users.json",b,0644)
	if err != nil {
		fmt.Println(err)
	}

}

func getProduct(element json.RawMessage) model.Product {
	var internetProduct model.InternetProduct
	err := json.Unmarshal(element,&internetProduct)
	if err == nil {
		return internetProduct
	}
	var normalProduct model.NormalProduct
	err = json.Unmarshal(element,&normalProduct)
	if err == nil && normalProduct.TypeNormal{
		return normalProduct
	}
	return nil
}

func ReadProducts() []model.InternetProduct {
	var products []model.InternetProduct
	res, err := http.Get("https://challenge.getsandbox.com/articles")
	if err != nil{
		return nil
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil
	}
	/*data, err := ioutil.ReadFile("src/data/products.json")

	if err != nil {
		fmt.Println(err)
		return nil
	}
	var slice []json.RawMessage
	err = json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range slice{
		products = append(products, getProduct(element))
	}*/
	return products
}

func ReadProductById(id string) model.Product {

	res, err := http.Get("https://challenge.getsandbox.com/articles/"+id)
	if err != nil{
		return nil
	}
	defer res.Body.Close()
	var product model.InternetProduct
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		return nil
	}
	return product
	/*data, err := ioutil.ReadFile("src/data/products.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var slice []json.RawMessage
	err = json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, element := range slice{
		product :=  getProduct(element)
		if product.GetId() == id{
			return product
		}
	}*/
}

func ReadUserById(id string) *model.User {
	users := ReadUsers()
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func PrintProducts(products []model.Product) string {
	var stringProducts string
	for _, element := range products{
		stringProducts += fmt.Sprintf("Title: %s --- Price $%.2f USD, $%.0f COP  \n",element.GetName(),
			element.GetPriceUSD(),element.GetPriceCOP())
	}
	return stringProducts
}

func DeleteCartUser(id string) bool {
	users := ReadUsers()
	for _, user := range users {
		if user.Id == id {
			user.Cart = make(map[string]uint)
			WriteUsers(users)
			return true
		}
	}
	return false
}

func RemoveProductCartUser(idUser,idProduct string) bool {
	users := ReadUsers()
	for _, user := range users {
		if user.Id == idUser {
			res := user.RemoveProductCart(idProduct)
			WriteUsers(users)
			return res
		}
	}
	return false
}

func AddProductCartUser(cart model.UserCart) bool {
	users := ReadUsers()
	for _, user := range users {
		if user.Id == cart.UserId {
			user.AddProductCart(cart)
			WriteUsers(users)
			return true
		}
	}
	return false
}