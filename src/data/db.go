package data

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golangbootcamp/src/model"
	"net/http"
)

var database *sql.DB
var errDB error

type IDatabase interface {
	OpenDB()
	NewUser(user model.User) bool
	ReadProducts() []model.InternetProduct
	ReadProductById(id string) model.Product
	ReadUserById(id string) *model.User
	ReadProductCartUser(idUser,idProduct string) (int,bool)
	UpdateProductCartUser(idUser,idProduct string, quantity uint) bool
	AddProductCartUser(idUser,idProduct string, quantity uint) bool
	RemoveOneProductCartUser(idUser,idProduct string) bool
	DeleteProductCartUser(idUser,idProduct string) bool
	InsertProductCartUser(cart model.UserCart) bool
	DeleteAllProductsCartUser(idUser string) bool
}
type Db struct {}

func (Db)OpenDB()  {
	database, errDB = sql.Open("mysql", "root:my-secret-pw@/bootcamp")
	if errDB != nil {
		panic(errDB.Error())
	}
}


func (Db) NewUser(user model.User) bool {
	_,err := database.Query("INSERT INTO USERS (id,name,currency) VALUES(?,?,?)",user.Id,user.Name,user.Currency)
	if err != nil{
		return false
	}
	return true
}


func (Db) ReadProducts() []model.InternetProduct {
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
	return products
}

func (Db) ReadProductById(id string) model.Product {

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
}

func (Db) ReadProductByIdGO(id string,q uint,channel chan model.ProductUser) {
	item := model.ProductUser{
		Quantity: q,
	}
	res, err := http.Get("https://challenge.getsandbox.com/articles/"+id)
	if err != nil{
		channel <- item
	}
	defer res.Body.Close()
	var product model.InternetProduct
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		channel <- item
	}
	item.Product = product
	channel <- item
}

func (db Db) ReadUserById(id string) *model.User {
	us := &model.User{Id: id}
	err := database.QueryRow("select name,currency from USERS where id=?", id).Scan(&us.Name, &us.Currency)
	if err != nil{
		return nil
	}
	rows, err := database.Query("select id_product,quantity from CART where id_user=?",id)
	channel := make(chan model.ProductUser)
	n := 0
	for rows.Next(){
		var idProduct string
		var quantity uint
		err = rows.Scan(&idProduct, &quantity)
		if err != nil{
			return nil
		}
		go db.ReadProductByIdGO(idProduct,quantity,channel)
		n++
	}
	for i := 0; i < n; i++{
		productUser := <- channel
		us.Products = append(us.Products, productUser)
	}
	return us
}
func (Db) ReadProductCartUser(idUser,idProduct string) (int,bool)  {
	var quantity int
	err := database.QueryRow("select quantity from CART where id_user=? AND id_product = ?",idUser, idProduct).Scan(&quantity)
	if err != nil{
		return 0,false
	}
	return quantity,true
}
func (Db) UpdateProductCartUser(idUser,idProduct string, quantity uint) bool {
	row := database.QueryRow("update CART set quantity = ? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func (Db) AddProductCartUser(idUser,idProduct string, quantity uint) bool {
	row := database.QueryRow("update CART set quantity = quantity+? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func (db Db) RemoveOneProductCartUser(idUser,idProduct string) bool {
	row := database.QueryRow("update CART set quantity = quantity-1 where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	quantity,_ := db.ReadProductCartUser(idUser,idProduct)
	if quantity <= 0{
		return db.DeleteProductCartUser(idUser,idProduct)
	}
	return true
}

func (Db) DeleteProductCartUser(idUser,idProduct string) bool {
	row := database.QueryRow("delete from CART where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func (Db) InsertProductCartUser(cart model.UserCart) bool {
	_,err := database.Query("INSERT INTO CART (id_user,id_product,quantity) VALUES(?,?,?)",cart.UserId,cart.ProductId,cart.Quantity)
	if err != nil{
		return false
	}
	return true
}

func (Db) DeleteAllProductsCartUser(idUser string) bool {
	row := database.QueryRow("delete from CART where id_user=? ",idUser)
	if row.Err() != nil{
		return false
	}
	return true
}

