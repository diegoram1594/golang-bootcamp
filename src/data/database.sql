create table USERS(
    id varchar(250) primary key ,
    name varchar(250),
    currency varchar(3)
);

create table CART(
    id_user varchar(250),
    id_product varchar(250),
    quantity int,
    FOREIGN KEY (id_user) REFERENCES USERS(id),
    UNIQUE (id_user,id_product)
);