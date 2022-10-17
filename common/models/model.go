package commonModel

import "gorm.io/gorm"

type BaseModel struct {
	// ID      uint  `gorm:"primarykey;autoIncrement"`
	// Created int64 `gorm:"autoCreateTime:milli"`
	// Updated int64 `gorm:"autoUpdateTime:milli"`
	gorm.Model
}
