package main

import (
	"fmt"
	"golangbootcamp/src/data"
	"golangbootcamp/src/model"
	"golangbootcamp/src/server"
)

func main() {
	products := data.ReadProducts()
	users := data.ReadUsers()
	//PrintProducts(products)
	user := findUserById("1",users)
	user.PrintCart(products)
	user = findUserById("2",users)
	data.WriteUsers(users)

	servidor := server.NewServer()
	err := servidor.Listen()
	if err!= nil {
		fmt.Println(err)
	}
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


