package model

import "router/utils"

type JwtBlacklist struct {
	utils.BaseModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
