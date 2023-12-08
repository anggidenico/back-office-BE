package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type FfsPeriode struct {
	PeriodeKey  uint64  `db:"periode_key"          json:"periode_key"`
	PeriodeDate string  `db:"periode_date"         json:"periode_date"`
	PeriodeName string  `db:"periode_name"         json:"periode_name"`
	DateOpened  *string `db:"date_opened"          json:"date_opened"`
	DateClosed  *string `db:"date_closed"          json:"date_closed"`
	RecStatus   uint8   `db:"rec_status"           json:"rec_status"`
	Instrument  *bool   `db:"instrument"           json:"instrument"`
	Sector      *bool   `db:"sector"               json:"sector"`
	Top10       *bool   `db:"top10"                json:"top10"`
}
type FfsPeriodeDetail struct {
	PeriodeKey  uint64  `db:"periode_key"          json:"periode_key"`
	PeriodeDate string  `db:"periode_date"         json:"periode_date"`
	PeriodeName string  `db:"periode_name"         json:"periode_name"`
	DateOpened  *string `db:"date_opened"          json:"date_opened"`
	DateClosed  *string `db:"date_closed"          json:"date_closed"`
	RecStatus   uint8   `db:"rec_status"           json:"rec_status"`
}

func GetFfsPeriodeModels(c *[]FfsPeriode) (int, error) {
	query := `SELECT periode_key,
	periode_date,
	periode_name,
	date_opened,
	date_closed 
	FROM ffs_periode 
	WHERE rec_status = 1 order by rec_order`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
	}
	for i := range *c {
		// Field Instrument diatur ke false jika nil
		if (*c)[i].Instrument == nil {
			instrumentDefault := false
			(*c)[i].Instrument = &instrumentDefault
		}

		// Field Sector diatur ke false jika nil
		if (*c)[i].Sector == nil {
			sectorDefault := false
			(*c)[i].Sector = &sectorDefault
		}

		// Field Top10 diatur ke false jika nil
		if (*c)[i].Top10 == nil {
			top10Default := false
			(*c)[i].Top10 = &top10Default
		}
	}
	return http.StatusOK, nil
}
func GetFfsPeriodeDetailModels(c *FfsPeriodeDetail, PeriodeKey string) (int, error) {
	query := `SELECT periode_key,
	periode_date,
	periode_name,
	date_opened,
	date_closed FROM ffs_periode
	WHERE rec_status = 1 
	AND periode_key =` + PeriodeKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Periode key not found")
			return http.StatusBadGateway, err
		}

		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

// 	err := db.Db.Get(c, query)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println(err.Error())
// 			return http.StatusBadGateway, err
// 		}
// 	}
// 	return http.StatusOK, nil
// }

func CreateFfsPeriode(params map[string]string) (int, error) {
	query := "INSERT INTO ffs_periode"
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
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
func UpdateFfsPeriode(PeriodeKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_periode SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "periode_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE periode_key = ?`
	values = append(values, PeriodeKey)

	log.Println("========== UpdateFFSPeriode ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
func DeleteFfsPeriode(PeriodeKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_periode SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "periode_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE periode_key = ?`
	values = append(values, PeriodeKey)

	log.Println("========== DeleteFfsPeriode ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
