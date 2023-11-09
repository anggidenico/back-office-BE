package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type ProductFeeRequest struct {
	RecPK                 *uint64                  `db:"rec_pk" json:"rec_pk"`
	RecAction             *string                  `db:"rec_action" json:"rec_action"`
	FeeKey                *uint64                  `db:"fee_key"               json:"fee_key"`
	ProductKey            *uint64                  `db:"product_key"           json:"product_key"`
	ProductName           *string                  `db:"product_name" json:"product_name"`
	FeeType               *uint64                  `db:"fee_type"              json:"fee_type"`
	FeeTypeName           *string                  `db:"fee_type_name"         json:"fee_type_name"`
	FeeCode               *string                  `db:"fee_code"              json:"fee_code"`
	FlagShowOntnc         *uint64                  `db:"flag_show_ontnc"       json:"flag_show_ontnc"`
	FeeAnnotation         *string                  `db:"fee_annotation"        json:"fee_annotation"`
	FeeDesc               *string                  `db:"fee_desc"              json:"fee_desc"`
	FeeDateStart          *string                  `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru           *string                  `db:"fee_date_thru"         json:"fee_date_thru"`
	FeeNominalType        *uint64                  `db:"fee_nominal_type"      json:"fee_nominal_type"`
	FeeNominalTypeName    *string                  `db:"fee_nominal_type_name" json:"fee_nominal_type_name"`
	EnabledMinAmount      *uint64                  `db:"enabled_min_amount"    json:"enabled_min_amount"`
	FeeMinAmount          *decimal.Decimal         `db:"fee_min_amount"        json:"fee_min_amount"`
	EnabledMaxAmount      *uint64                  `db:"enabled_max_amount"    json:"enabled_max_amount"`
	FeeMaxAmount          *decimal.Decimal         `db:"fee_max_amount"        json:"fee_max_amount"`
	FeeCalcMethod         *uint64                  `db:"fee_calc_method"       json:"fee_calc_method"`
	FeeCalcMethodName     *string                  `db:"fee_calc_method_name"       json:"fee_calc_method_name"`
	CalculationBaseon     *uint64                  `db:"calculation_baseon"    json:"calculation_baseon"`
	CalculationBaseonName *string                  `db:"calculation_baseon_name"    json:"calculation_baseon_name"`
	PeriodHold            *uint64                  `db:"period_hold"           json:"period_hold"`
	DaysInyear            *uint64                  `db:"days_inyear"           json:"days_inyear"`
	DaysInyearName        *string                  `db:"days_inyear_name"           json:"days_inyear_name"`
	RecCreatedDate        *string                  `db:"rec_created_date" json:"rec_created_date"`
	FeeItem               *[]ProductFeeItemRequest `db:"fee_item" json:"fee_item"`
}

type ProductFeeUpdateDetails struct {
	Existing ProductFeeRequest `json:"existing"`
	Updates  ProductFeeRequest `json:"updates"`
}

type ProductFeeItemRequest struct {
	RecPK             *uint64          `db:"rec_pk" json:"rec_pk"`
	RecAction         *string          `db:"rec_action" json:"rec_action"`
	ProductFeeItemKey *uint64          `db:"product_fee_item_key"  json:"product_fee_item_key"`
	ProductFeeKey     *uint64          `db:"product_fee_key"       json:"product_fee_key"`
	ItemSeqno         *uint64          `db:"item_seqno"            json:"item_seqno"`
	RowMax            *uint64          `db:"row_max"               json:"row_max"`
	PrincipleLimit    *decimal.Decimal `db:"principle_limit"       json:"principle_limit"`
	FeeValue          *decimal.Decimal `db:"fee_value"             json:"fee_value"`
	ItemNotes         *string          `db:"item_notes"            json:"item_notes"`
	RecOrder          *uint64          `db:"rec_order"             json:"rec_order"`
	RecStatus         *uint64          `db:"rec_status"            json:"rec_status"`
}

func GetProductFeeApprovalList() []ProductFeeRequest {
	query := `SELECT rec_pk, rec_action, fee_type, product_key, fee_desc, rec_created_date FROM ms_product_fee_request WHERE rec_status = 1 AND rec_approval_status IS NULL`
	var result []ProductFeeRequest

	var get []ProductFeeRequest
	err := db.Db.Select(&get, query)
	if err != nil {
		log.Println(err.Error())
	}

	for _, getdata := range get {
		data := getdata
		product_name := GetForeignKeyValue("ms_product", "product_name", "product_key", *data.ProductKey)
		data.ProductName = &product_name
		fee_type_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *data.FeeType)
		data.FeeTypeName = &fee_type_name
		result = append(result, data)
	}

	return result
}

func GetProductFeeApprovalDetail(rec_pk string) ProductFeeUpdateDetails {
	var result ProductFeeUpdateDetails

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	query := `SELECT rec_pk, rec_action, fee_type, product_key, fee_key, fee_code, flag_show_ontnc, fee_annotation, fee_desc, fee_date_start, fee_date_thru, fee_nominal_type, enabled_min_amount, fee_min_amount, enabled_max_amount, fee_max_amount, fee_calc_method, calculation_baseon, period_hold, days_inyear FROM ms_product_fee_request WHERE rec_pk = ` + rec_pk
	// log.Print(query)
	row := tx.QueryRow(query)
	err = row.Scan(&result.Updates.RecPK, &result.Updates.RecAction, &result.Updates.FeeType, &result.Updates.ProductKey, &result.Updates.FeeKey, &result.Updates.FeeCode, &result.Updates.FlagShowOntnc, &result.Updates.FeeAnnotation, &result.Updates.FeeDesc, &result.Updates.FeeDateStart, &result.Updates.FeeDateThru, &result.Updates.FeeNominalType, &result.Updates.EnabledMinAmount, &result.Updates.FeeMinAmount, &result.Updates.EnabledMaxAmount, &result.Updates.FeeMaxAmount, &result.Updates.FeeCalcMethod, &result.Updates.CalculationBaseon, &result.Updates.PeriodHold, &result.Updates.DaysInyear)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	product_name := GetForeignKeyValue("ms_product", "product_name", "product_key", *result.Updates.ProductKey)
	result.Updates.ProductName = &product_name

	fee_type_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Updates.FeeType)
	result.Updates.FeeTypeName = &fee_type_name

	fee_nominal_type_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Updates.FeeNominalType)
	result.Updates.FeeNominalTypeName = &fee_nominal_type_name

	fee_calc_method_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Updates.FeeCalcMethod)
	result.Updates.FeeCalcMethodName = &fee_calc_method_name

	calculation_baseon_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Updates.CalculationBaseon)
	result.Updates.CalculationBaseonName = &calculation_baseon_name

	days_inyear := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Updates.DaysInyear)
	result.Updates.DaysInyearName = &days_inyear

	if *result.Updates.RecAction == "UPDATE" {

		query2 := `SELECT rec_pk, rec_action, fee_type, product_key, fee_key, fee_code, flag_show_ontnc, fee_annotation, fee_desc, fee_date_start, fee_date_thru, fee_nominal_type, enabled_min_amount, fee_min_amount, enabled_max_amount, fee_max_amount, fee_calc_method, calculation_baseon, period_hold, days_inyear FROM ms_product_fee_request WHERE fee_key = ` + strconv.FormatUint(*result.Updates.FeeKey, 10)
		// log.Print(query)
		row := tx.QueryRow(query2)
		err = row.Scan(&result.Existing.RecPK, &result.Existing.RecAction, &result.Existing.FeeType, &result.Existing.ProductKey, &result.Existing.FeeKey, &result.Existing.FeeCode, &result.Existing.FlagShowOntnc, &result.Existing.FeeAnnotation, &result.Existing.FeeDesc, &result.Existing.FeeDateStart, &result.Existing.FeeDateThru, &result.Existing.FeeNominalType, &result.Existing.EnabledMinAmount, &result.Existing.FeeMinAmount, &result.Existing.EnabledMaxAmount, &result.Existing.FeeMaxAmount, &result.Existing.FeeCalcMethod, &result.Existing.CalculationBaseon, &result.Existing.PeriodHold, &result.Existing.DaysInyear)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
		}

		product_name := GetForeignKeyValue("ms_product", "product_name", "product_key", *result.Existing.ProductKey)
		result.Existing.ProductName = &product_name

		fee_type_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Existing.FeeType)
		result.Existing.FeeTypeName = &fee_type_name

		fee_nominal_type_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Existing.FeeNominalType)
		result.Existing.FeeNominalTypeName = &fee_nominal_type_name

		fee_calc_method_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Existing.FeeCalcMethod)
		result.Existing.FeeCalcMethodName = &fee_calc_method_name

		calculation_baseon_name := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Existing.CalculationBaseon)
		result.Existing.CalculationBaseonName = &calculation_baseon_name

		days_inyear := GetForeignKeyValue("gen_lookup", "lkp_name", "lookup_key", *result.Existing.DaysInyear)
		result.Existing.DaysInyearName = &days_inyear
	}

	return result
}

func ProductFeeCreateRequest(paramsFee map[string]string, feeItems []FeeItemData) (int, error) {
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	query := "INSERT INTO ms_product_fee_request"
	var fields, values string
	var bindvars []interface{}
	for key, value := range paramsFee {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
		log.Println(key)

	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	log.Println(query)

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

	queryItem := `INSERT INTO ms_product_fee_item_request(product_fee_key,item_seqno,row_max,principle_limit,fee_value,item_notes,rec_status,rec_created_date,rec_created_by,rec_action) 
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
		recAction := paramsFee["rec_action"]

		queryItem += `('` + productFeeKey + `','` + seqNo + `','` + rowMax + `','` + principleLimit + `','` + feeValue + `','` + itemNotes + `','` + recStatus + `','` + recCreatedDate + `','` + recCreatedBy + `','` + recAction + `'),`
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

func ProductFeeApprovalAction(params map[string]string) error {
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	recBy := params["rec_by"]
	recDate := params["rec_date"]

	qGetData := `SELECT rec_pk, rec_action, fee_type, product_key, fee_code, flag_show_ontnc, fee_annotation, fee_desc, fee_date_start, fee_date_thru, fee_nominal_type, enabled_min_amount, fee_min_amount, enabled_max_amount, fee_max_amount, fee_calc_method, calculation_baseon, period_hold, days_inyear
	FROM ms_product_fee_request WHERE rec_status = 1 AND rec_pk = ` + params["rec_pk"]

	var pf ProductFeeRequest

	row := tx.QueryRow(qGetData)
	err = row.Scan(&pf.RecPK, &pf.RecAction, &pf.FeeType, &pf.ProductKey, &pf.FeeCode, &pf.FlagShowOntnc, &pf.FeeAnnotation, &pf.FeeDesc, &pf.FeeDateStart, &pf.FeeDateThru, &pf.FeeNominalType, &pf.EnabledMinAmount, &pf.FeeMinAmount, &pf.EnabledMaxAmount, &pf.FeeMaxAmount, &pf.FeeCalcMethod, &pf.CalculationBaseon, &pf.PeriodHold, &pf.DaysInyear)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	query1 := `UPDATE ms_product_request SET rec_approval_status = ` + params["approval"] + ` , rec_approved_date = '` + recDate + `' , rec_approved_by = '` + recBy + `' , rec_attribute_id1 = '` + params["reason"] + `' WHERE rec_pk = ` + params["rec_pk"]

	// log.Println(query1)
	_, err = tx.Exec(query1)
	// log.Println(res.RowsAffected())
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	return nil
}
