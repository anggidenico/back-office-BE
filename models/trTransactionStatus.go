package models

import (
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type TrTransactionStatus struct {
	TransStatusKey    uint64  `db:"trans_status_key"        json:"trans_status_key"`
	StatusCode        *string `db:"status_code"             json:"status_code"`
	StatusDescription *string `db:"status_description"      json:"status_description"`
	StatusOrder       uint64  `db:"status_order"            json:"status_order"`
	StatusPhase       *string `db:"status_phase"            json:"status_phase"`
	RecOrder          *uint64 `db:"rec_order"            json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"           json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"           json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}

type TrTransactionStatusDropdown struct {
	TransStatusKey uint64  `json:"trans_status_key"`
	StatusCode     *string `json:"status_code"`
}

func GetMsTransactionStatusIn(c *[]TrTransactionStatus, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_transaction_status.* FROM 
				tr_transaction_status `
	query := query2 + " WHERE tr_transaction_status.rec_status = 1 AND tr_transaction_status." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionStatus(c *TrTransactionStatus, key string) (int, error) {
	query := `SELECT tr_transaction_status.* FROM tr_transaction_status WHERE tr_transaction_status.rec_status = 1 AND tr_transaction_status.trans_status_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAllMsTransactionStatus(c *[]TrTransactionStatus, params map[string]string) (int, error) {
	query := `SELECT
              tr_transaction_status.*
			  FROM tr_transaction_status
			  WHERE tr_transaction_status.rec_status = 1 AND tr_transaction_status.trans_status_key != 3`

	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction_status."+field+" = '"+value+"'")
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
