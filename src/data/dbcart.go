package data

import "golangbootcamp/src/model"

type DBCart interface {
	ReadProductCartUser(idUser,idProduct string) (int,error)
	UpdateProductCartUser(idUser,idProduct string, quantity uint) error
	AddProductCartUser(idUser,idProduct string, quantity uint) error
	RemoveOneProductCartUser(idUser,idProduct string) error
	DeleteProductCartUser(idUser,idProduct string) error
	InsertProductCartUser(cart model.UserCart) error
	DeleteAllProductsCartUser(idUser string) error
}

func (db DB) ReadProductCartUser(idUser,idProduct string) (int,error)  {
	var quantity int
	err := db.database.QueryRow("select quantity from CART where id_user=? AND id_product = ?",idUser, idProduct).Scan(&quantity)
	if err != nil{
		return 0,err
	}
	return quantity,nil
}

func (db DB) UpdateProductCartUser(idUser,idProduct string, quantity uint) error {
	row := db.database.QueryRow("update CART set quantity = ? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return row.Err()
	}
	return nil
}

func (db DB) AddProductCartUser(idUser,idProduct string, quantity uint) error {
	row := db.database.QueryRow("update CART set quantity = quantity+? where id_user=? AND id_product = ?",quantity,idUser,idProduct)
	if row.Err() != nil{
		return row.Err()
	}
	return nil
}

func (db DB) RemoveOneProductCartUser(idUser,idProduct string) error {
	row := db.database.QueryRow("update CART set quantity = quantity-1 where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return row.Err()
	}
	quantity,_ := db.ReadProductCartUser(idUser,idProduct)
	if quantity <= 0{
		return db.DeleteProductCartUser(idUser,idProduct)
	}
	return nil
}

func (db DB) DeleteProductCartUser(idUser,idProduct string) error {
	row := db.database.QueryRow("delete from CART where id_user=? AND id_product = ?",idUser,idProduct)
	if row.Err() != nil{
		return row.Err()
	}
	return nil
}

func (db DB) InsertProductCartUser(cart model.UserCart) error {
	_,err := db.database.Query("INSERT INTO CART (id_user,id_product,quantity) VALUES(?,?,?)",cart.UserId,cart.ProductId,cart.Quantity)
	if err != nil{
		return err
	}
	return nil
}

func (db DB) DeleteAllProductsCartUser(idUser string) error {
	row := db.database.QueryRow("delete from CART where id_user=? ",idUser)
	if row.Err() != nil{
		return row.Err()
	}
	return nil
}