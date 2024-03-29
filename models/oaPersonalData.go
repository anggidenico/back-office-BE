package models

import (
	"database/sql"
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"
)

type UserProfile struct {
	FullName       string            `json:"full_name"`
	SID            string            `json:"sid"`
	Email          string            `json:"email"`
	CIF            string            `json:"cif"`
	PhoneNumber    string            `json:"phone_number"`
	CifSuspendFlag bool              `json:"cif_suspend_flag"`
	RiskProfile    MsRiskProfileInfo `json:"risk_profile"`
	RecImage1      string            `json:"rec_image1"`
	BankAcc        BankAccount       `json:"bank_account"`
}

type OaPersonalData struct {
	PersonalDataKey        uint64  `db:"personal_data_key"          json:"personal_data_key"`
	OaRequestKey           uint64  `db:"oa_request_key"             json:"oa_request_key"`
	FullName               string  `db:"full_name"                  json:"email"`
	PlaceBirth             string  `db:"place_birth"                json:"place_birth"`
	DateBirth              string  `db:"date_birth"                 json:"date_birth"`
	Nationality            uint64  `db:"nationality"                json:"nationality"`
	IDcardType             uint64  `db:"idcard_type"                json:"idcard_type"`
	IDcardNo               string  `db:"idcard_no"                  json:"idcard_no"`
	IDcardExpiredDate      *string `db:"idcard_expired_date"        json:"idcard_expired_date"`
	IDcardNeverExpired     *uint8  `db:"idcard_never_expired"       json:"idcard_never_expired"`
	Gender                 *uint64 `db:"gender"                     json:"gender"`
	MaritalStatus          *uint64 `db:"marital_status"             json:"marital_status"`
	IDcardAddressKey       *uint64 `db:"idcard_address_key"         json:"idcard_address_key"`
	DomicileAddressKey     *uint64 `db:"domicile_address_key"       json:"domicile_address_key"`
	PhoneHome              *string `db:"phone_home"                 json:"phone_home"`
	PhoneMobile            string  `db:"phone_mobile"               json:"phone_mobile"`
	EmailAddress           string  `db:"email_address"              json:"email_address"`
	CorrespondAddress      *uint64 `db:"correspond_address"         json:"correspond_address"`
	Religion               *uint64 `db:"religion"                   json:"religion"`
	PicSelfie              *string `db:"pic_selfie"                 json:"pic_selfie"`
	PicKtp                 *string `db:"pic_ktp"                    json:"pic_ktp"`
	PicSelfieKtp           *string `db:"pic_selfie_ktp"             json:"pic_selfie_ktp"`
	GeolocName             *string `db:"geoloc_name"                json:"geoloc_name"`
	GeolocLongitude        *string `db:"geoloc_longitude"           json:"geoloc_longitude"`
	GeolocLatitude         *string `db:"geoloc_latitude"            json:"geoloc_latitude"`
	Education              *uint64 `db:"education"                  json:"education"`
	OccupJob               *uint64 `db:"occup_job"                  json:"occup_job"`
	OccupCompany           *string `db:"occup_company"              json:"occup_company"`
	OccupPosition          *uint64 `db:"occup_position"             json:"occup_position"`
	OccupAddressKey        *uint64 `db:"occup_address_key"          json:"occup_address_key"`
	OccupBusinessFields    *uint64 `db:"occup_business_fields"      json:"occup_business_fields"`
	OccupPhone             *string `db:"occup_phone"                json:"occup_phone"`
	OccupWebUrl            *string `db:"occup_web_url"              json:"occup_web_url"`
	Correspondence         *uint64 `db:"correspondence"             json:"correspondence"`
	AnnualIncome           *uint64 `db:"annual_income"              json:"annual_income"`
	SourceofFund           *uint64 `db:"sourceof_fund"              json:"sourceof_fund"`
	InvesmentObjectives    *uint64 `db:"invesment_objectives"       json:"invesment_objectives"`
	RelationType           *uint64 `db:"relation_type"              json:"relation_type"`
	RelationFullName       *string `db:"relation_full_name"         json:"relation_full_name"`
	RelationOccupation     *uint64 `db:"relation_occupation"        json:"relation_occupation"`
	RelationBusinessFields *uint64 `db:"relation_business_fields"   json:"relation_business_fields"`
	MotherMaidenName       string  `db:"mother_maiden_name"         json:"mother_maiden_name"`
	EmergencyFullName      *string `db:"emergency_full_name"        json:"emergency_full_name"`
	EmergencyRelation      *uint64 `db:"emergency_relation"         json:"emergency_relation"`
	EmergencyPhoneNo       *string `db:"emergency_phone_no"         json:"emergency_phone_no"`
	BeneficialFullName     *string `db:"beneficial_full_name"       json:"beneficial_full_name"`
	BeneficialRelation     *uint64 `db:"beneficial_relation"        json:"beneficial_relation"`
	BankAccountKey         *uint64 `db:"bank_account_key"           json:"bank_account_key"`
	PepStatus              *uint64 `db:"pep_status"                 json:"pep_status"`
	RecOrder               *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus              uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate         *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy           *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate        *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy          *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1              *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2              *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus      *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage       *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate        *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy          *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate         *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy           *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1        *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2        *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3        *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
	ReqBankAccKey          *uint64 `db:"req_bankacc_key" json:"req_bankacc_key"`
	PepName                *string `db:"pep_name" json:"pep_name"`
	PepPosition            *string `db:"pep_position" json:"pep_postiion"`
}

