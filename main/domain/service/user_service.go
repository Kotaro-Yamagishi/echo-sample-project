package service

import "echoProject/main/domain/entity"

type UserService interface {
    Add(u entity.User) 
    GetInfo() []entity.User
    Delete(id string) 
}
