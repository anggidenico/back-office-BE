package models

import (
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type OaInstitutionData struct {
	InstitutionDataKey        uint64           `db:"institution_data_key"                json:"institution_data_key"`
	OaRequestKey              uint64           `db:"oa_request_key"                      json:"oa_request_key"`
	Nationality               *uint64          `db:"nationality"                         json:"nationality"`
	FullName                  *string          `db:"full_name"                           json:"full_name"`
	ShortName                 *string          `db:"short_name"                          json:"short_name"`
	TinNumber                 *string          `db:"tin_number"                          json:"tin_number"`
	EstablishedCity           *string          `db:"established_city"                    json:"established_city"`
	EstablishedDate           *string          `db:"established_date"                    json:"established_date"`
	DeedNo                    *string          `db:"deed_no"                             json:"deed_no"`
	DeedDate                  *string          `db:"deed_date"                           json:"deed_date"`
	LsEstablishValidationNo   *string          `db:"ls_establish_validation_no"          json:"ls_establish_validation_no"`
	LsEstablishValidationDate *string          `db:"ls_establish_validation_date"        json:"ls_establish_validation_date"`
	LastChangeAaNo            *string          `db:"last_change_aa_no"                   json:"last_change_aa_no"`
	LastChangeAaDate          *string          `db:"last_change_aa_date"                 json:"last_change_aa_date"`
	LsLastChangeAaNo          *string          `db:"ls_last_change_aa_no"                json:"ls_last_change_aa_no"`
	LsLastChangeAaDate        *string          `db:"ls_last_change_aa_date"              json:"ls_last_change_aa_date"`
	ManagementDeedNo          *string          `db:"management_deed_no"                  json:"management_deed_no"`
	ManagementDeedDate        *string          `db:"management_deed_date"                json:"management_deed_date"`
	LsMgtChangeDeedNo         *string          `db:"ls_mgt_change_deed_no"               json:"ls_mgt_change_deed_no"`
	LsMgtChangeDeedDate       *string          `db:"ls_mgt_change_deed_date"             json:"ls_mgt_change_deed_date"`
	SkdLicenseNo              *string          `db:"skd_license_no"                      json:"skd_license_no"`
	SkdLicenseDate            *string          `db:"skd_license_date"                    json:"skd_license_date"`
	BizLicenseNo              *string          `db:"biz_license_no"                      json:"biz_license_no"`
	BizLicenseDate            *string          `db:"biz_license_date"                    json:"biz_license_date"`
	NibNo                     *string          `db:"nib_no"                              json:"nib_no"`
	NibDate                   *string          `db:"nib_date"                            json:"nib_date"`
	DomicileKey               *uint64          `db:"domicile_key"                        json:"domicile_key"`
	CorrespondenceKey         *uint64          `db:"correspondence_key"                  json:"correspondence_key"`
	PhoneNo                   *string          `db:"phone_no"                            json:"phone_no"`
	MobileNo                  *string          `db:"mobile_no"                           json:"mobile_no"`
	FaxNo                     *string          `db:"fax_no"                              json:"fax_no"`
	EmailAddress              *string          `db:"email_address"                       json:"email_address"`
	IntitutionType            *uint64          `db:"intitution_type"                     json:"intitution_type"`
	IntitutionClassification  *uint64          `db:"intitution_classification"           json:"intitution_classification"`
	IntitutionCharacteristic  *uint64          `db:"intitution_characteristic"           json:"intitution_characteristic"`
	IntitutionBusinessType    *uint64          `db:"intitution_business_type"            json:"intitution_business_type"`
	InstiAnnuallyIncome       *uint64          `db:"insti_annually_income"               json:"insti_annually_income"`
	InstiSourceOfIncome       *uint64          `db:"insti_source_of_income"              json:"insti_source_of_income"`
	InstiInvestmentPurpose    *uint64          `db:"insti_investment_purpose"            json:"insti_investment_purpose"`
	BoName                    *string          `db:"bo_name"                             json:"bo_name"`
	BoIdnumber                *string          `db:"bo_idnumber"                         json:"bo_idnumber"`
	BoBusiness                *string          `db:"bo_business"                         json:"bo_business"`
	BoIdaddress               *string          `db:"bo_idaddress"                        json:"bo_idaddress"`
	BoBusinessAddress         *string          `db:"bo_business_address"                 json:"bo_business_address"`
	BoAnnuallyIncome          *decimal.Decimal `db:"bo_annually_income"                  json:"bo_annually_income"`
	BoRelation                *uint64          `db:"bo_relation"                         json:"bo_relation"`
	AssetY1                   *decimal.Decimal `db:"asset_y1"                            json:"asset_y1"`
	AssetY2                   *decimal.Decimal `db:"asset_y2"                            json:"asset_y2"`
	AssetY3                   *decimal.Decimal `db:"asset_y3"                            json:"asset_y3"`
	OpsProfitY1               *decimal.Decimal `db:"ops_profit_y1"                       json:"ops_profit_y1"`
	OpsProfitY2               *decimal.Decimal `db:"ops_profit_y2"                       json:"ops_profit_y2"`
	OpsProfitY3               *decimal.Decimal `db:"ops_profit_y3"                       json:"ops_profit_y3"`
	BankAccountKey            *uint64          `db:"bank_account_key"                    json:"bank_account_key"`
	InstiRemarks              *string          `db:"insti_remarks"                       json:"insti_remarks"`
	InstitutionGroup          *string          `db:"institution_group"                   json:"institution_group"`
	InstiDocShipment          *uint64          `db:"insti_doc_shipment"                  json:"insti_doc_shipment"`
	DocShipmentDate           *string          `db:"doc_shipment_date"                   json:"doc_shipment_date"`
	DocShipmentEmail          *string          `db:"doc_shipment_email"                  json:"doc_shipment_email"`
	DocShipmentNotes          *string          `db:"doc_shipment_notes"                  json:"doc_shipment_notes"`
	RecOrder                  *uint64          `db:"rec_order"                           json:"rec_order"`
	RecStatus                 uint8            `db:"rec_status"                          json:"rec_status"`
	RecCreatedDate            *string          `db:"rec_created_date"                    json:"rec_created_date"`
	RecCreatedBy              *string          `db:"rec_created_by"                      json:"rec_created_by"`
	RecModifiedDate           *string          `db:"rec_modified_date"                   json:"rec_modified_date"`
	RecModifiedBy             *string          `db:"rec_modified_by"                     json:"rec_modified_by"`
	RecImage1                 *string          `db:"rec_image1"                          json:"rec_image1"`
	RecImage2                 *string          `db:"rec_image2"                          json:"rec_image2"`
	RecApprovalStatus         *uint8           `db:"rec_approval_status"                 json:"rec_approval_status"`
	RecApprovalStage          *uint64          `db:"rec_approval_stage"                  json:"rec_approval_stage"`
	RecApprovedDate           *string          `db:"rec_approved_date"                   json:"rec_approved_date"`
	RecApprovedBy             *string          `db:"rec_approved_by"                     json:"rec_approved_by"`
	RecDeletedDate            *string          `db:"rec_deleted_date"                    json:"rec_deleted_date"`
	RecDeletedBy              *string          `db:"rec_deleted_by"                      json:"rec_deleted_by"`
	RecAttributeID1           *string          `db:"rec_attribute_id1"                   json:"rec_attribute_id1"`
	RecAttributeID2           *string          `db:"rec_attribute_id2"                   json:"rec_attribute_id2"`
	RecAttributeID3           *string          `db:"rec_attribute_id3"                   json:"rec_attribute_id3"`
}

type AdminListOaInstitutionData struct {
	OaRequestKey     uint64  `db:"oa_request_key"         json:"oa_request_key"`
	BranchName       *string `db:"branch_name"            json:"branch_name"`
	AgentName        *string `db:"agent_name"             json:"agent_name"`
	OaStatus         *string `db:"oa_status"              json:"oa_status"`
	StatusData       *string `db:"status_data"            json:"status_data"`
	Npwp             *string `db:"npwp"                   json:"npwp"`
	FullName         *string `db:"full_name"              json:"full_name"`
	TanggalPendirian *string `db:"tanggal_pendirian"      json:"tanggal_pendirian"`
	TempatPendirian  *string `db:"tempat_pendirian"       json:"tempat_pendirian"`
	NomorAkta        *string `db:"nomor_akta"             json:"nomor_akta"`
	TanggalAkta      *string `db:"tanggal_akta"           json:"tanggal_akta"`
	NoIzinUsaha      *string `db:"no_izin_usaha"          json:"no_izin_usaha"`
}

type OaInstitutionDetail struct {
	OaRequestKey              uint64                             `json:"oa_request_key"`
	InstitutionDataKey        uint64                             `json:"institution_data_key"`
	SalesCode                 *string                            `json:"sales_code"`
	Check1Date                *string                            `json:"check1_date"`
	Check1Flag                *uint8                             `json:"check1_flag"`
	Check1References          *string                            `json:"check1_references"`
	Check1Notes               *string                            `json:"check1_notes"`
	Check2Date                *string                            `json:"check2_date"`
	Check2Flag                *uint8                             `json:"check2_flag"`
	Check2References          *string                            `json:"check2_references"`
	Check2Notes               *string                            `json:"check2_notes"`
	Nationality               *MsCountryList                     `json:"nationality"`
	FullName                  *string                            `json:"full_name"`
	ShortName                 *string                            `json:"short_name"`
	TinNumber                 *string                            `json:"tin_number"`
	EstablishedCity           *string                            `json:"established_city"`
	EstablishedDate           *string                            `json:"established_date"`
	DeedNo                    *string                            `json:"deed_no"`
	DeedDate                  *string                            `json:"deed_date"`
	LsEstablishValidationNo   *string                            `json:"ls_establish_validation_no"`
	LsEstablishValidationDate *string                            `json:"ls_establish_validation_date"`
	LastChangeAaNo            *string                            `json:"last_change_aa_no"`
	LastChangeDaDate          *string                            `json:"last_change_aa_date"`
	LsLastChangeAaNo          *string                            `json:"ls_last_change_aa_no"`
	LsLastChangeAaDate        *string                            `json:"ls_last_change_aa_date"`
	ManagementDeedNo          *string                            `json:"management_deed_no"`
	ManagementDeedDate        *string                            `json:"management_deed_date"`
	LsMgtChangeDeedNo         *string                            `json:"ls_mgt_change_deed_no"`
	LsMgtChangeDeedDate       *string                            `json:"ls_mgt_change_deed_date"`
	SkdLicenseNo              *string                            `json:"skd_license_no"`
	SkdLicenseDate            *string                            `json:"skd_license_date"`
	BizLicenseNo              *string                            `json:"biz_license_no"`
	BizLicenseDate            *string                            `json:"biz_license_date"`
	NibNo                     *string                            `json:"nib_no"`
	NibDate                   *string                            `json:"nib_date"`
	Domicile                  *AddressDetail                     `json:"domicile"`
	Correspondence            *AddressDetail                     `json:"correspondence"`
	PhoneNo                   *string                            `json:"phone_no"`
	MobileNo                  *string                            `json:"mobile_no"`
	FaxNo                     *string                            `json:"fax_no"`
	EmailAddress              *string                            `json:"email_address"`
	IntitutionType            *LookupTrans                       `json:"intitution_type"`
	IntitutionClassification  *LookupTrans                       `json:"intitution_classification"`
	IntitutionCharacteristic  *LookupTrans                       `json:"intitution_characteristic"`
	IntitutionBusinessType    *LookupTrans                       `json:"intitution_business_type"`
	InstiAnnuallyIncome       *LookupTrans                       `json:"insti_annually_income"`
	InstiSourceOfIncome       *LookupTrans                       `json:"insti_source_of_income"`
	InstiInvestmentPurpose    *LookupTrans                       `json:"insti_investment_purpose"`
	BoName                    *string                            `json:"bo_name"`
	BoIdnumber                *string                            `json:"bo_idnumber"`
	BoBusiness                *string                            `json:"bo_business"`
	BoIdaddress               *string                            `json:"bo_idaddress"`
	BoBusinessAddress         *string                            `json:"bo_business_address"`
	BoAnnuallyIncome          *decimal.Decimal                   `json:"bo_annually_income"`
	BoRelation                *LookupTrans                       `json:"bo_relation"`
	AssetY1                   *decimal.Decimal                   `json:"asset_y1"`
	AssetY2                   *decimal.Decimal                   `json:"asset_y2"`
	AssetY3                   *decimal.Decimal                   `json:"asset_y3"`
	OpsProfitY1               *decimal.Decimal                   `json:"ops_profit_y1"`
	OpsProfitY2               *decimal.Decimal                   `json:"ops_profit_y2"`
	OpsProfitY3               *decimal.Decimal                   `json:"ops_profit_y3"`
	InstiRemarks              *string                            `json:"insti_remarks"`
	InstitutionGroup          *string                            `json:"institution_group"`
	InstiDocShipment          *LookupTrans                       `json:"insti_doc_shipment"`
	DocShipmentDate           *string                            `json:"doc_shipment_date"`
	DocShipmentEmail          *string                            `json:"doc_shipment_email"`
	DocShipmentNotes          *string                            `json:"doc_shipment_notes"`
	RecCreatedDate            *string                            `json:"rec_created_date"`
	RecCreatedBy              *string                            `json:"rec_created_by"`
	RecModifiedDate           *string                            `json:"rec_modified_date"`
	RecModifiedBy             *string                            `json:"rec_modified_by"`
	RecApprovalStatus         *uint8                             `json:"rec_approval_status"`
	RecApprovalStage          *uint64                            `json:"rec_approval_stage"`
	RecApprovedDate           *string                            `json:"rec_approved_date"`
	RecApprovedBy             *string                            `json:"rec_approved_by"`
	OaStatus                  *LookupTrans                       `json:"oa_status"`
	Branch                    *MsBranchDropdown                  `json:"branch"`
	Agent                     *MsAgentDropdown                   `json:"agent"`
	BankRequest               *[]OaRequestByField                `json:"bank_request"`
	RiskProfile               *AdminOaRiskProfile                `json:"risk_profile"`
	RiskProfileQuiz           *[]RiskProfileQuiz                 `json:"risk_profile_quiz"`
	InstitutionDocs           *[]OaInstitutionDocsDetail         `json:"institution_docs"`
	InstitutionUserMaker      *[]OaInstitutionUserDetail         `json:"institution_user_maker"`
	InstitutionUserChecker    *[]OaInstitutionUserDetail         `json:"institution_user_checker"`
	InstitutionUserReleaser   *[]OaInstitutionUserDetail         `json:"institution_user_releaser"`
	InstitutionSharesHolder   *[]OaInstitutionSharesHolderDetail `json:"institution_shares_holder"`
	InstitutionAuthPerson     *[]OaInstitutionAuthPersonDetail   `json:"institution_auth_person"`
}

func GetOaInstitutionData(c *OaInstitutionData, key string, field string) (int, error) {
	query := `SELECT oa_institution_data.* FROM oa_institution_data 
	WHERE oa_institution_data.rec_status = 1 AND oa_institution_data.` + field + ` = ` + key
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateOaInstitutionData(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_institution_data"
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
		return http.StatusBadGateway, err, "0"
	}
	ret, err := tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func UpdateOaInstitutionData(params map[string]string) (int, error) {
	query := "UPDATE oa_institution_data SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "institution_data_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE institution_data_key = " + params["institution_data_key"]
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func AdminGetListOaInstitutionData(
	c *[]AdminListOaInstitutionData,
	oaValue []string,
	fieldNot string,
	valueNot string,
	params map[string]string,
	limit uint64,
	offset uint64,
	nolimit bool,
	searchLike string,
) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	query := `SELECT 
					o.oa_request_key,
					b.branch_name,
					CONCAT(a.agent_code, " - ", a.agent_name) AS agent_name,
					o.oa_status,
					oas.lkp_name AS status_data,
					oad.tin_number AS npwp,
					oad.full_name,
					DATE_FORMAT(oad.established_date, '%d %M %Y') AS tanggal_pendirian,
					oad.established_city AS tempat_pendirian,
					oad.deed_no AS nomor_akta,
					DATE_FORMAT(oad.deed_date, '%d %M %Y') AS tanggal_akta,
					oad.biz_license_no AS no_izin_usaha 
				FROM oa_request AS o
				INNER JOIN oa_institution_data AS oad ON o.oa_request_key = oad.oa_request_key
				INNER JOIN ms_branch AS b ON o.branch_key = b.branch_key
				INNER JOIN ms_agent AS a ON a.agent_key = o.agent_key 
				INNER JOIN gen_lookup AS oas ON oas.lookup_key = o.oa_status 
				WHERE o.rec_status = 1 AND oad.rec_status = 1`
	if len(oaValue) > 0 {
		inQuery := strings.Join(oaValue, ",")
		query += " AND o.oa_status IN(" + inQuery + ")"
	}

	if fieldNot != "" && valueNot != "" {
		query += " AND " + fieldNot + " != '" + valueNot + "'"
	}

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (o.oa_request_key like '%" + searchLike + "%' OR"
		condition += " b.branch_name like '%" + searchLike + "%' OR"
		condition += " CONCAT(a.agent_code, ' - ', a.agent_name) like '%" + searchLike + "%' OR"
		condition += " oas.lkp_name like '%" + searchLike + "%' OR"
		condition += " oad.tin_number like '%" + searchLike + "%' OR"
		condition += " oad.full_name like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(oad.established_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " oad.established_city like '%" + searchLike + "%' OR"
		condition += " oad.deed_no like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(oad.deed_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " oad.biz_license_no like '%" + searchLike + "%')"
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

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountGetListOaInstitutionData(
	c *OaRequestCountData,
	oaValue []string,
	fieldNot string,
	valueNot string,
	params map[string]string,
	searchLike string,
) (int, error) {
	var whereClause []string
	var condition string
	query := `SELECT 
					count(o.oa_request_key) as count_data 
				FROM oa_request AS o
				INNER JOIN oa_institution_data AS oad ON o.oa_request_key = oad.oa_request_key
				INNER JOIN ms_branch AS b ON o.branch_key = b.branch_key
				INNER JOIN ms_agent AS a ON a.agent_key = o.agent_key 
				INNER JOIN gen_lookup AS oas ON oas.lookup_key = o.oa_status 
				WHERE o.rec_status = 1 AND oad.rec_status = 1`
	if len(oaValue) > 0 {
		inQuery := strings.Join(oaValue, ",")
		query += " AND o.oa_status IN(" + inQuery + ")"
	}

	if fieldNot != "" && valueNot != "" {
		query += " AND " + fieldNot + " != '" + valueNot + "'"
	}

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (o.oa_request_key like '%" + searchLike + "%' OR"
		condition += " b.branch_name like '%" + searchLike + "%' OR"
		condition += " CONCAT(a.agent_code, ' - ', a.agent_name) like '%" + searchLike + "%' OR"
		condition += " oas.lkp_name like '%" + searchLike + "%' OR"
		condition += " oad.tin_number like '%" + searchLike + "%' OR"
		condition += " oad.full_name like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(oad.established_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " oad.established_city like '%" + searchLike + "%' OR"
		condition += " oad.deed_no like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(oad.deed_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " oad.biz_license_no like '%" + searchLike + "%')"
	}

	query += condition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaInstitutionDataIn(c *[]OaInstitutionData, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				oa_institution_data.* FROM 
				oa_institution_data `
	query := query2 + " WHERE oa_institution_data." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
