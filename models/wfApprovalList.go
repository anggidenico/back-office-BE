package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
)

type WfApprovalList struct {
	ApprovalHdrKey         uint64  `db:"approval_hdr_key"                json:"approval_hdr_key"`
	ApprovalAction         *uint64 `db:"approval_action"                 json:"approval_action"`
	ApprovalItem           *string `db:"approval_item"                   json:"approval_item"`
	ApprovalItemType       *string `db:"approval_item_type"              json:"approval_item_type"`
	ApprovalDescription    *string `db:"approval_description"            json:"approval_description"`
	ApprovalReferencesKey  *string `db:"approval_references_key"         json:"approval_references_key"`
	ApprovalControllerCode *string `db:"approval_controller_code"        json:"approval_controller_code"`
	ApprovalQuerystring    *string `db:"approval_querystring"            json:"approval_querystring"`
	DataContentBody        *string `db:"data_content_body"               json:"data_content_body"`
	RecOrder               *uint64 `db:"rec_order"                       json:"rec_order"`
	RecStatus              uint8   `db:"rec_status"                      json:"rec_status"`
	RecCreatedDate         *string `db:"rec_created_date"                json:"rec_created_date"`
	RecCreatedBy           *string `db:"rec_created_by"                  json:"rec_created_by"`
	RecModifiedDate        *string `db:"rec_modified_date"               json:"rec_modified_date"`
	RecModifiedBy          *string `db:"rec_modified_by"                 json:"rec_modified_by"`
	RecImage1              *string `db:"rec_image1"                      json:"rec_image1"`
	RecImage2              *string `db:"rec_image2"                      json:"rec_image2"`
	RecApprovalStatus      *uint8  `db:"rec_approval_status"             json:"rec_approval_status"`
	RecApprovalStage       *uint64 `db:"rec_approval_stage"              json:"rec_approval_stage"`
	RecApprovedDate        *string `db:"rec_approved_date"               json:"rec_approved_date"`
	RecApprovedBy          *string `db:"rec_approved_by"                 json:"rec_approved_by"`
	RecDeletedDate         *string `db:"rec_deleted_date"                json:"rec_deleted_date"`
	RecDeletedBy           *string `db:"rec_deleted_by"                  json:"rec_deleted_by"`
	RecAttributeID1        *string `db:"rec_attribute_id1"               json:"rec_attribute_id1"`
	RecAttributeID2        *string `db:"rec_attribute_id2"               json:"rec_attribute_id2"`
	RecAttributeID3        *string `db:"rec_attribute_id3"               json:"rec_attribute_id3"`
}

func CreateWfApprovalList(params map[string]string) (int, error, string) {
	query := "INSERT INTO wf_approval_list"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		if value == "NULL" {
			var s *string
			bindvars = append(bindvars, s)
		} else {
			bindvars = append(bindvars, value)
		}

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

func UpdateWfApprovalList(params map[string]string) (int, error) {
	query := "UPDATE wf_approval_list SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "approval_hdr_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE approval_hdr_key = " + params["approval_hdr_key"]
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

func UpdateWfApprovalListByApprovalItemAndKey(params map[string]string, appItem string, key string) (int, error) {
	query := "UPDATE wf_approval_list SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE approval_item = '" + appItem + "' AND approval_references_key = " + key
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

func GetWfApprovalListByApprovalItemAndKey(c *WfApprovalList, appItem string, key string) (int, error) {
	query := `SELECT * FROM wf_approval_list WHERE rec_status = 1 AND approval_references_key = "` + key + `" AND approval_item = "` + appItem + `" LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
