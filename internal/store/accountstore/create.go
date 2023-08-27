package accountstore

import (
	"context"
	"database/sql"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

func Insert(ctx context.Context, tx *sql.Tx, obj models.Account) error {

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO content.account (id, user_id, updated_at, created_at) VALUES($1,$2,$3,$4)")

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

	return nil

}
