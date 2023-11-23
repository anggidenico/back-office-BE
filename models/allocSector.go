package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
)

type AllocSector struct {
	AllocSectorKey uint8  `db:"alloc_sector_key" json:"alloc_sector_key"`
	ProductKey     uint64 `db:"product_key" json:"product_key"`
	PeriodeKey     uint64 `db:"periode_key" json:"periode_key"`
	SectorKey      uint64 `db:"sector_key" json:"sector_key"`
	SectorValue    uint64 `db:"sector_value" json:"sector_value"`
}

func GetAllocSectorModels(c *[]AllocSector) (int, error) {
	query := `SELECT a.alloc_sector_key, 
	a.product_key, 
	a.periode_key, 
	a.sector_key, 
	a.sector_value, 
	a.rec_order 
	FROM ffs_alloc_sector a 
	JOIN ffs_periode b ON a.periode_key = b.periode_key 
	JOIN ms_product c ON a.product_key = c.product_key 
	JOIN ms_securities_sector d ON a.sector_key = d.sector_key 
	WHERE rec STATUS =1 
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
