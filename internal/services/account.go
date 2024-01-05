package services

import (
	"context"

	"github.com/gofrs/uuid"

	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type AccountStorer interface {
	GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Account, error)
	Insert(ctx context.Context, obj *models.Account) error
}

type AccountService struct {
	AccountStore AccountStorer
}

func NewAccountService(s AccountStorer) *AccountService {

	return &AccountService{
		AccountStore: s,
	}
}

// GetOrSet(ctx context.Context, userID uuid.UUID) Выводит Account пользователя
// Если до этого,  Account не существовал, создает его и возвращает
func (s *AccountService) GetOrSet(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	args := map[string]interface{}{"user_id": userID}

	obj, err := s.AccountStore.GetFirstByParameters(ctx, args)

	if err != nil {
		return nil, lErrors.InternalServicesError
	}

	if obj == nil {

		obj, err := models.NewAccount()

		if err != nil {
			return nil, lErrors.InternalServicesError
		}

		obj.UserID = userID

		if err := s.AccountStore.Insert(ctx, obj); err != nil {
			return nil, lErrors.InternalServicesError
		}
	}

	return obj, nil
}
