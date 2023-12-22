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

type BenchmarkProduct struct {
	BenchProductKey int64           `db:"bench_prod_key" json:"bench_prod_key"`
	BenchmarkKey    int64           `db:"benchmark_key" json:"benchmark_key"`
	BenchmarkName   string          `db:"benchmark_name" json:"benchmark_name"`
	ProductKey      int64           `db:"product_key" json:"product_key"`
	ProductCode     string          `db:"product_code" json:"product_code"`
	ProductNameAlt  string          `db:"product_name_alt" json:"product_name_alt"`
	ProductName     string          `db:"product_name" json:"product_name"`
	BenchmarkRatio  decimal.Decimal `db:"benchmark_ratio" json:"benchmark_ratio"`
	RecStatus       int64           `db:"rec_status" json:"rec_status"`
}
type BenchmarkProdDetail struct {
	BenchProductKey  int64           `db:"bench_prod_key"  json:"bench_prod_key"`
	BenchmarkKey     int64           `db:"benchmark_key" json:"benchmark_key"`
	ProductKey       int64           `db:"product_key" json:"product_key"`
	ProductCode      string          `db:"product_code" json:"product_code"`
	ProductName      string          `db:"product_name" json:"product_name"`
	ProductNameAlt   string          `db:"product_name_alt" json:"product_name_alt"`
	BenchmarkRatio   decimal.Decimal `db:"benchmark_ratio" json:"benchmark_ratio"`
	BenchmarkRemarks *string         `db:"benchmark_remarks" json:"benchmark_remarks"`
	RecStatus        int8            `db:"rec_status" json:"rec_status"`
}

func GetBenchmarkProductModels(c *[]BenchmarkProduct) (int, error) {
	query := `SELECT a.bench_prod_key,
	a.benchmark_key,
	c.benchmark_name,
	a.product_key,
	b.product_code,
	b.product_name,
	b.product_name_alt,
	a.benchmark_ratio,
	a.rec_status
	FROM ffs_benchmark_product a 
	JOIN ms_product b 
	ON a.product_key = b.product_key 
	JOIN ffs_benchmark c
	ON a.benchmark_key = c.benchmark_key
	WHERE a.rec_status = 1 
	ORDER BY a.rec_order ASC`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
	}
	return http.StatusOK, nil
}
func GetBenchmarkProductDetailModels(c *BenchmarkProdDetail, BenchProdKey string) (int, error) {
	query := `SELECT a.bench_prod_key,
	a.benchmark_key,
	a.product_key,
	b.product_code,
	b.product_name,
	b.product_name_alt,
	a.benchmark_ratio,
	a.benchmark_remarks,
	a.rec_status
	FROM ffs_benchmark_product a 
	JOIN ms_product b 
	ON a.product_key = b.product_key 
	WHERE a.rec_status = 1 
	AND a.bench_prod_key =` + BenchProdKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func DeleteBenchmarkProduct(BenchProdKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_benchmark_product SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "bench_prod_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE bench_prod_key = ?`
	values = append(values, BenchProdKey)

	log.Println("========== DeleteBenchmarkProduct ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckDuplicateBenchmarkProd(productKey, benchmarkKey string) (bool, string, error) { //dari sini
	// Query to check for duplicates
	query := "SELECT bench_prod_key FROM ffs_benchmark_product WHERE product_key = ? OR benchmark_key = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, benchmarkKey, productKey).Scan(&key)

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

func CreateBenchmarkProd(params map[string]interface{}) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateBenchmarkProd(params["product_key"].(string), params["benchmark_key"].(string))
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

	query := "INSERT INTO ffs_benchmark_product (" + fields + ") VALUES (" + placeholders + ")"

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
func UpdateBenchmarkProd(benchProdKey string, params map[string]interface{}) (int, error) {
	query := `UPDATE ffs_benchmark_product SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE bench_prod_key = ?`
	values = append(values, benchProdKey)

	// log.Println("========== UpdateBenchmarkProduct ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetBenchmarkProdStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM ffs_benchmark_product WHERE bench_prod_key = ?"
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
