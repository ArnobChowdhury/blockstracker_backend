package config

import "gorm.io/gorm"

var GormConfig = &gorm.Config{
	TranslateError: true,
}
