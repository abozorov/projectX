package service

import (
	"context"

	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/package/errs"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.storage.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	// check id
	if id < 0 {
		return &models.User{}, errs.ErrInvalidUserId
	}
	return s.storage.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, u models.User) error {
	// check id
	if u.ID < 0 {
		return errs.ErrInvalidUserId
	}
	return s.storage.Create(ctx, u)
}

func (s *UserService) Update(ctx context.Context, id int, u models.User) error {
	// check path id
	if id < 0 || u.ID < 0 || id != u.ID {
		return errs.ErrBadRequest
	}
	return s.storage.Update(ctx, u)
}
