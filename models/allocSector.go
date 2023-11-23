package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type AllocSector struct {
	AllocSectorKey uint8  `db:"alloc_sector_key" json:"alloc_sector_key"`
	ProductKey     uint64 `db:"product_key" json:"product_key"`
	ProductName    string `db:"product_name" json:"product_name"`
	PeriodeKey     uint64 `db:"periode_key" json:"periode_key"`
	PeriodeName    string `db:"periode_name" json:"periode_name"`
	SectorKey      uint64 `db:"sector_key" json:"sector_key"`
	SectorName     string `db:"sector_name" json:"sector_name"`
	SectorCode     string `db:"sector_code" json:"sector_code"`
	SectorValue    uint64 `db:"sector_value" json:"sector_value"`
}
type AllocSectorDetail struct {
	AllocSectorKey uint8  `db:"alloc_sector_key" json:"alloc_sector_key"`
	ProductKey     uint64 `db:"product_key" json:"product_key"`
	ProductName    string `db:"product_name" json:"product_name"`
	PeriodeKey     uint64 `db:"periode_key" json:"periode_key"`
	PeriodeName    string `db:"periode_name" json:"periode_name"`
	SectorKey      uint64 `db:"sector_key" json:"sector_key"`
	SectorName     string `db:"sector_name" json:"sector_name"`
	SectorCode     string `db:"sector_code" json:"sector_code"`
	SectorValue    uint64 `db:"sector_value" json:"sector_value"`
	RecOrder       uint64 `db:"rec_order" json:"rec_order"`
}

func GetAllocSectorModels(c *[]AllocSector) (int, error) {
	query := `SELECT a.alloc_sector_key, 
	a.product_key, 
	c.product_name,
	a.periode_key, 
	b.periode_name,
	a.sector_key, 
	d.sector_name,
	d.sector_code,
	a.sector_value, 
	a.rec_order 
	FROM ffs_alloc_sector a 
	JOIN ffs_periode b ON a.periode_key = b.periode_key 
	JOIN ms_product c ON a.product_key = c.product_key 
	JOIN ms_securities_sector d ON a.sector_key = d.sector_key 
	WHERE a.rec_status =1 
	ORDER BY rec_order`
	// log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
	}
	return http.StatusOK, nil
}

func GetAllocSectorDetailModels(c *AllocSectorDetail, AllocSectorKey string) (int, error) {
	query := `SELECT a.alloc_sector_key, 
	a.product_key, 
	c.product_name,
	a.periode_key, 
	b.periode_name,
	a.sector_key, 
	d.sector_name,
	d.sector_code,
	a.sector_value, 
	a.rec_order 
	FROM ffs_alloc_sector a 
	JOIN ffs_periode b ON a.periode_key = b.periode_key 
	JOIN ms_product c ON a.product_key = c.product_key 
	JOIN ms_securities_sector d ON a.sector_key = d.sector_key WHERE a.rec_status = 1 AND a.alloc_sector_key =` + AllocSectorKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Periode key not found")
			return http.StatusBadGateway, err
		}
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func DeleteAllocSector(AllocSectorKey string, params map[string]string) (int, error) {
	query := `UPDATE ffs_alloc_sector SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "alloc_sector_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE alloc_sector_key = ?`
	values = append(values, AllocSectorKey)

	log.Println("========== Delete Portfolio Sector ==========>>>", query)

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

func CheckDuplicateAllocSector(periodeKey, productKey int64, sectorKey int64) (bool, string, error) { //dari sini
	// Query to check for duplicates
	query := "SELECT alloc_sector_key FROM ffs_alloc_sector WHERE product_key = ? OR periode_key = ? OR sector_key = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, periodeKey, productKey, sectorKey).Scan(&key)

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

func CreateAllocSector(params map[string]interface{}) (int, error) {
	periodeKey, _ := params["periode_key"].(int64)
	productKey, _ := params["product_key"].(int64)
	sectorKey, _ := params["sector_key"].(int64)
	// Check for duplicate records
	duplicate, key, err := CheckDuplicateAllocSector(periodeKey, productKey, sectorKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		status, err := UpdateAllocSector(key, params)
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

	query := "INSERT INTO ffs_alloc_sector (" + fields + ") VALUES (" + placeholders + ")"

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

func UpdateAllocSector(allocSectorKey string, params map[string]interface{}) (int, error) {
	query := `UPDATE ffs_alloc_sector SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE alloc_sector_key = ?`
	values = append(values, allocSectorKey)

	log.Println("========== UpdatePortfolioSector ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
