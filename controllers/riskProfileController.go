// handlers.go
package controllers

import (
	"log"
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

	risk_code := c.FormValue("risk_code")
	if risk_code == "" {

		return lib.CustomError(http.StatusBadRequest, "risk_code can not be blank", "risk_code can not be blank")
	}
	params["risk_code"] = risk_code

	risk_name := c.FormValue("risk_name")
	if risk_name == "" {

		return lib.CustomError(http.StatusBadRequest, "risk_name can not be blank", "risk_code can not be blank")
	}
	params["risk_name"] = risk_name

	risk_desc := c.FormValue("risk_desc")
	if risk_desc == "" {

		return lib.CustomError(http.StatusBadRequest, "risk_desc can not be blank", "risk_desc can not be blank")
	}
	params["risk_desc"] = risk_desc

	min_score := c.FormValue("min_score")
	if min_score == "" {

		return lib.CustomError(http.StatusBadRequest, "min_score can not be blank", "min_score can not be blank")
	}
	params["min_score"] = min_score

	max_score := c.FormValue("max_score")
	if max_score == "" {

		return lib.CustomError(http.StatusBadRequest, "max_score can not be blank", "min_score can not be blank")
	}
	params["max_score"] = max_score

	max_flag := c.FormValue("max_flag")
	if max_flag == "" {
		return lib.CustomError(http.StatusBadRequest, "max_flag can not be blank", "min_score can not be blank")
	}
	params["max_flag"] = max_flag

	rec_order := c.FormValue("rec_order")
	if rec_order == "" {

		return lib.CustomError(http.StatusBadRequest, "rec_order can not be blank", "min_score can not be blank")
	}

	params["rec_order"] = rec_order
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
	// key := c.Param("key")
	// if key == "" {
	// 	return lib.CustomError(http.StatusBadRequest, "Missing key", "Missing key")
	// }
	result := models.GetRiskProfileModels()
	// if len(result) < 1 {
	// 	return lib.CustomError(http.StatusInternalServerError, "Can not get risk profile", "Can not get risk profile")
	// }
	log.Println("data ga keambil")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
func GetDetailRiskProfileController(c echo.Context) error {

	riskProfileKey := c.Param("risk_profile_key")
	if riskProfileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question key", "Missing question key")
	}
	result := models.GetDetailRiskProfileModels(riskProfileKey)
	if len(result) < 1 {
		return lib.CustomError(http.StatusInternalServerError, "Can not get option list", "Can not get option list")
		// log.Println("Not Found")
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func UpdateRiskProfile(c echo.Context) error {
	var err error
	params := make(map[string]string)

	riskprofileKey := c.FormValue("risk_profile_key")
	if riskprofileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	}
	params["risk_profile_key"] = riskprofileKey

	risk_code := c.FormValue("risk_code")
	if risk_code == "" {

		return lib.CustomError(http.StatusBadRequest, "risk_code can not be blank", "risk_code can not be blank")
	}
	params["risk_code"] = risk_code

	risk_name := c.FormValue("risk_name")
	if risk_name == "" {

		return lib.CustomError(http.StatusBadRequest, "risk_name can not be blank", "risk_code can not be blank")
	}
	params["risk_name"] = risk_name

	risk_desc := c.FormValue("risk_desc")
	if risk_desc == "" {
		return lib.CustomError(http.StatusBadRequest, "risk_desc can not be blank", "risk_desc can not be blank")
	}
	params["risk_desc"] = risk_desc

	min_score := c.FormValue("min_score")
	if min_score == "" {
		return lib.CustomError(http.StatusBadRequest, "min_score can not be blank", "min_score can not be blank")
	}
	params["min_score"] = min_score

	max_score := c.FormValue("max_score")
	if max_score == "" {
		return lib.CustomError(http.StatusBadRequest, "max_score can not be blank", "min_score can not be blank")
	}
	params["max_score"] = max_score

	max_flag := c.FormValue("max_flag")
	if max_flag == "" {
		return lib.CustomError(http.StatusBadRequest, "max_flag can not be blank", "min_score can not be blank")
	}
	params["max_flag"] = max_flag

	rec_order := c.FormValue("rec_order")
	if rec_order == "" {

		return lib.CustomError(http.StatusBadRequest, "rec_order can not be blank", "min_score can not be blank")
	}
	params["rec_order"] = rec_order

	err = models.UpdateRiskProfile(riskprofileKey, params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed input data")
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

	err := models.DeleteRiskProfile(riskprofileKey, params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Risk Profile!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
