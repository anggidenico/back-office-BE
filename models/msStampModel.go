package models

import (
	"mf-bo-api/db"
	"net/http"

	"github.com/shopspring/decimal"
)

type TransStampData struct {
	CustomerKey      uint64           `db:"customer_key"         json:"customer_key"`
	NavDate          string           `db:"nav_date"             json:"nav_date"`
	HasStamp         bool             `db:"has_stamp"            json:"has_stamp"`
	StampFeeAmount   decimal.Decimal  `db:"stamp_fee_amount"     json:"stamp_fee_amount"`
	TransAmountTotal decimal.Decimal  `db:"trans_amount_total"   json:"trans_amount_total"`
	StampMessageInfo string           `db:"stamp_message_info"   json:"stamp_message_info"`
	StampDutyValues  []StampDutyValue `db:"stamp_duty_values"    json:"stamp_duty_values"`
}

type StampDutyValue struct {
	CurrencyKey       uint64          `db:"currency_key"         json:"currency_key"`
	CurrCode          string          `db:"currency_code"        json:"currency_code"`
	StampValue        decimal.Decimal `db:"stamp_value"          json:"stamp_value"`
	MinimumTransAccum decimal.Decimal `db:"minimum_trans_accum"  json:"minimum_trans_accum"`
}

func GetStampNominal(c *[]StampDutyValue, params map[string]string) (int, error) {
	//ambil nilai materai dari db => ini konstanta
	query := `SELECT a.currency_key, a.code AS currency_code, cast(b.stamp_value as decimal(10,2)) as stamp_value , cast(c.minimum_trans_accum as decimal(10,2)) as minimum_trans_accum 
	FROM ms_currency a
	INNER JOIN (
	
		SELECT '1' AS currency_key, app_config_value as stamp_value FROM sc_app_config WHERE app_config_code='TRX_STAMP_VALUE_IDR'
		UNION SELECT a.currency_key,  
		 (SELECT cast(app_config_value as DECIMAL(15,2)) FROM sc_app_config WHERE app_config_code='TRX_STAMP_VALUE_IDR')/a.rate_value AS stamp_value
		FROM tr_currency_rate a
		INNER JOIN (
		
			SELECT MAX(rate_date) AS rate_date, currency_key, MAX(curr_rate_key) AS curr_rate_key
			FROM tr_currency_rate
			WHERE rec_status=1 
			AND rate_type=293
			AND rate_date <= CURRENT_DATE()
			GROUP BY currency_key
		
		) b ON (a.currency_key = b.currency_key AND DATE_FORMAT(a.rate_date,'%Y-%m-%d') = DATE_FORMAT(b.rate_date,'%Y-%m-%d'))
		WHERE a.rec_status=1 
		AND a.rate_type=293
		AND a.currency_key=2		
			
	) b ON (a.currency_key = b.currency_key) 
	INNER JOIN (
	
		SELECT '1' AS currency_key, app_config_value as minimum_trans_accum FROM sc_app_config WHERE app_config_code='TRX_STAMP_MIN_VALUE_IDR'
		UNION SELECT a.currency_key,  
		 (SELECT cast(app_config_value as DECIMAL(15,2)) FROM sc_app_config WHERE app_config_code='TRX_STAMP_MIN_VALUE_IDR')/a.rate_value AS minimum_trans_accum
		FROM tr_currency_rate a
		INNER JOIN (
		
			SELECT MAX(rate_date) AS rate_date, currency_key, MAX(curr_rate_key) AS curr_rate_key
			FROM tr_currency_rate
			WHERE rec_status=1 
			AND rate_type=293
			AND rate_date <= CURRENT_DATE()
			GROUP BY currency_key
		
		) b ON (a.currency_key = b.currency_key AND DATE_FORMAT(a.rate_date,'%Y-%m-%d') = DATE_FORMAT(b.rate_date,'%Y-%m-%d'))
		WHERE a.rec_status=1 
		AND a.rate_type=293
		AND a.currency_key=2		

	) c ON (a.currency_key = c.currency_key)`

	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, " a."+field+" = '"+value+"'")
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

func GetTransactionStamps(c *TransStampData, params map[string]string) (int, error) {
	query := `SELECT
		c.customer_key,
		IFNull(t.stamp_fee_amount,0) AS stamp_fee_amount,
		case when t.entry_mode=139 then IFNULL(t.trans_unit,0) * n1.nav_value ELSE IFNull(t.trans_amount,0) END AS trans_amount_total
	FROM ms_customer c 
	LEFT JOIN tr_transaction t ON (c.customer_key = t.customer_key 
		AND cast(t.nav_date AS DATE) = CURRENT_DATE() 
		AND t.trans_status_key NOT IN (1,3) 
		AND t.rec_status=1)
	LEFT JOIN vwtr_last_nav n ON (n.product_key = t.product_key) 
	LEFT JOIN tr_nav n1 ON (n1.product_key = n.product_key AND n1.nav_date=n.nav_date) `

	var whereClause []string
	var condition string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, " c."+field+" = '"+value+"'")
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

	query += condition
	s_qry := `SELECT  
		CURRENT_DATE() as nav_date, 
		a.customer_key, 
		0 as has_stamp, 
		SUM(a.stamp_fee_amount) AS stamp_fee_amount, 
		SUM(a.trans_amount_total) AS trans_amount_total, 
		'Biaya materai' as stamp_message_info 
	FROM  ( ` + query + ` ) a 
	GROUP BY a.customer_key `

	// Main query
	// log.Info(s_qry)
	// log.Info(s_qry)
	err := db.Db.Get(c, s_qry)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
