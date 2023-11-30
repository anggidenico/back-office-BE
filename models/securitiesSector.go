package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
)

type SecuritiesSector struct {
	SectorKey         int64   `db:"sector_key" json:"sector_key"`
	SectorCode        *string `db:"sector_code" json:"sector_code"`
	SectorName        *string `db:"sector_name" json:"sector_name"`
	SectorDescription *string `db:"sector_description" json:"sector_description"`
	RecOrder          int64   `db:"rec_order" json:"rec_order"`
}

func GetSecuritiesSectorModels(c *[]SecuritiesSector) (int, error) {
	query := `SELECT sector_key,
	sector_code,
	sector_name,
	sector_description,
	rec_order FROM ms_securities_sector
	WHERE rec_status =1 
	ORDER BY sector_key DESC`
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

func GetSecuritiesSectorDetailModels(c *SecuritiesSector, SectorKey string) (int, error) {
	query := `SELECT sector_key,
	sector_code,
	sector_name,
	sector_description,
	rec_order FROM ms_securities_sector 
	WHERE rec_status = 1 AND sector_key =` + SectorKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("sector_key not found")
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
