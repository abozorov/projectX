package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/abozorov/projectX/internal/models"
	events "github.com/abozorov/projectX/internal/service/eventbus"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/package/errs"
)

type UserService struct {
	storage *storage.UserStorage
	bus     *events.Bus
}

func NewUserService(storage *storage.UserStorage, bus *events.Bus) *UserService {
	return &UserService{
		storage: storage,
		bus:     bus,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	// get all users
	users, err := s.storage.GetAll(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("s.storage.GetAll: %w", err)
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Get all users",
		ClientId: rand.Int(),
	})

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	// check id
	if id < 0 {
		return &models.User{}, errs.ErrInvalidUserId
	}

	// get user
	user, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return &models.User{}, fmt.Errorf("s.storage.GetByID: %w", err)
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Get user by id",
		ClientId: rand.Int(),
	})

	return user, nil
}

func (s *UserService) Create(ctx context.Context, u models.User) error {
	// check id
	if u.ID < 0 {
		return errs.ErrInvalidUserId
	}

	// creating
	err := s.storage.Create(ctx, u)
	if err != nil {
		return fmt.Errorf("s.storage.Create: %w", err)
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Create user",
		ClientId: rand.Int(),
	})

	return nil
}

func (s *UserService) Update(ctx context.Context, id int, u models.User) error {
	// check path id
	if id < 0 || u.ID < 0 || id != u.ID {
		return errs.ErrBadRequest
	}

	// updating
	err := s.storage.Update(ctx, u)
	if err != nil {
		return fmt.Errorf("s.storage.Update: %w", err)
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Update user",
		ClientId: rand.Int(),
	})

	return nil
}
