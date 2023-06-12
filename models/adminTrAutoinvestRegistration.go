package models

import (
	"mf-bo-api/db"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type AdminListAutoInvestRegistration struct {
	AutoinvestKey       uint64          `db:"autoinvest_key"                json:"autoinvest_key"`
	CustomerKey         uint64          `db:"customer_key" json:"customer_key"`
	FullName            string          `db:"full_name"                     json:"full_name"`
	ProductKey          uint64          `db:"product_key" json:"product_key"`
	ProductName         string          `db:"product_name"                  json:"product_name"`
	InvestAmount        decimal.Decimal `db:"invest_amount"                 json:"invest_amount"`
	InvestDateExecute   *uint64         `db:"invest_date_execute"           json:"invest_date_execute"`
	DateStart           string          `db:"date_start"                    json:"date_start"`
	DateThru            string          `db:"date_thru"                     json:"date_thru"`
	DateLastGenerate    *string         `db:"date_last_generate"            json:"date_last_generate"`
	SettleChannel       string          `db:"settle_channel"                json:"settle_channel"`
	AttemptCount        *uint64         `db:"attempt_count"                 json:"attempt_count"`
	SettlePaymentMethod string          `db:"settle_payment_method"         json:"settle_payment_method"`
	BankKey             uint64          `db:"bank_key" json:"bank_key"`
	BankName            *string         `db:"bank_name"                     json:"bank_name"`
	AccountNo           *string         `db:"account_no"                    json:"account_no"`
	AccountHolderName   *string         `db:"account_holder_name"           json:"account_holder_name"`
}

func GetAdminListAutoInvestRegistration(c *[]AdminListAutoInvestRegistration, params map[string]string, limit uint64, offset uint64, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

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

	query := `SELECT 
				a.autoinvest_key,
				c.customer_key,
				c.full_name,
				p.product_key,
				p.product_name_alt AS product_name,
				a.invest_amount,
				a.invest_date_execute,
				DATE_FORMAT(a.date_start, '%d %M %Y') AS date_start,
				DATE_FORMAT(a.date_thru, '%d %M %Y') AS date_thru,
				DATE_FORMAT(a.date_last_generate, '%d %M %Y') AS date_last_generate,
				a.attempt_count,
				setchannel.lkp_name AS settle_channel,
				setpayment.lkp_name AS settle_payment_method,
				b.bank_key,
				b.bank_name,
				ba.account_no,
				ba.account_holder_name 
			FROM tr_autoinvest_registration AS a 
			INNER JOIN tr_account AS ta ON ta.acc_key = a.acc_key 
			INNER JOIN ms_customer AS c ON c.customer_key = ta.customer_key 
			INNER JOIN ms_product AS p ON p.product_key = ta.product_key 
			INNER JOIN gen_lookup AS setchannel ON setchannel.lookup_key = a.settle_channel 
			INNER JOIN gen_lookup AS setpayment ON setpayment.lookup_key = a.settle_payment_method 
			LEFT JOIN ms_bank AS b ON b.bank_key = a.bank_key 
			LEFT JOIN ms_bank_account AS ba ON ba.bank_account_key = a.cust_bank_account_key 
			WHERE a.rec_status = 1 AND ta.rec_status = 1 AND (ta.sub_suspend_flag IS NULL OR ta.sub_suspend_flag = 0) 
			AND c.rec_status = 1 AND c.investor_type = 263 
			AND DATE_FORMAT(DATE(a.date_thru),'%Y-%m-%d') >= DATE_FORMAT(DATE(NOW()),'%Y-%m-%d')`

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += condition + orderCondition + limitOffset

	// Main query
	// log.Println("========== QUERY LIST AUTOINVEST ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAdminCountListAutoInvestRegistration(c *CountData, params map[string]string) (int, error) {
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

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

	query := `SELECT 
				COUNT(a.autoinvest_key) AS count_data 
			FROM tr_autoinvest_registration AS a 
			INNER JOIN tr_account AS ta ON ta.acc_key = a.acc_key 
			INNER JOIN ms_customer AS c ON c.customer_key = ta.customer_key 
			INNER JOIN ms_product AS p ON p.product_key = ta.product_key 
			INNER JOIN gen_lookup AS setchannel ON setchannel.lookup_key = a.settle_channel 
			INNER JOIN gen_lookup AS setpayment ON setpayment.lookup_key = a.settle_payment_method 
			LEFT JOIN ms_bank AS b ON b.bank_key = a.bank_key 
			LEFT JOIN ms_bank_account AS ba ON ba.bank_account_key = a.cust_bank_account_key 
			WHERE a.rec_status = 1 AND ta.rec_status = 1 AND (ta.sub_suspend_flag IS NULL OR ta.sub_suspend_flag = 0) 
			AND c.rec_status = 1 AND c.investor_type = 263 
			AND DATE_FORMAT(DATE(a.date_thru),'%Y-%m-%d') >= DATE_FORMAT(DATE(NOW()),'%Y-%m-%d')`

	query += orderCondition + limitOffset

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminValidateAccAndiInvestDateExecute(c *CountData, accKey string, dateExecute string, autoInvestKey string) (int, error) {
	query := `SELECT 
				COUNT(autoinvest_key) AS count_data 
			FROM tr_autoinvest_registration 
			WHERE rec_status = 1 
			AND DATE_FORMAT(DATE(date_thru),'%Y-%m-%d') >= DATE_FORMAT(DATE(NOW()),'%Y-%m-%d')
			AND acc_key = "` + accKey + `" AND invest_date_execute = "` + dateExecute + `"`

	if autoInvestKey != "" {
		query += " AND autoinvest_key != " + autoInvestKey
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type DetailAutoinvestRegistration struct {
	AutoinvestKey           uint64          `db:"autoinvest_key"               json:"autoinvest_key"`
	AccKey                  uint64          `db:"acc_key"                      json:"acc_key"`
	CustomerKey             uint64          `db:"customer_key"                 json:"customer_key"`
	FullName                string          `db:"full_name"                    json:"full_name"`
	ProductKey              uint64          `db:"product_key"                  json:"product_key"`
	SubSuspendFlag          *uint8          `db:"sub_suspend_flag"             json:"sub_suspend_flag"`
	ProductName             string          `db:"product_name"                 json:"product_name"`
	InvestAmount            decimal.Decimal `db:"invest_amount"                json:"invest_amount"`
	InvestFeeRates          decimal.Decimal `db:"invest_fee_rates"             json:"invest_fee_rates"`
	InvestFeeAmount         decimal.Decimal `db:"invest_fee_amount"            json:"invest_fee_amount"`
	InvestFeeCharges        decimal.Decimal `db:"invest_fee_charges"           json:"invest_fee_charges"`
	InvestDateExecute       uint8           `db:"invest_date_execute"          json:"invest_date_execute"`
	RecuringMode            *uint64         `db:"recuring_mode"                json:"recuring_mode"`
	RecuringModeName        *string         `db:"recuring_mode_name"           json:"recuring_mode_name"`
	InvestReference         *string         `db:"invest_references"            json:"invest_references"`
	InvestRemarks           *string         `db:"invest_remarks"               json:"invest_remarks"`
	DateStart               *string         `db:"date_start"                   json:"date_start"`
	DateStartOrigin         *string         `db:"date_start_origin"            json:"date_start_origin"`
	DateThru                *string         `db:"date_thru"                    json:"date_thru"`
	DateThruOrigin          *string         `db:"date_thru_origin"             json:"date_thru_origin"`
	DateLastGenerate        *string         `db:"date_last_generate"           json:"date_last_generate"`
	DateLastGenerateOrigin  *string         `db:"date_last_generate_origin"    json:"date_last_generate_origin"`
	AttempCount             uint64          `db:"attempt_count"                json:"attempt_count"`
	SettleChannel           uint64          `db:"settle_channel"               json:"settle_channel"`
	SettleChannelName       string          `db:"settle_channel_name"          json:"settle_channel_name"`
	SettlePaymentMethod     uint64          `db:"settle_payment_method"        json:"settle_payment_method"`
	SettlePaymentMethodName string          `db:"settle_payment_method_name"   json:"settle_payment_method_name"`
	BankKey                 *uint64         `db:"bank_key"                     json:"bank_key"`
	BankName                *string         `db:"bank_name"                    json:"bank_name"`
	CustBankAccKey          *uint64         `db:"cust_bank_account_key"        json:"cust_bank_account_key"`
	AccountNo               *string         `db:"account_no"                   json:"account_no"`
	AccountHolderName       *string         `db:"account_holder_name"          json:"account_holder_name"`
	BranchName              *string         `db:"branch_name"                  json:"branch_name"`
	ProdBankaccKey          *uint64         `db:"prod_bankacc_key"             json:"prod_bankacc_key"`
}

func AdminGetDetailTrAutoinvestRegistration(c *DetailAutoinvestRegistration, key string) (int, error) {
	query := `SELECT 
				a.autoinvest_key,
				a.acc_key,
				ta.customer_key,
				c.full_name,
				ta.product_key,
				ta.sub_suspend_flag,
				p.product_name_alt AS product_name,
				(CASE WHEN a.invest_amount IS NULL THEN "0" ELSE a.invest_amount END) AS invest_amount,
				(CASE WHEN a.invest_fee_rates IS NULL THEN "0" ELSE a.invest_fee_rates END) AS invest_fee_rates,
				(CASE WHEN a.invest_fee_amount IS NULL THEN "0" ELSE a.invest_fee_amount END) AS invest_fee_amount,
				(CASE WHEN a.invest_fee_charges IS NULL THEN "0" ELSE a.invest_fee_charges END) AS invest_fee_charges,
				a.invest_date_execute,
				a.recuring_mode,
				recuring.lkp_name AS recuring_mode_name,
				a.invest_references,
				a.invest_remarks,
				DATE_FORMAT(a.date_start, '%d %M %Y') AS date_start,
				a.date_start AS date_start_origin,
				DATE_FORMAT(a.date_thru, '%d %M %Y') AS date_thru,
				a.date_thru AS date_thru_origin,
				DATE_FORMAT(a.date_last_generate, '%d %M %Y') AS date_last_generate,
				a.date_last_generate AS date_last_generate_origin,
				(CASE WHEN a.attempt_count IS NULL THEN "0" ELSE a.attempt_count END) AS attempt_count,
				a.settle_channel,
				setchannel.lkp_name AS settle_channel_name,
				a.settle_payment_method,
				setpayment.lkp_name AS settle_payment_method_name,
				b.bank_key,
				b.bank_name,
				a.cust_bank_account_key,
				ba.account_no,
				ba.account_holder_name,
				ba.branch_name,
				pba.prod_bankacc_key 
			FROM tr_autoinvest_registration AS a 
			INNER JOIN tr_account AS ta ON ta.acc_key = a.acc_key 
			INNER JOIN ms_customer AS c ON c.customer_key = ta.customer_key 
			INNER JOIN ms_product AS p ON p.product_key = ta.product_key 
			LEFT JOIN gen_lookup AS recuring ON recuring.lookup_key = a.recuring_mode 
			INNER JOIN gen_lookup AS setchannel ON setchannel.lookup_key = a.settle_channel 
			INNER JOIN gen_lookup AS setpayment ON setpayment.lookup_key = a.settle_payment_method 
			LEFT JOIN ms_bank AS b ON b.bank_key = a.bank_key 
			LEFT JOIN ms_bank_account AS ba ON ba.bank_account_key = a.cust_bank_account_key 
			LEFT JOIN ms_product_bank_account AS pba ON pba.bank_account_key = ba.bank_account_key 
			AND pba.product_key = ta.product_key AND pba.rec_status = 1 AND pba.bank_account_purpose = 271 
			WHERE a.rec_status = 1 AND ta.rec_status = 1 
			AND c.rec_status = 1 AND a.autoinvest_key = ` + key
	// log.Println("========== QUERY GET AUTOINVEST REGIST DETAIL =========", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
