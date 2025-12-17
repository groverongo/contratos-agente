package handlers

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string `gorm:"primaryKey"` // StackAuth ID
	Email     string
	CreatedAt time.Time
}

type Contract struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title     string
	AuthorID  string
	Status    string `gorm:"default:'Pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Versions   []ContractVersion
	Recipients []ContractRecipient
}

type ContractVersion struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	ContractID    uuid.UUID
	VersionNumber int
	FilePath      string
	CreatedAt     time.Time
}

type ContractRecipient struct {
	ContractID     uuid.UUID `gorm:"primaryKey"`
	RecipientEmail string    `gorm:"primaryKey"`
	Status         string    `gorm:"default:'Pending'"`
	SignedAt       *time.Time
}

type ChatSession struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	ContractID uuid.UUID
	UserID     string
	CreatedAt  time.Time
	Messages   []ChatMessage `gorm:"foreignKey:SessionID"`
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	SessionID uuid.UUID
	Sender    string // user or ai
	Content   string
	CreatedAt time.Time
}
