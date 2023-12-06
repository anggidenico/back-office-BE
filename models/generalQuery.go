package models

import (
	"log"
	"mf-bo-api/db"
	"strconv"
	"strings"
)

func GenerateSelectQuery(tableName string, fieldsToShow []string, whereClause map[string]string) string {
	var fields string

	if len(fieldsToShow) > 1 {
		fields = strings.Join(fieldsToShow, ", ")
	} else {
		fields = `*`
	}

	query := `SELECT ` + fields + ` FROM ` + tableName

	if len(whereClause) > 0 {
		query += ` WHERE `
		i := 0
		for fieldName, fieldValue := range whereClause {
			query += fieldName + ` = '` + fieldValue + `' `
			if i < len(whereClause)-1 {
				query += ` AND `
			}
		}
	}

	return query
}

func GenerateInsertQuery(tableName string, params map[string]string) string {
	query := "INSERT INTO " + tableName
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + `, `
		if value == "" {
			values += ` NULL, `
		} else {
			values += ` '` + value + `', `
		}
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	log.Println("GenerateInsertQuery", query)

	return query
}

func GenerateUpdateQuery(tableName string, primaryKeyField string, params map[string]string) string {
	query := `UPDATE ` + tableName + ` SET `
	i := 0
	for key, value := range params {
		if key != primaryKeyField {
			if value == "" {
				query += key + " = NULL"
			} else {
				query += key + " = '" + value + "'"
			}
			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += ` WHERE ` + primaryKeyField + ` = ` + params[primaryKeyField]

	log.Println("GenerateUpdateQuery", query)

	return query
}

func GetForeignKeyValue(TableName string, FieldName string, FieldForeignKey string, ForeignKey uint64) string {
	FK := strconv.FormatUint(ForeignKey, 10)

	query := `SELECT ` + FieldName + ` FROM ` + TableName + ` WHERE ` + FieldForeignKey + ` = ` + FK
	var result string
	// log.Println(query)
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result

}
