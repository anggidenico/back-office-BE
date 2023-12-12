package models

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type AllocInstrument struct {
	AllocInstrumentKey int64           `db:"alloc_instrument_key" json:"alloc_instrument_key"`
	ProductKey         int64           `db:"product_key" json:"product_key"`
	ProductName        string          `db:"product_name" json:"product_name"`
	PeriodeKey         int64           `db:"periode_key" json:"periode_key"`
	PeriodeName        string          `db:"periode_name" json:"periode_name"`
	InstrumentKey      int64           `db:"instrument_key" json:"instrument_key"`
	InstrumentName     string          `db:"instrument_name" json:"instrument_name"`
	InstrumentValue    decimal.Decimal `db:"instrument_value" json:"instrument_value"`
	RecOrder           *int64          `db:"rec_order" json:"rec_order"`
}

func CheckDuplicateAllocInstrument(productKey, periodeKey, instrumentKey string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT alloc_instrument_key FROM ffs_alloc_instrument WHERE product_key = ? AND periode_key = ? AND instrument_key = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, productKey, periodeKey, instrumentKey).Scan(&key)

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

func CreateAllocInstrument(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateAllocInstrument(params["product_key"], params["periode_key"], params["instrument_key"])
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

	query := "INSERT INTO ffs_alloc_instrument (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdateAllocInstrument(allocInstrumentKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_alloc_instrument SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE alloc_instrument_key = ?`
	values = append(values, allocInstrumentKey)

	log.Println("========== UpdatePortfolioInstrument ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetAllocInstrumentStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM ffs_alloc_instrument WHERE alloc_instrument_key = ?"
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

func GetAllocInstrumentModels(c *[]AllocInstrument) (int, error) {
	query := `SELECT a.alloc_instrument_key,
	a.product_key,
	b.product_name,
	a.periode_key,
	c.periode_name,
	a.instrument_key,
	d.instrument_name,
	a.rec_order
	FROM ffs_alloc_instrument a
	JOIN ms_product b ON a.product_key = b.product_key
	JOIN ffs_periode c ON a.periode_key = c.periode_key
	JOIN ms_instrument d ON a.instrument_key = d.instrument_key
	WHERE a.rec_status = 1 ORDER BY a.alloc_instrument_key DESC` //order by

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
