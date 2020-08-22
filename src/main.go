package main

import (
	"fmt"
	"golangbootcamp/src/model"
)

func main() {
	products := model.Products
	users := model.Users
	PrintProducts(products)
	user := findUserById("1",users)
	user.AddProductCart(findProductByName("Water",products))
	user.AddProductCart(findProductByName("Water",products))
	user.AddProductCart(findProductByName("TV",products))
	user.RemoveProductCart(findProductByName("Water",products))
	user.PrintCart()
	user = findUserById("2",users)
}

func findProductByName(name string, products []model.Product) model.Product {
	for _, element := range products{
		if element.GetName() == name{
			return element
		}
	}
	return nil
}

func PrintProducts(products []model.Product)  {
	for _, element := range products{
		fmt.Printf("Name: %s --- Price $%.2f USD, $%.0f COP  \n",element.GetName(),element.GetPriceUSD(),element.GetPriceCOP())
	}
}

func findUserById(id string, users []model.User) model.User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return model.User{}
}