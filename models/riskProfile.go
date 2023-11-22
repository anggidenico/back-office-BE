// models.go
package models

import (
	"database/sql"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type RiskProfile struct {
	RiskProfileKey string `json:"risk_profile_key"  db:"risk_profile_key"`
	RiskCode       string `json:"risk_code" db:"risk_code"`
	RiskName       string `json:"risk_name" db:"risk_name"`
	RiskDesc       string `json:"risk_desc" db:"risk_desc"`
	MinScore       int    `json:"min_score" db:"min_score"`
	MaxScore       int    `json:"max_score" db:"max_score"`
	MaxFlag        bool   `json:"max_flag" db:"max_flag"`
	RecOrder       int    `json:"rec_order" db:"rec_order"`
	RecStatus      int    `json:"rec_status" db:"rec_status"`
}

type GetDetailRisk struct {
	RiskProfileKey string `json:"risk_profile_key"  db:"risk_profile_key"`
	RiskCode       string `json:"risk_code" db:"risk_code"`
	RiskName       string `json:"risk_name" db:"risk_name"`
	RiskDesc       string `json:"risk_desc" db:"risk_desc"`
	MinScore       int    `json:"min_score" db:"min_score"`
	MaxScore       int    `json:"max_score" db:"max_score"`
	MaxFlag        bool   `json:"max_flag" db:"max_flag"`
	RecOrder       int    `json:"rec_order" db:"rec_order"`
}

func GetRiskProfileModels(c *[]RiskProfile) (int, error) {
	query := `SELECT risk_profile_key,risk_code,risk_name,risk_desc,min_score,max_score,max_flag,rec_order,rec_status FROM ms_risk_profile
			  WHERE rec_status = 1 order by rec_order`
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
func GetDetailRiskProfileModels(c *GetDetailRisk, RiskProfileKey string) (int, error) {
	query := `SELECT risk_profile_key,
	risk_code,
	risk_name,
	risk_desc,
	min_score,
	max_score,
	max_flag,
	rec_order 
	FROM ms_risk_profile 
	WHERE risk_profile_key =` + RiskProfileKey

	// log.Println("====================>>>", query)
	err := db.Db.Get(c, query)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("PChannelKey not found")
			return http.StatusBadGateway, err
		}
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
func CreateRiskProfile(params map[string]string) (int, error) {
	query := "INSERT INTO ms_risk_profile"
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

func UpdateRiskProfile(RiskProfileKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_risk_profile SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "risk_profile_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE risk_profile_key = ?`
	values = append(values, RiskProfileKey)

	// log.Println("========== UpdateRiskProfile ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func DeleteRiskProfile(RiskProfileKey string, params map[string]string) (int, error) {
	query := `UPDATE ms_risk_profile SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "risk_profile_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE risk_profile_key = ?`
	values = append(values, RiskProfileKey)

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
