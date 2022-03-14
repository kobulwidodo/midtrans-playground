package repository

import (
	"go-midtrans/domain"
	"log"

	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *itemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item domain.Item) (domain.Item, error) {
	if err := r.db.Create(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func (r *itemRepository) GetById(id uint) (domain.Item, error) {
	var item domain.Item
	log.Println(r.db.Where("id = ?", id).Find(&item))
	if err := r.db.Where("id = ?", id).Find(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}
