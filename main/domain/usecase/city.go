package ucIF

import "echoProject/main/domain/entity"

type City interface {
	Select() []entity.City
}
