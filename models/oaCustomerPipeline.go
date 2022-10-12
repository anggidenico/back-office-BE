package models

import (
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type OaCustomerPipeline struct {
	PipelineKey            uint64  `db:"pipeline_key"                        json:"pipeline_key"`
	InvestorType           uint64  `db:"investor_type"                       json:"investor_type"`
	Nationality            *uint64 `db:"nationality"                         json:"nationality"`
	FullName               *string `db:"full_name"                           json:"full_name"`
	ShortName              *string `db:"short_name"                          json:"short_name"`
	TinNumber              *string `db:"tin_number"                          json:"tin_number"`
	ProvinceKey            *string `db:"province_key"                        json:"province_key"`
	KabupatenKey           *string `db:"kabupaten_key"                       json:"kabupaten_key"`
	KecamatanKey           *string `db:"kecamatan_key"                       json:"kecamatan_key"`
	FullAddress            *string `db:"full_address"                        json:"full_address"`
	PostalCode             *string `db:"postal_code"                         json:"postal_code"`
	PicName                *string `db:"pic_name"                            json:"pic_name"`
	PicEmailAddress        *string `db:"pic_email_address"                   json:"pic_email_address"`
	PicPhoneNo             *string `db:"pic_phone_no"                        json:"pic_phone_no"`
	PicPosition            *string `db:"pic_position"                        json:"pic_position"`
	PipelineStatus         *string `db:"pipeline_status"                     json:"pipeline_status"`
	HeadSalesNotified      *string `db:"head_sales_notified"                 json:"head_sales_notified"`
	HeadSalesEmailAddress  *string `db:"head_sales_email_address"            json:"head_sales_email_address"`
	HeadSalesEmailSent     *string `db:"head_sales_email_sent"               json:"head_sales_email_sent"`
	AgentKey               *string `db:"agent_key"                           json:"agent_key"`
	AssignmentDate         *string `db:"assignment_date"                     json:"assignment_date"`
	AssignmentNotes        *string `db:"assignment_notes"                    json:"assignment_notes"`
	AgentNotified          *string `db:"agent_notified"                      json:"agent_notified"`
	AgentEmailAddress      *string `db:"agent_email_address"                 json:"agent_email_address"`
	AssignmentAccepted     *string `db:"assignment_accepted"                 json:"assignment_accepted"`
	AcceptanceNotes        *string `db:"acceptance_notes"                    json:"acceptance_notes"`
	PipelineVisitedDate    *uint64 `db:"pipeline_visited_date"               json:"pipeline_visited_date"`
	PipelineVisitedRemarks *uint64 `db:"pipeline_visited_remarks"            json:"pipeline_visited_remarks"`
	OaRequestLink          *string `db:"oa_request_link"                     json:"oa_request_link"`
	ClosingNotes           *string `db:"closing_notes"                       json:"closing_notes"`
	ClosingDate            *string `db:"closing_date"                        json:"closing_date"`
	RecOrder               *uint64 `db:"rec_order"                           json:"rec_order"`
	RecStatus              uint8   `db:"rec_status"                          json:"rec_status"`
	RecCreatedDate         *string `db:"rec_created_date"                    json:"rec_created_date"`
	RecCreatedBy           *string `db:"rec_created_by"                      json:"rec_created_by"`
	RecModifiedDate        *string `db:"rec_modified_date"                   json:"rec_modified_date"`
	RecModifiedBy          *string `db:"rec_modified_by"                     json:"rec_modified_by"`
	RecImage1              *string `db:"rec_image1"                          json:"rec_image1"`
	RecImage2              *string `db:"rec_image2"                          json:"rec_image2"`
	RecApprovalStatus      *uint8  `db:"rec_approval_status"                 json:"rec_approval_status"`
	RecApprovalStage       *uint64 `db:"rec_approval_stage"                  json:"rec_approval_stage"`
	RecApprovedDate        *string `db:"rec_approved_date"                   json:"rec_approved_date"`
	RecApprovedBy          *string `db:"rec_approved_by"                     json:"rec_approved_by"`
	RecDeletedDate         *string `db:"rec_deleted_date"                    json:"rec_deleted_date"`
	RecDeletedBy           *string `db:"rec_deleted_by"                      json:"rec_deleted_by"`
	RecAttributeID1        *string `db:"rec_attribute_id1"                   json:"rec_attribute_id1"`
	RecAttributeID2        *string `db:"rec_attribute_id2"                   json:"rec_attribute_id2"`
	RecAttributeID3        *string `db:"rec_attribute_id3"                   json:"rec_attribute_id3"`
}

func GetOaCustomerPipeline(c *OaCustomerPipeline, key string, field string) (int, error) {
	query := `SELECT oa_customer_pipeline.* FROM oa_customer_pipeline 
	WHERE oa_customer_pipeline.rec_status = 1 AND oa_customer_pipeline.` + field + ` = ` + key
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateOaCustomerPipeline(params map[string]string) (int, error) {
	query := "INSERT INTO oa_customer_pipeline"
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
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
