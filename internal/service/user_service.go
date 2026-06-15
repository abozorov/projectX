package service

import (
	"context"
	"fmt"
	"math/rand"

	events "github.com/abozorov/projectX/internal/eventbus"
	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/internal/repo"
	"github.com/abozorov/projectX/pkg/errs"
)

type UserService struct {
	storage repo.UIUserRepo
	bus     *events.Bus
}

func NewUserService(storage repo.UIUserRepo, bus *events.Bus) *UserService {
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

	s.bus.Publish(events.Event{
		Type:     "Get all users",
		ClientId: rand.Int(),
	})

	return users, nil

}

func (s *UserService) Create(ctx context.Context, u models.User) error {
	// validation
	if !u.Validate(false) {
		return errs.ErrBadRequestBody
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

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	// check id
	if id < 1 {
		return &models.User{}, errs.ErrInvalidUserId
	}

	// get user
	user, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return &models.User{}, fmt.Errorf("s.storage.GetByID: %w", err)
	}

	// business logic

	// audit
	s.bus.Publish(events.Event{
		Type:     "Get user by id",
		ClientId: rand.Int(),
	})

	return user, nil
}

func (s *UserService) Update(ctx context.Context, u models.User) error {

	// validation
	if !u.Validate(true) {
		return errs.ErrBadRequestBody
	}

	// updating
	err := s.storage.Update(ctx, u)
	if err != nil {
		return fmt.Errorf("s.storage.Update: %w", err)
	}

	// business logic

	// audit
	s.bus.Publish(events.Event{
		Type:     "Update user",
		ClientId: rand.Int(),
	})

	return nil

}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// delete user
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("s.storage.DeleteUser: %w", err)
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Delete user",
		ClientId: rand.Int(),
	})

	// return autorization code
	return nil

}

func (s *UserService) Login(ctx context.Context, u models.User) (string, error) {
	// check user

	// get user
	user, err := s.storage.GetByLogin(ctx, u.Login)
	if err != nil {
		return "", fmt.Errorf("s.storage.Login: %w", err)
	}

	if user.Password != u.Password {
		return "", errs.ErrIncorrectLoginOrPassword
	}

	// audit
	s.bus.Publish(events.Event{
		Type:     "Login",
		ClientId: rand.Int(),
	})

	// return autorization code
	return "secret", nil

}
