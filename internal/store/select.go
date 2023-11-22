package store

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type OperatorFilterQuery string

const FilterIN OperatorFilterQuery = "IN"

type SortingParams struct {
	Key  string
	Desc bool
}

type FilterParams struct {
	Field    string
	Value    interface{}
	Operator OperatorFilterQuery
}
type OptionsQuery struct {
	Filter  []FilterParams
	Sorting []SortingParams
	Limit   int
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
	GetAllByParameters(ctx context.Context, options *OptionsQuery) ([]E, error)
}

func NewGetAllByParametersStore[E models.BaseModeler](baseStorage *Store, baseSelectQuery string) *GetAllByParametersStore[E] {
	return &GetAllByParametersStore[E]{baseStorage, baseSelectQuery}
}

// GetAllByParameters(ctx context.Context, args map[string]interface{}) (*E, error) общий метод
// возвращает из хранилища все строки удовлетворяющее параметрам
func (s *GetAllByParametersStore[E]) GetAllByParameters(ctx context.Context, options *OptionsQuery) ([]E, error) {
	var (
		obj           E
		stringTypeObj = reflect.TypeOf(obj).String()
		list          = []E{}
	)

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetAllByParameters...")

	query, args, err := preparationQueryAndArgs(s.baseSelectQuery, options)
	fmt.Println(*query)
	fmt.Println(args)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	logger.Log.Debug(*query)

	stmt, err := s.db.PrepareNamedContext(ctx, *query)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, *args); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	logger.Log.Debug("Хранилище: " + stringTypeObj + ": GetAllByParameters - ок")

	return list, nil

}
func preparationQueryAndArgs(baseQuery string, options *OptionsQuery) (*string, *map[string]interface{}, error) {
	var (
		query = baseQuery
		args  = map[string]interface{}{}
	)

	for i, filterParams := range options.Filter {

		if filterParams.Field != "" {

			if i == 0 {
				query += " WHERE "
			} else {
				query += " AND "
			}

			switch filterParams.Operator {

			case "":
				query += fmt.Sprintf("%s=:%s", filterParams.Field, filterParams.Field)
			case FilterIN:
				query += fmt.Sprintf("%s IN (:%s)", "status", "status")
			}

			args[filterParams.Field] = filterParams.Value
		}
	}
	for i, options := range options.Sorting {
		if options.Key != "" {

			if i == 0 {
				query += fmt.Sprintf(" ORDER BY  %s", options.Key)
			} else {
				query += fmt.Sprintf(", %s", options.Key)
			}

			if options.Desc {
				query += " DESC"
			}

		}

	}

	if options.Limit != 0 {
		query += " LIMIT " + strconv.Itoa(options.Limit)
	}

	query += " ;"

	return formatQuery(&query, &args)
}

func formatQuery(q *string, a *map[string]interface{}) (*string, *map[string]interface{}, error) {

	query, args, _ := sqlx.Named(*q, *a)

	query, args, _ = sqlx.In(query, args...)

	query = sqlx.Rebind(sqlx.NAMED, query)
	params := map[string]interface{}{}
	for i, arg := range args {
		key := fmt.Sprintf("arg%d", i+1)
		params[key] = arg
	}

	return &query, &params, nil
}
