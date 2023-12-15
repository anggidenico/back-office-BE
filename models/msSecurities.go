package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type MsSecurities struct {
	SecKey             uint64  `db:"sec_key"               json:"sec_key"`
	SecCode            string  `db:"sec_code"              json:"sec_code"`
	SecName            string  `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory uint64  `db:"securities_category"   json:"securities_category"`
	SecurityType       uint64  `db:"security_type"         json:"security_type"`
	DateIssued         *string `db:"date_issued"           json:"date_issued"`
	DateMatured        *string `db:"date_matured"          json:"date_matured"`
	CurrencyKey        *uint64 `db:"currency_key"          json:"currency_key"`
	SecurityStatus     uint64  `db:"security_status"       json:"security_status"`
	IsinCode           *string `db:"isin_code"             json:"isin_code"`
	SecClassification  uint64  `db:"sec_classification"    json:"sec_classification"`
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
type Securities struct {
	SecKey       uint64  `db:"sec_key"               json:"sec_key"`
	SecParentKey *uint64 `db:"sec_parent_key" json:"sec_parent_key"`
	// ParentSectorName       *string          `db:"parent_sector_name" json:"parent_sector_name"`
	SecCode                string           `db:"sec_code"              json:"sec_code"`
	SecName                string           `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory     uint64           `db:"securities_category"   json:"securities_category"`
	SecuritiesCategoryName string           `db:"securities_category_name" json:"securities_category_name"`
	SecurityType           uint64           `db:"security_type"         json:"security_type"`
	SecurityTypeName       string           `db:"security_type_name"         json:"security_type_name"`
	SectorKey              *uint64          `db:"sector_key" json:"sector_key"`
	SectorName             *string          `db:"sector_name" json:"sector_name"`
	DateIssued             *string          `db:"date_issued"           json:"date_issued"`
	DateMatured            *string          `db:"date_matured"          json:"date_matured"`
	SecTenorMonth          *int64           `db:"sec_tenor_month" json:"sec_tenor_month"`
	CurrencyKey            *uint64          `db:"currency_key"          json:"currency_key"`
	CurrencyCode           *string          `db:"currency_code"          json:"currency_code"`
	CurrencyName           *string          `db:"currency_name"          json:"currency_name"`
	SecurityStatus         *uint64          `db:"security_status"       json:"security_status"`
	SecurityStatusName     *string          `db:"security_status_name"       json:"security_status_name"`
	IsinCode               *string          `db:"isin_code"             json:"isin_code"`
	SecClassification      *uint64          `db:"sec_classification"    json:"sec_classification"`
	SecClassificationName  *string          `db:"sec_classification_name"    json:"sec_classification_name"`
	SecShares              *uint64          `db:"sec_shares" json:"sec_shares"`
	FlagSyariah            []byte           `db:"flag_syariah" json:"flag_syariah"`
	StockMarket            *uint64          `db:"stock_market" json:"stock_market"`
	StocKMarketName        *string          `db:"stock_market_name" json:"stock_market_name"`
	SecParates             *decimal.Decimal `db:"sec_pa_rates" json:"sec_pa_rates"`
	SecPrincipleValue      *decimal.Decimal `db:"sec_principle_value" json:"sec_principle_value"`
	TaxRates               *decimal.Decimal `db:"tax_rates" json:"tax_rates"`
	ParticipantKey         *int64           `db:"participant_key" json:"participant_key"`
	ParticipantName        *string          `db:"participant_name" json:"participant_name"`
	FlagHasCoupon          []byte           `db:"flag_has_coupon" json:"flag_has_coupon"`
	CouponType             *int64           `db:"coupon_type" json:"coupon_type"`
	CouponName             *string          `db:"coupon_name" json:"coupon_name"`
	FlagIsBreakable        []byte           `db:"flag_is_breakable" json:"flag_is_breakable"`
	RecOrder               *int64           `db:"rec_order" json:"rec_order"`
}

// Value converts Int64Bool to a database value.

