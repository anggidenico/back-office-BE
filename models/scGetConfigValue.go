package models

import (
	"mf-bo-api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScAppConfig1 struct {
	AppConfigValue *string `db:"app_config_value"           json:"app_config_value"`
}

func GetScAppConfigValueByCode(c *ScAppConfig1, code string) (int, error) {
	query := `SELECT sc_app_config.app_config_value 
	FROM sc_app_config 
	WHERE sc_app_config.rec_status = 1 
	AND sc_app_config.app_config_code ='` + code + `' `
	// log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
