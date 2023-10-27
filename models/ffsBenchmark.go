package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type Benchmark struct {
	BenchmarkKey       int64  `db:"benchmark_key"        json:"benchmark_key"`
	FundTypeKey        int64  `db:"fund_type_key"        json:"fund_type_key"`
	FundTypeName       string `db:"fund_type_name" json:"fund_type_name"`
	BenchmarkCode      string `db:"benchmark_code"  json:"benchmark_code"`
	BenchmarkName      string `db:"benchmark_name"  json:"benchmark_name"`
	BenchmarkShortName string `db:"benchmark_short_name" json:"benchmark_short_name"`
}
type BenchmarkDetail struct {
	FundTypeKey        int64  `db:"fund_type_key"        json:"fund_type_key"`
	FundTypeName       string `db:"fund_type_name" json:"fund_type_name"`
	BenchmarkCode      string `db:"benchmark_code"  json:"benchmark_code"`
	BenchmarkName      string `db:"benchmark_name"  json:"benchmark_name"`
	BenchmarkShortName string `db:"benchmark_short_name" json:"benchmark_short_name"`
	RecAttributeID1    string `db:"rec_attribute_id1" json:"rec_attribute_id1"`
	RecStatus          int8   `db:"rec_status"  json:"rec_status"`
}

func GetBenchmarkModels() (result []Benchmark) {
	query := `SELECT a.benchmark_key, 
	a.fund_type_key,
	b.fund_type_name, 
	a.benchmark_code, 
	a.benchmark_name, 
	a.benchmark_short_name 
	FROM ffs_benchmark a 
	JOIN ms_fund_type b 
	ON a.fund_type_key = b.fund_type_key WHERE a.rec_status = 1`

	log.Println("====================>>>", query)
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}

func GetBenchmarkDetailModels(BenchmarkKey string) (result BenchmarkDetail) {
	query := `SELECT a.fund_type_key,
	b.fund_type_name, 
	a.benchmark_code, 
	a.benchmark_name, 
	a.benchmark_short_name,
	a.rec_attribute_id1,
    a.rec_status
	FROM ffs_benchmark a 
	JOIN ms_fund_type b 
	ON a.fund_type_key = b.fund_type_key WHERE a.rec_status = 1 AND a.benchmark_key =` + BenchmarkKey

	log.Println("====================>>>", query)
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}

func DeleteBenchmark(BenchmarkKey string, params map[string]string) error {
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

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func CreateFfsBenchmark(params map[string]string) (int, error) {
	query := "INSERT INTO ffs_benchmark"
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

func UpdateFfsBenchmark(BenchmarkKey string, params map[string]string) error {
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

	log.Println("========== UpdateRiskProfile ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
