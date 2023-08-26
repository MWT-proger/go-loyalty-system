package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Withdrawal struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	Number    string    `json:"number,omitempty" db:"number"`
	UserID    uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Bonuses   int       `json:"bonuses,omitempty" db:"bonuses"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func (*Withdrawal) GetType() string {
	return "Withdrawal"
}

func (d *Withdrawal) GetArgsInsert() []any {

	return []any{d.ID, d.Number, d.UserID, d.UpdatedAt, d.CreatedAt}
}
