package services

import "echoProject/main/internal/models"

type UserServiceIF interface {
    Add(u domain.User) 
    GetInfo() []domain.User
    Delete(id string) 
}
