package store

import (
	"embed"

	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

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
