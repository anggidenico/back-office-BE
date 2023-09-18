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

func GetOaRequestListCustomerBuild(c echo.Context) error {
	var responseData []models.OaRequestListModelsResponse
	result := models.GetOaRequestCustomerBuildListQuery()

	if len(result) > 0 {
		for _, rData := range result {
			datas := rData

			requestKey := strconv.FormatUint(datas.OaRequestKey, 10)
			log.Println(requestKey)
			var personalData models.PengkinianPersonalDataResponse

			if *datas.OaRequestTypeInt == uint64(lib.OA_REQ_TYPE_PENGKINIAN_RISIKO_INT) { // JIKA PENGKINIAN PROFIL RISIKO
				requestKey1 := models.GetLastOaRequestHasPersonalData(strconv.FormatUint(*datas.UserLoginKey, 10))
				personalData = GetThePersonalDataDetails(requestKey1)

				layoutDateBirth := "02 Jan 2006"
				t1, _ := time.Parse(lib.TIMESTAMPFORMAT, *personalData.DateBirth)
				dateBirth := t1.Format(layoutDateBirth)
				datas.DateBirth = &dateBirth
				datas.FullName = personalData.FullName
				datas.Agent = personalData.Agent
				datas.Branch = personalData.Branch
				datas.IDCardNo = personalData.IdCardNo
			}

			layout := "02 Jan 2006 15:04"
			t2, _ := time.Parse(lib.TIMESTAMPFORMAT, *rData.OaDate)
			*datas.OaDate = t2.Format(layout)

			responseData = append(responseData, datas)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func RevertOAStatus(c echo.Context) error {

	OaRequestyKeys := c.QueryParam("oa_request_key")
	if OaRequestyKeys == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: oa_request_key", "Missing: oa_request_key")
	}

	err := models.SetOAStatusRevert(OaRequestyKeys)
	if err != nil {
		return lib.CustomError(http.StatusBadRequest, "Can not revert oa request", "Can not revert oa request")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
