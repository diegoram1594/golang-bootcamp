package data

import (
	"encoding/json"
	"fmt"
	"golangbootcamp/src/model"
	"io/ioutil"
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
	var basicProduct model.BasicProduct
	err := json.Unmarshal(element,&basicProduct)
	if err == nil && basicProduct.TypeBasic{
		return basicProduct
	}
	var normalProduct model.NormalProduct
	err = json.Unmarshal(element,&normalProduct)
	if err == nil && normalProduct.TypeNormal{
		return normalProduct
	}
	return nil
}

func ReadProducts() []model.Product {
	data, err := ioutil.ReadFile("src/data/products.json")
	var products []model.Product
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
	}
	return products
}

func ReadProductById(id string) model.Product {
	data, err := ioutil.ReadFile("src/data/products.json")
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
	}
	return nil
}

func FindUserById(id string) *model.User {
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
		stringProducts += fmt.Sprintf("Name: %s --- Price $%.2f USD, $%.0f COP  \n",element.GetName(),
			element.GetPriceUSD(),element.GetPriceCOP())
	}
	return stringProducts
}

func DeleteCartUser(id string) bool {
	users := ReadUsers()
	for _, user := range users {
		if user.Id == id {
			user.Cart = make(map[string]int)
			WriteUsers(users)
			return true
		}
	}
	return false
}