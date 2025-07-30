package userRepository

import (
	"echoProject/main/domain/entity"
)

type User interface {
	Store(u entity.User)
	Select() []entity.User
	Delete(id string)
}

// インメモリ