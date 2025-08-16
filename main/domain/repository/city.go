package repository

import "echoProject/main/domain/entity"

type City interface {
	Select() []entity.City
}
