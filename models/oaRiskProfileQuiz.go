package models

import (
	"database/sql"
	_ "database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type ParamsRiskProfileQuiz struct {
	OaRequestKey    string `db:"oa_request_key"`
	QuizQuestionKey string `db:"quiz_question_key"`
	QuizOptionKey   string `db:"quiz_option_key"`
	QuizOptionScore string `db:"quiz_option_score"`
	RecStatus       string `db:"rec_status"`
}

type OaRiskProfileQuiz struct {
	RiskProfileQuizKey uint64  `db:"risk_profile_quiz_key"   json:"risk_profile_quiz_key"`
	OaRequestKey       *uint64 `db:"oa_request_key"          json:"oa_request_key"`
	QuizQuestionKey    *uint64 `db:"quiz_question_key"       json:"quiz_question_key"`
	QuizOptionKey      *uint64 `db:"quiz_option_key"         json:"quiz_option_key"`
	QuizOptionScore    *uint64 `db:"quiz_option_score"       json:"quiz_option_score"`
	RecOrder           *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type AdminOaRiskProfileQuiz struct {
	RiskProfileQuizKey  uint64          `db:"risk_profile_quiz_key"      json:"risk_profile_quiz_key"`
	QuizQuestionKey     uint64          `db:"quiz_question_key"          json:"quiz_question_key"`
	QuizOptionKeyUser   uint64          `db:"quiz_option_key_user"       json:"quiz_option_key_user"`
	QuizOptionScoreUser decimal.Decimal `db:"quiz_option_score_user"     json:"quiz_option_score_user"`
	QuizTitle           string          `db:"quiz_title"                 json:"quiz_title"`
	HeaderQuizName      *string         `db:"header_quiz_name"           json:"header_quiz_name"`
	HeaderQuizDesc      *string         `db:"header_quiz_desc"           json:"header_quiz_desc"`
}

func CreateOaRiskProfileQuiz(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_risk_profile_quiz"
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

func CreateMultipleOaRiskProfileQuiz(params []interface{}) (int, error) {

	q := `INSERT INTO oa_risk_profile_quiz (
		oa_request_key, 
		quiz_question_key,
		quiz_option_key,
		quiz_option_score,
		rec_status) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	// log.Println("==========  ==========>>>", q)
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		// log.Error(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func AdminGetOaRiskProfileQuizByOaRequestKey(c *[]AdminOaRiskProfileQuiz, key string) (int, error) {
	query := `SELECT 
				oa_risk_profile_quiz.risk_profile_quiz_key AS risk_profile_quiz_key,
				oa_risk_profile_quiz.quiz_question_key AS quiz_question_key,
				oa_risk_profile_quiz.quiz_option_key AS quiz_option_key_user,
				oa_risk_profile_quiz.quiz_option_score AS quiz_option_score_user,
				cms_quiz_question.quiz_title AS quiz_title,
				cms_quiz_header.quiz_name AS header_quiz_name,
				cms_quiz_header.quiz_desc AS header_quiz_desc
			FROM oa_risk_profile_quiz AS oa_risk_profile_quiz 
			INNER JOIN cms_quiz_question AS cms_quiz_question ON cms_quiz_question.quiz_question_key = oa_risk_profile_quiz.quiz_question_key
			INNER JOIN cms_quiz_options AS cms_quiz_options ON cms_quiz_options.quiz_option_key = oa_risk_profile_quiz.quiz_option_key
			INNER JOIN cms_quiz_header AS cms_quiz_header ON cms_quiz_header.quiz_header_key = cms_quiz_question.quiz_header_key
			WHERE oa_risk_profile_quiz.rec_status = 1 AND oa_risk_profile_quiz.oa_request_key = ` + key +
		` ORDER BY oa_risk_profile_quiz.risk_profile_quiz_key DESC`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func RiskProfileQuizIfNull(c *[]OaRiskProfileQuiz, oarequestkey uint64) (int, error) {
	oa_key := strconv.FormatUint(uint64(oarequestkey), 10)
	query := `SELECT 
	0 AS risk_profile_quiz_key,
	'` + oa_key + `' AS oa_request_key,
	b.quiz_question_key, 
	0 AS quiz_option_key,
	0 AS quiz_option_score,
	b.rec_order,
	b.rec_status,
	NOW() as rec_created_date,
	'system' as rec_created_by,
	null as rec_modified_date,
	null as rec_modified_by,
	null as rec_image1,
	null as rec_image2,
	null as rec_approval_status,
	null as rec_approval_stage,
	null as rec_approved_date,
	null as rec_approved_by,
	null as rec_deleted_date,
	null as rec_deleted_by,
	null as rec_attribute_id1,
	null as rec_attribute_id2,
	NULL as rec_attribute_id3
	FROM cms_quiz_header a 
	INNER JOIN cms_quiz_question b ON (a.quiz_header_key=b.quiz_header_key)
	-- INNER JOIN cms_quiz_options c ON (c.quiz_question_key=b.quiz_question_key AND c.rec_status=1)
	WHERE a.rec_status=1
	AND a.quiz_type_key=2
	AND DATE(NOW()) BETWEEN DATE(a.quiz_published_start) AND DATE(a.quiz_published_thru)
	AND b.rec_status = 1
	ORDER BY b.rec_order;
	`

	// Main query
	// log.Println("===== GET ALL OA RISK PROFILE QUIZ IF NULL ===== >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllOaRiskProfileQuiz(c *[]OaRiskProfileQuiz, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              oa_risk_profile_quiz.* FROM 
			  oa_risk_profile_quiz `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_risk_profile_quiz."+field+" = '"+value+"'")
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
	// log.Println("===== GET ALL OA RISK PROFILE QUIZ ===== >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
func UpdateOaRiskProfileQuiz(params map[string]string) (int, error) {
	query := "UPDATE oa_risk_profile_quiz SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "risk_profile_quiz_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE risk_profile_quiz_key = " + params["risk_profile_quiz_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func DeleteOaRiskProfileQuiz(params map[string]string, riskKey []string, requestKey string) (int, error) {
	inQuery := strings.Join(riskKey, ",")
	query := "UPDATE oa_risk_profile_quiz SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	if len(riskKey) > 0 {
		query += " WHERE rec_status = 1 AND oa_request_key = '" + requestKey + "' AND risk_profile_quiz_key NOT IN(" + inQuery + ")"
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

func GetOaRiskProfileQuiz(c *OaRiskProfileQuiz, key string, field string) (int, error) {
	query := "SELECT oa_risk_profile_quiz.* FROM oa_risk_profile_quiz WHERE oa_risk_profile_quiz.rec_status = 1 AND oa_risk_profile_quiz." + field + " = " + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
