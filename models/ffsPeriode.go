package models

import (
	"database/sql"
	"errors"
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
	RecOrder    *uint8  `db:"rec_order"           json:"rec_order"`
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
	RecOrder    *uint8  `db:"rec_order"           json:"rec_order"`
}

func GetFfsPeriodeModels(c *[]FfsPeriode) (int, error) {
	query := `SELECT periode_key,
	periode_date,
	periode_name,
	date_opened,
	date_closed 
	FROM ffs_periode 
	WHERE rec_status = 1 order by rec_order DESC`

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
	date_closed,
	rec_order 
	FROM ffs_periode 
	WHERE rec_status = 1 
	AND periode_key = ` + PeriodeKey

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

func CheckDuplicatePeriode(PeriodeName string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT periode_key FROM ffs_periode WHERE periode_name = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, PeriodeName).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			// No duplicate found
			return false, "", nil
		}
		// Other error occurred
		return false, "", err
	}

	// Duplicate found
	return true, key, nil
}

func CreatePeriode(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicatePeriode(params["periode_name"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		return http.StatusBadRequest, errors.New("data duplikat ditemukan")
	}

	// Jika tidak ada duplikasi, buat data baru
	fields := ""
	placeholders := ""
	var bindvars []interface{}

	for key, value := range params {
		fields += key + `, `
		if value == "NULL" {
			placeholders += `NULL, `
		} else {
			placeholders += `?, `
			bindvars = append(bindvars, value)
		}
	}

	fields = fields[:len(fields)-2]
	placeholders = placeholders[:len(placeholders)-2]

	query := "INSERT INTO ffs_periode (" + fields + ") VALUES (" + placeholders + ")"

	tx, err := db.Db.Begin()
	if err != nil {
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	if err != nil {
		tx.Rollback()
		return http.StatusBadRequest, err
	}

	tx.Commit()

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

func GetPeriodeStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM ffs_periode WHERE periode_key = ?"
	var status int
	err := db.Db.QueryRow(query, key).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Data tidak ditemukan
			return 0, nil
		}
		// Terjadi error lain
		return 0, err
	}
	return status, nil
}
