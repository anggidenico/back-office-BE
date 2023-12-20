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
	paramOaRequest["oa_status"] = "444"
	paramOaRequest["oa_request_type"] = "127"
	paramOaRequest["oa_entry_start"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key != "" {
		paramOaRequest["oa_request_key"] = oa_request_key
	}

	oa_source := c.FormValue("oa_source")
	if oa_source == "" {
		return fmt.Errorf("Missing: oa_source"), oaRequestKey
	}
	paramOaRequest["oa_source"] = oa_source

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

	email := c.FormValue("email")
	if email == "" {
		return fmt.Errorf("Missing: email"), oaRequestKey
	}
	paramPersonalData["email_address"] = email

	mobile_phone_number := c.FormValue("mobile_phone_number")
	if mobile_phone_number == "" {
		return fmt.Errorf("Missing: mobile_phone_number"), oaRequestKey
	}
	paramPersonalData["phone_mobile"] = mobile_phone_number

	phone_number := c.FormValue("phone_number")
	if phone_number == "" {
		return fmt.Errorf("Missing: phone_number"), oaRequestKey
	}
	paramPersonalData["phone_home"] = phone_number

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

	paramsOther := make(map[string]string)

	education_other := c.FormValue("education_other")
	if education_other != "" {
		paramsOther[education] = education_other
	}

	religion_other := c.FormValue("religion_other")
	if religion_other != "" {
		paramsOther[religion] = religion_other
	}

	paramsOffice := make(map[string]string)

	err, oaRequestKey := models.CreateOrUpdateOAManual(paramOaRequest, paramPersonalData, paramsIDCard, paramsDomicile, paramsOffice, paramsOther)
	if err != nil {
		return err, oaRequestKey
	}

	return nil, oaRequestKey
}

func SaveStep2(c echo.Context) (error, int64) {
	var OaRequestKey int64

	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "2"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), OaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

	site_referer := c.FormValue("site_referer")
	if site_referer == "" {
		return fmt.Errorf("Missing: site_referer"), OaRequestKey
	}
	paramsOaRequest["site_referer"] = site_referer

	paramsPersonalData := make(map[string]string)
	paramsPersonalData["oa_request_key"] = oa_request_key

	occupation := c.FormValue("occupation")
	if occupation == "" {
		return fmt.Errorf("Missing: occupation"), OaRequestKey
	}
	paramsPersonalData["occup_job"] = occupation

	company := c.FormValue("company")
	if company == "" {
		return fmt.Errorf("Missing: company"), OaRequestKey
	}
	paramsPersonalData["occup_company"] = company

	position := c.FormValue("position")
	if position == "" {
		return fmt.Errorf("Missing: position"), OaRequestKey
	}
	paramsPersonalData["occup_position"] = position

	business_fields := c.FormValue("business_fields")
	if business_fields == "" {
		return fmt.Errorf("Missing: business_fields"), OaRequestKey
	}
	paramsPersonalData["occup_business_fields"] = business_fields

	annual_income := c.FormValue("annual_income")
	if annual_income == "" {
		return fmt.Errorf("Missing: annual_income"), OaRequestKey
	}
	paramsPersonalData["annual_income"] = annual_income

	source_of_income := c.FormValue("source_of_income")
	if source_of_income == "" {
		return fmt.Errorf("Missing: source_of_income"), OaRequestKey
	}
	paramsPersonalData["sourceof_fund"] = source_of_income

	investment_objectives := c.FormValue("investment_objectives")
	if investment_objectives == "" {
		return fmt.Errorf("Missing: investment_objectives"), OaRequestKey
	}
	paramsPersonalData["invesment_objectives"] = investment_objectives

	paramsOfficeAddr := make(map[string]string)
	paramsOfficeAddr["rec_status"] = "1"
	paramsOfficeAddr["rec_created_by"] = lib.UserIDStr
	paramsOfficeAddr["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	paramsOfficeAddr["address_type"] = "19"

	address := c.FormValue("address")
	if address == "" {
		return fmt.Errorf("Missing: address"), OaRequestKey
	}
	paramsOfficeAddr["address_line1"] = address

	paramsOther := make(map[string]string)

	occupation_other := c.FormValue("occupation_other")
	if occupation_other != "" {
		paramsOther[occupation] = occupation_other
	}

	position_other := c.FormValue("position_other")
	if position_other != "" {
		paramsOther[position] = position_other
	}

	business_fields_other := c.FormValue("business_fields_other")
	if business_fields_other != "" {
		paramsOther[business_fields] = business_fields_other
	}

	source_of_income_other := c.FormValue("source_of_income_other")
	if source_of_income_other != "" {
		paramsOther[source_of_income] = source_of_income_other
	}

	investment_objectives_other := c.FormValue("investment_objectives_other")
	if investment_objectives_other != "" {
		paramsOther[investment_objectives] = investment_objectives_other
	}

	paramsIdCard := make(map[string]string)
	paramsDomi := make(map[string]string)

	err, OaRequestKey := models.CreateOrUpdateOAManual(paramsOaRequest, paramsPersonalData, paramsIdCard, paramsDomi, paramsOfficeAddr, paramsOther)
	if err != nil {
		return err, OaRequestKey
	}

	return nil, OaRequestKey
}

