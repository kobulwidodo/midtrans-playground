package repository

import (
	"go-midtrans/domain"

	"gorm.io/gorm"
)

type midtransTrasactionRepository struct {
	db *gorm.DB
}

func NewMidtransTransactionRepository(db *gorm.DB) *midtransTrasactionRepository {
	return &midtransTrasactionRepository{db}
}

func (r *midtransTrasactionRepository) Create(mTrx domain.MidtransTrasaction) (domain.MidtransTrasaction, error) {
	if err := r.db.Create(&mTrx).Error; err != nil {
		return mTrx, err
	}

	return mTrx, nil
}

func (r *midtransTrasactionRepository) GetByOrderId(orderId string) (domain.MidtransTrasaction, error) {
	var mTrx domain.MidtransTrasaction
	if err := r.db.Where("order_code = ?", orderId).Find(&mTrx).Error; err != nil {
		return mTrx, err
	}

	return mTrx, nil
}

func (r *midtransTrasactionRepository) Update(mTrx domain.MidtransTrasaction) (domain.MidtransTrasaction, error) {
	if err := r.db.Save(&mTrx).Error; err != nil {
		return mTrx, err
	}

	return mTrx, nil
}
