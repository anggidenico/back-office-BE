package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type TrNavInfo struct {
	NavDate  string          `json:"nav_date"`
	NavValue decimal.Decimal `json:"nav_value"`
}

type TrNav struct {
	NavKey            uint64          `db:"nav_key"                   json:"nav_key"`
	ProductKey        uint64          `db:"product_key"               json:"product_key"`
	NavDate           string          `db:"nav_date"                  json:"nav_date"`
	NavValue          decimal.Decimal `db:"nav_value"                 json:"nav_value"`
	OriginalValue     decimal.Decimal `db:"original_value"            json:"original_value"`
	NavStatus         uint64          `db:"nav_status"                json:"nav_status"`
	ProdAumTotal      decimal.Decimal `db:"prod_aum_total"            json:"prod_aum_total"`
	ProdUnitTotal     decimal.Decimal `db:"prod_unit_total"           json:"prod_unit_total"`
	PublishMode       uint64          `db:"publish_mode"              json:"publish_mode"`
	RecOrder          *uint64         `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8           `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string         `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string         `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string         `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string         `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string         `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string         `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8          `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64         `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string         `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string         `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string         `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string         `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string         `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string         `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string         `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrNavList struct {
	NavKey        uint64  `json:"nav_key"`
	ProductName   string  `json:"product_name"`
	NavDate       string  `json:"nav_date"`
	NavValue      string  `json:"nav_value"`
	OriginalValue string  `json:"original_value"`
	NavStatus     *string `json:"nav_status"`
	ProdAumTotal  string  `json:"prod_aum_total"`
	ProdUnitTotal string  `json:"prod_unit_total"`
	PublishMode   *string `json:"publish_mode"`
}

type TrNavDetail struct {
	NavKey        uint64          `json:"nav_key"`
	Product       MsProductInfo   `json:"product"`
	NavDate       string          `json:"nav_date"`
	NavValue      decimal.Decimal `json:"nav_value"`
	OriginalValue decimal.Decimal `json:"original_value"`
	NavStatus     LookupTrans     `json:"nav_status"`
	ProdAumTotal  decimal.Decimal `json:"prod_aum_total"`
	ProdUnitTotal decimal.Decimal `json:"prod_unit_total"`
	PublishMode   LookupTrans     `json:"publish_mode"`
}

func GetAllTrNav(c *[]TrNav, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_nav."+field+" = '"+value+"'")
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
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrNav(c *TrNav, key string) (int, error) {
	query := `SELECT tr_nav.* FROM tr_nav WHERE tr_nav.product_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetLastNavIn(c *[]TrNav, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	// query2 := `SELECT t1.nav_key, t1.product_key, t1.nav_date, t1.nav_value FROM
	// 		   tr_nav t1 JOIN (SELECT MAX(nav_key) nav_key FROM tr_nav GROUP BY product_key) t2
	// 		   ON t1.nav_key = t2.nav_key`
	// query := query2 + " WHERE t1.product_key IN(" + inQuery + ") GROUP BY product_key"
	// inQuery = inQuery[:(len(inQuery) - 1)]

	query := `SELECT a.nav_key, a.product_key, a.nav_date, a.nav_value
		FROM tr_nav a
		INNER JOIN (
			SELECT product_key, MAX(nav_date) AS nav_date
			FROM tr_nav
			WHERE rec_status = 1
			AND publish_mode = 236 AND product_key IN(` + inQuery + `)
			GROUP BY product_key
		) b ON (a.product_key = b.product_key AND a.nav_date=b.nav_date)
		WHERE a.rec_status = 1
		AND a.publish_mode = 236 AND a.product_key IN(` + inQuery + `)
		ORDER BY a.nav_date DESC`

	// log.Println("========== QUERY SELECT LAST NAV =========", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTrNavBetween(c *[]TrNav, start string, end string, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	query += " WHERE tr_nav.product_key IN(" + inQuery + ") AND tr_nav.nav_date BETWEEN '" + start + "' AND '" + end + "'"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrNav1D(c *[]TrNav, productKey string) (int, error) {
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	query += " WHERE tr_nav.product_key=" + productKey + " ORDER BY tr_nav.nav_key DESC LIMIT 2 "

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrNavIn(c *[]TrNav, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_nav.* FROM 
				tr_nav `
	query := query2 + " WHERE tr_nav." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrNavByProductKeyAndNavDate(c *[]TrNav, productKey string, navDate string) (int, error) {
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	query += " WHERE tr_nav.rec_status = 1 AND tr_nav.product_key = '" + productKey + "' AND tr_nav.nav_date = '" + navDate + "'"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetNavByProductKeyAndNavDate(c *TrNav, productKey string, navDate string) (int, error) {
	query := `SELECT t1.* FROM tr_nav t1 WHERE t1.rec_status = 1 AND t1.product_key = ` + productKey + ` AND t1.nav_date = "` + navDate + `"`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTrNavCount(c *CountData, params map[string]string) (int, error) {
	query := `SELECT 
			  count(tr_nav.nav_key) as count_data 
              FROM tr_nav`
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_nav."+field+" = '"+value+"'")
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

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrNavByKey(c *TrNav, key string) (int, error) {
	query := `SELECT tr_nav.* FROM tr_nav WHERE tr_nav.nav_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateTrNav(params map[string]string) (int, error) {
	query := "UPDATE tr_nav SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "nav_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE nav_key = " + params["nav_key"]
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

func CreateTrNav(params map[string]string) (int, error) {
	query := "INSERT INTO tr_nav"
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

func GetNavByProductKeyAndNavDateExcept(c *TrNav, productKey string, navDate string, keyExcept string) (int, error) {
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	query += " WHERE tr_nav.rec_status =" +
		" 1 AND tr_nav.product_key = '" + productKey + "'" +
		" AND tr_nav.nav_date = '" + navDate + "' AND tr_nav.nav_key != '" + keyExcept + "'"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminLastNavValue(c *NavValue, productKey string, date string) (int, error) {
	query := `SELECT 
				nav_value 
			FROM tr_nav
			WHERE nav_date <= '` + date + `' AND product_key = '` + productKey + `' AND rec_status = '1' AND nav_status = '231'
			ORDER BY nav_date DESC LIMIT 1`

	// Main query
	// log.Println("========== QUERY GET LAST NAV VALUE ==========", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type TanggalBursa struct {
	TanggalBursa string `db:"date_result"                   json:"date_result"`
}

func GetTanggalBursa(c *TanggalBursa, date string, addDate string) (int, error) {
	query := `SELECT fn_AddDate('` + date + `', ` + addDate + `) AS date_result`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
