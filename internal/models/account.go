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
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Alias: (*Alias)(s),
	}

	if s.Current.Int64 > 0 {
		custumAccount.Current = float64(s.Current.Int64) / 100
	}
	if s.Withdrawn.Int64 > 0 {
		custumAccount.Withdrawn = float64(s.Withdrawn.Int64) / 100
	}

	return json.Marshal(custumAccount)

}

func NewAccount() (*Account, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	o := &Account{
		ID:        newUUID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	return o, nil
}
