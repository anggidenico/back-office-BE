package models

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/db"
	"strconv"
)

type ProductBankAccountRequest struct {
	RecPK                 *uint64 `db:"rec_pk" json:"rec_pk"`
	RecAction             *string `db:"rec_action" json:"rec_action"`
	ProductBankAccountKey *uint64 `db:"prod_bankacc_key" json:"prod_bankacc_key"`
	ProductKey            *uint64 `db:"product_key" json:"product_key"`
	BankKey               *uint64 `db:"bank_key" json:"bank_key"`
	BankAccountKey        *uint64 `db:"bank_account_key" json:"bank_account_key"`
	BankAccountPurpose    *uint64 `db:"bank_account_purpose" json:"bank_account_purpose"`
	SwiftCode             *string `db:"swift_code" json:"swift_code"`
	BankAccountName       *string `db:"bank_account_name" json:"bank_account_name"`
	AccountHolderName     *string `db:"account_holder_name" json:"account_holder_name"`
	// Foreign Key Value
	AccountNo              *string `db:"account_no" json:"account_no"`
	ProductName            *string `db:"product_name" json:"product_name"`
	BankName               *string `db:"bank_name" json:"bank_name"`
	BankAccountPurposeName *string `db:"bank_account_purpose_name" json:"bank_account_purpose_name"`
	// BankAccountName        *string `db:"bank_account_name" json:"bank_account_name"`
}

type ProductBankAccountDetail struct {
	Updates  ProductBankAccountRequest `json:"updates"`
	Existing ProductBankAccountRequest `json:"existing"`
}

func ProductBankAccountRequestList() []ProductBankAccountRequest {
	query := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no, t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.bank_key, t1.prod_bankacc_key, t4.swift_code, t1.bank_account_name, t4.account_holder_name

	FROM ms_product_bank_account_request t1
	INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
	INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
	INNER JOIN (
		SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code
		FROM ms_bank_account a1
		INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key
		WHERE a1.rec_status = 1
	) t4 ON t1.bank_account_key = t4.bank_account_key
	
	WHERE t1.rec_status = 1 AND t1.rec_approval_status IS NULL`

	var get []ProductBankAccountRequest
	err := db.Db.Select(&get, query)
	if err != nil {
		log.Println(err.Error())
	}

	return get

}

func ProductBankAccountRequestDetail(rec_pk string) ProductBankAccountDetail {
	var result ProductBankAccountDetail

	queryGetAction := `SELECT t1.rec_action FROM ms_product_bank_account_request t1 WHERE t1.rec_pk =` + rec_pk

	// var getUpdates ProductBankAccountRequest
	var getAction string
	err := db.Db.Get(&getAction, queryGetAction)
	if err != nil {
		log.Println(err.Error())
	}

	if getAction == "CREATE" {
		queryGetUpdate := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no, t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.bank_key, t1.prod_bankacc_key, t4.swift_code, t1.bank_account_name, t4.account_holder_name 
		FROM ms_product_bank_account_request t1
		INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
		INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
		INNER JOIN (
			SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code
			FROM ms_bank_account a1
			INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key
			WHERE a1.rec_status = 1
		) t4 ON t1.bank_account_key = t4.bank_account_key
		WHERE t1.rec_pk = ` + rec_pk

		var getUpdates ProductBankAccountRequest
		err = db.Db.Get(&getUpdates, queryGetUpdate)
		if err != nil {
			log.Println(err.Error())
		}

		result.Updates = getUpdates
	}

	if getAction == "UPDATE" {

		queryGetUpdate := `SELECT t1.rec_pk, t1.rec_action, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no, t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.bank_key, t1.prod_bankacc_key, t4.swift_code, t1.bank_account_name, t4.account_holder_name 
		FROM ms_product_bank_account_request t1
		INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
		INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
		INNER JOIN (
			SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code
			FROM ms_bank_account_request a1
			INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key
			WHERE a1.rec_status = 1
		) t4 ON t1.bank_account_key = t4.bank_account_key
		WHERE t1.rec_pk = ` + rec_pk
		// log.Println("queryGetUpdate", queryGetUpdate)
		var getUpdates ProductBankAccountRequest
		err = db.Db.Get(&getUpdates, queryGetUpdate)
		if err != nil {
			log.Println(err.Error())
		}

		result.Updates = getUpdates

		queryGetExisting := `SELECT t1.prod_bankacc_key, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no, t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.bank_key, t1.prod_bankacc_key, t4.swift_code, t1.bank_account_name, t4.account_holder_name 
		FROM ms_product_bank_account t1
		INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
		INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
		INNER JOIN (
			SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code 
			FROM ms_bank_account a1 
			INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key WHERE a1.rec_status = 1
		) t4 ON t1.bank_account_key = t4.bank_account_key		
		WHERE t1.prod_bankacc_key =` + strconv.FormatUint(*getUpdates.ProductBankAccountKey, 10)
		// log.Println("queryGetExisting", queryGetExisting)
		var getExisting ProductBankAccountRequest
		err := db.Db.Get(&getExisting, queryGetExisting)
		if err != nil {
			log.Println(err.Error())
		}

		result.Existing = getExisting
	}

	return result

}

