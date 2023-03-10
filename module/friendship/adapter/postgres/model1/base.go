package model

import (
	"time"

	"github.com/phantranhieunhan/s3-assignment/pkg/util"
)

type Condition map[string]interface{}

type Base struct {
	Id        string    `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (b *Base) BeforeCreate() {
	b.Id = util.GenUUID()
	b.CreatedAt = time.Now().UTC()
	b.UpdatedAt = time.Now().UTC()
}

func (b *Base) BeforeUpdate() {
	b.UpdatedAt = time.Now().UTC()
}
