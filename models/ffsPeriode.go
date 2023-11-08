package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
)

type FfsPeriode struct {
	PeriodeKey  uint64  `db:"periode_key"          json:"periode_key"`
	PeriodeDate string  `db:"periode_date"         json:"periode_date"`
	PeriodeName string  `db:"periode_name"         json:"periode_name"`
	DateOpened  *string `db:"date_opened"          json:"date_opened"`
	DateClosed  *string `db:"date_closed"          json:"date_closed"`
	RecStatus   uint8   `db:"rec_status"           json:"rec_status"`
}
type FfsPeriodeDetail struct {
	PeriodeKey  uint64  `db:"periode_key"          json:"periode_key"`
	PeriodeDate string  `db:"periode_date"         json:"periode_date"`
	PeriodeName string  `db:"periode_name"         json:"periode_name"`
	DateOpened  *string `db:"date_opened"          json:"date_opened"`
	DateClosed  *string `db:"date_closed"          json:"date_closed"`
	RecStatus   uint8   `db:"rec_status"           json:"rec_status"`
}

func GetFfsPeriodeModels(c *[]FfsPeriode) (int, error) {
	query := `SELECT periode_key,
	periode_date,
	periode_name,
	date_opened,
	date_closed 
	FROM ffs_periode 
	WHERE rec_status = 1`

	log.Println("====================>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func GetFfsPeriodeDetailModels(c *FfsPeriodeDetail, PeriodeKey string) (int, error) {
	query := `SELECT periode_key,
	periode_date,
	periode_name,
	date_opened,
	date_closed FROM ffs_periode
	WHERE rec_status = 1 
	AND periode_key =` + PeriodeKey

	log.Println("====================>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
