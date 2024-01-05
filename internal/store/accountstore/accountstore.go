package accountstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type AccountStore struct {
	*store.Store
	*store.InsertStore[*models.Account]
	*store.GetFirstByParametersStore[*models.Account]
}

type AccountStorer interface {
	store.Inserter[*models.Account]
	store.GetFirstByParameterser[*models.Account]
}

func New(baseStorage *store.Store) *AccountStore {
	insertQuery := "INSERT INTO content.account (id, user_id, updated_at, created_at) VALUES($1,$2,$3,$4)"
	baseSelectQueryFirst := "SELECT * FROM content.account "

	insertStore := store.NewInsertStore[*models.Account](baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.Account](baseStorage, baseSelectQueryFirst)

	return &AccountStore{baseStorage, insertStore, getFirst}
}
