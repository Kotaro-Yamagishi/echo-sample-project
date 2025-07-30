package datasource

import (
	"context"
	"database/sql"
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/entity"
	"echoProject/main/domain/model"
	"echoProject/main/infra/things/sqlboiler"
	"log"
)

type Ds struct {
	db        *sql.DB
	sqlBoiler sqlboiler.SQLBoiler
}

func NewUserDataSource(sqlBoiler sqlboiler.SQLBoiler) datasource.User {

	db := sqlBoiler.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	return &Ds{db: db, sqlBoiler: sqlBoiler}
}

func (ds *Ds) Store(u entity.User) {
	// ds.db.Create(&u)
}

func (ds *Ds) Select() model.UserSlice {
	users, err := model.Users().All(context.Background(), ds.db)
	if err != nil {
		log.Fatal(err)
	}
	return users
}

func (ds *Ds) Delete(id string) {
	// user := []entity.User{}
	// ds.db.DeleteById(&user, id)
}

func (ds *Ds) FindBetweenName(user entity.User) {
	// ds.db.Create(&user)
}

func (ds *Ds) FindBetweenEmail(user entity.User) {
	// ds.db.Create(&user)
}

func (ds *Ds) FindBetweenId(user entity.User) {
	// ds.db.Create(&user)
}
