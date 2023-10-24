package models

import (
	"log"
	"mf-bo-api/db"
)

func CreateProductUpdateRequest(params map[string]string) error {
	query := "INSERT INTO ms_product_request"
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += ` "` + value + `", `
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	// log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
