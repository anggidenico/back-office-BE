package models

import (
	"database/sql"
	"log"
	"mf-bo-api/db"
	"net/http"
)

type NurturingNotifSummary struct {
	Id              uint64  `db:"id"                   json:"id"`
	UserLogin       uint64  `db:"user_login"           json:"user_login"`
	TokenNotif      *string `db:"token_notif"          json:"token_notif"`
	FullName        *string `db:"full_name"            json:"full_name"`
	Email           *string `db:"email"                json:"email"`
	IdCategory      uint8   `db:"id_category"          json:"id_category"`
	IdCategoryTime  uint8   `db:"id_category_time"     json:"id_category_time"`
	IdForceType     uint8   `db:"id_force_type"        json:"id_force_type"`
	IdMessage       uint8   `db:"id_message"           json:"id_message"`
	MessageTitle    *string `db:"message_title"        json:"message_title"`
	MessageBody     *string `db:"message_body"         json:"message_body"`
	PushNotifTime   *string `db:"push_notif_time"      json:"push_notif_time"`
	PushNotifStatus *uint8  `db:"push_notif_status"    json:"push_notif_status"`
	PushNotifError  *string `db:"push_notif_error"     json:"push_notif_error"`
	RecStatus       *uint8  `db:"rec_status"           json:"rec_status"`
	RecCreatedDate  *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy    *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy   *string `db:"rec_modified_by"      json:"rec_modified_by"`
}

func CreateNurturingNotifSummary(params map[string]string) (int, error) {
	query := "INSERT INTO nurturing_notif_summary"
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
	log.Println("==========  ==========>>>", query)

	tx, err := db.DbDashboard.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateCmsPost(params map[string]string) (int, error) {
	query := "UPDATE nurturing_notif_summary SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "id" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE id = " + params["id"]
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
