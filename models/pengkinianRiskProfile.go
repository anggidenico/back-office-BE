package models

import (
	"log"
	"mf-bo-api/db"
)

type RiskProfileDetailResponse struct {
	RiskProfileQuizAnswer []RiskProfileQuizAnswerModels `json:"risk_profile_quiz_answer"`
	RiskProfileQuizResult RiskProfileQuizResultModels   `json:"risk_profile_quiz_result"`
}

type RiskProfileQuizAnswerModels struct {
	QuizTitle       string `db:"quiz_title" json:"quiz_title"`
	QuizOptionTitle string `db:"quiz_option_title" json:"quiz_option_title"`
	QuizOptionScore uint64 `db:"quiz_option_score" json:"quiz_option_score"`
}

type RiskProfileQuizResultModels struct {
	ScoreResult    uint64 `db:"score_result" json:"score_result"`
	RiskName       string `db:"risk_name" json:"risk_name"`
	RiskDesc       string `db:"risk_desc" json:"risk_desk"`
	RiskCode       string `db:"risk_code" json:"risk_code"`
	RiskProfileKey uint64 `db:"risk_profile_key" json:"risk_profile_key"`
}

func GetQuizQuestionAnswerQuery(OaRequestKey string) []RiskProfileQuizAnswerModels {
	query := `SELECT t3.quiz_title, t4.quiz_option_title, t4.quiz_option_score
	FROM oa_request t1
	INNER JOIN oa_risk_profile_quiz t2 ON t1.oa_request_key = t2.oa_request_key
	INNER JOIN cms_quiz_question t3 ON t2.quiz_question_key = t3.quiz_question_key
	INNER JOIN cms_quiz_options t4 ON t2.quiz_option_key = t4.quiz_option_key
	WHERE t1.rec_status = 1 AND t2.rec_status = 1 AND t2.oa_request_key = ` + OaRequestKey
	var result []RiskProfileQuizAnswerModels
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

func GetRiskProfileQuizResult(OaRequestKey string) RiskProfileQuizResultModels {
	query := `SELECT t2.score_result, t3.risk_name, t3.risk_desc, t3.risk_code, t3.risk_profile_key
	FROM oa_request t1 
	INNER JOIN oa_risk_profile t2 ON t2.oa_request_key = t1.oa_request_key
	INNER JOIN ms_risk_profile t3 ON t3.risk_profile_key = t2.risk_profile_key
	WHERE t1.rec_status = 1 AND t2.rec_status = 1 AND t2.oa_request_key = ` + OaRequestKey + `
	ORDER BY t1.oa_request_key DESC LIMIT 1`
	var result RiskProfileQuizResultModels
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}
