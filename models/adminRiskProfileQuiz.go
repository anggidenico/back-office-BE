package models

import (
	"log"
	"mf-bo-api/db"
)

type QuizQuestionData struct {
	QuizQuestionKey   uint64 `db:"quiz_question_key" json:"quiz_question_key"`
	QuizQuestionTitle string `db:"quiz_title" json:"quiz_title"`
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
	QuizOptionKey   uint64  `db:"quiz_option_key" json:"quiz_option_key"`
	QuizOptionLabel string  `db:"quiz_option_label" json:"quiz_option_label"`
	QuizOptionTitle string  `db:"quiz_option_title" json:"quiz_option_title"`
	QuizOptionScore uint64  `db:"quiz_option_score" json:"quiz_option_score"`
	QuizOptionOrder *uint64 `db:"quiz_option_order" json:"quiz_option_order"`
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

type QuizQuestionDetail struct {
	QuizQuestionKey    uint64  `db:"quiz_question_key" json:"quiz_question_key"`
	QuizHeaderKey      uint64  `db:"quiz_header_key" json:"quiz_header_key"`
	QuizHeaderName     string  `db:"quiz_header_name" json:"quiz_header_name"`
	QuizOptionTypeKey  uint64  `db:"quiz_option_type_key" json:"quiz_option_type_key"`
	QuizOptionTypeName string  `db:"quiz_option_type_name" json:"quiz_option_type_name"`
	QuizQuestionTitle  string  `db:"quiz_title" json:"quiz_title"`
	QuizQuestionOrder  *uint64 `db:"quiz_question_order" json:"quiz_question_order"`
}

func GetQuestionDetail(QuizQuestionKey string) (result QuizQuestionDetail) {
	query := `SELECT t1.quiz_question_key, t2.quiz_header_key, t2.quiz_name quiz_header_name, 
	t1.quiz_title, t1.rec_order quiz_question_order, t3.lookup_key quiz_option_type_key,
	t3.lkp_name quiz_option_type_name
	FROM cms_quiz_question t1
	INNER JOIN cms_quiz_header t2 ON t1.quiz_header_key = t2.quiz_header_key
	INNER JOIN gen_lookup t3 ON t3.lookup_key = t1.quiz_option_type
	WHERE t1.quiz_question_key = ` + QuizQuestionKey

	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func UpdateQuizQuestion(QuizQuestionKey string, params map[string]string) error {
	query := `UPDATE cms_quiz_question SET `
	i := 0
	for key, value := range params {
		if key != "quiz_question_key" {
			query += key + " = '" + value + "'"

			if (len(params) - 1) > i {
				query += ", "
			}
			i++
		}
	}
	query += ` WHERE quiz_question_key = ` + QuizQuestionKey

	log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func CreateQuizQuestion(params map[string]string) error {

	query := "INSERT INTO cms_quiz_question "

	var fields, values string
	var bindvars []interface{}

	for key, value := range params {
		fields += key + ", "
		values += ` "` + value + `", `
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

type QuizOptionDetail struct {
	QuizQuestionKey uint64 `db:"quiz_question_key" json:"quiz_question_key"`
	QuizOptionKey   uint64 `db:"quiz_option_key" json:"quiz_option_key"`
	QuizOptionLabel string `db:"quiz_option_label" json:"quiz_option_label"`
	QuizOptionTitle string `db:"quiz_option_title" json:"quiz_option_title"`
	QuizOptionScore uint64 `db:"quiz_option_score" json:"quiz_option_score"`
}

func GetOptionDetail(QuizOptionKey string) (result QuizOptionDetail) {
	query := `SELECT quiz_question_key, quiz_option_key, quiz_option_label, quiz_option_title, quiz_option_score 
	FROM cms_quiz_options WHERE rec_status = 1 AND quiz_option_key = ` + QuizOptionKey

	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func CreateQuizOption(params map[string]string) error {

	query := "INSERT INTO cms_quiz_options "

	var fields, values string
	var bindvars []interface{}

	for key, value := range params {
		fields += key + ", "
		values += ` "` + value + `", `
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func UpdateQuizOption(QuizOptionKey string, params map[string]string) error {
	query := `UPDATE cms_quiz_options SET `
	i := 0
	for key, value := range params {
		if key != "quiz_option_key" {
			query += key + " = '" + value + "'"

			if (len(params) - 1) > i {
				query += ", "
			}
			i++
		}
	}
	query += ` WHERE quiz_option_key = ` + QuizOptionKey

	log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