func CreateRequestProductBankAccount(paramsMsBankAccount map[string]string, paramsProductBankAccount map[string]string) error {

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	var ret sql.Result
	var lastID int64

	if _, ok := paramsMsBankAccount["bank_account_key"]; !ok {
		qInsertMsBankAccount := "INSERT INTO ms_bank_account"
		var fields, values string
		var bindvars []interface{}
		for key, value := range paramsMsBankAccount {
			if key != "rec_pk" && key != "rec_action" {
				fields += key + ", "
				values += "?, "
				bindvars = append(bindvars, value)
			}
		}
		fields = fields[:(len(fields) - 2)]
		values = values[:(len(values) - 2)]

		qInsertMsBankAccount += "(" + fields + ") VALUES(" + values + ")"
		// log.Println(qInsertMsBankAccount)
		// qInsertMsBankAccountReq += "(" + fields + ") VALUES(" + values + ")"

		ret, err = tx.Exec(qInsertMsBankAccount, bindvars...)
		// tx.Commit()
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err
		}

		lastID, err = ret.LastInsertId()
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err
		}

	} else {
		// log.Println("paramsMsBankAccount[bank_account_key] =", paramsMsBankAccount["bank_account_key"])
		lastID, _ = strconv.ParseInt(paramsMsBankAccount["bank_account_key"], 10, 64)
	}

	paramsMsBankAccount2 := paramsMsBankAccount
	paramsMsBankAccount2["bank_account_key"] = strconv.FormatInt(lastID, 10)
	var fields1, values1 string
	var bindvars1 []interface{}
	qInsertMsBankAccountReq := "INSERT INTO ms_bank_account_request"
	for key, value := range paramsMsBankAccount2 {
		fields1 += key + ", "
		values1 += "?, "
		bindvars1 = append(bindvars1, value)
	}
	fields1 = fields1[:(len(fields1) - 2)]
	values1 = values1[:(len(values1) - 2)]
	qInsertMsBankAccountReq += "(" + fields1 + ") VALUES(" + values1 + ")"
	// log.Println(qInsertMsBankAccountReq)
	_, err = tx.Exec(qInsertMsBankAccountReq, bindvars1...)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// if lastID > 0 {
	// msBankAccKey := strconv.FormatInt(lastID, 10)
	paramsProductBankAccount["bank_account_key"] = strconv.FormatInt(lastID, 10)
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
	// log.Println(qInsertMsBankAccount)
	_, err = tx.Exec(qInsertMsBankAccount, bindvars2...)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// }

	tx.Commit()

	return nil
}

