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

	err, oaRequestKey := models.CreateOrUpdateOAManual(paramOaRequest, paramPersonalData, paramsIDCard, paramsDomicile)
	if err != nil {
		return err, oaRequestKey
	}

	return nil, oaRequestKey
}
