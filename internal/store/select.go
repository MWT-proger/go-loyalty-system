package store

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type OptionsSelect struct {
	Args        map[string]interface{}
	Limit       int
	OrderBy     string
	DescOrderBy bool
}

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
	stringTypeObj := reflect.TypeOf(obj).String()
	list := []E{}

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetFirstByParameters...")
	var values []string

	for n := range args {

		params := fmt.Sprintf("%s=:%s", n, n)

		values = append(values, params)
	}

	query := s.baseSelectQuery + "WHERE " + strings.Join(values, " AND ") + ` LIMIT 1 ;`
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

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetFirstByParameters - ок")

	return obj, nil

}

type GetAllByParametersStore[E models.BaseModeler] struct {
	*Store
	baseSelectQuery string
}

type GetAllByParameterser[E models.BaseModeler] interface {
	GetAllByParameters(ctx context.Context, options *OptionsSelect) ([]E, error)
}

func NewGetAllByParametersStore[E models.BaseModeler](baseStorage *Store, baseSelectQuery string) *GetAllByParametersStore[E] {
	return &GetAllByParametersStore[E]{baseStorage, baseSelectQuery}
}

// GetAllByParameters(ctx context.Context, args map[string]interface{}) (*E, error) общий метод
// возвращает из хранилища все строки удовлетворяющее параметрам
func (s *GetAllByParametersStore[E]) GetAllByParameters(ctx context.Context, options *OptionsSelect) ([]E, error) {
	var obj E
	stringTypeObj := reflect.TypeOf(obj).String()
	argsBool := false
	list := []E{}

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetAllByParameters...")
	var values []string

	for n := range options.Args {
		argsBool = true
		params := fmt.Sprintf("%s=:%s", n, n)

		values = append(values, params)
	}

	query := s.baseSelectQuery

	if argsBool {
		query += "WHERE " + strings.Join(values, " AND ")
	}

	if options.OrderBy != "" {
		query += ` ORDER BY "` + options.OrderBy + `"`
		if options.DescOrderBy {
			query += ` DESC`
		}
	}

	if options.Limit != 0 {
		query += " LIMIT " + strconv.Itoa(options.Limit)
	}

	query += ` ;`

	logger.Log.Debug(query)
	stmt, err := s.db.PrepareNamedContext(ctx, query)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, options.Args); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetAllByParameters - ок")

	return list, nil

}
