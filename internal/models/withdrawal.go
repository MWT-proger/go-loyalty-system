package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Withdrawal struct {
	ID        uuid.UUID     `json:"-" db:"id"`
	Number    string        `json:"-" db:"number"`
	Status    StatusOrder   `json:"-" db:"status"`
	UserID    uuid.UUID     `json:"-" db:"user_id"`
	Bonuses   sql.NullInt64 `json:"-" db:"bonuses"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
}

func (*Withdrawal) GetType() string {
	return "Withdrawal"
}

func (d *Withdrawal) GetArgsInsert() []any {

	return []any{d.ID, d.Number, d.Bonuses, d.UserID, d.UpdatedAt, d.CreatedAt}
}
func (d *Withdrawal) MarshalJSON() ([]byte, error) {
	type Alias Withdrawal

	custumWithdrawal := &struct {
		*Alias
		Order     string `json:"order"`
		Sum       int64  `json:"sum"`
		CreatedAt string `json:"processed_at"`
	}{
		Alias:     (*Alias)(d),
		Order:     d.Number,
		Sum:       d.Bonuses.Int64,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}

	return json.Marshal(custumWithdrawal)

}

func NewWithdrawal() *Withdrawal {
	newUUID, _ := uuid.NewV4()
	o := &Withdrawal{
		ID:        newUUID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	return o
}
