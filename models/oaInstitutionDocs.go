package models

import (
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
)

type OaInstitutionDocs struct {
	InstiDocsKey              uint64  `db:"insti_docs_key"                      json:"insti_docs_key"`
	OaRequestKey              uint64  `db:"oa_request_key"                      json:"oa_request_key"`
	InstiDocumentType         *uint64 `db:"insti_document_type"                 json:"insti_document_type"`
	DocumentFileKey           *uint64 `db:"document_file_key"                   json:"document_file_key"`
	InstiDocumentName         *string `db:"insti_document_name"                 json:"insti_document_name"`
	InstiDocumentRemarks      *string `db:"insti_document_remarks"              json:"insti_document_remarks"`
	InstiDocValidationStatus  *uint64 `db:"insti_doc_validation_status"         json:"insti_doc_validation_status"`
	InstiDocValidationRemarks *string `db:"insti_doc_validation_remarks"        json:"insti_doc_validation_remarks"`
	InstiDocShipment          *uint64 `db:"insti_doc_shipment"                  json:"insti_doc_shipment"`
	DocShipmentDate           *string `db:"doc_shipment_date"                   json:"doc_shipment_date"`
	DocShipmentEmail          *string `db:"doc_shipment_email"                  json:"doc_shipment_email"`
	DocShipmentNotes          *string `db:"doc_shipment_notes"                  json:"doc_shipment_notes"`
	RecOrder                  *uint64 `db:"rec_order"                           json:"rec_order"`
	RecStatus                 uint8   `db:"rec_status"                          json:"rec_status"`
	RecCreatedDate            *string `db:"rec_created_date"                    json:"rec_created_date"`
	RecCreatedBy              *string `db:"rec_created_by"                      json:"rec_created_by"`
	RecModifiedDate           *string `db:"rec_modified_date"                   json:"rec_modified_date"`
	RecModifiedBy             *string `db:"rec_modified_by"                     json:"rec_modified_by"`
	RecImage1                 *string `db:"rec_image1"                          json:"rec_image1"`
	RecImage2                 *string `db:"rec_image2"                          json:"rec_image2"`
	RecApprovalStatus         *uint8  `db:"rec_approval_status"                 json:"rec_approval_status"`
	RecApprovalStage          *uint64 `db:"rec_approval_stage"                  json:"rec_approval_stage"`
	RecApprovedDate           *string `db:"rec_approved_date"                   json:"rec_approved_date"`
	RecApprovedBy             *string `db:"rec_approved_by"                     json:"rec_approved_by"`
	RecDeletedDate            *string `db:"rec_deleted_date"                    json:"rec_deleted_date"`
	RecDeletedBy              *string `db:"rec_deleted_by"                      json:"rec_deleted_by"`
	RecAttributeID1           *string `db:"rec_attribute_id1"                   json:"rec_attribute_id1"`
	RecAttributeID2           *string `db:"rec_attribute_id2"                   json:"rec_attribute_id2"`
	RecAttributeID3           *string `db:"rec_attribute_id3"                   json:"rec_attribute_id3"`
}

type OaInstitutionDocsDetail struct {
	InstiDocsKey          *uint64 `db:"insti_docs_key"                      json:"insti_docs_key"`
	InstiDocumentType     *uint64 `db:"insti_document_type"                 json:"insti_document_type"`
	InstiDocumentTypeName *string `db:"insti_document_type_name"            json:"insti_document_type_name"`
	DocumentFileName      *string `db:"document_file_name"                  json:"document_file_name"`
	InstiDocumentName     *string `db:"insti_document_name"                 json:"insti_document_name"`
	InstiDocumentRemarks  *string `db:"insti_document_remarks"              json:"insti_document_remarks"`
	Path                  *string `db:"path"                                json:"path"`
}

func GetOaInstitutionDocs(c *OaInstitutionDocs, key string, field string) (int, error) {
	query := `SELECT oa_institution_docs.* FROM oa_institution_docs 
	WHERE oa_institution_docs.rec_status = 1 AND oa_institution_docs.` + field + ` = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateOaInstitutionDocs(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_institution_docs"
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
		// log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	ret, err := tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetOaInstitutionDocsRequest(c *[]OaInstitutionDocsDetail, oaReqKey string) (int, error) {
	query := `SELECT 
				d.insti_docs_key AS insti_docs_key,
				ty.lookup_key AS insti_document_type,
				ty.lkp_name AS insti_document_type_name,
				f.file_name AS document_file_name,
				d.insti_document_name AS insti_document_name,
				d.insti_document_remarks AS insti_document_remarks 
			FROM gen_lookup AS ty
			LEFT JOIN oa_institution_docs AS d ON d.insti_document_type = ty.lookup_key AND d.rec_status = 1 AND d.oa_request_key = "` + oaReqKey + `" 
			LEFT JOIN ms_file AS f ON f.file_key = d.document_file_key AND f.ref_fk_domain = "oa_institution_docs" AND f.rec_status = 1
			WHERE ty.rec_status = 1 
			AND ty.lkp_group_key = "95"`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateOaInstitutionDocs(params map[string]string) (int, error) {
	query := "UPDATE oa_institution_docs SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "insti_docs_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE insti_docs_key = " + params["insti_docs_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}
