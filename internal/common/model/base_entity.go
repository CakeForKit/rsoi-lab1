package model

import (
	"reflect"
	"uuid"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID uuid.UUID `json:"id,omitempty" gorm:"type:UUID;primaryKey"`
}

func (baseEntity *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if reflect.ValueOf(baseEntity.ID).IsZero() {
		baseEntity.ID = uuid.New()
	}
	return nil
}
