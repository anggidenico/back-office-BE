package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type AutoInvestRegistration struct {
	AutoInvestKey           int64           `db:"autoinvest_key" json:"autoinvest_key"`
	ProductKey              *int64          `db:"product_key" json:"product_key"`
	ProductName             *string         `db:"product_name" json:"product_name"`
	AccKey                  int64           `db:"acc_key" json:"acc_key"`
	AccName                 *string         `db:"acc_name" json:"acc_name"`
	InvestAmount            decimal.Decimal `db:"invest_amount" json:"invest_amount"`
	InvestFeeRates          decimal.Decimal `db:"invest_fee_rates" json:"invest_fee_rates"`
	InvestFeeAmount         decimal.Decimal `db:"invest_fee_amount" json:"invest_fee_amount"`
	InvestFeeCharges        decimal.Decimal `db:"invest_fee_charges" json:"invest_fee_charges"`
	StampFeeAmount          decimal.Decimal `db:"stamp_fee_amount" json:"stamp_fee_amount"`
	TotalAmount             decimal.Decimal `db:"total_amount" json:"total_amount"`
	InvestDateExecute       int64           `db:"invest_date_execute" json:"invest_date_execute"`
	RecuringMode            *int64          `db:"recuring_mode" json:"recuring_mode"`
	InvestReferences        *string         `db:"invest_references" json:"invest_reference"`
	InvestRemarks           *string         `db:"invest_remarks" json:"invest_remarks"`
	DateStart               *string         `db:"date_start" json:"date_start"`
	DateThru                *string         `db:"date_thru" json:"date_thru"`
	DateLastGenerate        *string         `db:"date_last_generate" json:"date_last_generate"`
	AttemptCount            *int64          `db:"attempt_count" json:"attempt_count"`
	SettleChannel           int64           `db:"settle_channel" json:"settle_channel"`
	SettleChannelName       string          `db:"settle_channel_name" json:"settle_channel_name"`
	SettlePaymentMethod     int64           `db:"settle_payment_method" json:"settle_payment_method"`
	SettlePaymentMethodName string          `db:"settle_payment_method_name" json:"settle_payment_method_name"`
	BankKey                 *int64          `db:"bank_key" json:"bank_key"`
	BankName                *string         `db:"bank_name" json:"bank_name"`
	CustBankaccKey          *int64          `db:"cust_bankacc_key" json:"custom_bankacc_key"`
	BankAccountName         *string         `db:"bank_account_name" json:"bank_account_name"`
	RecOrder                *int64          `db:"rec_order" json:"rec_order"`
}

func GetAutoInvestRegistrationModels(c *[]AutoInvestRegistration) (int, error) {
	query := `SELECT a.autoinvest_key,
	a.product_key,
	b.product_name,
	a.acc_key,
	c.account_name acc_name,
	a.invest_amount,
	a.invest_fee_rates,
	a.invest_fee_amount,
	a.invest_fee_charges,
	a.stamp_fee_amount,
	a.total_amount,
	a.invest_date_execute,
	a.recuring_mode,
	a.invest_references,
	a.invest_remarks,
	a.date_start,
	a.date_thru,
	a.date_last_generate,
	a.attempt_count,
	a.settle_channel,
	f.lkp_name settle_channel_name,
	a.settle_payment_method,
	g.lkp_name settle_payment_method_name,
	a.bank_key,
	d.bank_name,
	a.cust_bankacc_key,
	e.bank_account_name,
	a.rec_order
FROM tr_autoinvest_registration a
LEFT JOIN ms_product b ON a.product_key = b.product_key
LEFT JOIN tr_account c ON a.acc_key = c.acc_key
LEFT JOIN ms_bank d ON a.bank_key = d.bank_key
LEFT JOIN ms_customer_bank_account e ON a.cust_bankacc_key = e.cust_bankacc_key
LEFT JOIN gen_lookup f ON a.settle_channel = f.lookup_key
LEFT JOIN gen_lookup g ON  a.settle_payment_method = g.lookup_key
WHERE a.rec_status = 1
ORDER BY a.autoinvest_key DESC`

	// log.Println(query)

	err := db.Db.Select(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAutoInvestRegisDetailModels(c *AutoInvestRegistration, AutoInvestKey string) (int, error) {
	query := `SELECT a.autoinvest_key,
	a.product_key,
	b.product_name,
	a.acc_key,
	c.account_name acc_name,
	a.invest_amount,
	a.invest_fee_rates,
	a.invest_fee_amount,
	a.invest_fee_charges,
	a.stamp_fee_amount,
	a.total_amount,
	a.invest_date_execute,
	a.recuring_mode,
	a.invest_references,
	a.invest_remarks,
	a.date_start,
	a.date_thru,
	a.date_last_generate,
	a.attempt_count,
	a.settle_channel,
	f.lkp_name settle_channel_name,
	a.settle_payment_method,
	g.lkp_name settle_payment_method_name,
	a.bank_key,
	d.bank_name,
	a.cust_bankacc_key,
	e.bank_account_name,
	a.rec_order
FROM tr_autoinvest_registration a
LEFT JOIN ms_product b ON a.product_key = b.product_key
LEFT JOIN tr_account c ON a.acc_key = c.acc_key
LEFT JOIN ms_bank d ON a.bank_key = d.bank_key
LEFT JOIN ms_customer_bank_account e ON a.cust_bankacc_key = e.cust_bankacc_key
LEFT JOIN gen_lookup f ON a.settle_channel = f.lookup_key
LEFT JOIN gen_lookup g ON  a.settle_payment_method = g.lookup_key
	WHERE a.rec_status = 1 
	AND a.autoinvest_key =` + AutoInvestKey

	// log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteAutoInvestRegis(AutoInvestKey string, params map[string]string) (int, error) {
	query := `UPDATE tr_autoinvest_registration SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "autoinvest_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE autoinvest_key = ?`
	values = append(values, AutoInvestKey)

	// log.Println("========== DeleteAutoInvestregistration ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
