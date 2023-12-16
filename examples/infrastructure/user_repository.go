package infrastructure

import (
	"context"
	"log"

	"github.com/hiroaki-u/ast-playground/examples/domain"
)

type userRepository struct {
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}
func (r *userRepository) FindById(ctx context.Context, id int) (*domain.User, error) {
	log.Default().Println("UserRepository.FindById")
	return nil, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	log.Default().Println("UserRepository.FindAll")
	return nil, nil
}

func (r *userRepository) Store(ctx context.Context, user *domain.User) error {
	log.Default().Println("UserRepository.Store")
	return nil
}
