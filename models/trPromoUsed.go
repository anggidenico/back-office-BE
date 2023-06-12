package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
)

type TrPromoUsed struct {
	PromoUsedKey      uint64  `db:"promo_used_key"          json:"promo_used_key"`
	UsedDate          string  `db:"used_date"               json:"used_date"`
	PromoKey          uint64  `db:"promo_key"               json:"promo_key"`
	UserLoginKey      uint64  `db:"user_login_key"          json:"user_login_key"`
	CustomerKey       uint64  `db:"customer_key"            json:"customer_key"`
	TransactionKey    uint64  `db:"transaction_key"         json:"transaction_key"`
	UsedStatus        uint64  `db:"used_status"             json:"used_status"`
	UsedNotes         *string `db:"used_notes"              json:"used_notes"`
	RecOrder          *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

func CreateTrPromoUsed(params map[string]string) (int, error) {
	query := "INSERT INTO tr_promo_used"
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

func AdminGetCountPromoUsed(c *CountData, promoKey *string, customerKey *string, transactionKey *string) (int, error) {
	query := `SELECT COUNT(promo_used_key) as count_data FROM tr_promo_used WHERE rec_status = 1`

	if promoKey != nil {
		query += " AND promo_key = '" + *promoKey + "'"
	}

	if customerKey != nil {
		query += " AND customer_key = '" + *customerKey + "'"
	}

	if transactionKey != nil {
		query += " AND transaction_key != '" + *transactionKey + "'"
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrPromoUsedByField(c *TrPromoUsed, field string, value string) (int, error) {
	query := `SELECT tr_promo_used.* FROM tr_promo_used WHERE tr_promo_used.rec_status = 1 
	AND tr_promo_used.` + field + ` = '` + value + `' LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateTrPromoUsed(params map[string]string) (int, error) {
	query := "UPDATE tr_promo_used SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "promo_used_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE tr_promo_used = " + params["promo_used_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
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
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
