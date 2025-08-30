package repository

import "echoProject/domain/entity"

type City interface {
	Select() []entity.City
}
