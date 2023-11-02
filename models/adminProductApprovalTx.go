package models

import (
	"log"
	"mf-bo-api/db"
)

func ApprovalAction(params map[string]string) error {
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// UpdProductRequest := `UPDATE ms_product_request SET rec_action = ?
	// WHERE rec_pk = ?`

	return nil
}
