package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type SecuritiesSector struct {
	SectorKey         int64   `db:"sector_key" json:"sector_key"`
	SectorParentKey   *int64  `db:"sector_parent_key" json:"sector_parent_key"`
	ParentSectorName  *string `db:"parent_sector_name" json:"parent_sector_name"`
	SectorCode        *string `db:"sector_code" json:"sector_code"`
	SectorName        *string `db:"sector_name" json:"sector_name"`
	SectorDescription *string `db:"sector_description" json:"sector_description"`
	RecOrder          *int64  `db:"rec_order" json:"rec_order"`
}

func GetSecuritiesSectorModels(c *[]SecuritiesSector) (int, error) {
	query := `SELECT
    a.sector_key,
    a.sector_parent_key,
    a.sector_code,
    a.sector_name,
    a.sector_description,
    a.rec_order,
    b.sector_name AS parent_sector_name
FROM
    ms_securities_sector a
LEFT JOIN
    ms_securities_sector b ON a.sector_parent_key = b.sector_key
	WHERE a.rec_status = 1 
	ORDER BY sector_key DESC`
	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func GetSecuritiesSectorDetailModels(c *SecuritiesSector, SectorKey string) (int, error) {
	query := `SELECT
    a.sector_key,
    a.sector_parent_key,
    a.sector_code,
    a.sector_name,
    a.sector_description,
    a.rec_order,
    b.sector_name AS parent_sector_name
FROM
    ms_securities_sector a
LEFT JOIN
    ms_securities_sector b ON a.sector_parent_key = b.sector_key
	WHERE a.rec_status = 1 AND a.sector_key =` + SectorKey + `ORDER BY a.sector_key DESC`

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("sector_key not found")
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func CheckDuplicateSecuritiesSector(SectorCode, SectorName string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT sector_key FROM ms_securities_sector WHERE sector_code = ? AND sector_name = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, SectorCode, SectorName).Scan(&key)

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

func CreateSecuritiesSector(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateSecuritiesSector(params["sector_code"], params["sector_name"])
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

	query := "INSERT INTO ms_securities_sector (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdateSecuritiesSector(SectorKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_securities_sector SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sector_key = ?`
	values = append(values, SectorKey)

	log.Println("========== UpdateMasterSecuritiesSector ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteSecuritiesSectorModels(SectorKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_securities_sector SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "sector_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE sector_key = ?`
	values = append(values, SectorKey)

	resultSQL, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}
	rows, _ := resultSQL.RowsAffected()
	if rows < 1 {
		log.Println("nothing rows affected")
		err2 := fmt.Errorf("nothing rows affected")
		return http.StatusNotFound, err2
	}

	return http.StatusOK, nil
}
