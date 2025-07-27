package service

import (
	"log"
	"stock-api/internal/repository"
)

type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(username, password string) error {
	err := s.userRepo.RegisterUser(username, password)
	if err != nil {
		log.Printf("failed to get reegister user %s: %v", username, err)
	}

	log.Printf("Successfully registered user")

	return nil
}


func (s *UserService) LoginUser(username, password string) (string, int, error) {
	token, userId, err := s.userRepo.LoginUser(username, password)
	if err != nil {
		log.Printf("failed to get login user %s: %v", username, err)
	}

	log.Printf("Successful Login")

	return token, userId, nil
}