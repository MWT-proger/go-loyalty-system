package orderstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type OrderStore struct {
	*store.Store
	*store.InsertStore[*models.Order]
	*store.GetFirstByParametersStore[*models.Order]
	*store.GetAllByParametersStore[*models.Order]
	*store.GetSumByParametersStore[*models.Order]
	*store.UpdateStore[*models.Order]
	*UpdateOrderPlusUserAccountStore[*models.Order]
}

type OrderStorer interface {
	store.Inserter[*models.Order]
	store.GetFirstByParameterser[*models.Order]
	store.GetAllByParameterser[*models.Order]
	store.GetSumByParameterser[*models.Order]
	store.Updateer[*models.Order]
	UpdateOrderPlusUserAccounter[*models.Order]
}

func New(baseStorage *store.Store) *OrderStore {
	insertQuery := "INSERT INTO content.order (id, number, user_id, updated_at, created_at) VALUES($1,$2,$3,$4,$5)"
	baseSelectQueryFirst := "SELECT * FROM content.order "
	sumSelectQueryFirst := "SELECT sum(bonuses) FROM content.order WHERE "
	updateQuery := "UPDATE content.order "

	insertStore := store.NewInsertStore[*models.Order](baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.Order](baseStorage, baseSelectQueryFirst)
	getAll := store.NewGetAllByParametersStore[*models.Order](baseStorage, baseSelectQueryFirst)
	sumNumber := store.NewGetSumByParametersStore[*models.Order](baseStorage, sumSelectQueryFirst)
	updateBatch := store.NewUpdateStore[*models.Order](baseStorage, updateQuery)
	update := NewUpdateOrderPlusUserAccountStore[*models.Order](baseStorage, updateQuery)

	return &OrderStore{baseStorage, insertStore, getFirst, getAll, sumNumber, updateBatch, update}
}
