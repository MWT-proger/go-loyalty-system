package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type StatusOrder string

const (
	New        StatusOrder = "NEW"
	Processing StatusOrder = "PROCESSING"
	Invaliud   StatusOrder = "INVALID"
	Processed  StatusOrder = "PROCESSED"
)

type BaseModeler interface {
	GetType() string
	GetArgsInsert() []any
}

type User struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	Login     string    `json:"login,omitempty" db:"login"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func (*User) GetType() string {
	return "User"
}

func (d *User) GetArgsInsert() []any {

	return []any{d.Login, d.Password, d.CreatedAt}
}

type Order struct {
	ID        uuid.UUID   `json:"id,omitempty" db:"id"`
	Number    string      `json:"number,omitempty" db:"number"`
	Status    StatusOrder `json:"status,omitempty" db:"status"`
	UserID    uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	Bonuses   int         `json:"bonuses,omitempty" db:"bonuses"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" db:"updated_at"`
	CreatedAt time.Time   `json:"created_at,omitempty" db:"created_at"`
}

type Withdrawal struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	Number    string    `json:"number,omitempty" db:"number"`
	UserID    uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Bonuses   int       `json:"bonuses,omitempty" db:"bonuses"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
