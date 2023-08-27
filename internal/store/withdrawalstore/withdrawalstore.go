package withdrawalstore

import (
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type WithdrawalStore struct {
	*store.Store
	*InsertStore
	*store.GetFirstByParametersStore[*models.Withdrawal]
	*store.GetAllByParametersStore[*models.Withdrawal]
	*store.GetSumByParametersStore[*models.Withdrawal]
}

type WithdrawalStorer interface {
	Inserter
	store.GetFirstByParameterser[*models.Withdrawal]
	store.GetAllByParameterser[*models.Withdrawal]
	store.GetSumByParameterser[*models.Withdrawal]
}

func New(baseStorage *store.Store) *WithdrawalStore {
	insertQuery := "INSERT INTO content.withdrawal (id, number, bonuses, user_id, updated_at, created_at) VALUES($1,$2,$3,$4,$5,$6)"
	baseSelectQueryFirst := "SELECT * FROM content.withdrawal WHERE "
	sumSelectQueryFirst := "SELECT sum(bonuses) FROM content.withdrawal WHERE "

	insertStore := NewInsertStore(baseStorage, insertQuery)
	getFirst := store.NewGetFirstByParametersStore[*models.Withdrawal](baseStorage, baseSelectQueryFirst)
	getAll := store.NewGetAllByParametersStore[*models.Withdrawal](baseStorage, baseSelectQueryFirst)
	sumNumber := store.NewGetSumByParametersStore[*models.Withdrawal](baseStorage, sumSelectQueryFirst)

	return &WithdrawalStore{baseStorage, insertStore, getFirst, getAll, sumNumber}
}