func GetAllOaPersonalData(c *[]OaPersonalData, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT oa_personal_data.* FROM oa_personal_data`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_personal_data."+field+" = '"+value+"'")
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

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaPersonalData(c *OaPersonalData, key string, field string) (int, error) {
	query := "SELECT oa_personal_data.* FROM oa_personal_data WHERE oa_personal_data." + field + " = " + key
	// log.Info("===== GET OA PERSONAL DATA =====>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateOaPersonalData(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_personal_data"
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

func GetOaPersonalDataByOaRequestKey(c *OaPersonalData, key string) (int, error) {
	query := `SELECT oa_personal_data.* FROM oa_personal_data WHERE oa_personal_data.oa_request_key = ` + key +
		` ORDER BY oa_personal_data.personal_data_key DESC LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetOaPersonalDataIn(c *[]OaPersonalData, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT oa_personal_data.* FROM oa_personal_data `
	query := query2 + " WHERE oa_personal_data." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ValidateUniquePersonalData(c *CountData, field string, value string, customerKey *string, requestKey *string) (int, error) {
	var query string
	query = `SELECT COUNT(a.oa_request_key) AS count_data
				FROM oa_request a
			INNER JOIN oa_personal_data b ON a.oa_request_key = b.oa_request_key
            where a.rec_status = '1' AND b.rec_status = '1' AND ` + field + ` = '` + value + `'`

	if requestKey != nil {
		query += ` AND (a.oa_request_key != '` + *requestKey + `')`
	}

	if customerKey != nil {
		query += ` AND (a.customer_key != '` + *customerKey + `' or a.customer_key IS NULL)`
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateOaPersonalData(params map[string]string) (int, error) {
	query := "UPDATE oa_personal_data SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "personal_data_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE personal_data_key = " + params["personal_data_key"]
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

type FullNameData struct {
	FullName string `db:"full_name"          json:"full_name"`
}

func GetFullName(c *FullNameData, userLoginKey string) (int, error) {
	query := `SELECT op.full_name as full_name  
	FROM oa_personal_data AS op
	INNER JOIN oa_request AS o ON op.oa_request_key = o.oa_request_key
	INNER JOIN sc_user_login AS u ON u.user_login_key = o.user_login_key 
	WHERE u.user_login_key = "` + userLoginKey + `" ORDER BY o.oa_request_key DESC LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateOaPersonalDataByKeyIn(params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := "UPDATE oa_personal_data SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	inQuery := strings.Join(valueIn, ",")
	query += " WHERE oa_personal_data." + fieldIn + " IN(" + inQuery + ")"

	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateRecStatusOaPersonalDocs(oareqkey string, doctype string, names string, datenows string) (int, error) {
	query := "UPDATE oa_personal_docs "
	query += "SET rec_status = 0, rec_modified_by = '" + names + "', rec_modified_date = '" + datenows + "'"
	query += "WHERE oa_request_key = " + oareqkey + " AND indi_document_type = " + doctype + " AND rec_status = 1 ORDER BY indi_docs_key DESC LIMIT 1"
	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	tx.Exec(query)
	tx.Commit()
	// log.Println("========== QUERY SET REC STATUS OA PERSONAL DOCS ==========>>>", query)
	return http.StatusOK, nil
}

func CreateOaDocsIndividu(params map[string]string) (int, error) {
	query := "INSERT INTO oa_personal_docs"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += " '" + value + "', "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Info("========== query insert oa personal docs ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	tx.Exec(query)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	// lastID, _ := ret.LastInsertId()
	// // log.Println(lastID)
	return http.StatusOK, nil
}

func UpdateDocPropertiesOaPersonalDocs(oareqkey string, docnames string, docremarks string, doctype string, names string, datenows string) (int, error) {
	query := "UPDATE oa_personal_docs "
	query += "SET indi_document_name = '" + docnames + "', indi_document_remarks ='" + docremarks + "', rec_modified_by = '" + names + "', rec_modified_date = '" + datenows + "'"
	query += "WHERE oa_request_key = " + oareqkey + " AND indi_document_type = " + doctype + " AND rec_status = 1"
	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	tx.Exec(query)
	tx.Commit()
	// log.Println("========== QUERY UPDATE DOC_NAME and DOC_REMARKS PERSONAL DOCS ==========>>>", query)
	return http.StatusOK, nil
}

type UdfOtherValueStruct struct {
	Values string `db:"udf_values"            json:"udf_values"`
}

func UdfOtherValueQuery(c *UdfOtherValueStruct, rowKey uint64, uikey string) (int, error) {
	row_key := strconv.FormatUint(rowKey, 10)
	query := `SELECT uv.udf_values
	FROM udf_value uv
	WHERE uv.udf_info_key = ` + uikey + ` AND uv.row_data_key = ` + row_key

	// log.Info("========= QUERY UDF OTHER VALUES ========== >>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
