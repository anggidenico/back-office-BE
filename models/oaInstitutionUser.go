package models

import (
	"database/sql"
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OaInstitutionUser struct {
	InstiUserKey       uint64  `db:"insti_user_key"                      json:"insti_user_key"`
	OaRequestKey       uint64  `db:"oa_request_key"                      json:"oa_request_key"`
	FullName           *string `db:"full_name"                           json:"full_name"`
	EmailAddress       *string `db:"email_address"                       json:"email_address"`
	EmailVerifiedFlag  uint8   `db:"email_verified_flag"                 json:"email_verified_flag"`
	EmailVerifiedDate  *string `db:"email_verified_date"                 json:"email_verified_date"`
	PhoneNumber        *string `db:"phone_number"                        json:"phone_number"`
	PhoneVerifiedFlag  uint8   `db:"phone_verified_flag"                 json:"phone_verified_flag"`
	PhoneVerifiedDate  *string `db:"phone_verified_date"                 json:"phone_verified_date"`
	UserPriority       uint8   `db:"user_priority"                       json:"user_priority"`
	UloginCreatedDate  *string `db:"ulogin_created_date"                 json:"ulogin_created_date"`
	UserLoginKey       *uint64 `db:"user_login_key"                      json:"user_login_key"`
	UserLoginKey1      *uint64 `db:"user_login_key1"                     json:"user_login_key1"`
	RoleKey            *uint64 `db:"role_key"                            json:"role_key"`
	RoleCategoryKey    *uint64 `db:"role_category_key"                   json:"role_category_key"`
	UloginDeletedDate  *string `db:"ulogin_deleted_date"                 json:"ulogin_deleted_date"`
	InstituteUserNotes *string `db:"institute_user_notes"                json:"institute_user_notes"`
	RecOrder           *uint64 `db:"rec_order"                           json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"                          json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"                    json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"                      json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"                   json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"                     json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"                          json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"                          json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"                 json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"                  json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"                   json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"                     json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"                    json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"                      json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"                   json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"                   json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"                   json:"rec_attribute_id3"`
}
type OaInstitutionUserDetail struct {
	InstiUserKey       uint64  `db:"insti_user_key"                      json:"insti_user_key"`
	FullName           *string `db:"full_name"                           json:"full_name"`
	EmailAddress       *string `db:"email_address"                       json:"email_address"`
	EmailVerifiedFlag  uint8   `db:"email_verified_flag"                 json:"email_verified_flag"`
	EmailVerifiedDate  *string `db:"email_verified_date"                 json:"email_verified_date"`
	PhoneNumber        *string `db:"phone_number"                        json:"phone_number"`
	PhoneVerifiedFlag  uint8   `db:"phone_verified_flag"                 json:"phone_verified_flag"`
	PhoneVerifiedDate  *string `db:"phone_verified_date"                 json:"phone_verified_date"`
	UserPriority       uint8   `db:"user_priority"                       json:"user_priority"`
	UserLoginKey       *uint64 `db:"user_login_key"                      json:"user_login_key"`
	Username           *uint64 `db:"username"                            json:"username"`
	RoleKey            *uint64 `db:"role_key"                            json:"role_key"`
	RoleName           *string `db:"role_name"                           json:"role_name"`
	InstituteUserNotes *string `db:"institute_user_notes"                json:"institute_user_notes"`
}

func CreateOaInstitutionUser(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_institution_user"
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
	log.Info(query)

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

func GetOaInstitutionUserRequest(c *[]OaInstitutionUserDetail, oaReqKey string, roleKey string) (int, error) {
	query := `SELECT 
					u.insti_user_key AS insti_user_key,
					u.full_name AS full_name,
					u.email_address AS email_address,
					u.email_verified_flag AS email_verified_flag,
					DATE_FORMAT(u.email_verified_date, '%d %M %Y') AS email_verified_date,
					u.phone_number AS phone_number,
					u.phone_verified_flag AS phone_verified_flag,
					DATE_FORMAT(u.phone_verified_date, '%d %M %Y') AS phone_verified_date,
					u.user_priority AS user_priority,
					u.user_login_key AS user_login_key,
					ul.ulogin_name AS username,
					u.role_key AS role_key,
					r.role_name AS role_name,
					u.institute_user_notes AS institute_user_notes  
				FROM oa_institution_user AS u 
				INNER JOIN sc_role AS r ON r.role_key = u.role_key 
				LEFT JOIN sc_user_login AS ul ON ul.user_login_key = u.user_login_key AND u.rec_status = "1"
				WHERE u.rec_status = "1" 
				AND u.role_key = "` + roleKey + `" 
				AND u.oa_request_key = "` + oaReqKey + `"`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ValidateUniqueInstitutionUser(c *CountData, field string, value string, instUserKey *string) (int, error) {
	var query string
	query = `SELECT
				count(insti_user_key) AS count_data
              FROM oa_institution_user where rec_status = '1' AND ` + field + ` = '` + value + `'`

	if instUserKey != nil {
		query += ` AND insti_user_key != '` + *instUserKey + `'`
	}

	// Main query
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaInstitutionUser(c *OaInstitutionUser, key string, field string) (int, error) {
	query := `SELECT oa_institution_user.* FROM oa_institution_user 
	WHERE oa_institution_user.rec_status = 1 AND oa_institution_user.` + field + ` = ` + key + ` LIMIT 1`
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateOaInstitutionUser(params map[string]string) (int, error) {
	query := "UPDATE oa_institution_user SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "insti_user_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE insti_user_key = " + params["insti_user_key"]
	log.Info(query)

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

func DeleteOaInstitutionUser(params map[string]string, userKey []string, requestKey string) (int, error) {
	inQuery := strings.Join(userKey, ",")
	query := "UPDATE oa_institution_user SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	if len(userKey) > 0 {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "' AND insti_user_key NOT IN(" + inQuery + ")"
	} else {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "'"
	}
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

type OaInstitutionUserGenerateLogin struct {
	InstiUserKey    uint64  `db:"insti_user_key"                      json:"insti_user_key"`
	Oarequestkey    uint64  `db:"oa_request_key"                      json:"oa_request_key"`
	FullName        *string `db:"full_name"                           json:"full_name"`
	EmailAddress    *string `db:"email_address"                       json:"email_address"`
	PhoneNumber     *string `db:"phone_number"                        json:"phone_number"`
	RoleKey         *uint64 `db:"role_key"                            json:"role_key"`
	RoleCategoryKey *uint64 `db:"role_category_key"                   json:"role_category_key"`
	OaRequestType   *uint64 `db:"oa_request_type"                     json:"oa_request_type"`
	OaStatus        *uint64 `db:"oa_status"                           json:"oa_status"`
	CustomerKey     *uint64 `db:"customer_key"                        json:"customer_key"`
	UserLoginKey    *uint64 `db:"user_login_key"                      json:"user_login_key"`
	InstName        *string `db:"inst_name"                           json:"inst_name"`
	ShortName       *string `db:"short_name"                          json:"short_name"`
	RecStatusEmail  *string `db:"rec_status_email"                    json:"rec_status_email"`
	StringToken     *string `db:"string_token"                        json:"string_token"`
}

func GetOaInstitutionUserGenerateLogin(c *[]OaInstitutionUserGenerateLogin) (int, error) {
	query := `SELECT 
				iu.insti_user_key,
				iu.oa_request_key,
				iu.full_name,
				iu.email_address,
				iu.phone_number,
				iu.role_key,
				iu.role_category_key,
				r.oa_request_type,
				r.oa_status,
				r.customer_key,
				od.full_name AS inst_name,
				REPLACE(LOWER(od.short_name), " ", "") AS short_name,
				iu.rec_attribute_id1 AS rec_status_email 
			FROM oa_institution_user AS iu
			INNER JOIN oa_request AS r ON iu.oa_request_key = r.oa_request_key
			INNER JOIN oa_institution_data AS od ON od.oa_request_key = r.oa_request_key
			WHERE iu.user_login_key IS NULL AND iu.rec_status = 1 AND od.rec_status = 1 
			AND r.rec_status = 1 AND r.oa_status IN (260, 261, 262) 
			ORDER BY r.rec_modified_date, r.oa_request_key, iu.role_key ASC LIMIT 6`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaInstitutionUserGenerateLoginFailedSendEmail(c *[]OaInstitutionUserGenerateLogin) (int, error) {
	query := `SELECT 
				iu.insti_user_key,
				iu.oa_request_key,
				iu.full_name,
				iu.email_address,
				iu.phone_number,
				iu.role_key,
				iu.role_category_key,
				r.oa_request_type,
				r.oa_status,
				r.customer_key,
				od.full_name AS inst_name,
				REPLACE(LOWER(od.short_name), " ", "") AS short_name,
				iu.rec_attribute_id1 AS rec_status_email,
				ul.string_token AS string_token 
			FROM oa_institution_user AS iu
			INNER JOIN oa_request AS r ON iu.oa_request_key = r.oa_request_key 
			INNER JOIN oa_institution_data AS od ON od.oa_request_key = r.oa_request_key 
			INNER JOIN sc_user_login AS ul ON ul.user_login_key = iu.user_login_key
			WHERE iu.user_login_key IS NOT NULL AND iu.rec_status = 1 AND od.rec_status = 1 
			AND ul.rec_status = 1 AND iu.rec_attribute_id1 = 0 
			AND r.rec_status = 1 AND r.oa_status IN (260, 261, 262) 
			ORDER BY r.rec_modified_date, r.oa_request_key, iu.role_key ASC LIMIT 6`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCountUserActive(c *OaRequestCountData, oaKey string) (int, error) {
	query := `SELECT 
				COUNT(insti_user_key) count_data
			FROM oa_institution_user
			WHERE oa_request_key = "` + oaKey + `" 
			AND user_login_key IS NOT NULL`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
