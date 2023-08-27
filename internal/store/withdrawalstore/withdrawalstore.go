package withdrawalstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type WithdrawalStore struct {
	*store.Store
	*store.InsertStore[*models.Withdrawal]
	*store.GetFirstByParametersStore[*models.Withdrawal]
	*store.GetAllByParametersStore[*models.Withdrawal]
}

type WithdrawalStorer interface {
	store.Inserter[*models.Withdrawal]
	store.GetFirstByParameterser[*models.Withdrawal]
	store.GetAllByParameterser[*models.Withdrawal]
}

func New(baseStorage *store.Store) *WithdrawalStore {
	insertQuery := "INSERT INTO content.withdrawal (id, number, bonuses, user_id, updated_at, created_at) VALUES($1,$2,$3,$4,$5,$6)"
	baseSelectQueryFirst := "SELECT * FROM content.withdrawal WHERE "

	insertStore := store.NewInsertStore[*models.Withdrawal](baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.Withdrawal](baseStorage, baseSelectQueryFirst)
	getAll := store.NewGetAllByParametersStore[*models.Withdrawal](baseStorage, baseSelectQueryFirst)

	return &WithdrawalStore{baseStorage, insertStore, getFirst, getAll}
}
