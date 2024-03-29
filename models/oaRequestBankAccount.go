package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OaRequestBankAccount struct {
	ReqBankaccKey     uint64  `db:"req_bankacc_key"        json:"req_bankacc_key"`
	OaRequestKey      uint64  `db:"oa_request_key"         json:"oa_request_key"`
	BankAccountKey    uint64  `db:"bank_account_key"       json:"bank_account_key"`
	FlagPriority      uint8   `db:"flag_priority"          json:"flag_priority"`
	BankAccountName   string  `db:"bank_account_name"      json:"bank_account_name"`
	RecOrder          *uint64 `db:"rec_order"              json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"             json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"             json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func CreateOaRequestBankAccount(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_request_bank_account"
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
		// log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func UpdateOaRequestBankAccount(params map[string]string) (int, error) {
	query := "UPDATE oa_request_bank_account SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "req_bankacc_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE req_bankacc_key = " + params["req_bankacc_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
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
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

type OaRequestByField struct {
	RecBankaccKey     uint64  `db:"req_bankacc_key"           json:"req_bankacc_key"`
	OaRequestKey      uint64  `db:"oa_request_key"            json:"oa_request_key"`
	BankAccountKey    uint64  `db:"bank_account_key"          json:"bank_account_key"`
	FlagPriority      uint64  `db:"flag_priority"             json:"flag_priority"`
	BankKey           uint64  `db:"bank_key"                  json:"bank_key"`
	BankFullname      *string `db:"bank_fullname"             json:"bank_fullname"`
	AccountNo         *string `db:"account_no"                json:"account_no"`
	AccountHolderName *string `db:"account_holder_name"       json:"account_holder_name"`
	BranchName        *string `db:"branch_name"               json:"branch_name"`
	CurrencyKey       *uint64 `db:"currency_key"              json:"currency_key,omitempty"`
	CurrencyName      *string `db:"currency_name"             json:"currency_name,omitempty"`
}

func GetOaRequestBankByField(c *[]OaRequestByField, field string, value string) (int, error) {
	query := `SELECT 
				br.req_bankacc_key AS req_bankacc_key,
				br.oa_request_key AS oa_request_key,
				br.bank_account_key AS bank_account_key,
				br.flag_priority AS flag_priority,
				ba.bank_key AS bank_key,
				b.bank_fullname AS bank_fullname,
				ba.account_no AS account_no,
				ba.account_holder_name AS account_holder_name,
				ba.branch_name AS branch_name,
				ba.currency_key AS currency_key, 
				cur.code AS currency_name 
			FROM oa_request_bank_account AS br 
			INNER JOIN ms_bank_account AS ba ON ba.bank_account_key = br.bank_account_key 
			INNER JOIN ms_bank AS b ON b.bank_key = ba.bank_key 
			LEFT JOIN ms_currency AS cur on cur.currency_key = ba.currency_key 
			WHERE br.rec_status = 1 AND ba.rec_status = 1 
			AND br.` + field + ` = '` + value + `'
			ORDER BY br.flag_priority DESC`

	// Main query
	// // log.Infoln("==================== query get bank account ====================")
	// log.Println(query)
	// // log.Infoln("==========================================================")
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMultipleOaRequestBankAccount(params []interface{}) (int, error) {

	q := `INSERT INTO oa_request_bank_account (
		oa_request_key, 
		bank_account_key,
		flag_priority,
		bank_account_name,
		rec_status,
		rec_created_date,
		rec_created_by) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	// log.Println("==========  ==========>>>", q)
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		// log.Error(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetOaRequestBankAccount(c *OaRequestBankAccount, key string, field string) (int, error) {
	query := `SELECT * 
				FROM oa_request_bank_account 
				WHERE rec_status = 1 
				AND ` + field + ` = "` + key + `"`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func DeleteOaRequestBankAccount(params map[string]string, bankKey []string, requestKey string) (int, error) {
	inQuery := strings.Join(bankKey, ",")
	query := "UPDATE oa_request_bank_account SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	if len(bankKey) > 0 {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "' AND req_bankacc_key NOT IN(" + inQuery + ")"
	} else {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "'"
	}
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
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
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

type BankAccountPriority struct {
	BankCode          string  `db:"bank_code"                json:"bank_code"`
	BankName          string  `db:"bank_name"                json:"bank_name"`
	Code              string  `db:"code"                     json:"code"`
	AccountNo         string  `db:"account_no"               json:"account_no"`
	AccountHolderName string  `db:"account_holder_name"      json:"account_holder_name"`
	BranchName        *string `db:"branch_name"              json:"branch_name"`
	SwiftCode         string  `db:"swift_code" json:"swift_code"`
}

func GetBankAccountPriority(c *BankAccountPriority, oaRequest string) (int, error) {
	query := `SELECT 
				b.bank_code,
				b.bank_name,
				cur.code,
				mba.account_no,
				mba.account_holder_name,
				mba.branch_name,
				b.swift_code 
			FROM oa_request_bank_account AS orba 
			INNER JOIN ms_bank_account AS mba ON mba.bank_account_key = orba.bank_account_key 
			INNER JOIN ms_bank AS b ON b.bank_key = mba.bank_key 
			INNER JOIN ms_currency AS cur ON cur.currency_key = mba.currency_key 
			WHERE orba.rec_status = 1 AND mba.rec_status = 1 
			AND orba.oa_request_key = "` + oaRequest + `" ORDER BY orba.flag_priority DESC LIMIT 1`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
