package orderstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type OrderStore struct {
	*store.Store
	*store.InsertStore[*models.Order]
	*store.GetFirstByParametersStore[*models.Order]
}

type OrderStorer interface {
	store.Inserter[*models.Order]
	store.GetFirstByParameterser[*models.Order]
}

func New(baseStorage *store.Store) *OrderStore {
	insertQuery := "INSERT INTO content.order (id, number, user_id, updated_at, created_at) VALUES($1,$2,$3,$4,$5)"
	baseSelectQueryFirst := "SELECT * FROM content.order WHERE "

	insertStore := store.NewInsertStore[*models.Order](baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.Order](baseStorage, baseSelectQueryFirst)

	return &OrderStore{baseStorage, insertStore, getFirst}
}
