package dao

import "gorm.io/gorm"

type Dao struct {
	DB *gorm.DB
}

type DaoInterface interface {
	TokenInterface
}
