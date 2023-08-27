package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type GetSumByParametersStore[E models.BaseModeler] struct {
	*Store
	baseSelectQuery string
}

type GetSumByParameterser[E models.BaseModeler] interface {
	GetSumByParameters(ctx context.Context, args map[string]interface{}) (int64, error)
}

func NewGetSumByParametersStore[E models.BaseModeler](baseStorage *Store, baseSelectQuery string) *GetSumByParametersStore[E] {
	return &GetSumByParametersStore[E]{baseStorage, baseSelectQuery}
}

// GetSumByParameters(ctx context.Context, args map[string]interface{}) (*E, error) общий метод
// возвращает сумму поля из строк удовлетворяющих условие
func (s *GetSumByParametersStore[E]) GetSumByParameters(ctx context.Context, args map[string]interface{}) (int64, error) {
	var sum int64

	logger.Log.Debug("Хранилище: GetSumByParameters...")
	var values []string

	for n := range args {

		params := fmt.Sprintf("%s=:%s", n, n)

		values = append(values, params)
	}

	query := s.baseSelectQuery + strings.Join(values, " AND ") + `;`
	logger.Log.Debug(query)
	stmt, err := s.db.PrepareNamedContext(ctx, query)

	if err != nil {
		logger.Log.Error(err.Error())
		return 0, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args)

	err = row.Scan(&sum)

	if err != nil {
		logger.Log.Error(err.Error())
		return 0, err
	}

	logger.Log.Debug("Хранилище: GetSumByParameters - ок")

	return sum, nil

}
