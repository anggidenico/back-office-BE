package models

import (
	"mf-bo-api/db"
	"net/http"
)

type MsParticipant struct {
	ParticipantKey      uint64  `db:"participant_key"       json:"participant_key"`
	ParticipantCode     string  `db:"participant_code"      json:"participant_code"`
	ParticipantName     string  `db:"participant_name"      json:"participant_name"`
	ParticipantCategory string  `db:"participant_category"  json:"participant_category"`
	RecOrder            *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus           uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate      *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy        *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate     *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy       *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1           *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2           *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus   *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage    *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate     *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy       *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate      *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy        *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1     *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2     *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3     *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type ListDropdownMsParticipant struct {
	ParticipantKey  uint64  `db:"participant_key"        json:"participant_key"`
	ParticipantCode *string `db:"participant_code"       json:"participant_code"`
	ParticipantName *string `db:"participant_name"       json:"participant_name"`
}

func AdminGetListDropdownMsParticipant(c *[]ListDropdownMsParticipant) (int, error) {
	query := `SELECT  
				p.participant_key,
				p.participant_code,
				p.participant_name 
			FROM ms_participant AS p
			WHERE p.rec_status = 1 `
	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
