package data

import "golangbootcamp/src/model"


func (db DB) NewUser(user model.User) error {
	_,err := db.database.Query("INSERT INTO USERS (id,name,currency) VALUES(?,?,?)",user.Id,user.Name,user.Currency)
	if err != nil{
		return err
	}
	return nil
}

func (db DB) ReadUserById(id string) (*model.User,error) {
	us := &model.User{Id: id}
	err := db.database.QueryRow("select name,currency from USERS where id=?", id).Scan(&us.Name, &us.Currency)
	if err != nil{
		return nil,err
	}
	rows, err := db.database.Query("select id_product,quantity from CART where id_user=?",id)
	channel := make(chan ResultProductUser)
	n := 0
	for rows.Next(){
		var idProduct string
		var quantity uint
		err = rows.Scan(&idProduct, &quantity)
		if err != nil{
			return nil,err
		}
		go db.ReadProductByIdGO(idProduct,quantity,channel)
		n++
	}
	for i := 0; i < n; i++{
		resProductUser := <- channel
		if resProductUser.Error != nil{
			return nil,err
		}
		us.Products = append(us.Products, resProductUser.ProductUser)
	}
	return us,nil
}


