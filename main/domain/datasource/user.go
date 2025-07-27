package datasource

import "echoProject/main/domain/entity"

type User interface{
	Store(u entity.User)
	Select() []entity.User
	Delete(id string)
	FindBetweenName(user entity.User)
	FindBetweenEmail(user entity.User)
	FindBetweenId(user entity.User)
}