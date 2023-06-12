package models

import (
	"database/sql"
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"
)

type OaInstitutionAuthPerson struct {
	OaInstitutionAuthPerson uint64  `db:"insti_auth_person_key"               json:"insti_auth_person_key"`
	OaRequestKey            uint64  `db:"oa_request_key"                      json:"oa_request_key"`
	FullName                *string `db:"full_name"                           json:"full_name"`
	PersonDob               *string `db:"person_dob"                          json:"person_dob"`
	Nationality             *uint64 `db:"nationality"                         json:"nationality"`
	IdcardType              *uint64 `db:"idcard_type"                         json:"idcard_type"`
	IdcardNo                *string `db:"idcard_no"                           json:"idcard_no"`
	IdcardExpiredDate       *string `db:"idcard_expired_date"                 json:"idcard_expired_date"`
	IdcardNeverExpired      uint8   `db:"idcard_never_expired"                json:"idcard_never_expired"`
	Position                *string `db:"position"                            json:"position"`
	PhoneNo                 *string `db:"phone_no"                            json:"phone_no"`
	PhoneNoAlt              *string `db:"phone_no_alt"                        json:"phone_no_alt"`
	EmailAddress            *string `db:"email_address"                       json:"email_address"`
	RecOrder                *uint64 `db:"rec_order"                           json:"rec_order"`
	RecStatus               uint8   `db:"rec_status"                          json:"rec_status"`
	RecCreatedDate          *string `db:"rec_created_date"                    json:"rec_created_date"`
	RecCreatedBy            *string `db:"rec_created_by"                      json:"rec_created_by"`
	RecModifiedDate         *string `db:"rec_modified_date"                   json:"rec_modified_date"`
	RecModifiedBy           *string `db:"rec_modified_by"                     json:"rec_modified_by"`
	RecImage1               *string `db:"rec_image1"                          json:"rec_image1"`
	RecImage2               *string `db:"rec_image2"                          json:"rec_image2"`
	RecApprovalStatus       *uint8  `db:"rec_approval_status"                 json:"rec_approval_status"`
	RecApprovalStage        *uint64 `db:"rec_approval_stage"                  json:"rec_approval_stage"`
	RecApprovedDate         *string `db:"rec_approved_date"                   json:"rec_approved_date"`
	RecApprovedBy           *string `db:"rec_approved_by"                     json:"rec_approved_by"`
	RecDeletedDate          *string `db:"rec_deleted_date"                    json:"rec_deleted_date"`
	RecDeletedBy            *string `db:"rec_deleted_by"                      json:"rec_deleted_by"`
	RecAttributeID1         *string `db:"rec_attribute_id1"                   json:"rec_attribute_id1"`
	RecAttributeID2         *string `db:"rec_attribute_id2"                   json:"rec_attribute_id2"`
	RecAttributeID3         *string `db:"rec_attribute_id3"                   json:"rec_attribute_id3"`
}

type OaInstitutionAuthPersonDetail struct {
	OaInstitutionAuthPerson uint64  `db:"insti_auth_person_key"               json:"insti_auth_person_key"`
	FullName                *string `db:"full_name"                           json:"full_name"`
	PersonDob               *string `db:"person_dob"                          json:"person_dob"`
	Nationality             *uint64 `db:"nationality"                         json:"nationality"`
	NationalityName         *string `db:"nationality_name"                    json:"nationality_name"`
	IdcardType              *uint64 `db:"idcard_type"                         json:"idcard_type"`
	IdcardTypeName          *string `db:"idcard_type_name"                    json:"idcard_type_name"`
	IdcardNo                *string `db:"idcard_no"                           json:"idcard_no"`
	IdcardExpiredDate       *string `db:"idcard_expired_date"                 json:"idcard_expired_date"`
	IdcardNeverExpired      uint8   `db:"idcard_never_expired"                json:"idcard_never_expired"`
	Position                *string `db:"position"                            json:"position"`
	PhoneNo                 *string `db:"phone_no"                            json:"phone_no"`
	EmailAddress            *string `db:"email_address"                       json:"email_address"`
}

func GetOaInstitutionAuthPerson(c *OaInstitutionAuthPerson, key string, field string) (int, error) {
	query := `SELECT oa_institution_auth_person.* FROM oa_institution_auth_person 
	WHERE oa_institution_auth_person.rec_status = 1 AND oa_institution_auth_person.` + field + ` = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateOaInstitutionAuthPerson(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_institution_auth_person"
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

func GetOaInstitutionAuthPersonRequest(c *[]OaInstitutionAuthPersonDetail, oaReqKey string) (int, error) {
	query := `SELECT
				ap.insti_auth_person_key AS insti_auth_person_key,
				ap.full_name AS full_name,
				DATE_FORMAT(ap.person_dob, '%d %M %Y') AS person_dob,
				ap.nationality AS nationality,
				c.cou_name AS nationality_name,
				ap.idcard_type AS idcard_type,
				idcard.lkp_name AS idcard_type_name,
				ap.idcard_no AS idcard_no,
				DATE_FORMAT(ap.idcard_expired_date, '%d %M %Y') AS idcard_expired_date,
				ap.idcard_never_expired AS idcard_never_expired,
				ap.position AS position,
				ap.phone_no AS phone_no,
				ap.email_address AS email_address 
			FROM oa_institution_auth_person AS ap
			LEFT JOIN ms_country AS c ON c.country_key = ap.nationality 
			LEFT JOIN gen_lookup AS idcard ON idcard.lookup_key = ap.idcard_type 
			WHERE ap.rec_status = "1" AND ap.oa_request_key = "` + oaReqKey + `"`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateOaInstitutionAuthPerson(params map[string]string) (int, error) {
	query := "UPDATE oa_institution_auth_person SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "insti_auth_person_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE insti_auth_person_key = " + params["insti_auth_person_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
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
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func DeleteOaInstitutionAuthPerson(params map[string]string, authKey []string, requestKey string) (int, error) {
	inQuery := strings.Join(authKey, ",")
	query := "UPDATE oa_institution_auth_person SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	if len(authKey) > 0 {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "' AND insti_auth_person_key NOT IN(" + inQuery + ")"
	} else {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "'"
	}
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
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
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
