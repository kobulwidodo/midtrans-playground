package repository

import (
	"go-midtrans/domain"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionReposity(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(trx domain.Transaction) (domain.Transaction, error) {
	if err := r.db.Create(&trx).Error; err != nil {
		return trx, err
	}

	return trx, nil
}
