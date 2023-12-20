package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mf-bo-api/db"
	"net/http"
	"strings"
)

type InvestPartner struct {
	InvestPartnerKey    int64   `db:"invest_partner_key" json:"invest_partner_key"`
	InvestPurposeKey    int64   `db:"invest_purpose_key" json:"invest_purpose_key"`
	PartnerCode         string  `db:"partner_code" json:"partner_code"`
	PurposeCode         *string `db:"purpose_code" json:"purpose_code"`
	PurposeName         *string `db:"purpose_name" json:"purpose_name"`
	PartnerBusinessName *string `db:"partner_business_name" json:"partner_business_name"`
	PartnerDesc         *string `db:"partner_desc" json:"partner_desc"`
	PartnerPicname      *string `db:"partner_picname" json:"partner_picname"`
	PartnerMobileNo     *string `db:"partner_mobileno" json:"partner_mobileno"`
	PartnerOfficeNo     *string `db:"partner_officeno" json:"partner_officeno"`
	PartnerEmail        *string `db:"partner_email" json:"partner_email"`
	PartnerCity         *string `db:"partner_city" json:"partner_city"`
	PartnerAddrress     *string `db:"partner_address" json:"partner_address"`
	PartnerUrl          *string `db:"partner_url" json:"partner_url"`
	PartnerDateStarted  *string `db:"partner_date_started" json:"partner_date_started"`
	PartnerDateExpired  *string `db:"partner_date_expired" json:"partner_date_expired"`
	PartnerBannerHits   *int64  `db:"partner_banner_hits" json:"partner_banner_hits"`
	RecImage1           *string `db:"rec_image1" json:"rec_image1"`
	RecOrder            *int64  `db:"rec_order" json:"recorder_order"`
}
type InvestPurpose struct {
	InvestPurposeKey int64   `db:"invest_purpose_key" json:"invest_purpose_key"`
	PurposeCode      *string `db:"purpose_code" json:"purpose_code"`
	PurposeName      *string `db:"purpose_name" json:"purpose_name"`
}

func GetInvestPurposeModels(c *[]InvestPurpose) (int, error) {
	query := `SELECT 
	invest_purpose_key,
	purpose_code,
	purpose_name   
	FROM cms_invest_purpose 
	WHERE rec_status =1 
	ORDER BY invest_purpose_key DESC`
	// log.Println("====================>>>", query)
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

func GetInvestPartnerModels(c *[]InvestPartner) (int, error) {
	query := `SELECT a.invest_partner_key,
	a.invest_purpose_key,
	b.purpose_code,
	b.purpose_name,
	a.partner_code,
	a.partner_business_name,
	a.partner_desc,
	a.partner_picname,
	a.partner_mobileno,
	a.partner_officeno,
	a.partner_email,
	a.partner_city,
	a.partner_address,
	a.partner_url,
	a.partner_date_started,
	a.partner_date_expired,
	a.partner_banner_hits,
	a.rec_image1,
	a.rec_order   
	FROM cms_invest_partner a
	JOIN cms_invest_purpose b 
	ON a.invest_purpose_key = b.invest_purpose_key
	WHERE a.rec_status =1 
	ORDER BY a.invest_partner_key DESC`
	// log.Println("====================>>>", query)
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

func GetInvestPartnerDetailModels(c *InvestPartner, InvestPartnerKey string) (int, error) {
	query := `SELECT a.invest_partner_key,
	a.invest_purpose_key,
	b.purpose_code,
	b.purpose_name,
	a.partner_code,
	a.partner_business_name,
	a.partner_desc,
	a.partner_picname,
	a.partner_mobileno,
	a.partner_officeno,
	a.partner_email,
	a.partner_city,
	a.partner_address,
	a.partner_url,
	a.partner_date_started,
	a.partner_date_expired,
	a.partner_banner_hits,
	a.rec_image1,
	a.rec_order   
	FROM cms_invest_partner a
	JOIN cms_invest_purpose b 
	ON a.invest_purpose_key = b.invest_purpose_key
	WHERE a.rec_status = 1 AND a.invest_partner_key =` + InvestPartnerKey
	err := db.Db.Get(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func DeleteInvestPartnerModels(InvestPartnerKey string, params map[string]string) (int, error) {
	query := `UPDATE cms_invest_partner SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		if key != "invest_partner_key" {
			setClauses = append(setClauses, key+" = ?")
			values = append(values, value)
		}
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE invest_partner_key = ?`
	values = append(values, InvestPartnerKey)

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
func CheckDuplicateInvestPartner(PartnerCode, PartnerBusName string) (bool, string, error) {
	// Query to check for duplicates
	query := "SELECT invest_partner_key FROM cms_invest_partner WHERE partner_code = ? AND partner_business_name = ? LIMIT 1"
	var key string
	err := db.Db.QueryRow(query, PartnerCode, PartnerBusName).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			// No duplicate found
			return false, "", nil
		}
		// Other error occurred
		return false, "", err
	}

	// Duplicate found
	return true, key, nil
}

func CreateInvestPartner(params map[string]string) (int, error) {
	// Check for duplicate records
	duplicate, _, err := CheckDuplicateInvestPartner(params["partner_code"], params["partner_business_name"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Jika duplikasi ditemukan, perbarui data yang sudah ada
	if duplicate {
		return http.StatusBadRequest, errors.New("data duplikat ditemukan")
	}

	// Jika tidak ada duplikasi, buat data baru
	fields := ""
	placeholders := ""
	var bindvars []interface{}

	for key, value := range params {
		fields += key + `, `
		if value == "" {
			placeholders += `NULL, `
		} else {
			placeholders += `?, `
			bindvars = append(bindvars, value)
		}
	}

	fields = fields[:len(fields)-2]
	placeholders = placeholders[:len(placeholders)-2]

	query := "INSERT INTO cms_invest_partner (" + fields + ") VALUES (" + placeholders + ")"

	tx, err := db.Db.Begin()
	if err != nil {
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	if err != nil {
		tx.Rollback()
		return http.StatusBadRequest, err
	}

	tx.Commit()

	return http.StatusOK, nil
}
func GetInvestPartnerStatusByKey(key string) (int, error) {
	query := "SELECT rec_status FROM cms_invest_partner WHERE invest_partner_key = ?"
	var status int
	err := db.Db.QueryRow(query, key).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Data tidak ditemukan
			return 0, nil
		}
		// Terjadi error lain
		return 0, err
	}
	return status, nil
}

func UpdateInvestPartner(investPartnerKey string, params map[string]string) (int, error) {
	query := `UPDATE cms_invest_partner SET `
	var setClauses []string
	var values []interface{}

	for key, value := range params {
		setClauses = append(setClauses, key+" = ?")
		values = append(values, value)
	}
	query += strings.Join(setClauses, ", ")
	query += ` WHERE invest_partner_key = ?`
	values = append(values, investPartnerKey)

	log.Println("========== UpdateInvestPartner ==========>>>", query)

	_, err := db.Db.Exec(query, values...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
