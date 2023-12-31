package store

import (
	"embed"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
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

type OptionsUpdateQuery struct {
	ListFieldValue map[string]interface{}
	Filter         []FilterParams
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

// migration() error - вызывается при запуске программы,
// проверяет новые миграции
// и при неообходимости обновляет БД,
// возвращает ошибку в случае неудачи
func (s *Store) migration() error {
	logger.Log.Info("Хранилище: Проверка и обновление миграций - ...")

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return err
	}

	logger.Log.Info("Хранилище: Проверка и обновление миграций - ок")

	return nil
}

// ping() error - вызывается при запуске программы,
// прверяет соединение и возвращает ошибку в случае неудачи
func (s *Store) ping() error {
	logger.Log.Info("Хранилище: Проверка соединения - ...")

	if err := s.db.Ping(); err != nil {
		return err
	}
	logger.Log.Info("Хранилище: Проверка соединения - ок")

	return nil
}

// PreparationQueryAndArgs(baseQuery string, options *OptionsQuery) (*string, *map[string]interface{}, error)
// обрабатывает OptionsQuery
// конструирует зпрос и параметры для БД
func PreparationQueryAndArgs(baseQuery string, options *OptionsQuery) (*string, *map[string]interface{}, error) {
	var (
		query = baseQuery
		args  = map[string]interface{}{}
	)

	query, args = AddWhereInQuery(query, args, options.Filter)

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

	return FormatQuery(&query, &args)
}

// FormatQuery(q *string, a *map[string]interface{}) (*string, *map[string]interface{}, error)
// форматирует в правильный вид зпрос и параметры для БД
func FormatQuery(q *string, a *map[string]interface{}) (*string, *map[string]interface{}, error) {

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

// AddWhereInQuery добавляет к SQL запросу (query) оператор WHERE с параметрами
// Возвращает новый запрос и карту (ключ - значение)
func AddWhereInQuery(query string, args map[string]interface{}, filterOptions []FilterParams) (string, map[string]interface{}) {
	for i, filterParams := range filterOptions {

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
				query += fmt.Sprintf("%s IN (:%s)", filterParams.Field, filterParams.Field)
			}

			args[filterParams.Field] = filterParams.Value
		}
	}
	return query, args

}

// AddSetInQuery добавляет к SQL запросу (query) оператор SET с параметрами
// Возвращает новый запрос и карту (ключ - значение)
func AddSetInQuery(query string, args map[string]interface{}, updetedFields map[string]interface{}) (string, map[string]interface{}) {
	var i int
	for key, value := range updetedFields {

		if i == 0 {
			query += " SET "
		} else {
			query += ", "
		}

		query += fmt.Sprintf("%s=(:%s)", key, key)

		args[key] = value
		i += 1
	}
	return query, args

}
