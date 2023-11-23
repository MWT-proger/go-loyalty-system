package models

import (
	"database/sql"
	"encoding/json"
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
	ID        uuid.UUID     `json:"-" db:"id"`
	Number    string        `json:"number,omitempty" db:"number"`
	Status    StatusOrder   `json:"status,omitempty" db:"status"`
	UserID    uuid.UUID     `json:"-" db:"user_id"`
	Bonuses   sql.NullInt64 `json:"-" db:"bonuses"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
}

func (*Order) GetType() string {
	return "Order"
}

func (s *Order) GetArgsInsert() []any {

	return []any{s.ID, s.Number, s.UserID, s.UpdatedAt, s.CreatedAt}
}

func (s *Order) MarshalJSON() ([]byte, error) {
	type Alias Order

	if s.Status != Processed {
		return json.Marshal(&struct {
			*Alias
			CreatedAt string `json:"uploaded_at"`
		}{
			Alias:     (*Alias)(s),
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
		})
	}

	custumOrder := &struct {
		*Alias
		Accural   float64 `json:"accrual"`
		CreatedAt string  `json:"uploaded_at"`
	}{

		Alias:     (*Alias)(s),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
	}
	if s.Bonuses.Int64 > 0 {
		custumOrder.Accural = float64(s.Bonuses.Int64) / 100
	}

	return json.Marshal(custumOrder)

}

func NewOrder() *Order {
	newUUID, _ := uuid.NewV4()
	o := &Order{
		ID:        newUUID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	return o
}
