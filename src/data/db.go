package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golangbootcamp/src/model"
	"io/ioutil"
	"net/http"
)

var database *sql.DB
var errDB error

func OpenDB()  {
	database, errDB = sql.Open("mysql", "root:my-secret-pw@/bootcamp")
	if errDB != nil {
		panic(errDB.Error())
	}
}

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

func NewUser(user model.User) bool {
	_,err := database.Query("INSERT INTO USERS (id,name,currency) VALUES(?,?,?)",user.Id,user.Name,user.Currency)
	if err != nil{
		return false
	}
	return true
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

func ReadProductByIdGO(id string,q uint,chanel chan model.ProductUser) {
	item := model.ProductUser{
		Quantity: q,
	}
	res, err := http.Get("https://challenge.getsandbox.com/articles/"+id)
	if err != nil{
		chanel <- item
	}
	defer res.Body.Close()
	var product model.InternetProduct
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		chanel <- item
	}
	item.Product = product
	chanel <- item
}

func ReadUserById(id string) *model.User {
	us := &model.User{Id: id}
	err := database.QueryRow("select name,currency from USERS where id=?", id).Scan(&us.Name, &us.Currency)
	if err != nil{
		return nil
	}
	rows, err := database.Query("select id_product,quantity from CART where id_user=?",id)
	chanel := make(chan model.ProductUser)
	n := 0
	for rows.Next(){
		var id_product string
		var quantity uint
		err = rows.Scan(&id_product, &quantity)
		if err != nil{
			return nil
		}
		go ReadProductByIdGO(id_product,quantity,chanel)
		n++
	}
	for i := 0; i < n; i++{
		productUser := <- chanel
		us.Products = append(us.Products, productUser)
	}
	return us
}
func ReadProductCartUser(idUser,idProduct string) (int,bool)  {
	var quantity int
	err := database.QueryRow("select quantity from CART where id_user=? AND id_product = ?",idUser, idProduct).Scan(&quantity)
	if err != nil{
		return 0,false
	}
	return quantity,true
}
func UpdateProductCartUser(idUser,idProduct string, quantity uint) bool {
	row := database.QueryRow("update CART set quantity = ? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func AddProductCartUser(idUser,idProduct string, quantity uint) bool {
	row := database.QueryRow("update CART set quantity = quantity+? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func RemoveOneProductCartUser(idUser,idProduct string) bool {
	row := database.QueryRow("update CART set quantity = quantity-1 where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	quantity,_ := ReadProductCartUser(idUser,idProduct)
	if quantity <= 0{
		return DeleteProductCartUser(idUser,idProduct)
	}
	return true
}

func DeleteProductCartUser(idUser,idProduct string) bool {
	row := database.QueryRow("delete from CART where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return false
	}
	return true
}

func InsertProductCartUser(cart model.UserCart) bool {
	_,err := database.Query("INSERT INTO CART (id_user,id_product,quantity) VALUES(?,?,?)",cart.UserId,cart.ProductId,cart.Quantity)
	if err != nil{
		return false
	}
	return true
}

func DeleteAllProductsCartUser(idUser string) bool {
	row := database.QueryRow("delete from CART where id_user=? ",idUser)
	if row.Err() != nil{
		return false
	}
	return true
}

/*func PrintProducts(products []model.Product) string {
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
}*/

/*func RemoveProductCartUser(idUser,idProduct string) bool {
	user := ReadUserById(idUser)
	if user == nil{
		return false
	}
	user.RemoveProductCart(idProduct)
	return false
}
*/
