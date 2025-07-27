package usecase

import (
	"echoProject/main/domain/entity"
	"echoProject/main/domain/repository"
	"echoProject/main/domain/service"
)

type UserService struct {
    repo userRepository.User
}

func NewUserService(repo userRepository.User) service.UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Add(u entity.User)  {
    s.repo.Store(u)
}

func (s *UserService) GetInfo() []entity.User {
    return s.repo.Select()
}

func (s *UserService) Delete(id string) {
    s.repo.Delete(id)
}
