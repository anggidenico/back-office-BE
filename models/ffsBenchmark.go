package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type Benchmark struct {
	BenchmarkKey       int64   `db:"benchmark_key"        json:"benchmark_key"`
	FundTypeKey        int64   `db:"fund_type_key"        json:"fund_type_key"`
	FundTypeName       string  `db:"fund_type_name" json:"fund_type_name"`
	BenchmarkCode      string  `db:"benchmark_code"  json:"benchmark_code"`
	BenchmarkName      string  `db:"benchmark_name"  json:"benchmark_name"`
	BenchmarkShortName string  `db:"benchmark_short_name" json:"benchmark_short_name"`
	RecAttributeID1    string  `db:"rec_attribute_id1" json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2" json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3" json:"rec_attribute_id3"`
}
type BenchmarkDetail struct {
	BenchmarkKey       int64   `db:"benchmark_key" json:"benchmark_key"`
	FundTypeKey        int64   `db:"fund_type_key"        json:"fund_type_key"`
	FundTypeName       string  `db:"fund_type_name" json:"fund_type_name"`
	BenchmarkCode      string  `db:"benchmark_code"  json:"benchmark_code"`
	BenchmarkName      string  `db:"benchmark_name"  json:"benchmark_name"`
	BenchmarkShortName string  `db:"benchmark_short_name" json:"benchmark_short_name"`
	RecAttributeID1    string  `db:"rec_attribute_id1" json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2" json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3" json:"rec_attribute_id3"`
	RecStatus          int8    `db:"rec_status"  json:"rec_status"`
}

func GetBenchmarkModels(c *[]Benchmark) (int, error) {
	query := `SELECT a.benchmark_key, 
	a.fund_type_key,
	b.fund_type_name, 
	a.benchmark_code, 
	a.benchmark_name, 
	a.benchmark_short_name,
	a.rec_attribute_id1,
	a.rec_attribute_id2,
	a.rec_attribute_id3
	FROM ffs_benchmark a 
	JOIN ms_fund_type b 
	ON a.fund_type_key = b.fund_type_key WHERE a.rec_status = 1` //order by

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error()) // sql.err
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetBenchmarkDetailModels(c *BenchmarkDetail, BenchmarkKey string) (int, error) {
	query := `SELECT 
	a.benchmark_key,
	a.fund_type_key,
	b.fund_type_name, 
	a.benchmark_code, 
	a.benchmark_name, 
	a.benchmark_short_name,
	a.rec_attribute_id1,
	a.rec_attribute_id2,
	a.rec_attribute_id3,
    a.rec_status
	FROM ffs_benchmark a 
	JOIN ms_fund_type b 
	ON a.fund_type_key = b.fund_type_key WHERE a.rec_status = 1 AND a.benchmark_key =` + BenchmarkKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error()) //sql
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteBenchmark(BenchmarkKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_benchmark SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "benchmark_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE benchmark_key = ?`
	values = append(values, BenchmarkKey)

	log.Println("========== DeleteBenchmark ==========>>>", query)

	resultSQL, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	rows, _ := resultSQL.RowsAffected()
	if rows < 1 {
		log.Println("nothing rows affected")
		err2 := fmt.Errorf("nothing rows affected")
		return http.StatusNotFound, err2
	}
	return http.StatusOK, nil
}

func CheckDuplicateFfsBenchmark(benchmarkCode, benchmarkName string) (bool, string, error) { //dari sini
	// Query to check for duplicates
	query := "SELECT benchmark_key FROM ffs_benchmark WHERE benchmark_code = ? OR benchmark_name = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, benchmarkCode, benchmarkName).Scan(&key)

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

func CreateFfsBenchmark(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, key, err := CheckDuplicateFfsBenchmark(params["benchmark_code"], params["benchmark_name"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := UpdateBenchmark(key, params)
		if err != nil {
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

	query := "INSERT INTO ffs_benchmark (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdateBenchmark(benchmarkKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_benchmark SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE benchmark_key = ?`
	values = append(values, benchmarkKey)

	log.Println("========== UpdateBenchmark ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
