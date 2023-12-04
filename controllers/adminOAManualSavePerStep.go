package controllers

import (
	"fmt"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"time"

	"github.com/labstack/echo"
)

func SaveStep1(c echo.Context) (error, int64) {
	var oaRequestKey int64

	paramOaRequest := make(map[string]string)
	paramOaRequest["rec_status"] = "1"
	paramOaRequest["rec_created_by"] = lib.UserIDStr
	paramOaRequest["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramOaRequest["oa_step"] = "1"
	paramOaRequest["oa_status"] = "258"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key != "" {
		paramOaRequest["oa_request_key"] = oa_request_key
	}

	branch_key := c.FormValue("branch_key")
	if branch_key == "" {
		return fmt.Errorf("Missing: branch_key"), oaRequestKey
	}
	paramOaRequest["branch_key"] = branch_key

	agent_key := c.FormValue("agent_key")
	if agent_key == "" {
		return fmt.Errorf("Missing: agent_key"), oaRequestKey
	}
	paramOaRequest["agent_key"] = agent_key

	// PERSONAL DATA

	paramPersonalData := make(map[string]string)
	paramPersonalData["rec_status"] = "1"
	paramPersonalData["rec_created_by"] = lib.UserIDStr
	paramPersonalData["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	nationality := c.FormValue("nationality")
	if nationality == "" {
		return fmt.Errorf("Missing: nationality"), oaRequestKey
	}
	paramPersonalData["nationality"] = nationality

	if nationality == "97" {
		paramPersonalData["idcard_type"] = "12"
	} else {
		paramPersonalData["idcard_type"] = "13"
	}

	idcard_no := c.FormValue("idcard_no")
	if idcard_no == "" {
		return fmt.Errorf("Missing: idcard_no"), oaRequestKey
	}
	paramPersonalData["idcard_no"] = idcard_no

	full_name := c.FormValue("full_name")
	if full_name == "" {
		return fmt.Errorf("Missing: full_name"), oaRequestKey
	}
	paramPersonalData["full_name"] = full_name

	gender := c.FormValue("gender")
	if gender == "" {
		return fmt.Errorf("Missing: gender"), oaRequestKey
	}
	paramPersonalData["gender"] = gender

	place_birth := c.FormValue("place_birth")
	if place_birth == "" {
		return fmt.Errorf("Missing: place_birth"), oaRequestKey
	}
	paramPersonalData["place_birth"] = place_birth

	date_birth := c.FormValue("date_birth")
	if date_birth == "" {
		return fmt.Errorf("Missing: date_birth"), oaRequestKey
	}
	paramPersonalData["date_birth"] = date_birth

	marital_status := c.FormValue("marital_status")
	if marital_status == "" {
		return fmt.Errorf("Missing: marital_status"), oaRequestKey
	}
	paramPersonalData["marital_status"] = marital_status

	religion := c.FormValue("religion")
	if religion == "" {
		return fmt.Errorf("Missing: religion"), oaRequestKey
	}
	paramPersonalData["religion"] = religion

	education := c.FormValue("education")
	if education == "" {
		return fmt.Errorf("Missing: education"), oaRequestKey
	}
	paramPersonalData["education"] = education

	// ADDRESS

	paramsIDCard := make(map[string]string)
	paramsIDCard["rec_status"] = "1"
	paramsIDCard["rec_created_by"] = lib.UserIDStr
	paramsIDCard["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsIDCard["address_type"] = "17"

	idcard_address := c.FormValue("idcard_address")
	if idcard_address == "" {
		return fmt.Errorf("Missing: idcard_address"), oaRequestKey
	}
	paramsIDCard["address_line1"] = idcard_address

	idcard_province := c.FormValue("idcard_province")
	if idcard_province == "" {
		return fmt.Errorf("Missing: idcard_province"), oaRequestKey
	}
	paramsIDCard["province_key"] = idcard_province

	idcard_city := c.FormValue("idcard_city")
	if idcard_city == "" {
		return fmt.Errorf("Missing: idcard_city"), oaRequestKey
	}
	paramsIDCard["kabupaten_key"] = idcard_city

	idcard_postal_code := c.FormValue("idcard_postal_code")
	if idcard_postal_code == "" {
		return fmt.Errorf("Missing: idcard_postal_code"), oaRequestKey
	}
	paramsIDCard["postal_code"] = idcard_postal_code

	// DOMICILE

	paramsDomicile := make(map[string]string)
	paramsDomicile["rec_status"] = "1"
	paramsDomicile["rec_created_by"] = lib.UserIDStr
	paramsDomicile["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsDomicile["address_type"] = "18"

	domicile_address := c.FormValue("domicile_address")
	if domicile_address == "" {
		return fmt.Errorf("Missing: domicile_address"), oaRequestKey
	}
	paramsDomicile["address_line1"] = domicile_address

	domicile_province := c.FormValue("domicile_province")
	if domicile_province == "" {
		return fmt.Errorf("Missing: domicile_province"), oaRequestKey
	}
	paramsDomicile["province_key"] = domicile_province

	domicile_city := c.FormValue("domicile_city")
	if domicile_city == "" {
		return fmt.Errorf("Missing: domicile_city"), oaRequestKey
	}
	paramsDomicile["kabupaten_key"] = domicile_city

	domicile_postal_code := c.FormValue("domicile_postal_code")
	if domicile_postal_code == "" {
		return fmt.Errorf("Missing: domicile_postal_code"), oaRequestKey
	}
	paramsDomicile["postal_code"] = domicile_postal_code

	paramsOffice := make(map[string]string)

	err, oaRequestKey := models.CreateOrUpdateOAManual(paramOaRequest, paramPersonalData, paramsIDCard, paramsDomicile, paramsOffice)
	if err != nil {
		return err, oaRequestKey
	}

	return nil, oaRequestKey
}

func SaveStep2(c echo.Context) (error, int64) {
	var oaRequestKey int64

	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "2"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), oaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

	site_referer := c.FormValue("site_referer")
	if site_referer == "" {
		return fmt.Errorf("Missing: site_referer"), oaRequestKey
	}
	paramsOaRequest["site_referer"] = site_referer

	paramsPersonalData := make(map[string]string)
	paramsPersonalData["oa_request_key"] = oa_request_key

	occupation := c.FormValue("occupation")
	if occupation == "" {
		return fmt.Errorf("Missing: occupation"), oaRequestKey
	}
	paramsPersonalData["occup_job"] = occupation

	company := c.FormValue("company")
	if company == "" {
		return fmt.Errorf("Missing: company"), oaRequestKey
	}
	paramsPersonalData["occup_company"] = company

	position := c.FormValue("position")
	if position == "" {
		return fmt.Errorf("Missing: position"), oaRequestKey
	}
	paramsPersonalData["occup_position"] = position

	business_fields := c.FormValue("business_fields")
	if business_fields == "" {
		return fmt.Errorf("Missing: business_fields"), oaRequestKey
	}
	paramsPersonalData["occup_business_fields"] = business_fields

	annual_income := c.FormValue("annual_income")
	if annual_income == "" {
		return fmt.Errorf("Missing: annual_income"), oaRequestKey
	}
	paramsPersonalData["annual_income"] = annual_income

	source_of_income := c.FormValue("source_of_income")
	if source_of_income == "" {
		return fmt.Errorf("Missing: source_of_income"), oaRequestKey
	}
	paramsPersonalData["sourceof_fund"] = source_of_income

	investment_objectives := c.FormValue("investment_objectives")
	if investment_objectives == "" {
		return fmt.Errorf("Missing: investment_objectives"), oaRequestKey
	}
	paramsPersonalData["investment_objectives"] = investment_objectives

	paramsOfficeAddr := make(map[string]string)
	paramsOfficeAddr["rec_status"] = "1"
	paramsOfficeAddr["rec_created_by"] = lib.UserIDStr
	paramsOfficeAddr["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsOfficeAddr["address_type"] = "19"

	address := c.FormValue("address")
	if address == "" {
		return fmt.Errorf("Missing: address"), oaRequestKey
	}
	paramsOfficeAddr["address_line1"] = address

	paramsIdCard := make(map[string]string)
	paramsDomi := make(map[string]string)

	err, oaRequestKey := models.CreateOrUpdateOAManual(paramsOaRequest, paramsPersonalData, paramsIdCard, paramsDomi, paramsOfficeAddr)
	if err != nil {
		return err, oaRequestKey
	}

	return nil, oaRequestKey
}
