package dao

import (
	"errors"
	"gorm.io/gorm"
	"router/global"
	"router/jwt/model"
)

func IsBlacklist(jwt string) bool {
	err := errors.Is(global.DB.Where("jwt = ?", jwt).First(&model.JwtBlacklist{}).Error, gorm.ErrRecordNotFound)
	return err
}

func JsonInBlacklist(jwtList model.JwtBlacklist) (err error) {
	return global.DB.Create(&jwtList).Error
}
