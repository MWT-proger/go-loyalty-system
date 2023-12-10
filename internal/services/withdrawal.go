package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"

	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type WithdrawalStorer interface {
	GetAllByParameters(ctx context.Context, options *store.OptionsQuery) ([]*models.Withdrawal, error)
	GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Withdrawal, error)
	Insert(ctx context.Context, obj *models.Withdrawal) error
}

type WithdrawalService struct {
	WithdrawalStore WithdrawalStorer
}

func NewWithdrawalService(s WithdrawalStorer) *WithdrawalService {

	return &WithdrawalService{
		WithdrawalStore: s,
	}
}

// Set(ctx context.Context, userID uuid.UUID, numberOrder string, sum int64) Создает новое списание
func (s *WithdrawalService) Set(ctx context.Context, userID uuid.UUID, numberOrder string, sum int64) error {
	args := map[string]interface{}{"number": numberOrder}
	obj, err := s.WithdrawalStore.GetFirstByParameters(ctx, args)

	if err != nil {
		return lErrors.InternalServicesError
	}

	if obj != nil {
		if obj.UserID != userID || obj.Bonuses.Int64 != sum {
			return lErrors.WithdrawalExistsOtherUserServicesError
		}
		return lErrors.WithdrawalExistsServicesError
	}

	newWithdrawal, err := models.NewWithdrawal()

	if err != nil {
		return lErrors.InternalServicesError
	}

	newWithdrawal.Number = numberOrder
	newWithdrawal.UserID = userID
	newWithdrawal.Bonuses = sql.NullInt64{Int64: sum, Valid: true}

	err = s.WithdrawalStore.Insert(ctx, newWithdrawal)

	if err != nil {

		if errors.Is(err, &lErrors.ErrorNotBonuses{}) {
			return lErrors.NotBonusesWithdrawalServicesError
		}

		return lErrors.InternalServicesError
	}
	return nil
}

// GetList(ctx context.Context, userID uuid.UUID) возвращает список списаний пользователя
func (s *WithdrawalService) GetList(ctx context.Context, userID uuid.UUID) ([]*models.Withdrawal, error) {

	filterParams := []store.FilterParams{
		{Field: "user_id", Value: userID},
	}

	objs, err := s.WithdrawalStore.GetAllByParameters(ctx, &store.OptionsQuery{Filter: filterParams})

	if err != nil {
		return nil, lErrors.InternalServicesError
	}

	if len(objs) == 0 {
		return nil, lErrors.ListWithdrawalsEmptyServicesError
	}

	return objs, nil
}
