package store

import (
	"context"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type InsertStore[E models.BaseModeler] struct {
	*Store
	insertQuery string
}

type Inserter[E models.BaseModeler] interface {
	Insert(ctx context.Context, obj E) error
}

func NewInsertStore[E models.BaseModeler](baseStorage *Store, insertQuery string) *InsertStore[E] {
	return &InsertStore[E]{baseStorage, insertQuery}
}

// Insert(ctx context.Context, obj E) error общий метод
// добавляет в хранилище новую строку
func (s *InsertStore[E]) Insert(ctx context.Context, obj E) error {
	logger.Log.Debug("Хранилище:" + obj.GetType() + ": Insert...")

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, s.insertQuery)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, obj.GetArgsInsert()...)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Debug("Хранилище:" + obj.GetType() + ": Insert - ok")

	return nil

}
