package models

import (
	"log"
	"mf-bo-api/db"
)

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
				query += ` AND ` + key + ` = '` + value + `'`
			}
		}
	}
	log.Println(query)
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}
