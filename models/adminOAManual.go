package models

import (
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
			resSQL, err := tx.Exec(qUpdateOAReq)
			if err != nil {
				tx.Rollback()
				log.Println(err.Error())
				return err, OaRequestKey
			}
			rowsAffected, _ := resSQL.RowsAffected()
			if rowsAffected == 0 {
				tx.Rollback()
				log.Println(fmt.Errorf("rowsAffected = 0"))
				return fmt.Errorf("rowsAffected = 0"), OaRequestKey
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
				if *udfInfoData.LookupKey == keyInt {
					var udfValueData UdfValue
					qCekExistUdf := `SELECT udf_value_key, udf_info_key, row_data_key FROM udf_value WHERE udf_info_key = ` + strconv.FormatUint(udfInfoData.UdfInfoKey, 10) + ` AND row_data_key = ` + strconv.FormatInt(PersonalDataKey, 10)
					err = db.Db.Get(&udfValueData, qCekExistUdf)
					if err != nil {
						tx.Rollback()
						log.Println(err.Error())
						return err, OaRequestKey
					}

					if udfValueData.RowDataKey > 0 { // DATA SUDAH ADA JADI UPDATE SAJA
						paramsUpdateUdf := make(map[string]string)
						paramsUpdateUdf["udf_value_key"] = strconv.FormatUint(udfValueData.UdfValueKey, 10)
						paramsUpdateUdf["udf_info_key"] = strconv.FormatUint(udfInfoData.UdfInfoKey, 10)
						paramsUpdateUdf["row_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
						paramsUpdateUdf["udf_values"] = value
						qUpdateUdf := GenerateUpdateQuery("udf_value", "udf_value_key", paramsUpdateUdf)
						_, err := tx.Exec(qUpdateUdf)
						if err != nil {
							tx.Rollback()
							log.Println(err.Error())
							return err, OaRequestKey
						}
					} else {
						paramsInsertUdf := make(map[string]string)
						paramsInsertUdf["udf_info_key"] = strconv.FormatUint(udfInfoData.UdfInfoKey, 10)
						paramsInsertUdf["row_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
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

		paramsPersonalData["personal_data_key"] = strconv.FormatInt(PersonalDataKey, 10)
		qUpdatePersonalData := GenerateUpdateQuery("oa_personal_data", "personal_data_key", paramsPersonalData)
		resSQL2, err := tx.Exec(qUpdatePersonalData)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		rowsAffected, _ := resSQL2.RowsAffected()
		if rowsAffected == 0 {
			tx.Rollback()
			log.Println(fmt.Errorf("rowsAffected = 0"))
			return fmt.Errorf("rowsAffected = 0"), OaRequestKey
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

func UpdateOAManual() (error, int64) {
	var OaRequestKey int64

	return nil, OaRequestKey
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
