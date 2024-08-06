package database

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                uuid.UUID `gorm:"type:uuid;primarykey;default(gen_random_uuid())"`
	UserID            string
	DescriptionHeader string
	Log               string
	ErrorMessage      string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
