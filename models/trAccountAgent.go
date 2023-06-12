package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type TrAccountAgent struct {
	AcaKey            uint64  `db:"aca_key"                   json:"aca_key"`
	AccKey            uint64  `db:"acc_key"                   json:"acc_key"`
	EffDate           *string `db:"eff_date"                  json:"eff_date"`
	BranchKey         *uint64 `db:"branch_key"                json:"branch_key"`
	AgentKey          uint64  `db:"agent_key"                 json:"agent_key"`
	BasketKey         *string `db:"basket_key" 					json:"basket_key"`
	RecOrder          *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func CreateTrAccountAgent(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_account_agent"
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

func GetAllTrAccountAgent(c *[]TrAccountAgent, params map[string]string) (int, error) {
	query := `SELECT
              tr_account_agent.* FROM 
			  tr_account_agent`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_account_agent."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
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

func GetTrAccountAgent(c *TrAccountAgent, key string) (int, error) {
	query := `SELECT tr_account_agent.* FROM tr_account_agent WHERE tr_account_agent.rec_status = 1 AND tr_account_agent.aca_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrAccountAgentIn(c *[]TrAccountAgent, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_account_agent.* FROM 
				tr_account_agent WHERE 
				tr_account_agent.rec_status = 1 `
	query := query2 + " AND tr_account_agent." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("========= QUERY GET TRX AGENT ========= >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type AumBalanceUnitStruct struct {
	BalanceUnit decimal.Decimal `db:"balance_unit" json:"balance_unit"`
}

func AumBalanceQuery(c *AumBalanceUnitStruct, acaKey string, date string) (int, error) {
	query := `SELECT SUM(trbal.balance_unit) AS balance_unit 
		FROM tr_balance trbal
		INNER JOIN (
			SELECT bal.*
			FROM tr_balance bal
			WHERE bal.rec_status = 1
			AND bal.aca_key = "` + acaKey + `"
			AND bal.balance_date <= " ` + date + `"
		) c ON c.aca_key = trbal.aca_key AND c.balance_date = trbal.balance_date`

	// log.Println("===== QUERY GET AUM BALANCE UNIT ===== >>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

type AumNavValueStruct struct {
	NavValue decimal.Decimal `db:"nav_value" json:"nav_value"`
}

func AumNavValueQuery(c *AumNavValueStruct, productKey string, date string) (int, error) {
	query := `SELECT nav.nav_value
		FROM tr_nav nav
		INNER JOIN (
			SELECT nab.*
			FROM tr_nav nab
			WHERE nab.rec_status = 1
			AND nab.publish_mode = 236
			AND nab.nav_status = 234
			AND nab.product_key = "` + productKey + `"
			AND nab.nav_date <= "` + date + `"
		) c ON c.product_key = nav.product_key AND c.nav_date = nav.nav_date
		ORDER BY nav.nav_date DESC LIMIT 1`

	// log.Println("===== QUERY GET Nav Value Aum report ===== >>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

type AumCurrencyRateStruct struct {
	RateValue decimal.Decimal `db:"rate_value" json:"rate_value"`
}

func AumCurrencyRateQuery(c *AumCurrencyRateStruct, currencyKey string, date string) (int, error) {
	query := `SELECT tcur.rate_value
		FROM tr_currency_rate tcur
		INNER JOIN (
			SELECT curt.*
			FROM tr_currency_rate curt
			WHERE curt.rec_status = 1 
			AND curt.currency_key  = "` + currencyKey + `" 
			AND curt.rate_date <= "` + date + `"
		) c ON c.rate_date = tcur.rate_date AND c.currency_key = tcur.currency_key
		ORDER BY tcur.rate_date DESC LIMIT 1`

	// log.Println("===== QUERY GET currency rate Aum report ===== >>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
