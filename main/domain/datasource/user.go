package datasource

import (
	"echoProject/main/domain/entity"
	"echoProject/main/domain/model"
)

type User interface{
	Store(u entity.User)
	Select() model.UserSlice
	Delete(id string)
	FindBetweenName(user entity.User)
	FindBetweenEmail(user entity.User)
	FindBetweenId(user entity.User)
}