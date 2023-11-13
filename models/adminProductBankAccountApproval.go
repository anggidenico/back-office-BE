package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"strconv"
)

func CreateRequestProductBankAccount(paramsMsBankAccount map[string]string, paramsProductBankAccount map[string]string) error {

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	qInsertMsBankAccount := "INSERT INTO ms_bank_account"
	var fields, values string
	var bindvars []interface{}
	for key, value := range paramsMsBankAccount {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	qInsertMsBankAccount += "(" + fields + ") VALUES(" + values + ")"

	var ret sql.Result
	ret, err = tx.Exec(qInsertMsBankAccount, bindvars...)
	tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}
	lastID, err := ret.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	if lastID > 0 {
		msBankAccKey := strconv.FormatInt(lastID, 10)
		paramsProductBankAccount["bank_account_key"] = msBankAccKey
		qInsertMsBankAccount := "INSERT INTO ms_product_bank_account_request"
		var fields2, values2 string
		var bindvars2 []interface{}
		for key2, value2 := range paramsProductBankAccount {
			fields2 += key2 + ", "
			values2 += "?, "
			bindvars2 = append(bindvars2, value2)
		}
		fields2 = fields2[:(len(fields2) - 2)]
		values2 = values2[:(len(values2) - 2)]
		qInsertMsBankAccount += "(" + fields2 + ") VALUES(" + values2 + ")"

		_, err = tx.Exec(qInsertMsBankAccount, bindvars2...)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err
		}

	}

	return nil
}
