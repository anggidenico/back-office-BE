package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
)

type AllocSecurity struct {
	AllocSecKey  uint64 `db:"alloc_security_key"          json:"alloc_security_key"`
	ProductKey   uint64 `db:"product_key"          json:"product_key"`
	ProductName  string `db:"product_name"          json:"product_name"`
	PeriodeName  string `db:"periode_name"          json:"periode_name"`
	SecurityKey  string `db:"sec_key"          json:"sec_key"`
	SecurityName string `db:"sec_name" json:"sec_name"`
}

func GetAllocSecModels(c *[]AllocSecurity) (int, error) {
	query := `SELECT a.alloc_security_key, 
	a.product_key, 
	b.product_name ,
	a.periode_key, 
	c.periode_name, 
	a.sec_key, d.sec_name 
	FROM ffs_alloc_security a 
	JOIN ms_product b ON a.product_key = b.product_key 
	JOIN ffs_periode c ON a.periode_key = c.periode_key 
	JOIN ms_securities d ON a.sec_key = d.sec_key 
	WHERE rec_status = 1`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
