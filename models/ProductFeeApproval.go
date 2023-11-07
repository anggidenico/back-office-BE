package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"
)

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

func ProductFeeUpdateRequest(paramsFee map[string]string, feeItems []FeeItemData) (int, error) {
	return http.StatusOK, nil
}
