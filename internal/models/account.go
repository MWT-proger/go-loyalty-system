package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Account struct {
	ID        uuid.UUID     `json:"-" db:"id"`
	UserID    uuid.UUID     `json:"-" db:"user_id"`
	Current   sql.NullInt64 `json:"-" db:"current"`
	Withdrawn sql.NullInt64 `json:"-" db:"withdrawn"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
}

func (*Account) GetType() string {
	return "Account"
}

func (s *Account) GetArgsInsert() []any {

	return []any{s.ID, s.UserID, s.UpdatedAt, s.CreatedAt}
}

func (s *Account) MarshalJSON() ([]byte, error) {
	type Alias Account

	custumAccount := &struct {
		*Alias
		Current   int64 `json:"current"`
		Withdrawn int64 `json:"withdrawn"`
	}{
		Alias:     (*Alias)(s),
		Withdrawn: s.Withdrawn.Int64,
		Current:   s.Current.Int64,
	}

	return json.Marshal(custumAccount)

}

func NewAccount() *Account {
	newUUID, _ := uuid.NewV4()
	o := &Account{
		ID:        newUUID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	return o
}
