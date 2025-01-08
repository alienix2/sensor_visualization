package mocks

import (
	"errors"
	"log"

	"github.com/alienix2/sensor_visualization/internal/models"
)

type MockAccountModel struct {
	Accounts []models.Account
}

func (m *MockAccountModel) Insert(username, email, password string) error {
	for _, acc := range m.Accounts {
		if acc.Email == email {
			return errors.New("email already exists")
		}
		if acc.Username == username {
			return errors.New("username already exists")
		}
	}

	m.Accounts = append(m.Accounts, models.Account{
		ID:           len(m.Accounts) + 1,
		Username:     username,
		Email:        email,
		PasswordHash: password,
		IsSuperuser:  false,
	})
	return nil
}

func (m *MockAccountModel) Authenticate(email, password string) (int, error) {
	log.Println("MockAccountModel.Authenticate")
	log.Println(email, password)
	for _, acc := range m.Accounts {
		if acc.Email == email && acc.PasswordHash == password {
			log.Println("MockAccountModel.Authenticate: found")
			return acc.ID, nil
		}
	}
	return 0, errors.New("invalid email or password")
}

func (m *MockAccountModel) GetUsername(id int) string {
	for _, acc := range m.Accounts {
		if acc.ID == id {
			return acc.Username
		}
	}
	return ""
}
