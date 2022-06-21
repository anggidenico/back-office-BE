package models

import (
	"mf-bo-api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type WfApprovalLog struct {
	ApprovalLogKey    uint64  `db:"approval_log_key"                json:"approval_log_key"`
	ApprovalHdrKey    uint64  `db:"approval_hdr_key"                json:"approval_hdr_key"`
	ApprovalNotes     *string `db:"approval_notes"                  json:"approval_notes"`
	RecOrder          *uint64 `db:"rec_order"                       json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                      json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"                json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"                  json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"               json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"                 json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                      json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                      json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"             json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"              json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"               json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"                 json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"                json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"                  json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"               json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"               json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"               json:"rec_attribute_id3"`
}

type LastUserApproveStatusApproved struct {
	UloginName   string  `db:"ulogin_name"                   json:"ulogin_name"`
	RecCreatedBy uint64  `db:"rec_created_by"                json:"rec_created_by"`
	CreatedDate  *string `db:"created_date"                  json:"created_date"`
}

func CreateWfApprovalLog(params map[string]string) (int, error) {
	query := "INSERT INTO wf_approval_log"
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
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetLastUserApproveStatusApproved(c *LastUserApproveStatusApproved, transactionKey string) (int, error) {
	query := `SELECT 
				u.ulogin_name,
				l.rec_created_by,
				DATE_FORMAT(l.rec_created_date, '%d %M %Y %H:%i') AS created_date  
			FROM wf_approval_log AS l 
			INNER JOIN wf_approval_list AS wfal ON wfal.approval_hdr_key = l.approval_hdr_key 
			INNER JOIN sc_user_login AS u ON u.user_login_key = l.rec_created_by
			WHERE wfal.approval_item = "tr_transaction" AND wfal.approval_references_key = "` + transactionKey + `" 
			AND l.rec_status = 1 AND l.rec_approval_status = "6" ORDER BY l.approval_log_key DESC LIMIT 1`

	// Main query
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
