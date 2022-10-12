package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type TrAutoinvestRegistration struct {
	AutoinvestKey       uint64           `db:"autoinvest_key"            json:"autoinvest_key"`
	AccKey              uint64           `db:"acc_key"                   json:"acc_key"`
	InvestAmount        *decimal.Decimal `db:"invest_amount"             json:"invest_amount"`
	InvestFeeRates      *decimal.Decimal `db:"invest_fee_rates"          json:"invest_fee_rates"`
	InvestFeeAmount     *decimal.Decimal `db:"invest_fee_amount"         json:"invest_fee_amount"`
	InvestFeeCharges    *decimal.Decimal `db:"invest_fee_charges"        json:"invest_fee_charges"`
	InvestDateExecute   uint8            `db:"invest_date_execute"       json:"invest_date_execute"`
	RecuringMode        *uint64          `db:"recuring_mode"             json:"recuring_mode"`
	InvestReference     *string          `db:"invest_references"         json:"invest_references"`
	InvestRemarks       *string          `db:"invest_remarks"            json:"invest_remarks"`
	DateStart           *string          `db:"date_start"                json:"date_start"`
	DateThru            *string          `db:"date_thru"                 json:"date_thru"`
	DateLastGenerate    *string          `db:"date_last_generate"        json:"date_last_generate"`
	AttempCount         *uint64          `db:"attempt_count"             json:"attempt_count"`
	SettleChannel       uint64           `db:"settle_channel"            json:"settle_channel"`
	SettlePaymentMethod uint64           `db:"settle_payment_method"     json:"settle_payment_method"`
	BankKey             *uint64          `db:"bank_key"                  json:"bank_key"`
	BankAccountKey      *uint64          `db:"bank_account_key"          json:"bank_account_key"`
	ProductKey          *uint64          `db:"product_key"               json:"product_key"`
	StampFeeAmount      *decimal.Decimal `db:"stamp_fee_amount"          json:"stamp_fee_amount"`
	TotalAmount         *decimal.Decimal `db:"total_amount"              json:"total_amount"`
	CustBankAccKey      *uint64          `db:"cust_bank_account_key"     json:"cust_bank_account_key"`
	RecOrder            *uint64          `db:"rec_order"                 json:"rec_order"`
	RecStatus           uint8            `db:"rec_status"                json:"rec_status"`
	RecCreatedDate      *string          `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy        *string          `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate     *string          `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy       *string          `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1           *string          `db:"rec_image1"                json:"rec_image1"`
	RecImage2           *string          `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus   *uint8           `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage    *uint64          `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate     *string          `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy       *string          `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate      *string          `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy        *string          `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1     *string          `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2     *string          `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3     *string          `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrAutoinvestRegistrationCount struct {
	CountData uint8 `db:"count_data"`
}

func CreateTrAutoinvestRegistration(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_autoinvest_registration"
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
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetTrAutoinvestRegistration(c *TrAutoinvestRegistration, key string) (int, error) {
	query := `SELECT tr_autoinvest_registration.* 
	FROM tr_autoinvest_registration 
	WHERE tr_autoinvest_registration.rec_status = "1" 
	AND tr_autoinvest_registration.autoinvest_key = ` + key
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrAutoinvestRegistrationCountData(c *TrAutoinvestRegistrationCount, acc string, date string, autoInvestKey string) (int, error) {
	query := `SELECT 
	COUNT(autoinvest_key) AS count_data 
   	FROM tr_autoinvest_registration
   	WHERE rec_status = 1 
	AND DATE_FORMAT(DATE(date_thru),'%Y-%m-%d') >= DATE_FORMAT(DATE(NOW()),'%Y-%m-%d')
   	AND acc_key = ` + acc + ` AND invest_date_execute = ` + date
	if autoInvestKey != "" {
		query += " AND autoinvest_key != " + autoInvestKey
	}
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateTrAutoinvestRegistration(params map[string]string) (int, error) {
	query := "UPDATE tr_autoinvest_registration SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "autoinvest_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE autoinvest_key = " + params["autoinvest_key"]
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetTrAutoinvestRegistrationIn(c *[]TrAutoinvestRegistration, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_autoinvest_registration.* FROM 
				tr_autoinvest_registration WHERE 
				tr_autoinvest_registration.rec_status = 1 `
	query := query2 + " AND tr_autoinvest_registration." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
