package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type CustomerProductModel struct {
	AccKey           uint64  `db:"acc_key" json:"acc_key"`
	CustomerKey      uint64  `db:"customer_key"  json:"customer_key"`
	UnitHolderIdno   string  `db:"unit_holder_idno" json:"unit_holder_idno"`
	CustomerName     string  `db:"customer_name" json:"customer_name"`
	SidNo            *string `db:"sid_no" json:"sid_no"`
	FundTypeKey      uint64  `db:"fund_type_key"   json:"fund_type_key"`
	ProductKey       uint64  `db:"product_key"   json:"product_key"`
	ProductName      string  `db:"product_name"  json:"product_name"`
	AccountNo        *string `db:"account_no"  json:"account_no"`
	AccountName      *string `db:"account_name"  json:"account_name"`
	IFUANo           *string `db:"ifua_no"  json:"ifua_no"`
	SuspendSubFlag   bool    `db:"suspend_sub_flag" json:"suspend_sub_flag"`
	SuspendSubReason *string `db:"suspend_sub_reason" json:"suspend_sub_reason"`
	SuspendSubDate   *string `db:"suspend_sub_date" json:"suspend_sub_date"`
	SuspendRedFlag   bool    `db:"suspend_red_flag" json:"suspend_red_flag"`
	SuspendRedReason *string `db:"suspend_red_reason" json:"suspend_red_reason"`
	SuspendRedDate   *string `db:"suspend_red_date" json:"suspend_red_date"`
	CurrencyKey      *uint64 `db:"currency_key" json:"currency_key"`
	SettlementPeriod *uint64 `db:"settlement_period" json:"settlement_period"`
}

func GetCustomerProductList(c *[]CustomerProductModel, CustomerKey string, FundTypeKey string) (int, error) {
	query := `SELECT ta.customer_key
	, c.unit_holder_idno
	, c.sid_no
	, c.full_name AS customer_name
	, p.fund_type_key
	, p.product_key
	, p.product_name_alt AS product_name
	, ta.account_name
	, ta.account_no
	, ta.ifua_no
	, IFNULL(ta.sub_suspend_flag, 0) AS suspend_sub_flag
	, ta.sub_suspend_reason AS suspend_sub_reason
	, IFNULL(ta.red_suspend_flag, 0) AS suspend_red_flag
	, ta.red_suspend_reason AS suspend_red_reason
	, p.currency_key
	, p.settlement_period
	FROM tr_account ta
	INNER JOIN ms_customer c ON (c.customer_key=ta.customer_key AND c.rec_status=1)
	INNER JOIN ms_product p ON (p.product_key = ta.product_key AND p.rec_status=1 AND p.flag_enabled = 1)
	WHERE ta.rec_status = 1`

	query += ` AND ta.customer_key = ` + CustomerKey

	if FundTypeKey != "" {
		query += ` AND p.fund_type_key = ` + FundTypeKey
	}

	query += ` ORDER BY ta.rec_order`

	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type CustomerBalanceModel struct {
	CustomerKey  uint64          `db:"customer_key"  json:"customer_key"`
	CustomerName string          `db:"customer_name" json:"customer_name"`
	ProductKey   uint64          `db:"product_key"   json:"product_key"`
	ProductName  string          `db:"product_name"  json:"product_name"`
	BalanceUnit  decimal.Decimal `db:"balance_unit"  json:"balance_unit"`
	AverageNav   decimal.Decimal `db:"average_nav" json:"average_nav"`
	InquiryDate  string          `db:"inquiry_date"  json:"inquiry_date"`
}

func GetCustomerBalance(c *CustomerBalanceModel, customer_key uint64, product_key uint64, balance_date string) (int, error) {
	query := `SELECT 
	ac.customer_key,
	c.full_name AS customer_name,
	ac.product_key,
	p.product_name_alt AS product_name,
	sum(a.balance_unit) AS balance_unit,
	'%v' as inquiry_date,
	a.avg_nav AS average_nav
	FROM tr_balance a
	INNER JOIN (
		SELECT
			x.tc_key, 
			cast(MAX(x.balance_date) AS DATE) AS balance_date
		FROM tr_balance x
		INNER JOIN tr_account_agent x2 ON (x2.aca_key=x.aca_key AND x2.rec_status=1)
		INNER JOIN tr_account x3 ON (x3.acc_key=x2.acc_key AND x3.rec_status=1)
		WHERE x.rec_status = 1 
		AND x3.customer_key = %v 
		AND x3.product_key = %v
		AND cast(x.balance_date AS DATE) <= '%v'
		AND x.balance_unit >= 1
		GROUP BY x.tc_key
	) b ON (a.tc_key=b.tc_key AND cast(a.balance_date AS DATE)=cast(b.balance_date AS DATE))
	INNER JOIN tr_account_agent aa ON (aa.aca_key=a.aca_key AND aa.rec_status=1)
	INNER JOIN tr_account ac ON (ac.acc_key=aa.acc_key AND ac.rec_status=1)
	INNER JOIN ms_product p ON (p.product_key=ac.product_key AND p.rec_status=1 AND p.flag_enabled = 1)
	INNER JOIN ms_currency cr ON (cr.currency_key=p.currency_key AND cr.rec_status=1)
	INNER JOIN ms_customer c ON (c.customer_key=ac.customer_key AND c.rec_status=1)
	WHERE a.rec_status = 1
	AND ac.customer_key = %v
	AND ac.product_key = %v 
	GROUP BY ac.customer_key, c.full_name, ac.product_key, p.product_name_alt
	ORDER BY ac.product_key`

	s_sql := fmt.Sprintf(query, balance_date, customer_key, product_key, balance_date, customer_key, product_key)
	// // log.Println("========== GetCustomerBalance ==========>>>", s_sql)

	err := db.Db.Get(c, s_sql)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type NavModel struct {
	ProductKey    uint32          `db:"product_key" json:"product_key"`
	NavDate       string          `db:"nav_date" json:"nav_date"`
	NavValue      decimal.Decimal `db:"nav_value" json:"nav_value"`
	ProdAumTotal  decimal.Decimal `db:"prod_aum_total" json:"prod_aum_total"`
	ProdUnitTotal decimal.Decimal `db:"prod_unit_total" json:"prod_unit_total"`
}

func GetLastNAV(c *NavModel, product_key uint64) (int, error) {
	query := `SELECT 
		a.product_key, 
		cast(a.nav_date as DATE) AS nav_date, 
		a.nav_value, 
		a.prod_aum_total, 
		a.prod_unit_total 
	FROM tr_nav a
	INNER JOIN (
		SELECT 
			m.product_key, 
			MAX(m.nav_date) AS nav_date
		FROM tr_nav m
		WHERE m.rec_status = 1 AND m.nav_status = 234 and m.publish_mode=236
		AND m.product_key = %v
		GROUP BY m.product_key
	) b ON (a.product_key=b.product_key AND cast(a.nav_date AS DATE)=cast(b.nav_date AS DATE))
	WHERE a.rec_status = 1 AND a.nav_status = 234 AND a.publish_mode = 236
	AND a.product_key = %v
	ORDER BY a.rec_order;`

	s_sql := fmt.Sprintf(query, product_key, product_key)
	// // log.Println("========== GetLastNAV ==========>>>", s_sql)

	err := db.Db.Get(c, s_sql)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil

}

func GetHolidayByDate(c *[]ListHoliday, p_date time.Time) (int, error) {
	layout := "2006-01-02"
	str_date := p_date.Format(layout)
	query := `SELECT holiday_key, stock_market, CAST(a.holiday_date AS DATE) AS holiday_date, holiday_name 
	FROM ms_holiday a 
	WHERE a.rec_status = 1 
	AND a.stock_market = 272 
	AND CAST(a.holiday_date AS DATE) = '%v'`
	s_sql := fmt.Sprintf(query, str_date)
	err := db.Db.Select(c, s_sql)
	if err != nil {
		if err != sql.ErrNoRows {
			return http.StatusBadGateway, err
		}
	}
	return http.StatusOK, nil
}

func IsHoliday(p_date time.Time) (bool, error) {
	weekday := p_date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return true, nil
	}
	var db_row []ListHoliday
	_, err := GetHolidayByDate(&db_row, p_date)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
	}
	if len(db_row) > 0 {
		return true, nil
	}
	return false, nil
}

