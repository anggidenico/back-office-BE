package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
)

type UserRegistrasiBaruBelumOA struct {
	UserLoginKey   uint64  `db:"user_login_key"       json:"user_login_key"`
	TokenNotif     *string `db:"token_notif"          json:"token_notif"`
	FullName       string  `db:"full_name"            json:"full_name"`
	UloginEmail    string  `db:"ulogin_email"         json:"ulogin_email"`
	IdCategory     uint64  `db:"id_category"          json:"id_category"`
	IdCategoryTime uint64  `db:"id_category_time"     json:"id_category_time"`
	IdForceType    uint64  `db:"id_force_type"        json:"id_force_type"`
	IdMessage      uint64  `db:"id_message"           json:"id_message"`
	MessageTitle   string  `db:"message_title"        json:"message_title"`
	MessageBody    string  `db:"message_body"         json:"message_body"`
	RecCreatedDate *string `db:"rec_created_date"    json:"rec_created_date,omitempty"`
}

func GetUserRegistrasiBaruBelumOA(c *[]UserRegistrasiBaruBelumOA, idCategory string, catTime string) (int, error) {
	query := `SELECT 
				u.user_login_key,
				u.token_notif,
				(CASE 
				WHEN op.full_name IS NOT NULL THEN op.full_name 
				ELSE u.ulogin_full_name 
				END) AS full_name,
				u.ulogin_email,
				nc.id AS id_category,
				nct.id AS id_category_time,
				nft.id AS id_force_type,
				nm.id AS id_message,
				nm.message_title,
				nm.message_body 
			FROM mam_core.sc_user_login AS u 
			LEFT JOIN mam_core.oa_request AS oa ON oa.user_login_key = u.user_login_key 
			LEFT JOIN mam_core.oa_personal_data AS op ON op.oa_request_key = oa.oa_request_key AND op.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_category AS nc ON 1 = 1 AND nc.id = ` + idCategory + ` 
			LEFT JOIN mam_dashboard.nurturing_category_time AS nct ON nct.id_category = nc.id AND nct.time = ` + catTime + ` AND nct.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_force_type AS nft ON 1 = 1 AND nft.id = 1 AND nft.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_message AS nm ON nm.id_force_type = nft.id AND nm.id_category_time = nct.id
			WHERE (u.role_key IS NULL OR u.role_key IN (1)) 
			AND u.user_category_key = 1 AND u.user_dept_key = 1 AND u.rec_status = 1 
			AND (oa.oa_request_key IS NULL OR (oa.oa_request_type = 127 AND oa.oa_status = 444 AND oa.customer_key IS NULL)) 
			AND DATE_FORMAT(DATE_ADD(u.rec_created_date, INTERVAL ` + catTime + ` DAY), '%Y-%m-%d') = DATE_FORMAT(NOW(),'%Y-%m-%d')`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetUserSudahCustomerBelumTransaksi(c *[]UserRegistrasiBaruBelumOA, idCategory string, catTime string) (int, error) {
	query := `SELECT 
				u.user_login_key,
				u.token_notif,
				c.full_name AS full_name,
				u.ulogin_email,
				nc.id AS id_category,
				nct.id AS id_category_time,
				nft.id AS id_force_type,
				nm.id AS id_message,
				nm.message_title,
				nm.message_body 
			FROM mam_core.ms_customer AS c 
			INNER JOIN mam_core.sc_user_login AS u ON u.customer_key = c.customer_key AND u.rec_status = 1 
			LEFT JOIN mam_core.tr_transaction AS t ON t.customer_key = c.customer_key AND t.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_category AS nc ON 1 = 1 AND nc.id = ` + idCategory + ` 
			LEFT JOIN mam_dashboard.nurturing_category_time AS nct ON nct.id_category = nc.id AND nct.time = ` + catTime + ` AND nct.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_force_type AS nft ON 1 = 1 AND nft.id = 1 AND nft.rec_status = 1 
			LEFT JOIN mam_dashboard.nurturing_message AS nm ON nm.id_force_type = nft.id AND nm.id_category_time = nct.id
			WHERE u.role_key IN (1) AND c.investor_type = 263 AND c.rec_status = 1 
			AND u.user_category_key = 1 AND u.user_dept_key = 1 AND u.rec_status = 1 
			AND t.transaction_key IS NULL 
			AND DATE_FORMAT(DATE_ADD(c.rec_created_date, INTERVAL ` + catTime + ` DAY), '%Y-%m-%d') = DATE_FORMAT(NOW(),'%Y-%m-%d') 
			GROUP BY c.customer_key`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetUserHanyaSubs1kali(c *[]UserRegistrasiBaruBelumOA, idCategory string, catTime string, firstCatTime string, forceType string) (int, error) {
	query := `SELECT dt.* FROM 
				(SELECT 
					t.rec_created_date,
					u.user_login_key,
					u.token_notif,
					c.full_name AS full_name,
					u.ulogin_email,
					nc.id AS id_category,
					nct.id AS id_category_time,
					nft.id AS id_force_type,
					nm.id AS id_message,
					nm.message_title,
					nm.message_body 
				FROM mam_core.ms_customer AS c 
				INNER JOIN mam_core.sc_user_login AS u ON u.customer_key = c.customer_key AND u.rec_status = 1 
				INNER JOIN mam_core.tr_transaction AS t ON t.customer_key = c.customer_key AND t.rec_status = 1 
				LEFT JOIN mam_dashboard.nurturing_category AS nc ON 1 = 1 AND nc.id = ` + idCategory + ` 
				LEFT JOIN mam_dashboard.nurturing_category_time AS nct ON nct.id_category = nc.id AND nct.time = ` + firstCatTime + ` AND nct.rec_status = 1 
				LEFT JOIN mam_dashboard.nurturing_force_type AS nft ON 1 = 1 AND nft.id = ` + forceType + ` AND nft.rec_status = 1 
				LEFT JOIN mam_dashboard.nurturing_message AS nm ON nm.id_force_type = nft.id AND nm.id_category_time = nct.id
				WHERE u.role_key IN (1) AND c.investor_type = 263 AND c.rec_status = 1 
				AND u.user_category_key = 1 AND u.user_dept_key = 1 AND u.rec_status = 1 AND t.rec_status = 1 AND t.trans_type_key = 1 
				GROUP BY c.customer_key HAVING COUNT(t.transaction_key) = 1) AS dt 
				WHERE DATE_FORMAT(DATE_ADD(dt.rec_created_date, INTERVAL ` + catTime + ` DAY), '%Y-%m-%d') = DATE_FORMAT(NOW(),'%Y-%m-%d')`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
