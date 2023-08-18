package controllers

import (
	"bytes"
	"html/template"
	"log"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

func GetPengkinianRiskProfileList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		return lib.CustomError(http.StatusUnauthorized, "You not allowed to access this page", "You not allowed to access this page")
	}
	var err error
	var responseData []models.RiskProfileListModels

	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	var getList []models.RiskProfileListModels
	pagination := models.GetPengkinianRiskProfileListQuery(&getList, lib.Profile.RoleKey, limit, offset)
	if len(getList) > 0 {
		// responseData = getList
		for _, getData := range getList {
			respData := getData
			layout := "02 Jan 2006 15:04"
			t2, _ := time.Parse(lib.TIMESTAMPFORMAT, getData.OaDate)
			respData.OaDate = t2.Format(layout)
			responseData = append(responseData, respData)
		}
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetPengkinianRiskProfileDetails(c echo.Context) error {
	// var err error
	var responseData models.RiskProfileDetailResponse
	layout := "02 Jan 2006 15:04"

	OaRequestKey := c.Param("key")
	if OaRequestKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing: oaRequestKey")
	}

	QnA_Array := models.GetQuizQuestionAnswerQuery(OaRequestKey)
	if len(QnA_Array) > 0 {
		responseData.RiskProfileQuizAnswer = QnA_Array
	}

	QuizResult := models.GetRiskProfileQuizResult(OaRequestKey)
	responseData.RiskProfileQuizResult = QuizResult

	t2, _ := time.Parse(lib.TIMESTAMPFORMAT, QuizResult.OaDate)
	OaDate := t2.Format(layout)
	responseData.OaDate = OaDate
	responseData.Cif = QuizResult.Cif
	responseData.FullName = QuizResult.FullName
	responseData.OaSource = QuizResult.OaSource
	responseData.OaStatus = QuizResult.OaStatus

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func SendEmailRejectRiskProfilePengkinianToCustomer(EmailTo string, Nama string) error {
	var err error

	paramEmail := make(map[string]string)
	paramEmail["FileUrl"] = config.ImageUrl + "/images/mail"
	paramEmail["Nama"] = Nama
	paramEmail["Judul"] = "Pengkinian Data Kamu Belum Dapat Disetujui"
	paramEmail["Keterangan"] = "Pengkinian Data yang kamu ajukan belum dapat dapat disetujui"

	emailTemplate := "email-pengkinian-risk-profile-rejected.html"
	t := template.New(emailTemplate)
	t, err = t.ParseFiles(config.BasePath + "/mail/" + emailTemplate)
	if err != nil {
		log.Println(err.Error())
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, paramEmail)
	if err != nil {
		log.Println(err.Error())
	}
	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("To", EmailTo)
	mailer.SetHeader("Subject", paramEmail["Judul"])
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