func FnAddDate(p_date time.Time, n_day int) (time.Time, error) {
	n_add := 0
	if n_day > 0 {
		n_add = 1
	} else if n_day < 0 {
		n_add = -1
	} else {
		n_add = 0
	}
	var err error
	var bHoliday bool
	new_date := p_date
	i_loop := n_day
	if i_loop == 0 {
		err = nil
	}
	for i_loop != 0 {
		new_date = new_date.AddDate(0, 0, n_add)
		bHoliday, err = IsHoliday(new_date)
		if !bHoliday { //loop terus jika masi libur sampai ketemu hari kerja (bukan libur).
			i_loop = i_loop - n_add
			if i_loop == 0 {
				err = nil
			}
		}
	}
	return new_date, err
}

func CekInProcessByUnit(productKey uint64, customerKey uint64) decimal.Decimal {
	result := decimal.NewFromInt(0)
	var unit_amount *decimal.Decimal

	query := `SELECT sum(t1.trans_unit) FROM tr_transaction t1 WHERE t1.rec_status = 1 AND t1.trans_type_key IN (2,3) 
	AND t1.trans_status_key IN (2,4,5,6,7,8) AND t1.trans_source IN (143,144)`
	prevDate, _ := FnAddDate(time.Now(), -1)
	query += ` AND CAST(t1.trans_date AS DATE) >= "` + prevDate.Format("2006-01-02") + `" 
	AND CAST(t1.trans_date AS DATE) <= "` + time.Now().Format("2006-01-02") + `" 
	AND t1.customer_key = ` + strconv.FormatUint(customerKey, 10) + ` AND t1.product_key = ` + strconv.FormatUint(productKey, 10)

	err := db.Db.Get(&unit_amount, query)
	if err != nil {
	}

	if unit_amount == nil || (unit_amount != nil && unit_amount.Cmp(decimal.NewFromInt(0)) <= 0) {
		var nominal_amount *decimal.Decimal
		query2 := `SELECT sum(t1.trans_amount)
		FROM tr_transaction t1
		WHERE t1.rec_status = 1 AND t1.trans_type_key IN (2,3) 
		AND t1.trans_status_key IN (2,4,5,6,7,8)
		AND t1.trans_source IN (143,144)`
		query2 += ` AND CAST(t1.trans_date AS DATE) >= "` + prevDate.Format("2006-01-02") + `" AND CAST(t1.trans_date AS DATE) <= "` + time.Now().Format("2006-01-02") + `" 
		AND t1.customer_key = ` + strconv.FormatUint(customerKey, 10) + ` AND t1.product_key = ` + strconv.FormatUint(productKey, 10)
		err := db.Db.Get(&nominal_amount, query2)
		if err != nil {
		}
		if nominal_amount != nil {
			var getNav NavModel
			_, err = GetLastNAV(&getNav, productKey)
			if err != nil {
			}
			re := nominal_amount.Div(getNav.NavValue)
			result = re
		}
	} else {
		result = *unit_amount
	}

	return result
}
