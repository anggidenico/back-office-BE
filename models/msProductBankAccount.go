package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"
)

type MsProductBankAccountInfo struct {
	ProdBankaccKey     uint64      `json:"prod_bankacc_key"`
	BankAccountName    string      `json:"bank_account_name"`
	BankAccountPurpose uint64      `json:"bank_account_purpose"`
	BankAccount        BankAccount `json:"bank_account"`
}

type MsProductBankAccount struct {
	ProdBankaccKey     uint64  `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	ProductKey         *uint64 `db:"product_key"           json:"product_key"`
	BankAccountKey     *uint64 `db:"bank_account_key"      json:"bank_account_key"`
	BankAccountName    string  `db:"bank_account_name"     json:"bank_account_name"`
	BankAccountPurpose uint64  `db:"bank_account_purpose"  json:"bank_account_purpose"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type AdminMsProductBankAccountList struct {
	ProdBankaccKey     uint64  `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	ProductKey         *uint64 `db:"product_key"           json:"product_key"`
	ProductCode        string  `db:"product_code"          json:"product_code"`
	ProductNameAlt     string  `db:"product_name_alt"      json:"product_name_alt"`
	BankAccountName    string  `db:"bank_account_name"     json:"bank_account_name"`
	BankAccountPurpose *string `db:"bank_account_purpose"  json:"bank_account_purpose"`
	BankFullname       *string `db:"bank_fullname"         json:"bank_fullname"`
	AccountNo          string  `db:"account_no"            json:"account_no"`
	AccountHolderName  string  `db:"account_holder_name"   json:"account_holder_name"`
	StatusUpdate       *bool   `db:"status_update" json:"status_update"`
}

type MsProductBankAccountDetailAdmin struct {
	ProdBankaccKey     uint64          `json:"prod_bankacc_key"`
	Product            *MsProductInfo  `json:"product"`
	Bank               *MsBankList     `json:"bank"`
	AccountNo          string          `json:"account_no"`
	AccountHolderName  string          `json:"account_holder_name"`
	BranchName         *string         `json:"branch_name"`
	Currency           *MsCurrencyInfo `json:"currency"`
	BankAccountType    LookupTrans     `json:"bank_account_type"`
	SwiftCode          *string         `json:"swift_code"`
	BankAccountName    string          `json:"bank_account_name"`
	BankAccountPurpose LookupTrans     `json:"bank_account_purpose"`
}

type MsProductBankAccountTransactionInfo struct {
	ProdBankaccKey uint64 `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	BankName       string `db:"bank_name"             json:"bank_name"`
	AccountNo      string `db:"account_no"            json:"account_no"`
	AccountName    string `db:"account_name"          json:"account_name"`
}

func GetAllMsProductBankAccount(c *[]MsProductBankAccount, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_bank_account.* FROM 
			  ms_product_bank_account WHERE  
			  ms_product_bank_account.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_bank_account."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ProductBankAccountStatusUpdate(prod_bankacc_key string) bool {
	query := `SELECT count(*) FROM ms_product_bank_account_request 
	WHERE rec_status = 1 AND rec_approval_status IS NULL 
	AND prod_bankacc_key = ` + prod_bankacc_key
	// log.Println(query)
	var count uint64
	err := db.Db.Get(&count, query)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	var result bool
	if count > 0 { // kalau ada request gantung maka false
		result = false
	} else {
		result = true
	}
	return result
}

func AdminGetAllMsProductBankAccount(c *[]AdminMsProductBankAccountList, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT 
				pba.prod_bankacc_key AS prod_bankacc_key,
				p.product_key AS product_key,
				p.product_code AS product_code,
				p.product_name_alt AS product_name_alt,
				pba.bank_account_name AS bank_account_name,
				bank_account_purpose.lkp_name AS bank_account_purpose,
				bank.bank_fullname AS bank_fullname,
				ba.account_no AS account_no,
				ba.account_holder_name AS account_holder_name
			FROM ms_product_bank_account pba
			INNER JOIN ms_product p ON pba.product_key = p.product_key
			INNER JOIN ms_bank_account ba ON pba.bank_account_key = ba.bank_account_key
			LEFT JOIN gen_lookup AS bank_account_purpose ON bank_account_purpose.lookup_key = pba.bank_account_purpose
			LEFT JOIN ms_bank AS bank ON bank.bank_key = ba.bank_key
			WHERE pba.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pba.prod_bankacc_key LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name_alt LIKE '%" + *searchLike + "%' OR"
		condition += " pba.bank_account_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank_account_purpose.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank.bank_fullname LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_no LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_holder_name LIKE '%" + *searchLike + "%')"
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountDataGetAllMsProductBankAccount(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT 
				count(pba.prod_bankacc_key) AS count_data 
			FROM ms_product_bank_account pba 
			INNER JOIN ms_product p ON pba.product_key = p.product_key 
			INNER JOIN ms_bank_account ba ON pba.bank_account_key = ba.bank_account_key 
			LEFT JOIN gen_lookup AS bank_account_purpose ON bank_account_purpose.lookup_key = pba.bank_account_purpose 
			LEFT JOIN ms_bank AS bank ON bank.bank_key = ba.bank_key 
			WHERE pba.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pba.prod_bankacc_key LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name_alt LIKE '%" + *searchLike + "%' OR"
		condition += " pba.bank_account_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank_account_purpose.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank.bank_fullname LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_no LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_holder_name LIKE '%" + *searchLike + "%')"
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsProductBankAccount(c *MsProductBankAccount, key string) (int, error) {
	query := `SELECT ms_product_bank_account.* FROM ms_product_bank_account WHERE ms_product_bank_account.rec_status = 1 AND ms_product_bank_account.prod_bankacc_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateMsProductBankAccount(params map[string]string) (int, error) {
	query := "UPDATE ms_product_bank_account SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "prod_bankacc_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE prod_bankacc_key = " + params["prod_bankacc_key"]
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

func CreateMsProductBankAccount(params map[string]string) (int, error) {
	query := "INSERT INTO ms_product_bank_account"
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

func GetAllMsProductBankAccountTransaction(c *[]MsProductBankAccountTransactionInfo, productKey string, transType string) (int, error) {
	query2 := `SELECT 
				ba.prod_bankacc_key AS prod_bankacc_key,
				bank.bank_fullname AS bank_name, 
				b.account_no AS account_no, 
				b.account_holder_name AS account_name 
			FROM ms_product_bank_account AS ba 
			INNER JOIN ms_bank_account AS b ON b.bank_account_key = ba.bank_account_key
			INNER JOIN ms_bank AS bank ON bank.bank_key = b.bank_key`
	query := query2 + " WHERE ba.rec_status = 1 AND ba.product_key = '" + productKey + "' AND ba.bank_account_purpose = '" + transType + "'"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllMsProductBankAccountTransactionAutoInvest(c *MsProductBankAccountTransactionInfo, productKey string) (int, error) {
	query := `SELECT 
				ba.prod_bankacc_key AS prod_bankacc_key,
				bank.bank_fullname AS bank_name, 
				b.account_no AS account_no, 
				b.account_holder_name AS account_name 
			FROM ms_product_bank_account AS ba 
			INNER JOIN ms_bank_account AS b ON b.bank_account_key = ba.bank_account_key
			INNER JOIN ms_bank AS bank ON bank.bank_key = b.bank_key`
	query += " WHERE ba.rec_status = 1 AND ba.product_key = '" + productKey + "' AND ba.bank_account_purpose = 271"
	query += " ORDER BY ba.prod_bankacc_key ASC LIMIT 1"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsProductBankAccountTransactionByKey(c *MsProductBankAccountTransactionInfo, prodBankaccKey string) (int, error) {
	query2 := `SELECT 
				ba.prod_bankacc_key AS prod_bankacc_key,
				bank.bank_fullname AS bank_name, 
				b.account_no AS account_no, 
				b.account_holder_name AS account_name 
			FROM ms_product_bank_account AS ba 
			INNER JOIN ms_bank_account AS b ON b.bank_account_key = ba.bank_account_key
			INNER JOIN ms_bank AS bank ON bank.bank_key = b.bank_key`
	query := query2 + " WHERE ba.rec_status = 1 AND ba.prod_bankacc_key = '" + prodBankaccKey + "'"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllMsProductBankAccountOrderByBank(c *[]MsProductBankAccount, params map[string]string) (int, error) {
	query := `SELECT
              mp.* FROM 
			  ms_product_bank_account as mp 
			  left join ms_bank_account as mb on mb.bank_account_key = mp.bank_account_key
			  left join ms_bank as b on b.bank_key = mb.bank_key 
			  WHERE mp.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
