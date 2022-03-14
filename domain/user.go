package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Email        string
	Password     string
	Transactions []Transaction
}

type AuthRepository interface {
	Create(user User) (User, error)
	GetByEmail(email string) (User, error)
}

type AuthUsecase interface {
	Register(name string, email string, password string) (User, error)
	Login(email string, password string) (User, error)
	ValidateToken(token string) (User, error)
}
