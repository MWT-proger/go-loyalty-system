package models

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	New        StatusOrder = "NEW"
	Processing StatusOrder = "PROCESSING"
	Invaliud   StatusOrder = "INVALID"
	Processed  StatusOrder = "PROCESSED"
)

type Order struct {
	ID        uuid.UUID   `json:"id,omitempty" db:"id"`
	Number    string      `json:"number,omitempty" db:"number"`
	Status    StatusOrder `json:"status,omitempty" db:"status"`
	UserID    uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	Bonuses   int         `json:"bonuses,omitempty" db:"bonuses"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" db:"updated_at"`
	CreatedAt time.Time   `json:"created_at,omitempty" db:"created_at"`
}

func (*Order) GetType() string {
	return "Order"
}

func (d *Order) GetArgsInsert() []any {

	return []any{d.ID, d.Number, d.UserID, d.UpdatedAt, d.CreatedAt}
}