func SaveStep3(c echo.Context) (error, int64) {
	var OaRequestKey int64

	paramsOaRequest := make(map[string]string)
	paramsOaRequest["oa_step"] = "3"

	oa_request_key := c.FormValue("oa_request_key")
	if oa_request_key == "" {
		return fmt.Errorf("Missing: oa_request_key"), OaRequestKey
	}
	paramsOaRequest["oa_request_key"] = oa_request_key

	paramsPersonalData := make(map[string]string)

	mother_maiden_name := c.FormValue("mother_maiden_name")
	if mother_maiden_name == "" {
		return fmt.Errorf("Missing: mother_maiden_name"), OaRequestKey
	}
	paramsPersonalData["mother_maiden_name"] = mother_maiden_name

	relation_type := c.FormValue("relation_type")
	if relation_type == "" {
		return fmt.Errorf("Missing: relation_type"), OaRequestKey
	}
	paramsPersonalData["relation_type"] = relation_type

	relation_full_name := c.FormValue("relation_full_name")
	if relation_full_name == "" {
		return fmt.Errorf("Missing: relation_full_name"), OaRequestKey
	}
	paramsPersonalData["relation_full_name"] = relation_full_name

	relation_occupation := c.FormValue("relation_occupation")
	if relation_occupation == "" {
		return fmt.Errorf("Missing: relation_occupation"), OaRequestKey
	}
	paramsPersonalData["relation_occupation"] = relation_occupation

	relation_business_fields := c.FormValue("relation_business_fields")
	if relation_business_fields == "" {
		return fmt.Errorf("Missing: relation_business_fields"), OaRequestKey
	}
	paramsPersonalData["relation_business_fields"] = relation_business_fields

	emergency_full_name := c.FormValue("emergency_full_name")
	if emergency_full_name == "" {
		return fmt.Errorf("Missing: emergency_full_name"), OaRequestKey
	}
	paramsPersonalData["emergency_full_name"] = emergency_full_name

	emergency_relation := c.FormValue("emergency_relation")
	if emergency_relation == "" {
		return fmt.Errorf("Missing: emergency_relation"), OaRequestKey
	}
	paramsPersonalData["emergency_relation"] = emergency_relation

	emergency_phone_no := c.FormValue("emergency_phone_no")
	if emergency_phone_no == "" {
		return fmt.Errorf("Missing: emergency_phone_no"), OaRequestKey
	}
	paramsPersonalData["emergency_phone_no"] = emergency_phone_no

	beneficial_full_name := c.FormValue("beneficial_full_name")
	if beneficial_full_name == "" {
		return fmt.Errorf("Missing: beneficial_full_name"), OaRequestKey
	}
	paramsPersonalData["beneficial_full_name"] = beneficial_full_name

	beneficial_relation := c.FormValue("beneficial_relation")
	if beneficial_relation == "" {
		return fmt.Errorf("Missing: beneficial_relation"), OaRequestKey
	}
	paramsPersonalData["beneficial_relation"] = beneficial_relation

	pep_status := c.FormValue("pep_status")
	if pep_status == "" {
		return fmt.Errorf("Missing: pep_status"), OaRequestKey
	}
	paramsPersonalData["pep_status"] = pep_status

	pep_name := c.FormValue("pep_name")
	if pep_name == "" {
		// return fmt.Errorf("Missing: pep_name"), OaRequestKey
	}
	paramsPersonalData["pep_name"] = pep_name

	pep_position := c.FormValue("pep_position")
	if pep_position == "" {
		// return fmt.Errorf("Missing: pep_position"), OaRequestKey
	}
	paramsPersonalData["pep_position"] = pep_position

	paramsOther := make(map[string]string)

	relation_occupation_other := c.FormValue("relation_occupation_other")
	if relation_occupation_other != "" {
		paramsOther[relation_occupation] = relation_occupation_other
	}

	relation_business_fields_other := c.FormValue("relation_business_fields_other")
	if relation_business_fields_other != "" {
		paramsOther[relation_business_fields] = relation_business_fields_other
	}

	beneficial_relation_other := c.FormValue("beneficial_relation_other")
	if beneficial_relation_other != "" {
		paramsOther[beneficial_relation] = beneficial_relation_other
	}

	paramsIdCard := make(map[string]string)
	paramsDomi := make(map[string]string)
	paramsOfficeAddr := make(map[string]string)

	err, OaRequestKey := models.CreateOrUpdateOAManual(paramsOaRequest, paramsPersonalData, paramsIdCard, paramsDomi, paramsOfficeAddr, paramsOther)
	if err != nil {
		return err, OaRequestKey
	}

	return nil, OaRequestKey
}
