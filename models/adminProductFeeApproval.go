package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"reflect"
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

	qFeeItem := `SELECT item_seqno, row_max, principle_limit, fee_value, item_notes FROM ms_product_fee_item_request WHERE product_fee_key = ` + strconv.FormatUint(*result.Updates.FeeKey, 10)
	rows, err := tx.Query(qFeeItem)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}
	// var FeeItems []ProductFeeItemRequest
	for rows.Next() {
		var Item ProductFeeItemRequest
		if err := rows.Scan(&Item.ItemSeqno, &Item.RowMax, &Item.PrincipleLimit, &Item.FeeValue, &Item.ItemNotes); err != nil {
			tx.Rollback()
			log.Println(err.Error())
		}
		*result.Updates.FeeItem = append(*result.Updates.FeeItem, Item)
	}

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

		qFeeItem2 := `SELECT item_seqno, row_max, principle_limit, fee_value, item_notes FROM ms_product_fee_item_request WHERE product_fee_key = ` + strconv.FormatUint(*result.Existing.RecPK, 10)
		rows, err := tx.Query(qFeeItem2)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
		}
		// var FeeItems []ProductFeeItemRequest
		for rows.Next() {
			var Item ProductFeeItemRequest
			if err := rows.Scan(&Item.ItemSeqno, &Item.RowMax, &Item.PrincipleLimit, &Item.FeeValue, &Item.ItemNotes); err != nil {
				tx.Rollback()
				log.Println(err.Error())
			}
			*result.Existing.FeeItem = append(*result.Existing.FeeItem, Item)
		}
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

	// log.Println(query)

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

	qGetData := `SELECT rec_pk, rec_action, fee_key, fee_type, product_key, fee_code, flag_show_ontnc, fee_annotation, fee_desc, fee_date_start, fee_date_thru, fee_nominal_type, enabled_min_amount, fee_min_amount, enabled_max_amount, fee_max_amount, fee_calc_method, calculation_baseon, period_hold, days_inyear
	FROM ms_product_fee_request WHERE rec_status = 1 AND rec_pk = ` + params["rec_pk"]

	var pf ProductFeeRequest

	row := tx.QueryRow(qGetData)
	err = row.Scan(&pf.RecPK, &pf.RecAction, &pf.FeeKey, &pf.FeeType, &pf.ProductKey, &pf.FeeCode, &pf.FlagShowOntnc, &pf.FeeAnnotation, &pf.FeeDesc, &pf.FeeDateStart, &pf.FeeDateThru, &pf.FeeNominalType, &pf.EnabledMinAmount, &pf.FeeMinAmount, &pf.EnabledMaxAmount, &pf.FeeMaxAmount, &pf.FeeCalcMethod, &pf.CalculationBaseon, &pf.PeriodHold, &pf.DaysInyear)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	qGetFeeItemReq := `SELECT item_seqno, row_max, principle_limit, fee_value, item_notes FROM ms_product_fee_item_request WHERE product_fee_key = ` + params["rec_pk"]
	rowsItem, err := tx.Query(qGetFeeItemReq)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}
	for rowsItem.Next() {
		var Item ProductFeeItemRequest
		if err := rowsItem.Scan(&Item.ItemSeqno, &Item.RowMax, &Item.PrincipleLimit, &Item.FeeValue, &Item.ItemNotes); err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err
		}
		*pf.FeeItem = append(*pf.FeeItem, Item)
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

	if params["approval"] == "1" {
		if *pf.RecAction == "CREATE" {
			insertMsProductFee := make(map[string]string)

			var reflectValue = reflect.ValueOf(pf)
			if reflectValue.Kind() == reflect.Ptr {
				reflectValue = reflectValue.Elem()
			}
			var reflectType = reflectValue.Type()
			for i := 0; i < reflectValue.NumField(); i++ {
				columnName := reflectType.Field(i).Tag.Get("db")
				value := reflectValue.Field(i).Interface()
				if val, ok := value.(*uint64); ok {
					if val != nil {
						columnValue := strconv.FormatUint(*val, 10)
						insertMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*uint8); ok {
					if val != nil {
						columnValue := strconv.FormatUint(uint64(*val), 10)
						insertMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*string); ok {
					if val != nil {
						columnValue := *val
						insertMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*decimal.Decimal); ok {
					if val != nil {
						columnValue := val.String()
						insertMsProductFee[columnName] = columnValue
					}
				}
			}
			insertMsProductFee["rec_status"] = "1"
			insertMsProductFee["rec_created_by"] = recBy
			insertMsProductFee["rec_created_date"] = recDate
			var fields, values string
			var bindvars []interface{}
			for key, value := range insertMsProductFee {
				if key != "rec_pk" && key != "rec_action" {
					fields += key + ", "
					values += ` "` + value + `", `
					bindvars = append(bindvars, value)
				}
			}
			fields = fields[:(len(fields) - 2)]
			values = values[:(len(values) - 2)]
			query := `INSERT INTO ms_product_fee (` + fields + `) VALUES(` + values + `)`

			// log.Println(query)
			res, err := tx.Exec(query)
			// log.Println(res.LastInsertId())
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}
			lastKey, err := res.LastInsertId()
			productFeeKey := strconv.FormatInt(lastKey, 10)

			// INPUT KE MS FEE ITEMS
			if len(*pf.FeeItem) > 0 {
				queryItem := `INSERT INTO ms_product_fee_item(product_fee_key,item_seqno,row_max,principle_limit,fee_value,item_notes,rec_status,rec_created_date,rec_created_by) VALUES`
				for i, data := range *pf.FeeItem {
					principleLimit := data.PrincipleLimit.String()
					feeValue := data.FeeValue.String()
					itemNotes := *data.ItemNotes
					seqNo := strconv.FormatInt(int64(i), 10)
					rowMax := "0"
					if i == len(*pf.FeeItem)-1 {
						rowMax = "1"
					}
					recStatus := "1"
					recCreatedDate := recDate
					recCreatedBy := recBy
					queryItem += `('` + productFeeKey + `','` + seqNo + `','` + rowMax + `','` + principleLimit + `','` + feeValue + `','` + itemNotes + `','` + recStatus + `','` + recCreatedDate + `','` + recCreatedBy + `'),`
				}
				queryItem = queryItem[0 : len(queryItem)-1]

				_, err = tx.Exec(queryItem)
				if err != nil {
					tx.Rollback()
					log.Println(err.Error())
					return err
				}
			}

		}

		if *pf.RecAction == "UPDATE" {
			updtMsProductFee := make(map[string]string)

			var reflectValue = reflect.ValueOf(pf)
			if reflectValue.Kind() == reflect.Ptr {
				reflectValue = reflectValue.Elem()
			}
			var reflectType = reflectValue.Type()
			for i := 0; i < reflectValue.NumField(); i++ {
				columnName := reflectType.Field(i).Tag.Get("db")
				value := reflectValue.Field(i).Interface()
				if val, ok := value.(*uint64); ok {
					if val != nil {
						columnValue := strconv.FormatUint(*val, 10)
						updtMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*uint8); ok {
					if val != nil {
						columnValue := strconv.FormatUint(uint64(*val), 10)
						updtMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*string); ok {
					if val != nil {
						columnValue := *val
						updtMsProductFee[columnName] = columnValue
					}
				}
				if val, ok := value.(*decimal.Decimal); ok {
					if val != nil {
						columnValue := val.String()
						updtMsProductFee[columnName] = columnValue
					}
				}
			}
			updtMsProductFee["rec_modified_by"] = recBy
			updtMsProductFee["rec_modified_date"] = recDate
			var fields, values string
			var bindvars []interface{}
			for key, value := range updtMsProductFee {
				if key != "rec_pk" && key != "rec_action" {
					fields += key + ", "
					values += ` "` + value + `", `
					bindvars = append(bindvars, value)
				}
			}
			fields = fields[:(len(fields) - 2)]
			values = values[:(len(values) - 2)]
			qUpdMsProductFee := "UPDATE ms_product_fee SET "
			i := 0
			for key, value := range updtMsProductFee {
				if key != "fee_key" && key != "rec_pk" && key != "rec_action" {
					qUpdMsProductFee += key + " = '" + value + "'"
					if (len(updtMsProductFee) - 4) > i {
						qUpdMsProductFee += ", "
					}
					i++
				}
			}
			qUpdMsProductFee += " WHERE fee_key = " + updtMsProductFee["fee_key"]
			ret, err := tx.Exec(qUpdMsProductFee)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}
			countRows, err := ret.RowsAffected()
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err
			}
			if countRows > 0 {
				if len(*pf.FeeItem) > 0 {
					var queryItem string
					for _, data := range *pf.FeeItem {
						// itemKey :=
						feeKey := strconv.FormatUint(*data.ProductFeeKey, 10)
						principleLimit := data.PrincipleLimit.String()
						feeValue := data.FeeValue.String()
						itemNotes := *data.ItemNotes
						recstatus := strconv.FormatUint(*data.RecStatus, 10)
						seqNo := strconv.FormatInt(int64(i), 10)
						rowMax := "0"
						if i == len(*pf.FeeItem)-1 {
							rowMax = "1"
						}

						if *data.ProductFeeItemKey > 0 {
							// JIKA ADA FEE ITEM KEY MAKA UPDATE
							queryItem = `UPDATE ms_product_fee_item SET 
							principle_limit = ` + principleLimit + `, 
							fee_value = ` + feeValue + `, 
							item_notes = '` + itemNotes + `',
							item_seqno = '` + seqNo + `',
							row_max = '` + rowMax + `', 
							rec_status = '` + recstatus + `', 
							rec_modified_date = '` + recDate + `', 
							rec_modified_by = ` + recBy + `
							WHERE product_fee_item_key = ` + strconv.FormatUint(*data.ProductFeeItemKey, 10)
						} else {
							// JIKA TIDAK ADA FEE ITEM KEY MAKA CREATE
							recstatus = `1`
							queryItem = `INSERT INTO ms_product_fee_item(product_fee_key,item_seqno,row_max,principle_limit,fee_value,item_notes,rec_status,rec_created_date,rec_created_by) 
							VALUES('` + feeKey + `','` + seqNo + `','` + rowMax + `','` + principleLimit + `','` + feeValue + `','` + itemNotes + `','` + recstatus + `','` + params["rec_modified_date"] + `','` + params["rec_modified_by"] + `'),`
						}
						// log.Println("query update fee item:", queryItem)
						_, err = tx.Exec(queryItem)
						if err != nil {
							tx.Rollback()
							log.Println(err.Error())
							return err
						}
					}
				}

			}
		}
	}

	return nil
}