type SecuritiesResponse struct {
	SecKey       uint64  `db:"sec_key"               json:"sec_key"`
	SecParentKey *uint64 `db:"sec_parent_key" json:"sec_parent_key"`
	// ParentSectorName       *string          `db:"parent_sector_name" json:"parent_sector_name"`
	SecCode                string           `db:"sec_code"              json:"sec_code"`
	SecName                string           `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory     uint64           `db:"securities_category"   json:"securities_category"`
	SecuritiesCategoryName string           `db:"securities_category_name" json:"securities_category_name"`
	SecurityType           uint64           `db:"security_type"         json:"security_type"`
	SecurityTypeName       string           `db:"security_type_name"         json:"security_type_name"`
	SectorKey              *uint64          `db:"sector_key" json:"sector_key"`
	SectorName             *string          `db:"sector_name" json:"sector_name"`
	DateIssued             *string          `db:"date_issued"           json:"date_issued"`
	DateMatured            *string          `db:"date_matured"          json:"date_matured"`
	SecTenorMonth          *int64           `db:"sec_tenor_month" json:"sec_tenor_month"`
	CurrencyKey            *uint64          `db:"currency_key"          json:"currency_key"`
	CurrencyCode           *string          `db:"currency_code"          json:"currency_code"`
	CurrencyName           *string          `db:"currency_name"          json:"currency_name"`
	SecurityStatus         *uint64          `db:"security_status"       json:"security_status"`
	SecurityStatusName     *string          `db:"security_status_name"       json:"security_status_name"`
	IsinCode               *string          `db:"isin_code"             json:"isin_code"`
	SecClassification      *uint64          `db:"sec_classification"    json:"sec_classification"`
	SecClassificationName  *string          `db:"sec_classification_name"    json:"sec_classification_name"`
	SecShares              *uint64          `db:"sec_shares" json:"sec_shares"`
	FlagSyariah            *bool            `db:"flag_syariah" json:"flag_syariah"`
	StockMarket            *uint64          `db:"stock_market" json:"stock_market"`
	StocKMarketName        *string          `db:"stock_market_name" json:"stock_market_name"`
	SecParates             *decimal.Decimal `db:"sec_pa_rates" json:"sec_pa_rates"`
	SecPrincipleValue      *decimal.Decimal `db:"sec_principle_value" json:"sec_principle_value"`
	TaxRates               *decimal.Decimal `db:"tax_rates" json:"tax_rates"`
	ParticipantKey         *int64           `db:"participant_key" json:"participant_key"`
	ParticipantName        *string          `db:"participant_name" json:"participant_name"`
	FlagHasCoupon          *bool            `db:"flag_has_coupon" json:"flag_has_coupon"`
	CouponType             *int64           `db:"coupon_type" json:"coupon_type"`
	CouponName             *string          `db:"coupon_name" json:"coupon_name"`
	FlagIsBreakable        *bool            `db:"flag_is_breakable" json:"flag_is_breakable"`
	RecOrder               *int64           `db:"rec_order" json:"rec_order"`
}

type SecuritiesDetail struct {
	SecKey       uint64  `db:"sec_key"               json:"sec_key"`
	SecParentKey *uint64 `db:"sec_parent_key" json:"sec_parent_key"`
	// ParentSectorName       *string          `db:"parent_sector_name" json:"parent_sector_name"`
	SecCode                string           `db:"sec_code"              json:"sec_code"`
	SecName                string           `db:"sec_name"              json:"sec_name"`
	SecuritiesCategory     uint64           `db:"securities_category"   json:"securities_category"`
	SecuritiesCategoryName string           `db:"securities_category_name" json:"securities_category_name"`
	SecurityType           uint64           `db:"security_type"         json:"security_type"`
	SecurityTypeName       string           `db:"security_type_name"         json:"security_type_name"`
	SectorKey              *uint64          `db:"sector_key" json:"sector_key"`
	SectorName             *string          `db:"sector_name" json:"sector_name"`
	DateIssued             *string          `db:"date_issued"           json:"date_issued"`
	DateMatured            *string          `db:"date_matured"          json:"date_matured"`
	SecTenorMonth          *int64           `db:"sec_tenor_month" json:"sec_tenor_month"`
	CurrencyKey            *uint64          `db:"currency_key"          json:"currency_key"`
	CurrencyCode           *string          `db:"currency_code"          json:"currency_code"`
	CurrencyName           *string          `db:"currency_name"          json:"currency_name"`
	SecurityStatus         *uint64          `db:"security_status"       json:"security_status"`
	SecurityStatusName     *string          `db:"security_status_name"       json:"security_status_name"`
	IsinCode               *string          `db:"isin_code"             json:"isin_code"`
	SecClassification      *uint64          `db:"sec_classification"    json:"sec_classification"`
	SecClassificationName  *string          `db:"sec_classification_name"    json:"sec_classification_name"`
	SecShares              *uint64          `db:"sec_shares" json:"sec_shares"`
	FlagSyariah            []byte           `db:"flag_syariah" json:"flag_syariah"`
	StockMarket            *uint64          `db:"stock_market" json:"stock_market"`
	StocKMarketName        *string          `db:"stock_market_name" json:"stock_market_name"`
	SecParates             *decimal.Decimal `db:"sec_pa_rates" json:"sec_pa_rates"`
	SecPrincipleValue      *decimal.Decimal `db:"sec_principle_value" json:"sec_principle_value"`
	TaxRates               *decimal.Decimal `db:"tax_rates" json:"tax_rates"`
	ParticipantKey         *int64           `db:"participant_key" json:"participant_key"`
	ParticipantName        *string          `db:"participant_name" json:"participant_name"`
	FlagHasCoupon          []byte           `db:"flag_has_coupon" json:"flag_has_coupon"`
	CouponType             *int64           `db:"coupon_type" json:"coupon_type"`
	CouponName             *string          `db:"coupon_name" json:"coupon_name"`
	FlagIsBreakable        []byte           `db:"flag_is_breakable" json:"flag_is_breakable"`
	RecOrder               *int64           `db:"rec_order" json:"rec_order"`
}
type ParticipantList struct {
	ParticipantKey  int64  `db:"participant_key" json:"participant_key"`
	ParticipantCode string `db:"participant_code" json:"participant_code"`
	ParticipantName string `db:"participant_name" json:"participant"`
}

func GetParticipantListModels(c *[]ParticipantList) (int, error) {
	query := `SELECT participant_key,
	participant_code,
	participant_name
	FROM ms_participant
WHERE rec_status = 1
ORDER BY participant_key ASC`

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

func GetSecuritiesModels(c *[]Securities) (int, error) {
	query := `SELECT
    a.sec_key,
    a.sec_parent_key,
    a.sec_code,
    a.sec_name,
    a.securities_category,
    b.lkp_name securities_category_name,
    a.security_type,
    c.lkp_name security_type_name,
    a.sector_key,
    g.sector_name,
    a.date_issued,
    a.date_matured,
    a.sec_tenor_month,
    a.currency_key,
    e.code currency_code,
    e.name currency_name,
    a.security_status,
    d.lkp_name security_status_name,
    a.isin_code,
    a.sec_classification,
    f.lkp_name sec_classification_name,
    a.sec_shares,
    a.flag_syariah,
    a.stock_market,
    h.lkp_name stock_market_name,
    a.sec_pa_rates,
    a.sec_principle_value,
    a.tax_rates,
    a.participant_key,
    i.participant_name,
    a.flag_has_coupon,
    a.coupon_type,
    j.lkp_name coupon_name,
    a.flag_is_breakable,
    a.rec_order
FROM ms_securities a
LEFT JOIN gen_lookup b ON a.securities_category = b.lookup_key
LEFT JOIN gen_lookup c ON a.security_type = c.lookup_key
LEFT JOIN gen_lookup d ON a.security_status = d.lookup_key
LEFT JOIN ms_currency e ON a.currency_key = e.currency_key
LEFT JOIN gen_lookup f ON a.sec_classification = f.lookup_key
LEFT JOIN ms_securities_sector g ON a.sector_key = g.sector_key
LEFT JOIN gen_lookup h ON a.stock_market = h.lookup_key
LEFT JOIN ms_participant i ON a.participant_key = i.participant_key
LEFT JOIN gen_lookup j ON a.coupon_type = j.lookup_key
WHERE a.rec_status = 1
ORDER BY a.sec_key DESC;`

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

func DeleteMsSecurities(SecKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_securities SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "sec_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sec_key = ?`
	values = append(values, SecKey)

	// log.Println("========== DeleteMsSecurities ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
func GetMsSecuritiesDetailModels(c *SecuritiesDetail, SecKey string) (int, error) {
	query := `SELECT
    a.sec_key,
    a.sec_parent_key,
    a.sec_code,
    a.sec_name,
    a.securities_category,
    b.lkp_name securities_category_name,
    a.security_type,
    c.lkp_name security_type_name,
    a.sector_key,
    g.sector_name,
    a.date_issued,
    a.date_matured,
    a.sec_tenor_month,
    a.currency_key,
    e.code currency_code,
    e.name currency_name,
    a.security_status,
    d.lkp_name security_status_name,
    a.isin_code,
    a.sec_classification,
    f.lkp_name sec_classification_name,
    a.sec_shares,
    a.flag_syariah,
    a.stock_market,
    h.lkp_name stock_market_name,
    a.sec_pa_rates,
    a.sec_principle_value,
    a.tax_rates,
    a.participant_key,
    i.participant_name,
    a.flag_has_coupon,
    a.coupon_type,
    j.lkp_name coupon_name,
    a.flag_is_breakable,
    a.rec_order
FROM ms_securities a
LEFT JOIN gen_lookup b ON a.securities_category = b.lookup_key
LEFT JOIN gen_lookup c ON a.security_type = c.lookup_key
LEFT JOIN gen_lookup d ON a.security_status = d.lookup_key
LEFT JOIN ms_currency e ON a.currency_key = e.currency_key
LEFT JOIN gen_lookup f ON a.sec_classification = f.lookup_key
LEFT JOIN ms_securities_sector g ON a.sector_key = g.sector_key
LEFT JOIN gen_lookup h ON a.stock_market = h.lookup_key
LEFT JOIN ms_participant i ON a.participant_key = i.participant_key
LEFT JOIN gen_lookup j ON a.coupon_type = j.lookup_key
	WHERE a.rec_status = 1 
	AND a.sec_key =` + SecKey

	// log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func UpdateMsSecurities(SecKey string, params map[string]interface{}) (int, error) {
	query := `UPDATE ms_securities SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "sec_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sec_key = ?`
	values = append(values, SecKey)

	// log.Println("========== UpdateMsSecurities ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func CreateMsSecurities(params map[string]string) (int, error) {
	// Check for duplicate records

	// log.Println(params)

	tx, err := db.Db.Begin()
	if err != nil {
		return http.StatusBadGateway, err
	}

	QueryCekDuplicate := `SELECT COUNT(*) FROM ms_securities WHERE rec_status = 1 AND ( sec_code = '` + params["sec_code"] + `' OR sec_name = '` + params["sec_name"] + `' OR security_type = '` + params["security_type"] + `' )`

	var CountDup int64
	err = db.Db.Get(&CountDup, QueryCekDuplicate)
	if err != nil {
		tx.Rollback()
		return http.StatusBadGateway, err
	}
	if CountDup > 0 {

		queryUpdate := GenerateUpdateQuery("ms_securities", "sec_code", params)
		_, err = tx.Exec(queryUpdate)
		if err != nil {
			tx.Rollback()
			return http.StatusBadGateway, err
		}

		queryUpdate = GenerateUpdateQuery("ms_securities", "sec_name", params)
		_, err = tx.Exec(queryUpdate)
		if err != nil {
			tx.Rollback()
			return http.StatusBadGateway, err
		}

		queryUpdate = GenerateUpdateQuery("ms_securities", "security_type", params)
		_, err = tx.Exec(queryUpdate)
		if err != nil {
			tx.Rollback()
			return http.StatusBadGateway, err
		}

	} else {

		queryInsert := GenerateInsertQuery("ms_securities", params)
		_, err = tx.Exec(queryInsert)
		if err != nil {
			tx.Rollback()
			return http.StatusBadGateway, err
		}

	}

	// duplicate, _, err := CheckDuplicateSecurities(params["sec_code"].(string), params["sec_name"].(string), params["security_type"].(string))
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// }

	// // Jika duplikasi ditemukan, perbarui data yang sudah ada
	// if duplicate {
	// 	return http.StatusBadRequest, errors.New("data duplikat ditemukan")
	// }

	// Jika tidak ada duplikasi, buat data baru
	// fields := ""
	// placeholders := ""
	// var bindvars []interface{}

	// for key, value := range params {
	// 	fields += key + `, `
	// 	if value == "" {
	// 		placeholders += `NULL, `
	// 	} else {
	// 		placeholders += `?, `
	// 		bindvars = append(bindvars, value)
	// 	}
	// }

	// fields = fields[:len(fields)-2]
	// placeholders = placeholders[:len(placeholders)-2]

	// query := "INSERT INTO ms_securities (" + fields + ") VALUES (" + placeholders + ")"

	// _, err = tx.Exec(query, bindvars...)
	// if err != nil {
	// 	tx.Rollback()
	// 	return http.StatusBadRequest, err
	// }

	tx.Commit()

	return http.StatusOK, nil
}

func CheckDuplicateSecurities(SecCode, SecName, SecType string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT sec_key FROM ms_securities WHERE sec_code = ? AND sec_name = ? AND security_type = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, SecCode, SecName, SecType).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			// No duplicate found
			return false, "", nil
		}
		// Other error occurred
		return false, "", err
	}

	// Duplicate found
	return true, key, nil
}

func GetSecuritiesStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM ms_securities WHERE sec_key = ?"
	var status int
	err := db.Db.QueryRow(query, key).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Data tidak ditemukan
			return 0, nil
		}
		// Terjadi error lain
		return 0, err
	}
	return status, nil
}