func ProductBankAccountApprovalAction(params map[string]string) error {
	var resultSQL sql.Result
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	recBy := params["rec_by"]
	recDate := params["rec_date"]

	queryGetAction := `SELECT t1.rec_action FROM ms_product_bank_account_request t1 WHERE t1.rec_pk = ` + params["rec_pk"]
	var getAction string
	err = db.Db.Get(&getAction, queryGetAction)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	qGetBankAccKey := `SELECT bank_account_key FROM ms_product_bank_account_request WHERE rec_pk = ` + params["rec_pk"]
	var getBankAccKey string
	err = db.Db.Get(&getBankAccKey, qGetBankAccKey)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	qGetLastMsBankAccReq := `SELECT rec_pk FROM ms_bank_account_request 
	WHERE bank_account_key = ` + getBankAccKey + ` ORDER BY rec_pk DESC LIMIT 1`
	log.Println(qGetLastMsBankAccReq)
	var getLastMsBankAccReq string
	err = db.Db.Get(&getLastMsBankAccReq, qGetLastMsBankAccReq)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	query1 := `UPDATE ms_bank_account_request SET rec_approval_status = ` + params["approval"] + ` , rec_approved_date = '` + recDate + `' , rec_approved_by = '` + recBy + `' , rec_attribute_id1 = '` + params["reason"] + `' WHERE rec_pk = ` + getLastMsBankAccReq
	resultSQL, err = tx.Exec(query1)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}
	rowsAffected, err := resultSQL.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		err = errors.New("No Rows Affected")
		return err
	}

	query2 := `UPDATE ms_product_bank_account_request SET rec_approval_status = ` + params["approval"] + ` , rec_approved_date = '` + recDate + `' , rec_approved_by = '` + recBy + `' , rec_attribute_id1 = '` + params["reason"] + `' WHERE rec_pk = ` + params["rec_pk"]
	resultSQL, err = tx.Exec(query2)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}
	rowsAffected, err = resultSQL.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		err = errors.New("No Rows Affected")
		return err
	}

	if params["approval"] == "1" {

		if getAction == "CREATE" {

			qGetProdBankAccReq := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t4.bank_key, t1.bank_account_key, t4.bank_name, t4.account_no , t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.swift_code, t1.bank_account_name, t4.account_holder_name
			FROM ms_product_bank_account_request t1
			INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
			INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
			INNER JOIN (
				SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code FROM ms_bank_account a1 INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key WHERE a1.rec_status = 1
			) t4 ON t1.bank_account_key = t4.bank_account_key
			WHERE t1.rec_pk =` + params["rec_pk"]

			var getReqProductBankAcc ProductBankAccountRequest
			err = db.Db.Get(&getReqProductBankAcc, qGetProdBankAccReq)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}

			inputProductBankAccount := make(map[string]string)
			inputProductBankAccount["rec_status"] = "1"
			inputProductBankAccount["rec_created_by"] = recBy
			inputProductBankAccount["rec_created_date"] = recDate
			inputProductBankAccount["product_key"] = strconv.FormatUint(*getReqProductBankAcc.ProductKey, 10)
			inputProductBankAccount["bank_account_key"] = strconv.FormatUint(*getReqProductBankAcc.BankAccountKey, 10)
			inputProductBankAccount["bank_account_name"] = *getReqProductBankAcc.BankAccountPurposeName + " " + *getReqProductBankAcc.ProductName
			inputProductBankAccount["bank_account_purpose"] = strconv.FormatUint(*getReqProductBankAcc.BankAccountPurpose, 10)

			var fields, values string
			var bindvars []interface{}
			for key, value := range inputProductBankAccount {
				if key != "rec_pk" && key != "rec_action" {
					fields += key + ", "
					values += ` "` + value + `", `
					bindvars = append(bindvars, value)
				}
			}
			fields = fields[:(len(fields) - 2)]
			values = values[:(len(values) - 2)]
			query := `INSERT INTO ms_product_bank_account (` + fields + `) VALUES(` + values + `)`

			resultSQL, err = tx.Exec(query)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}

		}

		if getAction == "UPDATE" {

			qGetProdBankAccReq := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t4.bank_key, t1.bank_account_key, t4.bank_name, t4.account_no , t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name, t4.swift_code, t1.bank_account_name, t4.account_holder_name
			FROM ms_product_bank_account_request t1
			INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
			INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
			INNER JOIN (
				SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name, a1.swift_code 
				FROM ms_bank_account_request a1 
				INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key WHERE a1.rec_status = 1
			) t4 ON t1.bank_account_key = t4.bank_account_key
			WHERE t1.rec_pk = ` + params["rec_pk"]

			var getReqProductBankAcc ProductBankAccountRequest
			err = db.Db.Get(&getReqProductBankAcc, qGetProdBankAccReq)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}

			updMsBankAccount := make(map[string]string)
			updMsBankAccount["rec_modified_by"] = recBy
			updMsBankAccount["rec_modified_date"] = recDate
			updMsBankAccount["bank_account_key"] = strconv.FormatUint(*getReqProductBankAcc.BankAccountKey, 10)
			updMsBankAccount["account_no"] = *getReqProductBankAcc.AccountNo
			updMsBankAccount["bank_key"] = strconv.FormatUint(*getReqProductBankAcc.BankKey, 10)
			updMsBankAccount["swift_code"] = *getReqProductBankAcc.SwiftCode
			updMsBankAccount["account_holder_name"] = *getReqProductBankAcc.AccountHolderName
			queryMsBankAccount := `UPDATE ms_bank_account SET `
			i := 0
			for key, value := range updMsBankAccount {
				if key != "bank_account_key" {
					queryMsBankAccount += key + " = '" + value + "'"
					if (len(updMsBankAccount) - 2) > i {
						queryMsBankAccount += ", "
					}
					i++
				}
			}
			queryMsBankAccount += " WHERE bank_account_key = " + strconv.FormatUint(*getReqProductBankAcc.BankAccountKey, 10)
			log.Println(queryMsBankAccount)
			resultSQL, err = tx.Exec(queryMsBankAccount)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}
			rowsAffected, err := resultSQL.RowsAffected()
			if rowsAffected == 0 {
				tx.Rollback()
				err = errors.New("No Rows Affected UPDATE ms_bank_account")
				return err
			}

			updProductBankAccount := make(map[string]string)
			updProductBankAccount["rec_modified_by"] = recBy
			updProductBankAccount["rec_modified_date"] = recDate
			updProductBankAccount["prod_bankacc_key"] = strconv.FormatUint(*getReqProductBankAcc.ProductBankAccountKey, 10)
			updProductBankAccount["product_key"] = strconv.FormatUint(*getReqProductBankAcc.ProductKey, 10)
			updProductBankAccount["bank_account_key"] = strconv.FormatUint(*getReqProductBankAcc.BankAccountKey, 10)
			updProductBankAccount["bank_account_name"] = *getReqProductBankAcc.BankAccountName
			updProductBankAccount["bank_account_purpose"] = strconv.FormatUint(*getReqProductBankAcc.BankAccountPurpose, 10)

			queryProdBankAcc := "UPDATE ms_product_bank_account SET "
			i = 0
			for key, value := range updProductBankAccount {
				if key != "prod_bankacc_key" {
					queryProdBankAcc += key + " = '" + value + "'"
					if (len(updProductBankAccount) - 2) > i {
						queryProdBankAcc += ", "
					}
					i++
				}
			}
			queryProdBankAcc += " WHERE prod_bankacc_key = " + updProductBankAccount["prod_bankacc_key"]
			log.Println(queryProdBankAcc)
			resultSQL, err = tx.Exec(queryProdBankAcc)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}
			rowsAffected, err = resultSQL.RowsAffected()
			if rowsAffected == 0 {
				tx.Rollback()
				err = errors.New("No Rows Affected UPDATE ms_product_bank_account")
				return err
			}
		}
	}

	tx.Commit()
	return nil
}
