package domain

import (
	"gorm.io/gorm"
)

type MidtransTrasaction struct {
	gorm.Model
	OrderCode   string `gorm:"unique"`
	PaymentType string
	Amount      int64
	Status      string
	PaymentData string
	MidtransId  string
}

type MidtransTrasactionRepository interface {
	Create(MidtransTrasaction) (MidtransTrasaction, error)
	Update(MidtransTrasaction) (MidtransTrasaction, error)
	GetByOrderId(orderId string) (MidtransTrasaction, error)
}

type MidtransTransactionUsecase interface {
	Handler(orderId string)
}
