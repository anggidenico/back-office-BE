package models

import (
	"database/sql"
	"errors"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type ProductInfoValue struct {
	ProductInfoValueKey int64   `db:"product_info_value_key" json:"product_info_value_key"`
	ProductKey          int64   `db:"product_key" json:"product_key"`
	ProductName         string  `db:"product_name" json:"product_name"`
	ProductInfoKey      int64   `db:"product_info_key" json:"product_info_key"`
	ProductInfoName     string  `db:"product_info_name" json:"product_info_name"`
	RecValue            *string `db:"rec_value" json:"rec_value"`
	RecOrder            *uint64 `db:"rec_order" json:"rec_order"`
}
type ProductInfoKey struct {
	ProductInfoKey int64   `db:"product_info_key" json:"product_info_key"`
	RecKey         *string `db:"rec_key" json:"rec_key"`
	RecCaption     *string `db:"rec_caption" json:"rec_caption"`
	RecDesc        *string `db:"rec_desc" json:"rec_desc"`
}

func GetProductInfoKeyModels(c *[]ProductInfoKey) (int, error) {
	query := `SELECT product_info_key,
	rec_key,
	rec_caption,
	rec_desc
	FROM ffs_product_info   
WHERE rec_status = 1
ORDER BY product_info_key ASC`

	log.Println(query)

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
func GetProductInfoValueModels(c *[]ProductInfoValue) (int, error) {
	query := `SELECT
	a.product_info_value_key,
	a.product_key,
	b.product_name,
	a.product_info_key,
	c.rec_caption product_info_name,
	a.rec_value,
	a.rec_order
	FROM ffs_product_info_value a
	JOIN ms_product b ON a.product_key = b.product_key
	JOIN ffs_product_info c ON a.product_info_key = c.product_info_key   
WHERE a.rec_status = 1
ORDER BY a.product_info_value_key DESC`

	log.Println(query)

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

func GetProductInfoValueDetailModels(c *ProductInfoValue, ProdInfoValKey string) (int, error) {
	query := `SELECT
	a.product_info_value_key,
	a.product_key,
	b.product_name,
	a.product_info_key,
	c.rec_caption product_info_name,
	a.rec_value,
	a.rec_order
	FROM ffs_product_info_value a
	JOIN ms_product b ON a.product_key = b.product_key
	JOIN ffs_product_info c ON a.product_info_key = c.product_info_key   
WHERE a.rec_status = 1
	AND a.product_info_value_key =` + ProdInfoValKey

	// log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func CheckDuplicateProductInfoValue(ProductKey, ProductInfoKey string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT product_info_value_key FROM ffs_product_info_value WHERE product_key = ? AND product_info_key = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, ProductKey, ProductInfoKey).Scan(&key)

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

func GetProductInfoValueStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM ffs_product_info_value WHERE product_info_value_key = ?"
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

func CreateProductInfoValue(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateProductInfoValue(params["product_key"], params["product_info_key"])
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

	query := "INSERT INTO ffs_product_info_value (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdateProductInfoValue(ProdInfoValKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_product_info_value SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "product_info_value_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE product_info_value_key = ?`
	values = append(values, ProdInfoValKey)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteProductInfoValue(ProdInfoValKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_product_info_value SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "product_info_value_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE product_info_value_key = ?`
	values = append(values, ProdInfoValKey)

	// log.Println("========== DeleteProductInfoValue ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
