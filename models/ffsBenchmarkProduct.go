package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type BenchmarkProduct struct {
	BenchProductKey int64  `db:"bench_prod_key"        json:"bench_prod_key"`
	ProductKey      int64  `db:"product_key" json:"product_key"`
	ProductCode     string `db:"product_code" json:"product_code"`
	ProductName     string `db:"product_name" json:"product_name"`
	BenchmarkRatio  int64  `db:"benchmark_ratio" json:"benchmark_ratio"`
}
type BenchmarkProdDetail struct {
	BenchProductKey int64  `db:"bench_prod_key"  json:"bench_prod_key"`
	ProductKey      int64  `db:"product_key" json:"product_key"`
	ProductCode     string `db:"product_code" json:"product_code"`
	ProductName     string `db:"product_name" json:"product_name"`
	ProductNameAlt  string `db:"product_name_alt" json:"product_name_alt"`
	BenchmarkRatio  int64  `db:"benchmark_ratio" json:"benchmark_ratio"`
	RecStatus       int8   `db:"rec_status" json:"rec_status"`
}

func CreateFfsProductBenchmark(params map[string]string) (int, error) {
	query := "INSERT INTO ffs_benchmark_product"
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
func GetBenchmarkProductModels(c *[]BenchmarkProduct) (int, error) {
	query := `SELECT a.bench_prod_key,
	a.product_key,
	b.product_code,
	b.product_name,
	a.benchmark_ratio 
	FROM ffs_benchmark_product a 
	JOIN ms_product b 
	ON a.product_key = b.product_key 
	WHERE a.rec_status =1`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetBenchmarkProductDetailModels(c *BenchmarkProdDetail, BenchProdKey string) (int, error) {
	query := `SELECT a.bench_prod_key,
	a.product_key,
	b.product_code,
	b.product_name,
	b.product_name_alt,
	a.benchmark_ratio,
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

func UpdateBenchmarkProduct(BenchProdKey string, params map[string]string) (int, error) {
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

	log.Println("========== UpdateFFsBenchmarkProduct ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
