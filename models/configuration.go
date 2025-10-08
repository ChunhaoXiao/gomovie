package models

import "gorm.io/gorm"

type Configuration struct {
	gorm.Model

	KeyName string
	Value   string
}
