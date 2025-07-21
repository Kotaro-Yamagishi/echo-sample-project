package repositories

import "echoProject/main/internal/models"


type UserRepository struct {
	SqlHandler
}


func (db *UserRepository) Store(u models.User) {
	db.Create(&u)
}

func (db *UserRepository) Select() []models.User {
	user := []models.User{}
	db.FindAll(&user)
	return user
}
func (db *UserRepository) Delete(id string) {
	user := []models.User{}
	db.DeleteById(&user, id)
}

// コンストラクタ
func NewUserRepository(sqlHandler SqlHandler) UserRepositoryIF {
    return &UserRepository{sqlHandler}
}