package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type GetFirstByParametersStore[E models.BaseModeler] struct {
	*Store
	baseSelectQuery string
}

type GetFirstByParameterser[E models.BaseModeler] interface {
	GetFirstByParameters(ctx context.Context, args map[string]interface{}) (E, error)
}

func NewGetFirstByParametersStore[E models.BaseModeler](baseStorage *Store, baseSelectQuery string) *GetFirstByParametersStore[E] {
	return &GetFirstByParametersStore[E]{baseStorage, baseSelectQuery}
}

// GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*E, error) общий метод
// возвращает первую строку из хранилища удовлетворяющею параметрам
func (s *GetFirstByParametersStore[E]) GetFirstByParameters(ctx context.Context, args map[string]interface{}) (E, error) {
	var obj E
	list := []E{}

	logger.Log.Debug("Хранилище: GetFirstByParameters...")
	var values []string

	for n := range args {

		params := fmt.Sprintf("%s=:%s", n, n)

		values = append(values, params)
	}

	query := s.baseSelectQuery + strings.Join(values, " AND ") + ` LIMIT 1 ;`
	logger.Log.Debug(query)
	stmt, err := s.db.PrepareNamedContext(ctx, query)

	if err != nil {
		logger.Log.Error(err.Error())
		return obj, err
	}

	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, args); err != nil {
		logger.Log.Error(err.Error())
		return obj, err
	}

	if len(list) > 0 {
		obj = list[0]
	}

	logger.Log.Debug("Хранилище: GetFirstByParameters - ок")

	return obj, nil

}

type GetAllByParametersStore[E models.BaseModeler] struct {
	*Store
	baseSelectQuery string
}

type GetAllByParameterser[E models.BaseModeler] interface {
	GetAllByParameters(ctx context.Context, args map[string]interface{}) ([]E, error)
}

func NewGetAllByParametersStore[E models.BaseModeler](baseStorage *Store, baseSelectQuery string) *GetAllByParametersStore[E] {
	return &GetAllByParametersStore[E]{baseStorage, baseSelectQuery}
}

// GetAllByParameters(ctx context.Context, args map[string]interface{}) (*E, error) общий метод
// возвращает из хранилища все строки удовлетворяющее параметрам
func (s *GetAllByParametersStore[E]) GetAllByParameters(ctx context.Context, args map[string]interface{}) ([]E, error) {
	list := []E{}

	logger.Log.Debug("Хранилище: GetAllByParameters...")
	var values []string

	for n := range args {
		params := fmt.Sprintf("%s=:%s", n, n)

		values = append(values, params)
	}

	query := s.baseSelectQuery + strings.Join(values, " AND ") + ` ;`
	logger.Log.Debug(query)
	stmt, err := s.db.PrepareNamedContext(ctx, query)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, args); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	logger.Log.Debug("Хранилище: GetAllByParameters - ок")

	return list, nil

}
