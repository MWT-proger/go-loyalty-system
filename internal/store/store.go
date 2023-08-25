package store

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

type Store struct {
	db *sqlx.DB
}

// Init(ctx context.Context) error - вызывается при запуске программы,
// инициализирует соединение
// и возвращает ошибку в случае не удачи
func (s *Store) Init(ctx context.Context) error {
	logger.Log.Info("Хранилище: Подключение...")

	var (
		conf    = configs.GetConfig()
		db, err = sqlx.Open("pgx", conf.DatabaseDSN)
	)

	if err != nil {
		return err
	}

	s.db = db

	if err := s.ping(); err != nil {
		return err
	}

	if err := s.migration(); err != nil {
		return err
	}
	logger.Log.Info("Хранилище: Соединение установленно")

	return nil

}

func (s *Store) GetDB() *sqlx.DB {
	return s.db
}

// Close() error - вызывается при завершение программы,
// закрывает соединение и возвращает ошибку в случае не удачи
func (s *Store) Close() error {
	logger.Log.Info("Хранилище: Закрытие соединения...")

	if err := s.db.Close(); err != nil {
		return err
	}
	logger.Log.Info("Хранилище: Соединение успешно закрыто")

	return nil
}
