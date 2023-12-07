package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"strconv"
)

func CreateOrUpdateOAManual(paramsOARequest map[string]string, paramsPersonalData map[string]string, paramsIDCard map[string]string, paramsDomicile map[string]string, paramsOfficeAddress map[string]string, paramsOther map[string]string) (error, int64) {

	var OaRequestKey int64

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if _, ok := paramsOARequest["oa_request_key"]; ok {

		// UPDATE EXISTING OA
		OaRequestKey, _ = strconv.ParseInt(paramsOARequest["oa_request_key"], 10, 64)

		if len(paramsOARequest) > 1 {
			qUpdateOAReq := GenerateUpdateQuery("oa_request", "oa_request_key", paramsOARequest)
			_, err := tx.Exec(qUpdateOAReq)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
		}

		// GET PERSONAL DATA KEY
		var PersonalDataKey int64
		qPersonalDataKey := `SELECT personal_data_key FROM oa_personal_data WHERE oa_request_key = ` + paramsOARequest["oa_request_key"] + ` ORDER BY personal_data_key DESC LIMIT 1`
		err = db.Db.Get(&PersonalDataKey, qPersonalDataKey)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}

		// GET ID CARD ADDRESS KEY
		var IdCardAddrKey *int64
		qIdCardAddrKey := `SELECT idcard_address_key FROM oa_personal_data WHERE personal_data_key = ` + strconv.FormatInt(PersonalDataKey, 10) + ` ORDER BY personal_data_key DESC LIMIT 1`
		err = db.Db.Get(&IdCardAddrKey, qIdCardAddrKey)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		if len(paramsIDCard) > 0 && IdCardAddrKey == nil {
			// INSERT ID CARD ADDRESS
			qInsertAddrKTP := GenerateInsertQuery("oa_postal_address", paramsIDCard)
			resSQL3, err := tx.Exec(qInsertAddrKTP)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			KTPAddrKey, _ := resSQL3.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
			UpdAddrToPrsnData["idcard_address_key"] = strconv.FormatInt(KTPAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			resSQL5, err := tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			a, _ := resSQL5.RowsAffected()
			if a == 0 {
				tx.Rollback()
				log.Println(fmt.Errorf("rowsAffected = 0"))
				return fmt.Errorf("rowsAffected = 0"), OaRequestKey
			}
			IdCardAddrKey = &KTPAddrKey
		}

		// GET DOMICILE ADDRESS KEY
		var DomicileAddrKey *int64
		qDomicileAddrKey := `SELECT domicile_address_key FROM oa_personal_data WHERE personal_data_key = ` + strconv.FormatInt(PersonalDataKey, 10) + ` ORDER BY personal_data_key DESC LIMIT 1`
		err = db.Db.Get(&DomicileAddrKey, qDomicileAddrKey)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		if len(paramsDomicile) > 0 && DomicileAddrKey == nil {
			// INSERT DOMICILE ADDRESS
			qInsertAddrDomicile := GenerateInsertQuery("oa_postal_address", paramsDomicile)
			resSQL4, err := tx.Exec(qInsertAddrDomicile)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			DomAddrKey, _ := resSQL4.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
			UpdAddrToPrsnData["domicile_address_key"] = strconv.FormatInt(DomAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			_, err = tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}

			DomicileAddrKey = &DomAddrKey
		}

		// GET OFFICE ADDRESS KEY
		var OfficeAddrKey *int64
		qOfficeAddrKey := `SELECT occup_address_key FROM oa_personal_data WHERE personal_data_key = ` + strconv.FormatInt(PersonalDataKey, 10) + ` ORDER BY personal_data_key DESC LIMIT 1`
		err = db.Db.Get(&OfficeAddrKey, qOfficeAddrKey)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		if len(paramsOfficeAddress) > 0 && OfficeAddrKey == nil {
			// INSERT OFFICE ADDRESS
			qInsertAddrOffice := GenerateInsertQuery("oa_postal_address", paramsOfficeAddress)
			resSQL4, err := tx.Exec(qInsertAddrOffice)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			OffAddrKey, _ := resSQL4.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
			UpdAddrToPrsnData["occup_address_key"] = strconv.FormatInt(OffAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			resSQL5, err := tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			a, _ := resSQL5.RowsAffected()
			if a == 0 {
				tx.Rollback()
				log.Println(fmt.Errorf("rowsAffected = 0"))
				return fmt.Errorf("rowsAffected = 0"), OaRequestKey
			}
			OfficeAddrKey = &OffAddrKey
		}

		// INSERT OR UPDATE OTHERS UDF VALUE
		if len(paramsOther) > 0 {
			for key, value := range paramsOther {

				keyInt, _ := strconv.ParseUint(key, 10, 64)
				var udfInfoData UdfInfo
				qCekOthers := `SELECT udf_info_key, lookup_key, udf_info_name, udf_info_code FROM udf_info WHERE lookup_key = ` + key
				err = db.Db.Get(&udfInfoData, qCekOthers)
				if err != nil {
					tx.Rollback()
					log.Println(err.Error())
					return err, OaRequestKey
				}
				if *udfInfoData.LookupKey == keyInt { // JIKA TERMASUK DI LOOKUPKEY PADA TABEL UDF INFO

					var udfValueData UdfValue
					qCekExistUdf := `SELECT udf_value_key, udf_info_key, row_data_key FROM udf_value WHERE udf_info_key = ` + strconv.FormatUint(udfInfoData.UdfInfoKey, 10) + ` AND row_data_key = ` + strconv.FormatInt(PersonalDataKey, 10)
					err = db.Db.Get(&udfValueData, qCekExistUdf)
					if err == sql.ErrNoRows { // JIKA NO ROWS MAKA INPUT BARU
						paramsInsertUdf := make(map[string]string)
						paramsInsertUdf["udf_info_key"] = strconv.FormatUint(udfInfoData.UdfInfoKey, 10)
						paramsInsertUdf["row_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
						paramsInsertUdf["udf_values"] = value
						qInsertUdf := GenerateInsertQuery("udf_value", paramsInsertUdf)
						_, err = tx.Exec(qInsertUdf)
						if err != nil {
							tx.Rollback()
							log.Println(err.Error())
							return err, OaRequestKey
						}
					}
					if err == nil {

						paramsUpdateUdf := make(map[string]string)
						paramsUpdateUdf["udf_value_key"] = strconv.FormatUint(udfValueData.UdfValueKey, 10)
						paramsUpdateUdf["udf_values"] = value
						qUpdateUdf := GenerateUpdateQuery("udf_value", "udf_value_key", paramsUpdateUdf)
						_, err = tx.Exec(qUpdateUdf)
						if err != nil {
							tx.Rollback()
							log.Println(err.Error())
							return err, OaRequestKey
						}
					}
					if err != nil && err != sql.ErrNoRows {
						tx.Rollback()
						log.Println(err.Error())
						return err, OaRequestKey
					}

				}
			}
		}

		paramsPersonalData["personal_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
		qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", paramsPersonalData)
		_, err := tx.Exec(qUpdatePersonalData)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
	} else {

		// CREATE NEW OA

		qInsertOAReq := GenerateInsertQuery("oa_request", paramsOARequest)
		resSQL, err := tx.Exec(qInsertOAReq)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		OaRequestKey, _ = resSQL.LastInsertId()

		paramsPersonalData["oa_request_key"] = strconv.FormatInt(OaRequestKey, 10)
		qInsertPersonalData := GenerateInsertQuery("oa_personal_data", paramsPersonalData)
		resSQL2, err := tx.Exec(qInsertPersonalData)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		OaPersonalDataKey, _ := resSQL2.LastInsertId()

		if len(paramsIDCard) > 0 {
			qInsertAddrKTP := GenerateInsertQuery("oa_postal_address", paramsIDCard)
			resSQL3, err := tx.Exec(qInsertAddrKTP)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			KTPAddrKey, _ := resSQL3.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(OaPersonalDataKey, 10)
			UpdAddrToPrsnData["idcard_address_key"] = strconv.FormatInt(KTPAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			resSQL5, err := tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			a, _ := resSQL5.RowsAffected()
			if a == 0 {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
		}

		if len(paramsDomicile) > 0 {
			qInsertAddrDomicile := GenerateInsertQuery("oa_postal_address", paramsDomicile)
			resSQL4, err := tx.Exec(qInsertAddrDomicile)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			DomAddrKey, _ := resSQL4.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(OaPersonalDataKey, 10)
			UpdAddrToPrsnData["domicile_address_key"] = strconv.FormatInt(DomAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			resSQL5, err := tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			a, _ := resSQL5.RowsAffected()
			if a == 0 {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
		}

		if len(paramsOfficeAddress) > 0 {
			qInsertAddrOffice := GenerateInsertQuery("oa_postal_address", paramsOfficeAddress)
			resSQL4, err := tx.Exec(qInsertAddrOffice)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			OfficeAddrKey, _ := resSQL4.LastInsertId()
			UpdAddrToPrsnData := make(map[string]string)
			UpdAddrToPrsnData["personal_data_key"] = strconv.FormatInt(OaPersonalDataKey, 10)
			UpdAddrToPrsnData["occup_address_key"] = strconv.FormatInt(OfficeAddrKey, 10)
			qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", UpdAddrToPrsnData)
			resSQL5, err := tx.Exec(qUpdatePersonalData)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			a, _ := resSQL5.RowsAffected()
			if a == 0 {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
		}

		// INSERT OR UPDATE OTHERS UDF VALUE
		if len(paramsOther) > 0 {
			for key, value := range paramsOther {
				keyInt, _ := strconv.ParseUint(key, 10, 64)
				var udfInfoData UdfInfo
				qCekOthers := `SELECT udf_info_key, lookup_key, udf_info_name, udf_info_code FROM udf_info WHERE lookup_key = ` + key
				err = db.Db.Get(&udfInfoData, qCekOthers)
				if err != nil {
					tx.Rollback()
					log.Println(err.Error())
					return err, OaRequestKey
				}
				if *udfInfoData.LookupKey == keyInt {
					paramsInsertUdf := make(map[string]string)
					paramsInsertUdf["udf_info_key"] = strconv.FormatUint(udfInfoData.UdfInfoKey, 10)
					paramsInsertUdf["row_data_key"] = strconv.FormatInt(OaPersonalDataKey, 10)
					paramsInsertUdf["udf_values"] = value
					qInsertUdf := GenerateInsertQuery("udf_value", paramsInsertUdf)
					_, err := tx.Exec(qInsertUdf)
					if err != nil {
						tx.Rollback()
						log.Println(err.Error())
						return err, OaRequestKey
					}
				}
			}
		}

	}

	tx.Commit()
	return nil, OaRequestKey
}

func CreateOaBankAccount(paramsBankAccount map[string]string) (error, int64) {
	OaRequestKey, _ := strconv.ParseInt(paramsBankAccount["oa_request_key"], 10, 64)

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
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
	CRUDOaRiskProfile["rec_created_by"] = paramsOaRequest["rec_created_by"]
	CRUDOaRiskProfile["rec_created_date"] = paramsOaRequest["rec_created_date"]
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
		qCRUDOaRiskProfile = GenerateInsertQuery("oa_risk_profile", CRUDOaRiskProfile)
	} else {
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

func GetCountryList(paramSearch map[string]string) []CountryModels {
	var result []CountryModels

	query := `SELECT country_key, country_code, country_name, calling_code FROM ms_country WHERE rec_status = 1`

	if len(paramSearch) > 0 {
		for key, value := range paramSearch {
			if value != "" {
				if key == "country_name" {
					query += ` AND ` + key + ` LIKE '%` + value + `%'`
				} else {
					query += ` AND ` + key + ` = '` + value + `'`
				}
			}
		}
	}

	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}

func GetCityList(paramSearch map[string]string) []CityModels {
	var result []CityModels

	query := `SELECT city_key, city_code, city_name, postal_code FROM ms_city WHERE rec_status = 1`

	if len(paramSearch) > 0 {
		for key, value := range paramSearch {
			if value != "" {
				if key == "city_name" {
					query += ` AND ` + key + ` LIKE '%` + value + `%'`
				} else {
					query += ` AND ` + key + ` = '` + value + `'`
				}
			}
		}
	}
	// log.Println(query)
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}

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

type CountryModels struct {
	CountryKey  uint64  `db:"country_key" json:"country_key"`
	CountryCode string  `db:"country_code" json:"country_code"`
	CountryName string  `db:"country_name" json:"country_name"`
	CallingCode *string `db:"calling_code" json:"calling_code"`
}

type CityModels struct {
	CityKey    uint64  `db:"city_key" json:"city_key"`
	CityCode   string  `db:"city_code" json:"city_code"`
	CityName   string  `db:"city_name" json:"city_name"`
	PostalCode *string `db:"postal_code" json:"postal_code"`
}
