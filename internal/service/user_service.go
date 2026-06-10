package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/abozorov/projectX/internal/models"
	events "github.com/abozorov/projectX/internal/service/eventbus"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/pkg/errs"
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

func validation(u *models.User, create bool) error {
	if u.ID < 0 {
		return errs.ErrInvalidUserId
	}

	if u.Name, u.Password = strings.TrimSpace(u.Name), strings.TrimSpace(u.Password); u.Name == "" ||
		(u.Password == "" && create) {
		return errs.ErrBadRequestBody
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// delete user
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("s.storage.DeleteUser: %w", err)
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("s.storage.DeleteUser: %w", errs.ErrTimeoutExceeded)
	default:
		// audit
		s.bus.Publish(events.Event{
			Type:     "Login",
			ClientId: rand.Int(),
		})

		// return autorization code
		return nil
	}
}

func (s *UserService) Login(ctx context.Context, u models.User) (string, error) {
	// check user
	if u.Password = strings.TrimSpace(u.Password); u.ID < 0 || u.Password == "" {
		return "", errs.ErrIncorrectLoginOrPassword
	}

	// get user
	user, err := s.storage.GetByID(ctx, u.ID)
	if err != nil {
		return "", fmt.Errorf("s.storage.Login: %w", err)
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("s.storage.Login: %w", errs.ErrTimeoutExceeded)
	default:
		log.Print("\"", user.Password, "\" \"", u.Password, "\"")
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
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	// get all users
	users, err := s.storage.GetAll(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("s.storage.GetAll: %w", err)
	}
	select {
	case <-ctx.Done():
		return []models.User{}, fmt.Errorf("s.storage.GetAll: %w", errs.ErrTimeoutExceeded)
	default:
		// delete non active users
		for i := 0; i < len(users); i++ {
			if !users[i].IsActive {
				users = append(users[:i], users[i+1:]...)
				i--
			}
		}
		// audit
		s.bus.Publish(events.Event{
			Type:     "Get all users",
			ClientId: rand.Int(),
		})

		return users, nil
	}
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

	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("s.storage.GetByID: %w", errs.ErrTimeoutExceeded)
	default:
		// audit
		s.bus.Publish(events.Event{
			Type:     "Get user by id",
			ClientId: rand.Int(),
		})

		return user, nil
	}
}

func (s *UserService) Create(ctx context.Context, u models.User) error {
	// validation
	if err := validation(&u, true); err != nil {
		return err
	}

	// creating
	err := s.storage.Create(ctx, u)
	if err != nil {
		return fmt.Errorf("s.storage.Create: %w", err)
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("s.storage.Create: %w", errs.ErrTimeoutExceeded)
	default:

		// audit
		s.bus.Publish(events.Event{
			Type:     "Create user",
			ClientId: rand.Int(),
		})

		return nil
	}
}

func (s *UserService) Update(ctx context.Context, u models.User) error {
	// check path id
	// if id < 0 || u.ID < 0 || id != u.ID {
	// 	return errs.ErrBadRequest
	// }

	// validation
	if err := validation(&u, false); err != nil {
		return err
	}

	// updating
	err := s.storage.Update(ctx, u)
	if err != nil {
		return fmt.Errorf("s.storage.Update: %w", err)
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("s.storage.Update: %w", errs.ErrTimeoutExceeded)
	default:

		// audit
		s.bus.Publish(events.Event{
			Type:     "Update user",
			ClientId: rand.Int(),
		})

		return nil
	}
}
