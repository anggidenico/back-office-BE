package models

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

func CheckDuplicateAllocInstrument(periodeKey int64, productKey int64, instrumentKey int64) (bool, string, error) { //dari sini
	// Query to check for duplicates
	query := "SELECT alloc_instrument_key FROM ffs_alloc_instrument WHERE product_key = ? OR periode_key = ? OR instrument_key = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, periodeKey, productKey, instrumentKey).Scan(&key)

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

func CreateAllocInstrument(params map[string]interface{}) (int, error) {
	periodeKey, ok := params["periode_key"].(int64)
	if !ok {
		return http.StatusBadRequest, errors.New("invalid periode_key")
	}
	productKey, ok := params["product_key"].(int64)
	if !ok {
		return http.StatusBadRequest, errors.New("invalid product_key")
	}
	instrumentKey, ok := params["instrument_key"].(int64)
	if !ok {
		return http.StatusBadRequest, errors.New("invalid instrument_key")
	}

	// Check for duplicate records
	duplicate, key, err := CheckDuplicateAllocInstrument(periodeKey, productKey, instrumentKey)
	log.Println("Error checking for duplicates:", err)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := UpdateAllocInstrument(key, params)
		if err != nil {
			log.Println("Failed to update data:", err)
			return status, err
		}
		return http.StatusOK, nil
	}

	// Jika tidak ada duplikasi, buat data baru
	fields := ""
	placeholders := ""
	var bindvars []interface{}

	for key, value := range params {
		fields += key + ", "
		placeholders += "?, "
		bindvars = append(bindvars, value)
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

func UpdateAllocInstrument(allocInstrumentKey string, params map[string]interface{}) (int, error) {
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
