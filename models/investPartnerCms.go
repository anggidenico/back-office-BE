package models

import (
	"database/sql"
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
	RecOrder            *int64  `db:"rec_order" json:"recorder_order"`
}
type InvestPurpose struct {
	InvestPartnerKey int64   `db:"invest_partner_key" json:"invest_partner_key"`
	InvestPurposeKey int64   `db:"invest_purpose_key" json:"invest_purpose_key"`
	PurposeCode      *string `db:"purpose_code" json:"purpose_code"`
	PurposeName      *string `db:"purpose_name" json:"purpose_name"`
}

func GetInvestPurposeModels(c *[]InvestPurpose) (int, error) {
	query := `SELECT a.invest_partner_key,
	a.invest_purpose_key,
	b.purpose_code,
	b.purpose_name   
	FROM cms_invest_partner a
	JOIN cms_invest_purpose b 
	ON a.invest_purpose_key = b.invest_purpose_key
	WHERE a.rec_status =1 
	ORDER BY a.invest_partner_key DESC`
	log.Println("====================>>>", query)
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
	a.rec_order   
	FROM cms_invest_partner a
	JOIN cms_invest_purpose b 
	ON a.invest_purpose_key = b.invest_purpose_key
	WHERE a.rec_status =1 
	ORDER BY a.invest_partner_key DESC`
	log.Println("====================>>>", query)
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
