package models

import (
	"log"
	"mf-bo-api/db"
	"strconv"
)

func CreateOrUpdateOAManual(paramsOARequest map[string]string, paramsPersonalData map[string]string, paramsIDCard map[string]string, paramsDomicile map[string]string) (error, int64) {

	var OaRequestKey int64

	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err, OaRequestKey
	}

	if _, ok := paramsOARequest["oa_request_key"]; ok { // JIKA TRUE MAKA UPDATE

	} else { // JIKA FALSE MAKA CREATE

		// var resSQL sql.Result

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

		qInsertAddrKTP := GenerateInsertQuery("oa_postal_address", paramsIDCard)
		resSQL3, err := tx.Exec(qInsertAddrKTP)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return err, OaRequestKey
		}
		KTPAddrKey, _ := resSQL3.LastInsertId()

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
		UpdAddrToPrsnData["idcard_address_key"] = strconv.FormatInt(KTPAddrKey, 10)
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

	tx.Commit()
	return nil, OaRequestKey
}

func GenerateInsertQuery(tableName string, params map[string]string) string {
	query := "INSERT INTO " + tableName
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + `, `
		values += ` '` + value + `', `
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	return query
}

func GenerateUpdateQuery(tableName string, primaryKeyField string, params map[string]string) string {
	query := `UPDATE ` + tableName + ` SET `
	i := 0
	for key, value := range params {
		if key != primaryKeyField {
			query += key + " = '" + value + "'"
			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += ` WHERE ` + primaryKeyField + ` = ` + params[primaryKeyField]

	return query
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
