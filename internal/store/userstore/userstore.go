package userstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type UserStore struct {
	*store.Store
	*store.InsertStore[*models.User]
}

type UserStorer interface {
	store.Inserter[*models.User]
}

func New(baseStorage *store.Store) *UserStore {
	insertQuery := "INSERT INTO auth.user (login, password, created_at) VALUES($1,$2,$3)"

	insertStore := store.NewInsertStore[*models.User](baseStorage, insertQuery)

	return &UserStore{baseStorage, insertStore}
}
