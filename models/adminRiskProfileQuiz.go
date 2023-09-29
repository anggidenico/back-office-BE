package models

import (
	"log"
	"mf-bo-api/db"
)

type QuizQuestionData struct {
	QuizQuestionKey   uint64 `db:"quiz_question_key" json:"quiz_question_key"`
	QuizQuestionTitle string `db:"quiz_title" json:"quiz_question_title"`
}

func GetQuizQuestion() (result []QuizQuestionData) {
	query := `SELECT t1.quiz_question_key, t1.quiz_title FROM cms_quiz_question t1
	WHERE t1.rec_status = 1 AND t1.quiz_header_key = 2`

	// log.Println("========== GetQuizQuestion ==========>>>", query)

	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}

type QuizOptionData struct {
	QuizOptionKey   uint64 `db:"quiz_option_key" json:"quiz_option_key"`
	QuizOptionLabel string `db:"quiz_option_label" json:"quiz_option_label"`
	QuizOptionTitle string `db:"quiz_option_title" json:"quiz_option_title"`
	QuizOptionScore uint64 `db:"quiz_option_score" json:"quiz_option_score"`
	QuizOptionOrder uint64 `db:"quiz_option_order" json:"quiz_option_order"`
}

func GetQuizOption(QuizQuestionKey string) (result []QuizOptionData) {
	query := `SELECT 
	t1.quiz_option_key,
	t1.quiz_option_label, 
	t1.quiz_option_title, 
	t1.quiz_option_score,
	t1.rec_order as quiz_option_order
	FROM cms_quiz_options t1 WHERE t1.rec_status = 1 AND t1.quiz_question_key = ` + QuizQuestionKey

	// log.Println("========== GetQuizOption ==========>>>", query)

	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err)
	}

	return
}
