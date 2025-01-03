package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;autoCreateTime"`
	Username     string    `gorm:"size:255;unique;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;unique;not null"`
	ID           int       `gorm:"primaryKey;autoIncrement"`
	IsSuperuser  bool      `gorm:"type:tinyint(1);default:0"`
}

func (Account) TableName() string {
	return "account"
}

type AccountModel struct {
	DB *gorm.DB
}

func (a *AccountModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	return a.DB.Create(&Account{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}).Error
}

func (a *AccountModel) Authenticate(email, password string) (int, error) {
	var account Account

	err := a.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("account with email '%s' not found", email)
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}
