package withdrawalstore

import (
	"context"
	"time"

	"github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
)

type InsertStore struct {
	*store.Store
	insertQuery string
}

type Inserter interface {
	Insert(ctx context.Context, obj *models.Withdrawal) error
}

func NewInsertStore(baseStorage *store.Store, insertQuery string) *InsertStore {
	return &InsertStore{baseStorage, insertQuery}
}

// Insert(ctx context.Context, obj E) error общий метод
// добавляет в хранилище новую строку
func (s *InsertStore) Insert(ctx context.Context, obj *models.Withdrawal) error {
	logger.Log.Debug("Хранилище:" + obj.GetType() + ": Insert...")

	tx, err := s.GetDB().BeginTx(ctx, nil)
	var account models.Account

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	// Блокировка строк для определенного user_id
	row := tx.QueryRowContext(ctx, "SELECT id, current, withdrawn FROM content.account WHERE user_id = $1 FOR UPDATE ", obj.UserID)
	err = row.Scan(&account.ID, &account.Current, &account.Withdrawn)

	if err != nil {
		account = *models.NewAccount()
		account.UserID = obj.UserID

		if err := accountstore.Insert(ctx, tx, account); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	current := account.Current.Int64 - obj.Bonuses.Int64
	if current < 0 {
		err := errors.ErrorNotBonuses{}
		return &err
	}
	withdrawn := account.Withdrawn.Int64 + obj.Bonuses.Int64

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

	_, err = tx.ExecContext(ctx,
		"UPDATE content.account SET current=$1 , withdrawn=$2, updated_at=$3 WHERE id=$4",
		current, withdrawn, time.Now(), account.ID)

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
