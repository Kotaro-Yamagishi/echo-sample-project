package controller

import (
	"echoProject/domain/entity"
)

type City interface {
	GetCity() []entity.City
}
