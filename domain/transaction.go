package domain

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ItemID               uint
	Item                 Item
	MidtransTrasactionID string             `gorm:"unique"`
	MidtransTrasaction   MidtransTrasaction `gorm:"foreignKey:MidtransTrasactionID; references:OrderCode"`
	UserID               uint
	Price                int64
}

type TransactionRepository interface {
	Create(Transaction) (Transaction, error)
}

type TransactionUsecase interface {
	Order(itemId uint, paymentType string, user User) (Transaction, error)
}
