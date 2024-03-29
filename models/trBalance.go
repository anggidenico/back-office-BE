package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type TrBalance struct {
	BalanceKey        uint64           `db:"balance_key"               json:"balance_key"`
	AcaKey            uint64           `db:"aca_key"                   json:"aca_key"`
	TcKey             uint64           `db:"tc_key"                    json:"tc_key"`
	BalanceDate       string           `db:"balance_date"              json:"balance_date"`
	BalanceUnit       decimal.Decimal  `db:"balance_unit"              json:"balance_unit"`
	AvgNav            *decimal.Decimal `db:"avg_nav"                   json:"avg_nav"`
	TcKeyRed          *uint64          `db:"tc_key_red"                json:"tc_key_red"`
	RecOrder          *uint64          `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8            `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string          `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string          `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string          `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string          `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string          `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string          `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8           `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64          `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string          `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string          `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string          `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string          `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string          `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string          `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string          `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrBalanceCustomerProduk struct {
	BalanceKey     uint64          `db:"balance_key"               json:"balance_key"`
	AcaKey         uint64          `db:"aca_key"                   json:"aca_key"`
	BalanceUnit    decimal.Decimal `db:"balance_unit"              json:"balance_unit"`
	TcKey          uint64          `db:"tc_key"                    json:"tc_key"`
	TransactionKey uint64          `db:"transaction_key"           json:"transaction_key"`
	NavDate        string          `db:"nav_date"                  json:"nav_date"`
}

type AvgNav struct {
	AvgNav *decimal.Decimal `db:"avg_nav"                   json:"avg_nav"`
}

func GetLastBalanceIn(c *[]TrBalance, acaKey []string) (int, error) {
	inQuery := strings.Join(acaKey, ",")
	query2 := `SELECT t1.balance_key, t1.aca_key, t1.tc_key, t1.balance_date, t1.balance_unit, t1.avg_nav, t1.tc_key_red 
	FROM tr_balance t1 JOIN (SELECT MAX(balance_date) balance_date, tc_key FROM tr_balance WHERE rec_status=1 AND balance_date < NOW()  GROUP BY tc_key) t2
	ON (t1.balance_date = t2.balance_date AND t1.tc_key = t2.tc_key)`
	query := query2 + " WHERE t1.rec_status = 1 AND t1.aca_key IN(" + inQuery + ") GROUP BY tc_key ORDER BY t1.balance_key DESC"

	// log.Println("========= QUERY GET LAST BALANCE ========= >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateTrBalance(params map[string]string) (int, error) {
	query := "INSERT INTO tr_balance"
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

func GetLastBalanceCustomerByProductKey(c *[]TrBalanceCustomerProduk, customerKey string, productKey string) (int, error) {
	query := `SELECT 
				tb.balance_key as balance_key, 
				tb.aca_key as aca_key, 
				tb.balance_unit as balance_unit, 
				tc.tc_key as tc_key, 
				tr.transaction_key as transaction_key, 
				tr.nav_date as nav_date 
				FROM tr_balance AS tb
				JOIN (SELECT MAX(balance_key) balance_key FROM tr_balance where rec_status = 1 GROUP BY tc_key) AS t2 ON tb.balance_key = t2.balance_key 
				INNER JOIN tr_transaction_confirmation AS tc ON tb.tc_key = tc.tc_key
				INNER JOIN tr_transaction AS tr ON tc.transaction_key = tr.transaction_key
				WHERE tr.customer_key = ` + customerKey +
		` AND tr.product_key = ` + productKey +
		` AND tr.trans_status_key = 9 AND tr.rec_status = 1 AND tb.rec_status = 1 
		  AND tc.rec_status = 1 AND tr.trans_type_key = 1 AND tb.balance_unit > 0 
				GROUP BY tb.tc_key  ORDER BY tc.tc_key ASC`

	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetLastTrBalanceByTcRed(c *TrBalance, tcKeyRed string) (int, error) {
	query := `SELECT * FROM tr_balance WHERE rec_status = 1 AND tc_key_red = ` + tcKeyRed + ` ORDER BY rec_order DESC LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

/*
ambil AvgNav terakhir dari product yang dimiliki customer
*/
func GetLastAvgNavTrBalanceCustomerByProductKey(c *AvgNav, customerKey string, productKey string) (int, error) {
	/*
		d.customer_key,
		d.product_key,
		count(c.acc_key) AS count_acc,
		count(a.aca_key) AS count_aca,
		CURRENT_DATE() as inquiry_date,
	*/
	query := `
	SELECT a.avg_nav FROM tr_balance a
	INNER JOIN (
		SELECT x.tc_key, cast(MAX(x.balance_date) AS DATE) AS balance_date
		FROM tr_balance x
		INNER JOIN tr_account_agent x2 ON (x2.aca_key = x.aca_key AND x2.rec_status = 1)
		INNER JOIN tr_account x3 ON (x3.acc_key = x2.acc_key AND x3.rec_status = 1)
		WHERE x.rec_status = 1 
		AND x3.customer_key = ` + customerKey + `
		AND x3.product_key = ` + productKey + `
		AND cast(x.balance_date AS DATE) <= CURRENT_DATE()
		AND x.balance_unit >= 1
		GROUP BY x.tc_key
	) b ON (a.tc_key = b.tc_key AND cast(a.balance_date AS DATE) = cast(b.balance_date AS DATE))
	INNER JOIN tr_account_agent aa ON (aa.aca_key = a.aca_key AND aa.rec_status = 1)
	INNER JOIN tr_account ac ON (ac.acc_key = aa.acc_key AND ac.rec_status = 1)
	INNER JOIN ms_product p ON (p.product_key = ac.product_key AND p.rec_status = 1)
	INNER JOIN ms_currency cr ON (cr.currency_key = p.currency_key AND cr.rec_status = 1)
	INNER JOIN ms_customer c ON (c.customer_key = ac.customer_key AND c.rec_status = 1)
	WHERE a.rec_status = 1
	AND ac.customer_key = ` + customerKey + `
	AND ac.product_key = ` + productKey + `
	GROUP BY ac.customer_key, c.full_name, ac.product_key, p.product_name_alt
	ORDER BY ac.product_key`

	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err != sql.ErrNoRows {
			// log.Println(err)
			return http.StatusBadGateway, err
		}
	}

	return http.StatusOK, nil
}

