package models

import (
	"mf-bo-api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type WfPath struct {
	PathKey           uint64  `db:"path_key"               json:"path_key"`
	WfProfileKey      uint64  `db:"wf_profile_key"         json:"wf_profile_key"`
	PathSeqno         uint64  `db:"path_seqno"             json:"path_seqno"`
	StageCurrent      uint64  `db:"stage_current"          json:"stage_current"`
	StageNext         *uint64 `db:"stage_next"             json:"stage_next"`
	StagePrev         *uint64 `db:"stage_prev"             json:"stage_prev"`
	NotifyTonext      uint8   `db:"notify_tonext"          json:"notify_tonext"`
	NotiftToprev      uint8   `db:"notift_toprev"          json:"notift_toprev"`
	NotifyToall       uint8   `db:"notify_toall"           json:"notify_toall"`
	RecOrder          *uint64 `db:"rec_order"              json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"             json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"             json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func GetAllWfPath(c *[]WfPath, params map[string]string) (int, error) {
	query := `SELECT
              wf_path.* FROM 
			  wf_path `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "wf_path."+field+" = '"+value+"'")
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
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetWfPathByProfilAndStageCurrent(c *WfPath, profileKey string, stageCurrent string) (int, error) {
	query := `SELECT wf_path.* 
			FROM wf_path 
			WHERE wf_path.rec_status = "1" AND wf_path.wf_profile_key = "` + profileKey + `" 
			AND wf_path.stage_current = "` + stageCurrent + `"`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
