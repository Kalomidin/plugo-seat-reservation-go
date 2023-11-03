package common

import (
	"time"
)

type CreatedDeleted struct {
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
