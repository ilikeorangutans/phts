package db

import (
	"database/sql"
	"fmt"
)

func checkResult(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("No row updated")
	}
	return nil
}
