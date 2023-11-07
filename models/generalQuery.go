package models

import (
	"log"
	"mf-bo-api/db"
	"strconv"
)

func GetForeignKeyValue(TableName string, FieldName string, FieldForeignKey string, ForeignKey uint64) string {
	FK := strconv.FormatUint(ForeignKey, 10)
	
	query := `SELECT ` + FieldName + ` FROM ` + TableName + ` WHERE ` + FieldForeignKey + ` = ` + FK
	var result string
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result

}
