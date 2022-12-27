package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Request struct {
	ID        uuid.UUID `gorm:"default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Headers   datatypes.JSONMap
}
