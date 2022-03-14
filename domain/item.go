package domain

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string
	Price int64
}

type ItemRepository interface {
	Create(item Item) (Item, error)
	GetById(id uint) (Item, error)
}
