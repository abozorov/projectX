package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/abozorov/projectX/models"
	"github.com/abozorov/projectX/package/errs"
)

type UserStorage struct {
	mu       sync.Mutex
	fileName string
}

func NewUSerStorage(fileName string) *UserStorage {
	return &UserStorage{
		mu:       sync.Mutex{},
		fileName: fileName,
	}
}

func (s *UserStorage) GetAll() ([]models.User, error) {
	// open file
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0644)
	if err != nil {
		return []models.User{}, err
	}
	defer file.Close()

	// read data
	users := make([]models.User, 0)
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return []models.User{}, err
	}

	// return data
	return users, nil
}

func (s *UserStorage) GetByID(id int) (*models.User, error) {
	// load all users
	users, err := s.GetAll()
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
func writeData(s *UserStorage, data []models.User) error {
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

func (s *UserStorage) Create(user models.User) error {
	// load users
	users, err := s.GetAll()
	if err != nil {
		return err
	}

	// check for exist
	if _, err := s.GetByID(user.ID); err == nil {
		return errs.ErrUserAlreadyExists
	}

	// add user
	users = append(users, user)

	// write data
	return writeData(s, users)
}

func (s *UserStorage) Update(id int, user models.User) error {
	// load users
	users, err := s.GetAll()
	if err != nil {
		return err
	}

	// update with id
	ok := false
	for k, v := range users {
		if v.ID == id {
			ok = true
			users[k] = user
			break
		}
	}
	if !ok {
		return errs.ErrUserIDNotFound
	}

	// write data
	return writeData(s, users)
}
