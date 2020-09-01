package model

type User struct {
	Name string
	Currency string
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


