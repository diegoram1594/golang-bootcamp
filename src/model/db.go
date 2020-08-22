package model

var Products = []Product{
	BasicProduct{Name: "Water", PriceCOP: 2000},
	BasicProduct{Name: "Rice", PriceCOP: 15000},
	BasicProduct{Name: "Milk", PriceCOP: 4000},
	BasicProduct{Name: "Meat", PriceCOP: 8000},
	NormalProduct{Name: "TV", PriceCOP: 2000000},
	NormalProduct{Name: "MacBook", PriceCOP: 40000000},
}

var Users = []User{
	{
		Name:     "Diego",
		Currency: "COP",
		Cart:     make(map[Product]int),
		Id: "1",
	},
	{
		Name:     "Lorena",
		Currency: "USD",
		Cart:     make(map[Product]int),
		Id: "2",
	},
}