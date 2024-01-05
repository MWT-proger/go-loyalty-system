package models

import (
	"time"

	"github.com/gofrs/uuid"
)

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

	return []any{d.ID, d.Login, d.Password, d.CreatedAt}
}

func NewUser() (*User, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	o := &User{ID: newUUID, CreatedAt: time.Now()}
	return o, nil
}
