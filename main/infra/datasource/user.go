package datasource

import (
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/entity"
	"echoProject/main/infra/things/mysql"
)

type Ds struct {
	db mysql.SqlHandler
}


func NewUserDataSource(db mysql.SqlHandler) datasource.User {
	return &Ds{db: db}
}

func (ds *Ds) Store(u entity.User) {
	ds.db.Create(&u)
}

func (ds *Ds) Select() []entity.User {
	user := []entity.User{}
	ds.db.FindAll(&user)
	return user
}
func (ds *Ds) Delete(id string) {
	user := []entity.User{}
	ds.db.DeleteById(&user, id)
}

func (ds *Ds) FindBetweenName(user entity.User) {
	ds.db.Create(&user)
}

func (ds *Ds) FindBetweenEmail(user entity.User) {
	ds.db.Create(&user)
}

func (ds *Ds) FindBetweenId(user entity.User) {
	ds.db.Create(&user)
}