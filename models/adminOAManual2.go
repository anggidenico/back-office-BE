package models

import (
	"log"
	"mf-bo-api/db"
	"strconv"
)

type OaQuizAnswer struct {
	QuestionKey int64 `db:"quiz_question_key" json:"quiz_question_key"`
	OptionKey   int64 `db:"quiz_option_key" json:"quiz_option_key"`
}

type OaRequestBankAccountDetails struct {
	BankAccountKey  *uint64 `db:"bank_account_key" json:"bank_account_key"`
	BankKey         *uint64 `db:"bank_key" json:"bank_key"`
	BankAccountNo   *string `db:"account_no" json:"account_no"`
	BankAccountName *string `db:"account_holder_name" json:"account_holder_name"`
	BankBranchName  *string `db:"branch_name" json:"branch_name"`
	CurrencyKey     *uint64 `db:"currency_key" json:"currency_key"`
	FlagPriority    *uint64 `db:"flag_priority" json:"flag_priority"`
}

func CreateOaBankAccount(paramsOARequest map[string]string, paramsBankAccount map[string]string) (error, int64) {
	OaRequestKey, _ := strconv.ParseInt(paramsBankAccount["oa_request_key"], 10, 64)

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if len(paramsOARequest) > 1 {
		qUpdateOAReq := GenerateUpdateQuery("oa_request", "oa_request_key", paramsOARequest)
		_, err := tx.Exec(qUpdateOAReq)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	var ListOaBankAccount []OaRequestBankAccountDetails
	qGetListOaBankAccount := `SELECT a2.bank_key, a2.bank_account_key, a2.account_no , a2.account_holder_name , a2.branch_name, a1.flag_priority FROM oa_request_bank_account a1 JOIN ms_bank_account a2 ON a1.bank_account_key = a2.bank_account_key AND a2.rec_status = 1 JOIN ms_bank a3 ON a2.bank_key = a3.bank_key WHERE a1.rec_status = 1 AND a1.oa_request_key = ` + strconv.FormatInt(OaRequestKey, 10)
	err = db.Db.Select(&ListOaBankAccount, qGetListOaBankAccount)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if len(ListOaBankAccount) == 0 { // JIKA TIDAK ADA SAMA SEKALI MAKA CREATE

		insertMsBankAccount := make(map[string]string)
		insertMsBankAccount["rec_status"] = paramsBankAccount["rec_status"]
		insertMsBankAccount["rec_created_by"] = paramsBankAccount["rec_created_by"]
		insertMsBankAccount["rec_created_date"] = paramsBankAccount["rec_created_date"]
		insertMsBankAccount["bank_key"] = paramsBankAccount["bank_key"]
		insertMsBankAccount["account_no"] = paramsBankAccount["account_no"]
		insertMsBankAccount["account_holder_name"] = paramsBankAccount["account_holder_name"]
		insertMsBankAccount["branch_name"] = paramsBankAccount["branch_name"]
		insertMsBankAccount["currency_key"] = paramsBankAccount["currency_key"]
		insertMsBankAccount["bank_account_type"] = "129"
		insertMsBankAccount["rec_domain"] = "131"

		qInsertMsBankAccount := GenerateInsertQuery("ms_bank_account", insertMsBankAccount)
		resSQL, err := tx.Exec(qInsertMsBankAccount)
		// _, err = tx.Exec(qInsertMsBankAccount)

		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}

		MsBankAccountKey, err := resSQL.LastInsertId()
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		// log.Println(MsBankAccountKey)

		insertOaRequestBankAccount := make(map[string]string)
		insertOaRequestBankAccount["rec_status"] = paramsBankAccount["rec_status"]
		insertOaRequestBankAccount["rec_created_by"] = paramsBankAccount["rec_created_by"]
		insertOaRequestBankAccount["rec_created_date"] = paramsBankAccount["rec_created_date"]
		insertOaRequestBankAccount["oa_request_key"] = strconv.FormatInt(OaRequestKey, 10)
		insertOaRequestBankAccount["bank_account_key"] = strconv.FormatInt(MsBankAccountKey, 10)
		insertOaRequestBankAccount["flag_priority"] = paramsBankAccount["flag_priority"]
		insertOaRequestBankAccount["bank_account_name"] = paramsBankAccount["account_holder_name"]

		qInsertOaRequestBankAccount := GenerateInsertQuery("oa_request_bank_account", insertOaRequestBankAccount)
		_, err = tx.Exec(qInsertOaRequestBankAccount)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	tx.Commit()

	return nil, OaRequestKey
}

func CreateOrUpdateOaRiskProfileQuiz(paramsOaRequest map[string]string, paramsQuizAnswer []OaQuizAnswer) (error, int64) {
	OaRequestKey, _ := strconv.ParseInt(paramsOaRequest["oa_request_key"], 10, 64)

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if len(paramsOaRequest) > 1 {
		qUpdateOaRequest := GenerateUpdateQuery("oa_request", "oa_request_key", paramsOaRequest)
		_, err := tx.Exec(qUpdateOaRequest)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	var total_score int64

	for _, v := range paramsQuizAnswer {

		qGetScore := `SELECT t1.quiz_option_score FROM cms_quiz_options t1 WHERE t1.rec_status = 1 AND t1.quiz_option_key = ` + strconv.FormatInt(v.OptionKey, 10) + ` AND t1.quiz_question_key = ` + strconv.FormatInt(v.QuestionKey, 10)

		var getScore int64
		err = db.Db.Get(&getScore, qGetScore)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}

		CRUDOaRiskProfileQuiz := make(map[string]string)
		CRUDOaRiskProfileQuiz["rec_status"] = paramsOaRequest["rec_status"]
		CRUDOaRiskProfileQuiz["rec_created_by"] = paramsOaRequest["rec_created_by"]
		CRUDOaRiskProfileQuiz["rec_created_date"] = paramsOaRequest["rec_created_date"]
		CRUDOaRiskProfileQuiz["oa_request_key"] = strconv.FormatInt(OaRequestKey, 10)
		CRUDOaRiskProfileQuiz["quiz_question_key"] = strconv.FormatInt(v.QuestionKey, 10)
		CRUDOaRiskProfileQuiz["quiz_option_key"] = strconv.FormatInt(v.OptionKey, 10)
		CRUDOaRiskProfileQuiz["quiz_option_score"] = strconv.FormatInt(getScore, 10)

		total_score += getScore

		// CEK ROW SUDAH ADA  ATAU BELUM
		qCekRows := `SELECT count(*) FROM oa_risk_profile_quiz t1 WHERE t1.rec_status = 1 AND t1.oa_request_key = ` + strconv.FormatInt(OaRequestKey, 10) + ` AND t1.quiz_question_key = ` + strconv.FormatInt(v.QuestionKey, 10)
		var countRows uint64
		err := db.Db.Get(&countRows, qCekRows)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}

		qCRUDOaRiskProfileQuiz := ``

		if countRows == 0 { // JIKA TIDAK ADA MAKA CREATE
			qCRUDOaRiskProfileQuiz = GenerateInsertQuery("oa_risk_profile_quiz", CRUDOaRiskProfileQuiz)
		} else {
			qCRUDOaRiskProfileQuiz = "UPDATE oa_risk_profile_quiz SET "
			i := 0
			for key, value := range CRUDOaRiskProfileQuiz {
				if key != "oa_request_key" {
					if value == "" {
						qCRUDOaRiskProfileQuiz += key + " = NULL"
					} else {
						qCRUDOaRiskProfileQuiz += key + " = '" + value + "'"
					}
					if (len(CRUDOaRiskProfileQuiz) - 2) > i {
						qCRUDOaRiskProfileQuiz += ", "
					}
					i++
				}
			}
			qCRUDOaRiskProfileQuiz += " WHERE oa_request_key = " + CRUDOaRiskProfileQuiz["oa_request_key"] + " AND quiz_question_key = " + CRUDOaRiskProfileQuiz["quiz_question_key"]
		}

		_, err = tx.Exec(qCRUDOaRiskProfileQuiz)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	// GET RISK PROFILE KEY
	var RiskProfileKey int64
	qRiskProfileKey := `SELECT risk_profile_key FROM ms_risk_profile WHERE rec_status = 1 
	AND ` + strconv.FormatInt(total_score, 10) + ` >= min_score 
	AND ` + strconv.FormatInt(total_score, 10) + ` <= max_score 
	ORDER BY risk_profile_key ASC LIMIT 1`
	err = db.Db.Get(&RiskProfileKey, qRiskProfileKey)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	CRUDOaRiskProfile := make(map[string]string)
	CRUDOaRiskProfile["rec_status"] = paramsOaRequest["rec_status"]
	CRUDOaRiskProfile["oa_request_key"] = strconv.FormatInt(OaRequestKey, 10)
	CRUDOaRiskProfile["score_result"] = strconv.FormatInt(total_score, 10)
	CRUDOaRiskProfile["risk_profile_key"] = strconv.FormatInt(RiskProfileKey, 10)

	qCekRows := `SELECT count(oa_risk_profile_key) FROM oa_risk_profile WHERE rec_status = 1 AND oa_request_key = ` + strconv.FormatInt(OaRequestKey, 10)
	var countRows uint64
	err = db.Db.Get(&countRows, qCekRows)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	qCRUDOaRiskProfile := ``

	if countRows == 0 { // JIKA TIDAK ADA MAKA CREATE
		CRUDOaRiskProfile["rec_created_by"] = paramsOaRequest["rec_created_by"]
		CRUDOaRiskProfile["rec_created_date"] = paramsOaRequest["rec_created_date"]
		qCRUDOaRiskProfile = GenerateInsertQuery("oa_risk_profile", CRUDOaRiskProfile)
	} else { // JIKA SUDAH ADA MAKA UPDATE
		CRUDOaRiskProfile["rec_modified_by"] = paramsOaRequest["rec_created_by"]
		CRUDOaRiskProfile["rec_modified_date"] = paramsOaRequest["rec_created_date"]
		qCRUDOaRiskProfile = GenerateUpdateQuery("oa_risk_profile", "oa_request_key", CRUDOaRiskProfile)
	}

	_, err = tx.Exec(qCRUDOaRiskProfile)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	tx.Commit()
	return nil, OaRequestKey
}

func CreateOrUpdateFileOaManual(paramsOaRequest map[string]string, paramsFile map[string]string) (error, int64) {
	OaRequestKey, _ := strconv.ParseInt(paramsOaRequest["oa_request_key"], 10, 64)

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if len(paramsOaRequest) > 1 {
		qUpdateOaRequest := GenerateUpdateQuery("oa_request", "oa_request_key", paramsOaRequest)
		_, err := tx.Exec(qUpdateOaRequest)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	if len(paramsFile) > 0 {
		qUpdateFile := GenerateInsertQuery("ms_file", paramsFile)
		_, err := tx.Exec(qUpdateFile)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	}

	tx.Commit()
	return nil, OaRequestKey
}
