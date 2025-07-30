package user

// TODO: package名どうするか
import (
	"echoProject/main/domain/datasource"
	"echoProject/main/domain/entity"
	userRepository "echoProject/main/domain/repository"
)

type User_impl struct {
	ds datasource.User
}

// コンストラクタ
func NewUserRepository(ds datasource.User) userRepository.User {
	return &User_impl{ds: ds}
}

func (db *User_impl) Store(u entity.User) {
	db.ds.Store(u)
}

func (db *User_impl) Select() []entity.User {
	users := db.ds.Select()
	var entities []entity.User
	for _, u := range users {
		entities = append(entities, entity.User{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	return entities
}

func (db *User_impl) Delete(id string) {
	db.ds.Delete(id)
}
