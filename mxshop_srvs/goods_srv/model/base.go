package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey" json:"id"`
	CreateAt  *time.Time     `gorm:"column:add_time" json:"create_at"`
	UpdateAt  *time.Time     `gorm:"column:update_time" json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	IsDeleted bool           `json:"is_deleted"`
}
