package orderstore

import (
	"context"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
)

type UpdateOrderPlusUserAccountStore[E models.BaseModeler] struct {
	*store.Store
	UpdateQuery string
}

type UpdateOrderPlusUserAccounter[E models.BaseModeler] interface {
	UpdateOrderPlusUserAccount(ctx context.Context, options *store.OptionsUpdateQuery, userID uuid.UUID, bonuses int64) error
}

func NewUpdateOrderPlusUserAccountStore[E models.BaseModeler](baseStorage *store.Store, UpdateQuery string) *UpdateOrderPlusUserAccountStore[E] {
	return &UpdateOrderPlusUserAccountStore[E]{baseStorage, UpdateQuery}
}

// UpdateOrderPlusUserAccount(ctx context.Context, options *OptionsQuery) error общий метод
// Обновляет в хранилище у списка строк
// отфильтрованных по OptionsUpdateQuery.Filter,
// определенные поля из OptionsUpdateQuery.ListFieldValue
func (s *UpdateOrderPlusUserAccountStore[E]) UpdateOrderPlusUserAccount(ctx context.Context, options *store.OptionsUpdateQuery, userID uuid.UUID, bonuses int64) error {

	var (
		obj           E
		stringTypeObj = reflect.TypeOf(obj).String()
		account       models.Account
		order         models.Order
	)

	logger.Log.Debug("Хранилище:" + stringTypeObj + ": Update...")

	options.Filter = append(options.Filter, store.FilterParams{Field: "user_id", Value: userID})

	tx, err := s.GetDB().BeginTxx(ctx, nil)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	if order.Status == models.Processed {
		logger.Log.Debug(
			"У Заказа уже конечный статус ",
			zap.String("Заказ", order.Number),
		)
		return nil
	}

	// Старт - Блокировка строк для определенного user_id
	row := tx.QueryRowContext(ctx, "SELECT id, current FROM content.account WHERE user_id = $1 FOR UPDATE ", userID)
	err = row.Scan(&account.ID, &account.Current)

	if err != nil {
		logger.Log.Error(err.Error())
		account = *models.NewAccount()
		account.UserID = userID

		if err := accountstore.Insert(ctx, tx.Tx, account); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	// Конец - Блокировка строк для определенного user_id

	// Старт - Обновление заказа
	query, args := store.AddSetInQuery(s.UpdateQuery, map[string]interface{}{}, options.ListFieldValue)
	query, args = store.AddWhereInQuery(query, args, options.Filter)

	queryF, argsF, _ := store.FormatQuery(&query, &args)

	*queryF += " ;"

	stmt, err := tx.PrepareNamedContext(ctx, *queryF)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, *argsF)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	// Конец - Обновление заказа

	// Старт - Обновление пользователя
	_, err = tx.ExecContext(ctx,
		"UPDATE content.account SET current=$1, updated_at=$2 WHERE id=$3",
		account.Current.Int64+bonuses, time.Now(), account.ID)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	// Конец - Обновление пользователя

	if err := tx.Commit(); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Debug("Хранилище:" + stringTypeObj + ": Update - ok")

	return nil

}
