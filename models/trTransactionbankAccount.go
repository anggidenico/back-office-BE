package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type TrTransactionBankAccount struct {
	TransBankaccKey   uint64  `db:"trans_bankacc_key"     json:"trans_bankacc_key"`
	TransactionKey    uint64  `db:"transaction_key"       json:"transaction_key"`
	ProdBankaccKey    uint64  `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	CustBankaccKey    uint64  `db:"cust_bankacc_key"      json:"cust_bankacc_key"`
	RecOrder          *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type TransactionPoductBankAccount struct {
	TransactionKey    uint64  `db:"transaction_key"         json:"transaction_key"`
	AccountNo         *string `db:"account_no"              json:"account_no"`
	AccountHolderName *string `db:"account_holder_name"     json:"account_holder_name"`
	BranchName        *string `db:"branch_name"             json:"branch_name"`
	BankName          *string `db:"bank_name"               json:"bank_name"`
}

func CreateTrTransactionBankAccount(params map[string]string) (int, error) {
	query := "INSERT INTO tr_transaction_bank_account"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateTrTransactionBankAccount(params map[string]string, value string, field string) (int, error) {
	query := "UPDATE tr_transaction_bank_account SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " = " + value
	// log.Println("========== UpdateTrTransactionBankAccount ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetTrTransactionBankAccountByField(c *TrTransactionBankAccount, value string, field string) (int, error) {
	query := `SELECT
              tr_transaction_bank_account.* FROM 
			  tr_transaction_bank_account 
			  where tr_transaction_bank_account.rec_status = 1 and tr_transaction_bank_account.` + field + ` = ` + value + ` limit 1`
	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionBankAccountIn(c *[]TransactionPoductBankAccount, value []string) (int, error) {
	inQuery := strings.Join(value, ",")
	query := `SELECT
	t.transaction_key,
	a.account_no,
	a.account_holder_name,
	a.branch_name,
	b.bank_name 
	FROM tr_transaction_bank_account t
	LEFT JOIN ms_product_bank_account p ON t.prod_bankacc_key = p.prod_bankacc_key
	LEFT JOIN ms_bank_account a ON p.bank_account_key = a.bank_account_key 
	LEFT JOIN ms_bank b ON a.bank_key = b.bank_key
	WHERE t.transaction_key IN(` + inQuery + `)`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
