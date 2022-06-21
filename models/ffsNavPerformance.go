package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type FfsNavPerformanceInfo struct {
	NavDate string `json:"nav_date,omitempty"`
	D1      string `json:"1d,omitempty"`
	MTD     string `json:"mtd,omitempty"`
	M1      string `json:"1m,omitempty"`
	M3      string `json:"3m,omitempty"`
	M6      string `json:"6m,omitempty"`
	Y1      string `json:"1y,omitempty"`
	Y3      string `json:"3y,omitempty"`
	Y5      string `json:"5y,omitempty"`
	YTD     string `json:"ytd,omitempty"`
	CAGR    string `json:"cagr,omitempty"`
	ALL     string `json:"all,omitempty"`
}

type PublicNavPerformance struct {
	LastDate       string          `db:"last_date"             json:"last_date"`
	NavPerformKey  uint64          `db:"nav_perform_key"       json:"nav_perform_key"`
	FundTypeKey    uint64          `db:"fund_type_key"         json:"fund_type_key"`
	FundTypeCode   string          `db:"fund_type_code"        json:"fund_type_code"`
	FundTypeName   string          `db:"fund_type_name"        json:"fund_type_name"`
	ProductKey     uint64          `db:"product_key"           json:"product_key"`
	ProductName    string          `db:"product_name"          json:"product_name"`
	ProductNameAlt string          `db:"product_name_alt"      json:"product_name_alt"`
	NavDate        string          `db:"nav_date"              json:"nav_date"`
	CurrencyCode   string          `db:"currency_code"         json:"currency_code"`
	NavValue       decimal.Decimal `db:"nav_value"             json:"nav_value"`
	PerformD1      string          `db:"perform_d1"            json:"perform_d1"`
	PerformMtd     string          `db:"perform_mtd"           json:"perform_mtd"`
	PerformM1      string          `db:"perform_m1"            json:"perform_m1"`
	PerformM3      string          `db:"perform_m3"            json:"perform_m3"`
	PerformM6      string          `db:"perform_m6"            json:"perform_m6"`
	PerformYtd     string          `db:"perform_ytd"           json:"perform_ytd"`
	PerformY1      string          `db:"perform_y1"            json:"perform_y1"`
	PerformY3      string          `db:"perform_y3"            json:"perform_y3"`
	PerformY5      string          `db:"perform_y5"            json:"perform_y5"`
	PerformCagr    string          `db:"perform_cagr"          json:"perform_cagr"`
	PerformAll     string          `db:"perform_all"           json:"perform_all"`
}

