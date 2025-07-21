package services

import "echoProject/main/internal/models"

type UserServiceIF interface {
    Add(u models.User) 
    GetInfo() []models.User
    Delete(id string) 
}
