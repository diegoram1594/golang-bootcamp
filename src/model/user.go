package model

type User struct {
	Name string
	Currency string
	//Cart map[string]uint
	Id string
	Products []ProductUser
}
type ProductUser struct {
	Product InternetProduct
	Quantity uint
}
type UserCart struct {
	UserId, ProductId string
	Quantity          uint
	Total             bool
}

/*func (u *User) AddProductCart(cart UserCart)  {
	if cart.Total{
		u.Cart[cart.ProductId] = cart.Quantity
		return
	}
	value,ok := u.Cart[cart.ProductId]
	if ok{
		u.Cart[cart.ProductId] = value + cart.Quantity
	}else{
		u.Cart[cart.ProductId] = cart.Quantity
	}
}
func (u *User) RemoveProductCart(productId string) bool {
	value,ok := u.Cart[productId]
	if !ok{
		return false
	}
	if value >1 {
		u.Cart[productId] = value - 1
	}else{
		delete(u.Cart, productId)
	}
	return true
}
func findProductById(id string, products []Product) Product {
	for _, element := range products{
		if element.GetId() == id {
			return element
		}
	}
	return nil
}
func getPrice(u User,product Product)  float64{
	if u.Currency == "COP"{
		return product.GetPriceCOP()
	}else{
		return product.GetPriceUSD()
	}
}*/
/*func (u User) PrintCart(products []Product)  {
	fmt.Printf("%s's cart\n",u.Name)
	for id,quantity:= range u.Cart{
		product := findProductById(id,products)
		price := getPrice(u,product)
		fmt.Printf("Title: %s --- Price $%.2f USD Quantity: %d  \n", product.GetName(), price, quantity)
	}
}*/


