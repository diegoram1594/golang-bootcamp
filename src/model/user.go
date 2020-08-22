package model

import "fmt"

type User struct {
	Name string
	Currency string
	Cart map[Product]int
	Id string
}

func (u User) AddProductCart(product Product)  {
	value,ok := u.Cart[product]
	if ok{
		u.Cart[product] = value + 1
	}else{
		u.Cart[product] = 1
	}
}
func (u User) RemoveProductCart(product Product)  {
	value,ok := u.Cart[product]
	if ok && value >1 {
		u.Cart[product] = value - 1
	}else{
		delete(u.Cart, product)
	}
}

func (u User) PrintCart()  {
	fmt.Printf("%s's cart\n",u.Name)
	for product,quantity:= range u.Cart{
		var price float64
		if u.Currency == "COP"{
			price = product.GetPriceCOP()
		}else{
			price = product.GetPriceUSD()
		}
		fmt.Printf("Name: %s --- Price $%.2f USD Quantity: %d  \n", product.GetName(), price, quantity)
	}
}
