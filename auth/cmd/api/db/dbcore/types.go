package dbcore

import "github.com/google/uuid"

// our schemas to manage the accounts

type Account struct {
	ID   uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	User User      `gorm:"constraint:OnDelete:Cascade;"`
	Role Role      `gorm:"constraint:OnDelete:Cascade;"`
}

type User struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"column:name;not null"`
	Email        string    `gorm:"column:email;not null;unique"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	AccountID    uuid.UUID `gorm:"column:account_id;not null;unique"`
}

type Role struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	IsVerified bool      `gorm:"column:is_verified;not null;default:false"`
	AboType    string    `gorm:"column:abo_type;not null;default:'NONE'"`
	AccountID  uuid.UUID `gorm:"column:account_id;not null;unique"`
}
