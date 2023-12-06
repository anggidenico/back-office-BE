package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type PriceList struct {
	PriceKey      int64            `db:"price_key" json:"price_key"`
	BenchmarkKey  int64            `db:"benchmark_key" json:"benchmark_key"`
	BenchmarkName string           `db:"benchmark_name" json:"benchmark_name"`
	PriceType     *int64           `db:"price_type" json:"price_type"`
	PriceName     string           `db:"price_name" json:"price_name"`
	PriceDate     *string          `db:"price_date" json:"price_date"`
	PriceValue    *decimal.Decimal `db:"price_value" json:"price_value"`
	PriceRemarks  *string          `db:"price_remarks" json:"price_remarks"`
	RecOrder      *int64           `db:"rec_order" json:"rec_order"`
}
type PriceType struct {
	PriceKey  int64  `db:"price_key" json:"price_key"`
	PriceType *int64 `db:"price_type" json:"price_type"`
	PriceName string `db:"price_name" json:"price_name"`
}

// func GetPriceListModels(c *[]PriceList) (int, error) {
// 	query := `SELECT a.price_key,
// 	a.benchmark_key,
// 	b.benchmark_name,
//     a.price_type,
// 	c.lkp_name price_name,
//     a.price_date,
//     a.price_value,
//     a.price_remarks,
// 	a.rec_order
// 	FROM ffs_price a
// 	JOIN ffs_benchmark b
// 	ON a.benchmark_key = b.benchmark_key
// 	JOIN gen_lookup c
// 	ON a.price_type = c.lookup_key
// 	WHERE a.rec_status =1
// 	ORDER BY a.rec_created_date DESC`
// 	// log.Println("====================>>>", query)
// 	err := db.Db.Select(c, query)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println(err.Error())
// 			return http.StatusBadGateway, err
// 		}
// 		return http.StatusNotFound, err
// 	}
// 	return http.StatusOK, nil
// }

// GetPriceListModels mengambil semua data harga dari database
func GetPriceListModels(priceLists *[]PriceList) (int, error) {
	query := `
		SELECT
			a.price_key,
			a.benchmark_key,
			b.benchmark_name,
			a.price_type,
			c.lkp_name price_name,
			a.price_date,
			a.price_value,
			a.price_remarks,
			a.rec_order
		FROM ffs_price a
		JOIN ffs_benchmark b ON a.benchmark_key = b.benchmark_key
		JOIN gen_lookup c ON a.price_type = c.lookup_key
		WHERE a.rec_status = 1
		ORDER BY a.rec_created_date DESC
	`

	if err := db.Db.Select(priceLists, query); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetPriceDetailModels(c *PriceList, PriceKey string) (int, error) {
	query := `SELECT a.price_key,
	a.benchmark_key,
	b.benchmark_name,
    a.price_type,
	c.lkp_name price_name,
    a.price_date,
    a.price_value,
    a.price_remarks,
	a.rec_order
	FROM ffs_price a
	JOIN ffs_benchmark b 
	ON a.benchmark_key = b.benchmark_key
	JOIN gen_lookup c
	ON a.price_type = c.lookup_key WHERE a.rec_status = 1 AND a.price_key =` + PriceKey

	// log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("price_key not found")
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func CheckDuplicateFfsPrice(benchkmarkKey, priceType, priceDate string) (bool, string, error) { //dari sini
	// Query to check for duplicates
	query := "SELECT benchmark_key FROM ffs_price WHERE benchmark_key = ? OR price_type = ? OR price_date = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, benchkmarkKey, priceType, priceDate).Scan(&key)

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

func CreateFfsPrice(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, key, err := CheckDuplicateFfsPrice(params["benchmark_key"], params["price_type"], params["price_date"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := UpdatePrice(key, params)
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

	query := "INSERT INTO ffs_price (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdatePrice(priceKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_price SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE price_key = ?`
	values = append(values, priceKey)

	log.Println("========== UpdatePrice ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeletePriceModels(PriceKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_price SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "price_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE price_key = ?`
	values = append(values, PriceKey)

	// log.Println("========== UpdateRiskProfile ==========>>>", query)

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

// FilterByBenchmarkAndDateModels melakukan filter pada PriceList berdasarkan kunci benchmark dan rentang tanggal
func FilterByBenchmarkAndDateModels(benchmarkKey int64, startDate, endDate time.Time, priceLists []PriceList) []PriceList {
	var filteredPriceLists []PriceList

	for _, priceList := range priceLists {
		// Filter by benchmark key
		if priceList.BenchmarkKey != benchmarkKey {
			continue
		}

		// Parse price date string to time.Time
		priceDate, err := time.Parse("2006-01-02", *priceList.PriceDate)
		if err != nil {
			// Handle parsing error
			continue
		}

		// Filter by date range
		if priceDate.After(startDate) && priceDate.Before(endDate) {
			filteredPriceLists = append(filteredPriceLists, priceList)
		}
	}

	return filteredPriceLists
}

func GetPriceListFilterModels(priceLists *[]PriceList, startDate string, endDate string, benchmarkKey string) (int, error) {
	query := `
		SELECT
			a.price_key,
			a.benchmark_key,
			b.benchmark_name,
			a.price_type,
			c.lkp_name price_name,
			DATE_FORMAT(a.price_date,'%d %M %Y') as price_date,
			a.price_value,
			a.price_remarks,
			a.rec_order
		FROM ffs_price a
		JOIN ffs_benchmark b ON a.benchmark_key = b.benchmark_key
		JOIN gen_lookup c ON a.price_type = c.lookup_key
		WHERE a.rec_status = 1
		AND CAST(a.price_date AS DATE) >= ? 
		AND CAST(a.price_date AS DATE) <= ?`

	// Jika benchmarkKey tidak kosong, tambahkan kondisi ke query
	if benchmarkKey != "" {
		query += " AND a.benchmark_key=?"
	}

	query += " ORDER BY a.price_key DESC"

	var err error
	if benchmarkKey != "" {
		err = db.Db.Select(priceLists, query, startDate, endDate, benchmarkKey)
	} else {
		err = db.Db.Select(priceLists, query, startDate, endDate)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
