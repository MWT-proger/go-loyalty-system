package store

import (
	"context"
	"reflect"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type UpdateStore[E models.BaseModeler] struct {
	*Store
	UpdateQuery string
}

type Updateer[E models.BaseModeler] interface {
	UpdateBatch(ctx context.Context, options *OptionsUpdateQuery) error
}

func NewUpdateStore[E models.BaseModeler](baseStorage *Store, UpdateQuery string) *UpdateStore[E] {
	return &UpdateStore[E]{baseStorage, UpdateQuery}
}

// UpdateBatch(ctx context.Context, options *OptionsQuery) error общий метод
// Обновляет в хранилище у списка строк
// отфильтрованных по OptionsUpdateQuery.Filter,
// определенные поля из OptionsUpdateQuery.ListFieldValue
func (s *UpdateStore[E]) UpdateBatch(ctx context.Context, options *OptionsUpdateQuery) error {

	var (
		obj           E
		stringTypeObj = reflect.TypeOf(obj).String()
	)

	logger.Log.Debug("Хранилище:" + stringTypeObj + ": Update...")

	query, args := addSetInQuery(s.UpdateQuery, map[string]interface{}{}, options.ListFieldValue)
	query, args = addWhereInQuery(query, args, options.Filter)

	queryF, argsF, _ := formatQuery(&query, &args)

	*queryF += " ;"

	stmt, err := s.db.PrepareNamedContext(ctx, *queryF)

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

	logger.Log.Debug("Хранилище:" + stringTypeObj + ": Update - ok")

	return nil

}
