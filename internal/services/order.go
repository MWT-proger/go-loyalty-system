package services

import (
	"context"

	"github.com/gofrs/uuid"

	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type OrderStorer interface {
	GetAllByParameters(ctx context.Context, options *store.OptionsQuery) ([]*models.Order, error)
	GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Order, error)
	Insert(ctx context.Context, obj *models.Order) error
}

type OrderService struct {
	OrderStore OrderStorer
}

func NewOrderService(s OrderStorer) *OrderService {

	return &OrderService{
		OrderStore: s,
	}
}

// Set(ctx context.Context, userID uuid.UUID, numberOrder string) создает новый заказ
// с numberOrder для пользователя с userID
func (s *OrderService) Set(ctx context.Context, userID uuid.UUID, numberOrder string) error {

	args := map[string]interface{}{"number": numberOrder}
	obj, err := s.OrderStore.GetFirstByParameters(ctx, args)

	if err != nil {
		return lErrors.InternalServicesError
	}

	if obj != nil {

		if obj.UserID != userID {
			return lErrors.OrderExistsOtherUserServicesError
		}

		return lErrors.OrderExistsServicesError
	}

	newOrder, err := models.NewOrder()

	if err != nil {
		return lErrors.InternalServicesError
	}
	newOrder.Number = numberOrder
	newOrder.UserID = userID

	err = s.OrderStore.Insert(ctx, newOrder)

	return err
}

// GetList(ctx context.Context, userID uuid.UUID) Возвращает список заказов пользователя
func (s *OrderService) GetList(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {

	filterParams := []store.FilterParams{
		{Field: "user_id", Value: userID},
	}
	objs, err := s.OrderStore.GetAllByParameters(ctx, &store.OptionsQuery{
		Filter: filterParams, Sorting: []store.SortingParams{{Key: "updated_at"}}})

	if err != nil {
		return nil, lErrors.InternalServicesError
	}

	if len(objs) == 0 {
		return nil, lErrors.ListOrdersEmptyServicesError
	}

	return objs, nil
}
