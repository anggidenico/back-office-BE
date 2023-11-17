package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"strconv"
)

type ProductBankAccountRequest struct {
	RecPK                 uint64  `db:"rec_pk" json:"rec_pk"`
	RecAction             string  `db:"rec_action" json:"rec_action"`
	ProductBankAccountKey *uint64 `db:"product_bankacc_key" json:"product_bankacc_key"`
	ProductKey            *uint64 `db:"product_key" json:"product_key"`
	BankAccountKey        *uint64 `db:"bank_account_key" json:"bank_account_key"`
	BankAccountPurpose    *uint64 `db:"bank_account_purpose" json:"bank_account_purpose"`
	// Foreign Key Value
	AccountNo              *string `db:"account_no" json:"account_no"`
	ProductName            *string `db:"product_name" json:"product_name"`
	BankName               *string `db:"bank_name" json:"bank_name"`
	BankAccountPurposeName *string `db:"bank_account_purpose_name" json:"bank_account_purpose_name"`
	// BankAccountName        *string `db:"bank_account_name" json:"bank_account_name"`
}

func ProductBankAccountRequestList() []ProductBankAccountRequest {
	query := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no, t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name 

	FROM ms_product_bank_account_request t1
	INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
	INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
	INNER JOIN (
		SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name
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

func ProductBankAccountRequestDetail(rec_pk string) ProductBankAccountRequest {
	query := `SELECT t1.rec_pk, t1.rec_action, t1.prod_bankacc_key, t1.product_key, t3.product_name, t1.bank_account_key, t4.bank_name, t4.account_no , t1.bank_account_purpose, t2.lkp_name bank_account_purpose_name 

	FROM ms_product_bank_account t1
	INNER JOIN gen_lookup t2 ON t2.lookup_key = t1.bank_account_purpose
	INNER JOIN ms_product t3 ON t3.product_key = t1.product_key
	INNER JOIN (
		SELECT a1.bank_account_key, a1.bank_account_type, a1.bank_key, a2.bank_name, a1.account_no, a1.account_holder_name
		FROM ms_bank_account a1
		INNER JOIN ms_bank a2 ON a2.bank_key = a1.bank_key
		WHERE a1.rec_status = 1
	) t4 ON t1.bank_account_key = t4.bank_account_key
	
	WHERE t1.rec_pk =` + rec_pk

	var get ProductBankAccountRequest
	err := db.Db.Get(&get, query)
	if err != nil {
		log.Println(err.Error())
	}

	return get

}

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
