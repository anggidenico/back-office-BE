package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
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
func GetBenchmarkProductModels() (result []BenchmarkProduct) {
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
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}

func GetBenchmarkProductDetailModels(BenchProdKey string) (result BenchmarkProdDetail) {
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
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
		// return http.StatusBadGateway, err
	}
	return
}