type FfsNavPerformance struct {
	NavPerformKey     uint64          `db:"nav_perform_key"       json:"nav_perform_key"`
	ProductKey        uint64          `db:"product_key"           json:"product_key"`
	PeriodeKey        *uint64         `db:"periode_key"           json:"periode_key"`
	NavDate           string          `db:"nav_date"              json:"nav_date"`
	NavD0             decimal.Decimal `db:"nav_d0"                json:"nav_d0"`
	NavD1             decimal.Decimal `db:"nav_d1"                json:"nav_d1"`
	NavM0             decimal.Decimal `db:"nav_m0"                json:"nav_m0"`
	NavM1             decimal.Decimal `db:"nav_m1"                json:"nav_m1"`
	NavM3             decimal.Decimal `db:"nav_m3"                json:"nav_m3"`
	NavM6             decimal.Decimal `db:"nav_m6"                json:"nav_m6"`
	NavYtd            decimal.Decimal `db:"nav_ytd"               json:"nav_ytd"`
	Navy1             decimal.Decimal `db:"nav_y1"                json:"nav_y1"`
	Navy3             decimal.Decimal `db:"nav_y3"                json:"nav_y3"`
	Navy5             decimal.Decimal `db:"nav_y5"                json:"nav_y5"`
	PerformD1         decimal.Decimal `db:"perform_d1"            json:"perform_d1"`
	PerformMtd        decimal.Decimal `db:"perform_mtd"           json:"perform_mtd"`
	PerformM1         decimal.Decimal `db:"perform_m1"            json:"perform_m1"`
	PerformM3         decimal.Decimal `db:"perform_m3"            json:"perform_m3"`
	PerformM6         decimal.Decimal `db:"perform_m6"            json:"perform_m6"`
	PerformYtd        decimal.Decimal `db:"perform_ytd"           json:"perform_ytd"`
	PerformY1         decimal.Decimal `db:"perform_y1"            json:"perform_y1"`
	PerformY3         decimal.Decimal `db:"perform_y3"            json:"perform_y3"`
	PerformY5         decimal.Decimal `db:"perform_y5"            json:"perform_y5"`
	PerformCagr       decimal.Decimal `db:"perform_cagr"          json:"perform_cagr"`
	PerformAll        decimal.Decimal `db:"perform_all"           json:"perform_all"`
	RecOrder          *uint64         `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8           `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string         `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string         `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string         `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string         `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string         `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string         `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8          `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64         `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string         `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string         `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string         `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string         `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string         `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string         `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string         `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

func GetAllFfsNavPerformance(c *[]FfsNavPerformance, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ffs_nav_performance.* FROM 
			  ffs_nav_performance`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ffs_nav_performance."+field+" = '"+value+"'")
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

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetLastNavPerformanceIn(c *[]FfsNavPerformance, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query2 := `SELECT 
				t1.nav_perform_key, 
				t1.product_key, nav_date, 
				t1.perform_d1,
				t1.perform_mtd,
				t1.perform_m1,
				t1.perform_m3,
				t1.perform_m6,
				t1.perform_ytd,
				t1.perform_y1,
				t1.perform_y3,
				t1.perform_y5,
				t1.perform_cagr,
				t1.perform_all FROM
				ffs_nav_performance t1 JOIN (SELECT MAX(nav_perform_key) nav_perform_key FROM ffs_nav_performance GROUP BY product_key) t2
				ON t1.nav_perform_key = t2.nav_perform_key`
	query := query2 + " WHERE t1.product_key IN(" + inQuery + ")"

	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllLastNavPerformance(c *[]FfsNavPerformance, params map[string]string) (int, error) {
	query := `SELECT 
			t1.nav_perform_key, 
			t1.product_key, nav_date, 
			t1.perform_d1,
			t1.perform_mtd,
			t1.perform_m1,
			t1.perform_m3,
			t1.perform_m6,
			t1.perform_ytd,
			t1.perform_y1,
			t1.perform_y3,
			t1.perform_y5,
			t1.perform_cagr,
			t1.perform_all FROM
			ffs_nav_performance t1 JOIN (SELECT MAX(nav_perform_key) nav_perform_key FROM ffs_nav_performance GROUP BY product_key) t2
			ON t1.nav_perform_key = t2.nav_perform_key`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ffs_nav_performance."+field+" = '"+value+"'")
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
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetPublicNavPerformance(c *[]PublicNavPerformance) (int, error) {
	query := `SELECT 
				t6.last_date AS last_date,
				t1.nav_perform_key, 
				t4.fund_type_key,
				t4.fund_type_code,
				t4.fund_type_name,
				t1.product_key, 
				t3.product_name,
				t3.product_name_alt,
				cast(t1.nav_date AS DATE) AS nav_date,
				t5.code AS currency_code,
				t1.nav_d0 as nav_value, 
				CONCAT(truncate(t1.perform_d1,2),'%') AS perform_d1,
				CONCAT(truncate(t1.perform_mtd,2),'%') AS perform_mtd,
				CONCAT(truncate(t1.perform_m1,2),'%') AS perform_m1,
				CONCAT(truncate(t1.perform_m3,2),'%') AS perform_m3,
				CONCAT(truncate(t1.perform_m6,2),'%') AS  perform_m6,
				CONCAT(truncate(t1.perform_ytd,2),'%') AS perform_ytd,
				CONCAT(truncate(t1.perform_y1,2),'%') AS perform_y1,
				CONCAT(truncate(t1.perform_y3,2),'%') AS perform_y3,
				CONCAT(truncate(t1.perform_y5,2),'%') AS perform_y5,
				CONCAT(truncate(t1.perform_cagr,2),'%') AS perform_cagr,
				CONCAT(truncate(t1.perform_all,2),'%') AS perform_all
				FROM ffs_nav_performance t1 
				INNER JOIN (
					SELECT MAX(nav_date) as nav_date, product_key 
					FROM ffs_nav_performance 
					WHERE rec_status = 1
					AND cast(nav_date AS DATE) <= CURRENT_DATE 
					GROUP BY product_key
				) t2 ON (t1.nav_date = t2.nav_date AND t1.product_key = t2.product_key)
				INNER JOIN ms_product t3 ON (t3.product_key = t2.product_key)
				INNER JOIN ms_currency t5 ON (t5.currency_key = t3.currency_key)
				INNER JOIN ms_fund_type t4 ON (t4.fund_type_key = t3.fund_type_key)
				CROSS JOIN (
					SELECT MAX(cast(nav_date AS DATE)) as last_date
					FROM ffs_nav_performance 
					WHERE cast(nav_date AS DATE) <= CURRENT_DATE 
				) t6
			WHERE t3.rec_status = 1
			AND t3.flag_enabled = 1
			AND t3.product_phase = 177
			ORDER BY t4.rec_order, t1.product_key`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
