package controllers

import (
	"database/sql"
	"errors"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func CreateRiskProfile(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	riskCode := c.FormValue("risk_code")
	if riskCode != "" {
		if len(riskCode) > 30 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 30", "risk_code should be <= 30")
		}
	} else {
		if riskCode == "" {
			return lib.CustomError(http.StatusBadRequest, "risk_code can not be blank", "risk_code can not be blank")
		}
	}
	params["risk_code"] = riskCode

	riskName := c.FormValue("risk_name")
	if riskName != "" {
		if len(riskName) > 50 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 50", "risk_code should be <= 50")
		}
	}
	params["risk_name"] = riskName

	riskDesc := c.FormValue("risk_desc")
	if riskDesc != "" {
		if len(riskDesc) > 1000 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 1000", "risk_code should be <= 1000")
		}
	}
	params["risk_desc"] = riskDesc

	min_score := c.FormValue("min_score")
	if min_score == "" {
		return lib.CustomError(http.StatusBadRequest, "min_score can not be blank", "min_score can not be blank")
	}
	params["min_score"] = min_score

	max_score := c.FormValue("max_score")
	if max_score == "" {
		return lib.CustomError(http.StatusBadRequest, "max_score can not be blank", "max_score can not be blank")
	}
	params["max_score"] = max_score

	maxFlag := c.FormValue("max_flag")
	if maxFlag != "" {
		if len(maxFlag) > 1 {
			return lib.CustomError(http.StatusBadRequest, "max_flag should be <= 1", "max_flag should be <= 1")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "max_flag can not be blank", "max_flag can not be blank")
	}
	params["max_flag"] = maxFlag

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		value, err := strconv.Atoi(recOrder)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be numeric", "rec_order should be numeric")
		}
		if len(recOrder) > 11 {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be <= 11 characters", "rec_order should be <= 11 characters")
		}
		params["rec_order"] = strconv.Itoa(value)
	}

	params["rec_status"] = "1"

	status, err = models.CreateRiskProfile(params)
	if err != nil {

		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func GetriskProfileController(c echo.Context) error {
	var riskprofile []models.RiskProfile
	status, err := models.GetRiskProfileModels(&riskprofile)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var response lib.Response
	if len(riskprofile) < 1 {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = []models.RiskProfile{}
	} else {
		response.Status.Code = http.StatusOK
		response.Status.MessageServer = "OK"
		response.Status.MessageClient = "OK"
		response.Data = riskprofile
	}
	return c.JSON(http.StatusOK, response)
}

func GetDetailRiskProfileController(c echo.Context) error {

	riskProfileKey := c.Param("risk_profile_key")
	if riskProfileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Risk profile not found", "Risk profile not found")
	}
	var detailrisk models.GetDetailRisk
	status, err := models.GetDetailRiskProfileModels(&detailrisk, riskProfileKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "Periode key not found", "Periode key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}
	if detailrisk.RiskProfileKey == "" {
		return c.NoContent(http.StatusOK)
	}
	var response lib.Response

	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailrisk
	return c.JSON(http.StatusOK, response)
}

func UpdateRiskProfile(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	riskprofileKey := c.FormValue("risk_profile_key")
	if riskprofileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	}
	params["risk_profile_key"] = riskprofileKey

	riskCode := c.FormValue("risk_code")
	if riskCode != "" {
		if len(riskCode) > 30 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 30", "risk_code should be <= 30")
		}
	} else {
		if riskCode == "" {
			return lib.CustomError(http.StatusBadRequest, "risk_code can not be blank", "risk_code can not be blank")
		}
	}
	params["risk_code"] = riskCode

	riskName := c.FormValue("risk_name")
	if riskName != "" {
		if len(riskName) > 50 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 50", "risk_code should be <= 50")
		}
	}
	params["risk_name"] = riskName

	riskDesc := c.FormValue("risk_desc")
	if riskDesc != "" {
		if len(riskDesc) > 1000 {
			return lib.CustomError(http.StatusBadRequest, "risk_code should be <= 1000", "risk_code should be <= 1000")
		}
	}
	params["risk_desc"] = riskDesc

	min_score := c.FormValue("min_score")
	if min_score == "" {
		return lib.CustomError(http.StatusBadRequest, "min_score can not be blank", "min_score can not be blank")
	}
	params["min_score"] = min_score

	max_score := c.FormValue("max_score")
	if max_score == "" {
		return lib.CustomError(http.StatusBadRequest, "max_score can not be blank", "max_score can not be blank")
	}
	params["max_score"] = max_score

	maxFlag := c.FormValue("max_flag")
	if maxFlag != "" {
		if len(maxFlag) > 1 {
			return lib.CustomError(http.StatusBadRequest, "max_flag should be <= 1", "max_flag should be <= 1")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "max_flag can not be blank", "max_flag can not be blank")
	}
	params["max_flag"] = maxFlag

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		value, err := strconv.Atoi(recOrder)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be numeric", "rec_order should be numeric")
		}
		if len(recOrder) > 11 {
			return lib.CustomError(http.StatusBadRequest, "rec_order should be <= 11 characters", "rec_order should be <= 11 characters")
		}
		params["rec_order"] = strconv.Itoa(value)
	}

	// params["rec_status"] = "1"
	status, err = models.UpdateRiskProfile(riskprofileKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func DeleteRiskProfile(c echo.Context) error {
	params := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	riskprofileKey := c.FormValue("risk_profile_key")
	if riskprofileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	}

	status, err := models.DeleteRiskProfile(riskprofileKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Risk Profile!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
