package userstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type UserStore struct {
	*store.Store
	*store.InsertStore[*models.User]
	*store.GetFirstByParametersStore[*models.User]
}

type UserStorer interface {
	store.Inserter[*models.User]
	store.GetFirstByParameterser[*models.User]
}

func New(baseStorage *store.Store) *UserStore {
	insertQuery := "INSERT INTO auth.user (id, login, password, created_at) VALUES($1,$2,$3,$4)"
	baseSelectQueryFirst := "SELECT * FROM auth.user "

	insertStore := store.NewInsertStore[*models.User](baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.User](baseStorage, baseSelectQueryFirst)

	return &UserStore{baseStorage, insertStore, getFirst}
}
