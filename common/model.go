package common

import (
	"time"

	"gorm.io/gorm"
)

type CreatedDeleted struct {
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
