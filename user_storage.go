package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/pkg/errs"
)

type UserStorage struct {
	mu       sync.Mutex
	fileName string
}

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

func newUser(v models.User) *user {
	return &user{
		ID:       v.ID,
		Name:     v.Name,
		Password: v.Password,
		IsActive: v.IsActive,
	}
}

func NewUserStorage(fileName string) *UserStorage {
	return &UserStorage{
		mu:       sync.Mutex{},
		fileName: fileName,
	}
}

func (s *UserStorage) DeleteUser(ctx context.Context, id int) error {
	// load users
	users, err := s.GetAll(ctx)
	if err != nil {
		return err
	}

	// update with id
	ok := false
	for k, v := range users {
		if v.ID == id && v.IsActive {
			ok = true
			users[k].IsActive = false
			break
		}
	}
	if !ok {
		return errs.ErrUserIDNotFound
	}

	// write data
	resp := make([]user, 0, len(users))
	for _, v := range users {
		resp = append(resp, *newUser(v))
	}
	return writeData(s, resp)
}

func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	// // sleep
	// time.Sleep(time.Second * 15)

	// open file
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0644)
	if err != nil {
		return []models.User{}, err
	}
	defer file.Close()

	// read data
	users := make([]user, 0)
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return []models.User{}, err
	}

	// return data
	resp := make([]models.User, 0, len(users))
	for _, v := range users {
		resp = append(resp, models.User{
			ID:       v.ID,
			Name:     v.Name,
			IsActive: v.IsActive,
		})
	}
	return resp, nil
}

func (s *UserStorage) GetByID(ctx context.Context, id int) (*models.User, error) {
	// load all users
	users, err := s.GetAll(ctx)
	if err != nil {
		return &models.User{}, err
	}

	// find by id
	for _, v := range users {
		if v.ID == id {
			return &v, nil
		}
	}
	return &models.User{}, errs.ErrUserIDNotFound
}

// write data
func writeData(s *UserStorage, data []user) error {
	// open file
	s.mu.Lock()
	defer s.mu.Unlock()
	writeFile, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer writeFile.Close()

	// write
	encoder := json.NewEncoder(writeFile)
	encoder.SetIndent("", "	")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) Create(ctx context.Context, usr models.User) error {
	// load users
	users, err := s.GetAll(ctx)
	if err != nil {
		return err
	}

	// check for exist
	if _, err := s.GetByID(ctx, usr.ID); err == nil {
		return errs.ErrUserAlreadyExists
	}

	// add user
	users = append(users, usr)

	// write data
	resp := make([]user, 0, len(users))
	for _, v := range users {
		resp = append(resp, user{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return writeData(s, resp)
}

func (s *UserStorage) Update(ctx context.Context, usr models.User) error {
	// load users
	users, err := s.GetAll(ctx)
	if err != nil {
		return err
	}

	// update with id
	ok := false
	for k, v := range users {
		if v.ID == usr.ID {
			ok = true
			usr.IsActive = users[k].IsActive
			users[k] = usr
			break
		}
	}
	if !ok {
		return errs.ErrUserIDNotFound
	}

	// write data
	resp := make([]user, 0, len(users))
	for _, v := range users {
		resp = append(resp, *newUser(v))
	}
	return writeData(s, resp)
}
