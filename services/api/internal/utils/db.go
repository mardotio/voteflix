package utils

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"log"
)

func TxnRollback(tx *bun.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Println(err)
	}
}
