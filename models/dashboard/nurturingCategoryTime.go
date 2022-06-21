package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
)

type NurturingCategoryTime struct {
	Id         uint64 `db:"id"                  json:"id"`
	IdCategory uint64 `db:"id_category"         json:"id_category"`
	Time       uint64 `db:"time"                json:"time"`
	Action     string `db:"action"              json:"action"`
	Frequency  uint64 `db:"frequency"           json:"frequency"`
	IsLooping  *uint8 `db:"is_looping"          json:"is_looping"`
	RecStatus  uint8  `db:"rec_status"          json:"rec_status"`
}

func GetNurturingCategoryTime(c *NurturingCategoryTime, key string) (int, error) {
	query := `SELECT 
				nurturing_category_time.* 
			FROM nurturing_category_time 
			WHERE nurturing_category_time.rec_status = 1 
			AND cms_post.post_key = ` + key
	log.Println(query)
	err := db.DbDashboard.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAllNurturingCategoryTime(c *[]NurturingCategoryTime, params map[string]string) (int, error) {
	query := `SELECT
				nurturing_category_time.* 
			FROM nurturing_category_time `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "nurturing_category_time."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	// Main query
	err := db.DbDashboard.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
