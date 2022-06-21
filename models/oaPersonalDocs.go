package models

import (
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type OaPersonalDocuments struct {
	OaRequestKey uint64                 `json:"oa_request_key"`
	PersonalDocs []OaPersonalDocsDetail `json:"personal_docs"`
}

type OaPersonalDocsDetail struct {
	IndiDocsKey          *uint64 `db:"indi_docs_key"                      json:"personal_docs_key"`
	OaRequestKey         *uint64 `db:"oa_request_key"                     json:"oa_request_key"`
	IndiDocumentType     *uint64 `db:"indi_document_type"                 json:"personal_document_type"`
	IndiDocumentTypeName *string `db:"indi_document_type_name"            json:"personal_document_type_name"`
	DocumentFileName     *string `db:"document_file_name"                 json:"document_file_name"`
	IndiDocumentName     *string `db:"indi_document_name"                 json:"personal_document_name"`
	IndiDocumentRemarks  *string `db:"indi_document_remarks"              json:"personal_document_remarks"`
	IndiDocumentFullname *string `db:"indi_document_fullname"             json:"personal_document_fullname"`
	Path                 *string `db:"path"                               json:"path"`
}

func GetOaPersonalDocsRequest(c *[]OaPersonalDocsDetail, oaReqKey string) (int, error) {
	query := `SELECT 
		d.oa_request_key AS oa_request_key,
		d.indi_docs_key AS indi_docs_key,
		ty.lookup_key AS indi_document_type,
		ty.lkp_name AS indi_document_type_name,
		d.indi_document_name AS indi_document_name,
		d.indi_document_remarks AS indi_document_remarks,
		d.rec_image1 AS indi_document_fullname
	FROM gen_lookup AS ty
	LEFT JOIN oa_personal_docs AS d ON d.indi_document_type = ty.lookup_key 
	AND d.rec_status = 1
	AND d.oa_request_key = "` + oaReqKey + `"
	WHERE ty.rec_status = 1
	AND ty.lkp_group_key = 105`

	// Main query
	log.Println("========== query cek personal data ==========", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