func UpdateTrBalance(params map[string]string, value string, field string) (int, error) {
	query := "UPDATE tr_balance SET "
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
	// // log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// // log.Error(err)
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
		// // log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

type SumBalanceUnit struct {
	AcaKey uint64          `db:"aca_key"           json:"aca_key"`
	Unit   decimal.Decimal `db:"unit"              json:"unit"`
}

func GetSumBalanceUnit(c *[]SumBalanceUnit, acaKeys []string) (int, error) {
	inQuery := strings.Join(acaKeys, ",")
	query := `SELECT
					t.aca_key as aca_key,
					SUM(t.balance_unit) AS unit
				FROM
					(
						SELECT 
							*
						FROM tr_balance 
						WHERE balance_key IN (SELECT MAX(balance_key) AS id 
						FROM tr_balance WHERE aca_key IN(` + inQuery + `)
						GROUP BY tc_key
					) ORDER BY balance_key
				) AS t
				GROUP BY t.aca_key`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetBalanceUnitByCustomerAndProduct(c *SumBalanceUnit, customerKey string, productKey string) (int, error) {
	query := `SELECT
				t.aca_key AS aca_key,
				SUM(t.balance_unit) AS unit
			FROM
				(
					SELECT 
						*
					FROM tr_balance
					WHERE balance_key IN 
					(
						SELECT MAX(t.balance_key) AS id 
						FROM tr_balance AS t
						INNER JOIN tr_account_agent AS aca ON aca.aca_key = t.aca_key
						INNER JOIN tr_account AS ta ON ta.acc_key = aca.acc_key
						WHERE ta.customer_key = '` + customerKey + `' AND ta.product_key = '` + productKey + `'
						GROUP BY tc_key
					) ORDER BY balance_key
				) AS t
			GROUP BY t.aca_key`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type BeginningEndingBalance struct {
	Tanggal     string           `db:"tgl"             json:"tgl"`
	Description string           `db:"description"     json:"description"`
	Amount      *decimal.Decimal `db:"amount"          json:"amount"`
	NavValue    *decimal.Decimal `db:"nav_value"       json:"nav_value"`
	Unit        *decimal.Decimal `db:"unit"            json:"unit"`
	AvgNav      *decimal.Decimal `db:"avg_nav"         json:"avg_nav"`
	Fee         *decimal.Decimal `db:"fee"             json:"fee"`
	DecNav      *int32           `db:"dec_nav"         json:"dec_nav"`
	DecUnit     *int32           `db:"dec_unit"        json:"dec_unit"`
	DecAmount   *int32           `db:"dec_amount" json:"dec_amount"`
}

func GetBeginningEndingBalanceAcc(c *BeginningEndingBalance, desc string, date string, accKey string, productKey string) (int, error) {
	query := `SELECT
				DATE_FORMAT('` + date + `', '%d %b %Y') AS tgl,
				'` + desc + `' AS description,
				(nv.nav_value * SUM(t.balance_unit)) AS amount,
				nv.nav_value,
				SUM(t.balance_unit) AS unit,
				t.avg_nav,
				0 AS fee,
				msp.dec_nav,
				msp.dec_unit,
				msp.dec_amount
			FROM
				(
					SELECT 
						*
					FROM tr_balance 
					WHERE balance_key IN (
					SELECT MAX(balance_key) AS id 
						FROM tr_balance WHERE aca_key IN (SELECT aca_key FROM tr_account_agent WHERE acc_key = '` + accKey + `') 
						AND rec_status = 1 AND balance_date <= '` + date + `'
						GROUP BY tc_key
					) AND rec_status = 1 ORDER BY balance_key
				) AS t 
			LEFT JOIN (
					SELECT 
						*
					FROM tr_nav
					WHERE nav_key = (
					SELECT MAX(nav_key) AS id 
						FROM tr_nav WHERE nav_date <= '` + date + `' AND product_key = '` + productKey + `'
					) ORDER BY nav_key
				) AS nv ON 1=1
			INNER JOIN ms_product AS msp ON nv.product_key = msp.product_key
			GROUP BY nv.product_key`

	// Main query
	// log.Println("========== QUERY GET BEGINNING AND ENDING BALANCE ACCOUNT ==========", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetBeginningEndingBalanceAca(c *BeginningEndingBalance, desc string, date string, acaKey string, productKey string) (int, error) {
	query := `SELECT
				DATE_FORMAT('` + date + `', '%d %b %Y') AS tgl,
				'` + desc + `' AS description,
				(nv.nav_value * SUM(t.balance_unit)) AS amount,
				nv.nav_value,
				SUM(t.balance_unit) AS unit,
				t.avg_nav,
				0 AS fee,
				msp.dec_nav,
				msp.dec_unit,
				msp.dec_amount
			FROM
				(
					SELECT 
						*
					FROM tr_balance 
					WHERE balance_key IN (
					SELECT MAX(balance_key) AS id 
						FROM tr_balance WHERE aca_key  = '` + acaKey + `' 
						AND balance_date <= '` + date + `' 
						AND rec_status = 1 
						GROUP BY tc_key
					) AND rec_status = 1 ORDER BY balance_key
				) AS t 
			LEFT JOIN (
					SELECT 
						*
					FROM tr_nav
					WHERE nav_key = (
					SELECT MAX(nav_key) AS id 
						FROM tr_nav WHERE nav_date <= '` + date + `' AND product_key = '` + productKey + `'
					) ORDER BY nav_key
				) AS nv ON 1=1
			INNER JOIN ms_product AS msp ON nv.product_key = msp.product_key
			GROUP BY nv.product_key`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
