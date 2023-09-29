package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetRiskProfileQuestionList(c echo.Context) error {

	result := models.GetQuizQuestion()
	if len(result) < 1 {
		return lib.CustomError(http.StatusInternalServerError, "Can not get question list", "Can not get question list")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func GetOptionListPerQuestion(c echo.Context) error {

	questionKey := c.Param("question_key")
	if questionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question key", "Missing question key")
	}

	result := models.GetQuizOption(questionKey)
	if len(result) < 1 {
		return lib.CustomError(http.StatusInternalServerError, "Can not get option list", "Can not get option list")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
