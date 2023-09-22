package models

import (
	"log"
	"math"
	"mf-bo-api/db"
	"strconv"
)

type RiskProfileListModels struct {
	OaRequestKey uint64  `db:"oa_request_key" json:"oa_request_key"`
	CustomerKey  uint64  `db:"customer_key" json:"customer_key"`
	Cif          string  `db:"cif" json:"cif"`
	FullName     string  `db:"full_name" json:"full_name"`
	Email        string  `db:"email_address" json:"email_address"`
	OaStatus     string  `db:"oa_status" json:"oa_status"`
	OaDate       string  `db:"oa_date" json:"oa_date"`
	OaSource     *string `db:"oa_source" json:"oa_source"`
}

type RiskProfileDetailResponse struct {
	Cif                   string                        `json:"cif"`
	FullName              string                        `json:"full_name"`
	OaStatus              string                        `json:"oa_status"`
	OaDate                string                        `json:"oa_date"`
	OaSource              *string                       `json:"oa_source"`
	RiskProfileQuizAnswer []RiskProfileQuizAnswerModels `json:"risk_profile_quiz_answer"`
	RiskProfileQuizResult RiskProfileQuizResultModels   `json:"risk_profile_quiz_result"`
}

type RiskProfileQuizAnswerModels struct {
	QuizTitle       string `db:"quiz_title" json:"quiz_title"`
	QuizOptionTitle string `db:"quiz_option_title" json:"quiz_option_title"`
	QuizOptionScore uint64 `db:"quiz_option_score" json:"quiz_option_score"`
}

type RiskProfileQuizResultModels struct {
	OaStatus       string  `db:"oa_status" json:"oa_status"`
	OaDate         string  `db:"oa_date" json:"oa_date"`
	OaSource       *string `db:"oa_source" json:"oa_source"`
	Cif            string  `db:"cif" json:"cif"`
	FullName       string  `db:"full_name" json:"full_name"`
	ScoreResult    uint64  `db:"score_result" json:"score_result"`
	RiskName       string  `db:"risk_name" json:"risk_name"`
	RiskDesc       string  `db:"risk_desc" json:"risk_desk"`
	RiskCode       string  `db:"risk_code" json:"risk_code"`
	RiskProfileKey uint64  `db:"risk_profile_key" json:"risk_profile_key"`
}

func GetPengkinianRiskProfileListQuery(c *[]RiskProfileListModels, backOfficeRole uint64, limit uint64, offset uint64) int {
	query := `SELECT t1.oa_request_key, t4.customer_key, t3.ulogin_email AS email_address, 
	t2.lkp_name AS oa_status,  t1.oa_entry_start AS oa_date, t4.unit_holder_idno AS cif, 
	t4.full_name, t5.lkp_name AS oa_source
	FROM oa_request t1
	INNER JOIN gen_lookup t2 ON t1.oa_status = t2.lookup_key
	LEFT JOIN gen_lookup t5 ON t1.oa_source = t5.lookup_key
	INNER JOIN sc_user_login t3 ON t3.user_login_key = t1.user_login_key
	INNER JOIN ms_customer t4 ON t4.customer_key = t3.customer_key
	WHERE t1.rec_status = 1 AND t1.oa_request_type = 128`

	if backOfficeRole == 11 {
		query += ` AND t1.oa_status = 258`
	}
	if backOfficeRole == 12 {
		query += ` AND t1.oa_status = 259`
	}

	queryCountPage := `SELECT count(*) FROM
	( ` + query + `) t1`

	// log.Println(limit)
	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// EXECUTE DATANYA
	// log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
	}

	// EXECUTE PAGING
	var pagination int
	var count uint64
	// log.Println(queryPage)
	err = db.Db.Get(&count, queryCountPage)
	if err != nil {
		log.Println(err.Error())
	}

	if limit > 0 {
		if count < limit {
			pagination = 1
		} else {
			calc := math.Ceil(float64(count) / float64(limit))
			pagination = int(calc)
		}
	}

	return pagination
}

func GetQuizQuestionAnswerQuery(OaRequestKey string) []RiskProfileQuizAnswerModels {
	query := `SELECT t3.quiz_title, t4.quiz_option_title, t4.quiz_option_score
	FROM oa_request t1
	INNER JOIN oa_risk_profile_quiz t2 ON t1.oa_request_key = t2.oa_request_key
	INNER JOIN cms_quiz_question t3 ON t2.quiz_question_key = t3.quiz_question_key
	INNER JOIN cms_quiz_options t4 ON t2.quiz_option_key = t4.quiz_option_key
	WHERE t1.rec_status = 1 AND t2.rec_status = 1 AND t2.oa_request_key = ` + OaRequestKey + `
	ORDER BY t3.quiz_question_key ASC`
	var result []RiskProfileQuizAnswerModels
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

func GetRiskProfileQuizResult(OaRequestKey string) RiskProfileQuizResultModels {
	query := `SELECT t2.score_result, t3.risk_name, t3.risk_desc, t3.risk_code, t3.risk_profile_key,
	t5.unit_holder_idno AS cif, t5.full_name, src.lkp_name AS oa_source, t1.oa_entry_start AS oa_date, 
	stt.lkp_name AS oa_status
	FROM oa_request t1
	LEFT JOIN gen_lookup src ON t1.oa_source = src.lookup_key
	LEFT JOIN gen_lookup stt ON t1.oa_status = stt.lookup_key
	INNER JOIN oa_risk_profile t2 ON t2.oa_request_key = t1.oa_request_key
	INNER JOIN ms_risk_profile t3 ON t3.risk_profile_key = t2.risk_profile_key
	INNER JOIN sc_user_login t4 ON t1.user_login_key = t4.user_login_key
	INNER JOIN ms_customer t5 ON t4.customer_key = t5.customer_key
	WHERE t1.rec_status = 1 AND t2.rec_status = 1 AND t2.oa_request_key = ` + OaRequestKey + `
	ORDER BY t1.oa_request_key DESC LIMIT 1`
	var result RiskProfileQuizResultModels
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}
