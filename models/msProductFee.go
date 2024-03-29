package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type MsProductFeeInfo struct {
	FeeAnnotation string                 `json:"fee_annotation"`
	FeeDesc       string                 `json:"fee_desc"`
	FeeCode       string                 `json:"fee_code"`
	FeeType       uint64                 `json:"fee_type"`
	FlagShowOntnc uint8                  `json:"flag_show_ontnc"`
	FeeItem       []MsProductFeeItemInfo `json:"fee_item"`
}

type MsProductFee struct {
	FeeKey            uint64           `db:"fee_key"               json:"fee_key"`
	ProductKey        uint64           `db:"product_key"           json:"product_key"`
	FeeType           *uint64          `db:"fee_type"              json:"fee_type"`
	FeeCode           *string          `db:"fee_code"              json:"fee_code"`
	FlagShowOntnc     *uint8           `db:"flag_show_ontnc"       json:"flag_show_ontnc"`
	FeeAnnotation     *string          `db:"fee_annotation"        json:"fee_annotation"`
	FeeDesc           *string          `db:"fee_desc"              json:"fee_desc"`
	FeeDateStart      *string          `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru       *string          `db:"fee_date_thru"         json:"fee_date_thru"`
	FeeNominalType    *uint64          `db:"fee_nominal_type"      json:"fee_nominal_type"`
	EnabledMinAmount  uint8            `db:"enabled_min_amount"    json:"enabled_min_amount"`
	FeeMinAmount      *decimal.Decimal `db:"fee_min_amount"        json:"fee_min_amount"`
	EnabledMaxAmount  uint8            `db:"enabled_max_amount"    json:"enabled_max_amount"`
	FeeMaxAmount      *decimal.Decimal `db:"fee_max_amount"        json:"fee_max_amount"`
	FeeCalcMethod     *uint64          `db:"fee_calc_method"       json:"fee_calc_method"`
	CalculationBaseon *uint64          `db:"calculation_baseon"    json:"calculation_baseon"`
	PeriodHold        uint64           `db:"period_hold"           json:"period_hold"`
	DaysInyear        *uint64          `db:"days_inyear"           json:"days_inyear"`
	RecOrder          *uint64          `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8            `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string          `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string          `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string          `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string          `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string          `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string          `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8           `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64          `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string          `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string          `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string          `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string          `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string          `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string          `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string          `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type MsProductFeeDetailAdmin struct {
	FeeKey            uint64                        `json:"fee_key"`
	Product           MsProductListDropdown         `json:"product"`
	FeeType           *LookupTrans                  `json:"fee_type"`
	FeeCode           *string                       `json:"fee_code"`
	FlagShowOntnc     bool                          `json:"flag_show_ontnc"`
	FeeAnnotation     *string                       `json:"fee_annotation"`
	FeeDesc           *string                       `json:"fee_desc"`
	FeeDateStart      *string                       `json:"fee_date_start"`
	FeeDateThru       *string                       `json:"fee_date_thru"`
	FeeNominalType    *LookupTrans                  `json:"fee_nominal_type"`
	EnabledMinAmount  bool                          `json:"enabled_min_amount"`
	FeeMinAmount      *decimal.Decimal              `json:"fee_min_amount"`
	EnabledMaxAmount  bool                          `json:"enabled_max_amount"`
	FeeMaxAmount      *decimal.Decimal              `json:"fee_max_amount"`
	FeeCalcMethod     *LookupTrans                  `json:"fee_calc_method"`
	CalculationBaseon *LookupTrans                  `json:"calculation_baseon"`
	PeriodHold        uint64                        `json:"period_hold"`
	DaysInyear        *LookupTrans                  `json:"days_inyear"`
	ProductFeeItems   *[]MsProductFeeItemDetailList `json:"product_fee_items"`
}

type AdminListMsProductFee struct {
	FeeKey       uint64  `db:"fee_key"               json:"fee_key"`
	FeeCode      *string `db:"fee_code"              json:"fee_code"`
	ProductKey   uint64  `db:"product_key"           json:"product_key"`
	ProductCode  string  `db:"product_code"          json:"product_code"`
	ProductName  string  `db:"product_name"          json:"product_name"`
	FeeTypeName  *string `db:"feetypename"           json:"feetypename"`
	FeeDateStart *string `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru  *string `db:"fee_date_thru"         json:"fee_date_thru"`
	PeriodHold   uint64  `db:"period_hold"           json:"period_hold"`
	StatusUpdate bool    `db:"status_update" json:"status_update"`
}

func GetAllMsProductFee(c *[]MsProductFee, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_fee.* FROM 
			  ms_product_fee WHERE 
			  ms_product_fee.fee_date_start <= NOW() AND 
			  ms_product_fee.fee_date_thru > NOW() AND 
			  ms_product_fee.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_fee."+field+" = '"+value+"'")
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

func ProductFeeStatusUpdate(FeeKey string) bool {

	query := `SELECT count(*) FROM ms_product_fee_request WHERE rec_approval_status IS NULL AND fee_key = ` + FeeKey
	var get int
	err := db.Db.Get(&get, query)
	if err != nil {
		log.Println(err.Error())
	}
	result := true
	if get > 0 {
		result = false
	}

	return result
}

func AdminGetAllMsProductFee(c *[]AdminListMsProductFee, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT
				pf.fee_key AS fee_key, 
				pf.fee_code AS fee_code, 
				p.product_key AS product_key, 
				p.product_code AS product_code, 
				p.product_name AS product_name, 
				feetype.lkp_name AS feetypename, 
				DATE_FORMAT(pf.fee_date_start, '%d %M %Y') AS fee_date_start, 
				DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') AS fee_date_thru,
				pf.period_hold AS period_hold 
			  FROM ms_product_fee AS pf
			  INNER JOIN ms_product AS p ON p.product_key = pf.product_key
			  LEFT JOIN gen_lookup AS feetype ON feetype.lookup_key = pf.fee_type
			  WHERE pf.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pf.fee_key LIKE '%" + *searchLike + "%' OR"
		condition += " pf.fee_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name LIKE '%" + *searchLike + "%' OR"
		condition += " feetype.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_start, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " pf.period_hold LIKE '%" + *searchLike + "%')"
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
	// log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountDataGetAllMsProductFee(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT
				count(pf.fee_key) AS count_data
			  FROM ms_product_fee AS pf
			  INNER JOIN ms_product AS p ON p.product_key = pf.product_key
			  LEFT JOIN gen_lookup AS feetype ON feetype.lookup_key = pf.fee_type
			  WHERE pf.rec_status = 1`
	var whereClause []string
	var condition string

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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pf.fee_key LIKE '%" + *searchLike + "%' OR"
		condition += " pf.fee_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name LIKE '%" + *searchLike + "%' OR"
		condition += " feetype.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_start, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " pf.period_hold LIKE '%" + *searchLike + "%')"
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

func GetMsProductFee(c *MsProductFee, key string) (int, error) {
	query := `SELECT ms_product_fee.* FROM ms_product_fee WHERE ms_product_fee.rec_status = 1 AND ms_product_fee.fee_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateMsProductFee(params map[string]string) (int, error) {
	query := "UPDATE ms_product_fee SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "fee_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE fee_key = " + params["fee_key"]
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

func CreateMsProductFee(params map[string]string) (int, error, string) {
	query := "INSERT INTO ms_product_fee"
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

type ProductFeeValueSubscription struct {
	FeeNominalType uint64          `db:"fee_nominal_type"    json:"fee_nominal_type"`
	FeeValue       decimal.Decimal `db:"fee_value"           json:"fee_value"`
}

func GetProductFeeValueSubscription(c *ProductFeeValueSubscription, productKey string) (int, error) {
	query := `SELECT 
				pf. fee_nominal_type,
				pfi.fee_value 
			FROM ms_product_fee AS pf 
			LEFT JOIN ms_product_fee_item AS pfi ON pfi.product_fee_key = pf.fee_key AND pfi.rec_status = 1 
			WHERE pf.rec_status = 1 AND pf.fee_type = 183 AND pf.product_key = "` + productKey + `" LIMIT 1`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type FeeItemData struct {
	ProductFeeItemKey uint64          `db:"product_fee_item_key" json:"product_fee_item_key"`
	PrincipleLimit    decimal.Decimal `db:"principle_limit" json:"principle_limit"`
	FeeValue          decimal.Decimal `db:"fee_value" json:"fee_value"`
	ItemNotes         string          `db:"item_notes" json:"item_notes"`
	RecStatus         uint64          `db:"rec_status" json:"rec_status"`
	// ProductFeeKey uint64 `db:"product_fee_key" json:"product_fee_key"`
	// ItemSeqno      string          `db:"item_seqno" json:"item_seqno"`
	// RowMax         uint64          `db:"row_max" json:"row_max"`
	// RecCreatedDate string `db:"rec_created_date" json:"rec_created_date"`
	// RecCreatedBy string `db:"rec_created_by" json:"rec_created_by"`
}

func CreateProductFeeSettings(paramsFee map[string]string, feeItems []FeeItemData) (int, error) {
	query := "INSERT INTO ms_product_fee"
	var fields, values string
	var bindvars []interface{}
	for key, value := range paramsFee {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	// tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}
	lastKey, _ := ret.LastInsertId()

	productFeeKey := strconv.FormatInt(lastKey, 10)

	queryItem := `INSERT INTO ms_product_fee_item(product_fee_key,item_seqno,row_max,principle_limit,fee_value,item_notes,rec_status,rec_created_date,rec_created_by) 
	VALUES`
	for i, data := range feeItems {
		principleLimit := data.PrincipleLimit.String()
		feeValue := data.FeeValue.String()
		itemNotes := data.ItemNotes
		seqNo := strconv.FormatInt(int64(i), 10)
		rowMax := "0"
		if i == len(feeItems)-1 {
			rowMax = "1"
		}
		recStatus := "1"
		recCreatedDate := paramsFee["rec_created_date"]
		recCreatedBy := paramsFee["rec_created_by"]

		queryItem += `('` + productFeeKey + `','` + seqNo + `','` + rowMax + `','` + principleLimit + `','` + feeValue + `','` + itemNotes + `','` + recStatus + `','` + recCreatedDate + `','` + recCreatedBy + `'),`
	}
	queryItem = queryItem[0 : len(queryItem)-1]

	_, err = tx.Exec(queryItem)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}

	tx.Commit()
	return http.StatusOK, nil
}

func UpdateProductFeeSettings(params map[string]string, feeItems []FeeItemData) (int, error) {
	query := "UPDATE ms_product_fee SET "
	i := 0
	for key, value := range params {
		if key != "fee_key" {
			query += key + " = '" + value + "'"
			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE fee_key = " + params["fee_key"]

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	count, err := ret.RowsAffected()

	if count > 0 {
		var queryItem string
		for _, data := range feeItems {
			// itemKey :=
			principleLimit := data.PrincipleLimit.String()
			feeValue := data.FeeValue.String()
			itemNotes := data.ItemNotes
			recstatus := strconv.FormatUint(data.RecStatus, 10)
			seqNo := strconv.FormatInt(int64(i), 10)
			rowMax := "0"
			if i == len(feeItems)-1 {
				rowMax = "1"
			}

			if data.ProductFeeItemKey > 0 {
				// JIKA ADA FEE ITEM KEY MAKA UPDATE
				queryItem = `UPDATE ms_product_fee_item SET 
				principle_limit = ` + principleLimit + `, 
				fee_value = ` + feeValue + `, 
				item_notes = '` + itemNotes + `',
				item_seqno = '` + seqNo + `',
				row_max = '` + rowMax + `', 
				rec_status = '` + recstatus + `', 
				rec_modified_date = '` + params["rec_modified_date"] + `', 
				rec_modified_by = ` + params["rec_modified_by"] + `
				WHERE product_fee_item_key = ` + strconv.FormatUint(data.ProductFeeItemKey, 10)
			} else {
				// JIKA TIDAK ADA FEE ITEM KEY MAKA CREATE
				recstatus = `1`
				queryItem = `INSERT INTO ms_product_fee_item(product_fee_key,item_seqno,row_max,principle_limit,fee_value,item_notes,rec_status,rec_created_date,rec_created_by) 
				VALUES('` + params["product_fee_key"] + `','` + seqNo + `','` + rowMax + `','` + principleLimit + `','` + feeValue + `','` + itemNotes + `','` + recstatus + `','` + params["rec_modified_date"] + `','` + params["rec_modified_by"] + `'),`
			}
			// log.Println("query update fee item:", queryItem)
			_, err = tx.Exec(queryItem)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return http.StatusBadGateway, err
			}
		}
		// log.Println(queryItem)
	} else {
		tx.Rollback()
	}

	tx.Commit()

	return http.StatusOK, nil
}
