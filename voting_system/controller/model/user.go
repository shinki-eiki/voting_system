package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Age       int
	Telephone string
	Password  string
	Poll      int `default:"10"`
}

func (User) TableName() string {
	return "user"
}
