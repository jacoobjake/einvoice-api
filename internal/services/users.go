package services

import (
	"context"

	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/jacoobjake/einvoice-api/internal/repositories"
	"github.com/jacoobjake/einvoice-api/pkg"
)

type UserService struct {
	repo *repositories.UserRepository
}

// CreateUser creates a user and returns the user and the original (plain) password if generated.
// If the password was provided, the returned plain password will be empty.
func (s *UserService) CreateUser(ctx context.Context, user models.UserSetter) (*models.User, string, error) {
	pw, isset := user.Password.Get()
	var plainPw string

	if !isset {
		// Generate random password and hash it
		randPw, hashedPw, err := pkg.GenerateAndHashPassword(12)
		if err != nil {
			return nil, "", err
		}
		plainPw = randPw
		user.Password.Set(string(hashedPw))
	} else {
		// Hash the provided password
		hashedPw, err := pkg.HashPassword(pw)
		if err != nil {
			return nil, "", err
		}
		user.Password.Set(string(hashedPw))
	}

	createdUser, err := s.repo.Create(ctx, &user)
	if err != nil {
		return nil, "", err
	}
	return createdUser, plainPw, nil
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}
