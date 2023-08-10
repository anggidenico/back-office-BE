package models

import (
	"log"
	"mf-bo-api/db"
)

type InfoCustomerModels struct {
	UserLoginKey *uint64 `db:"user_login_key" json:"user_login_key"`
	CustomerKey  *uint64 `db:"customer_key" json:"customer_key"`
	Email        *string `db:"email" json:"email"`
	Phone        *string `db:"phone" json:"phone"`
	Cif          *string `db:"cif" json:"cif"`
	Sid          *string `db:"sid" json:"sid"`
}

func GetCustomerInfobyCustomerKey(CustomerKey string) interface{} {
	query := `SELECT t1.user_login_key, t2.customer_key, t1.ulogin_email AS email, 
	t1.ulogin_mobileno AS phone, t2.unit_holder_idno AS cif, t2.sid_no 
	FROM sc_user_login t1
	INNER JOIN ms_customer t2 ON t1.customer_key = t2.customer_key
	WHERE t2.customer_key = ` + CustomerKey
	var result InfoCustomerModels
	err := db.Db.Get(&result, query)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}
