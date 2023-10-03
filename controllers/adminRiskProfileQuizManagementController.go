package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

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

func GetQuestionDetail(c echo.Context) error {

	questionKey := c.Param("question_key")
	if questionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question key", "Missing question key")
	}

	result := models.GetQuestionDetail(questionKey)

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

func UpdateQuizQuestion(c echo.Context) error {
	UpdateMaps := make(map[string]string)
	UpdateMaps["rec_modified_by"] = lib.UserIDStr
	UpdateMaps["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	quizQuestionKey := c.FormValue("quiz_question_key")
	if quizQuestionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question quiz_question_key", "Missing question quiz_question_key")
	}

	quizHeaderKey := c.FormValue("quiz_header_key")
	if quizHeaderKey != "" {
		UpdateMaps["quiz_header_key"] = quizHeaderKey
	}

	quizTitle := c.FormValue("quiz_title")
	if quizTitle == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_title", "Missing quiz_title")
	} else {
		UpdateMaps["quiz_title"] = quizTitle
	}

	quizOptionType := c.FormValue("quiz_option_type")
	if quizOptionType == "" {
		UpdateMaps["quiz_option_type"] = "256" // PILIHAN GANDA
	} else {
		UpdateMaps["quiz_option_type"] = quizOptionType
	}

	err := models.UpdateQuizQuestion(quizQuestionKey, UpdateMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil ubah question!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func DeleteQuizQuestion(c echo.Context) error {
	UpdateMaps := make(map[string]string)
	UpdateMaps["rec_status"] = "0"
	UpdateMaps["rec_deleted_by"] = lib.UserIDStr
	UpdateMaps["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	quizQuestionKey := c.FormValue("quiz_question_key")
	if quizQuestionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question quiz_question_key", "Missing question quiz_question_key")
	}

	err := models.UpdateQuizQuestion(quizQuestionKey, UpdateMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus question!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func CreateQuizQuestion(c echo.Context) error {
	insertMaps := make(map[string]string)
	insertMaps["rec_status"] = "1"
	insertMaps["rec_created_by"] = lib.UserIDStr
	insertMaps["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	quizHeaderKey := c.FormValue("quiz_header_key")
	if quizHeaderKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_header_key", "Missing quiz_header_key")
	} else {
		insertMaps["quiz_header_key"] = quizHeaderKey
	}

	quizTitle := c.FormValue("quiz_title")
	if quizTitle == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_title", "Missing quiz_title")
	} else {
		insertMaps["quiz_title"] = quizTitle
	}

	quizOptionType := c.FormValue("quiz_option_type")
	if quizOptionType == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_type", "Missing quiz_option_type")
	} else {
		insertMaps["quiz_option_type"] = quizOptionType
	}

	err := models.CreateQuizQuestion(insertMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil tambah question!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func GetOptionDetail(c echo.Context) error {

	quizOptionKey := c.Param("quiz_option_key")
	if quizOptionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_key", "Missing quiz_option_key")
	}

	result := models.GetOptionDetail(quizOptionKey)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func CreateQuizOption(c echo.Context) error {
	insertMaps := make(map[string]string)
	insertMaps["rec_status"] = "1"
	insertMaps["rec_created_by"] = lib.UserIDStr
	insertMaps["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	quizquestionKey := c.FormValue("quiz_question_key")
	if quizquestionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_question_key", "Missing quiz_question_key")
	} else {
		insertMaps["quiz_question_key"] = quizquestionKey
	}

	quizOptionLabel := c.FormValue("quiz_option_label")
	if quizOptionLabel == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_label", "Missing quiz_option_label")
	} else {
		insertMaps["quiz_option_label"] = quizOptionLabel
	}

	quizOptionTitle := c.FormValue("quiz_option_title")
	if quizOptionTitle == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_title", "Missing quiz_option_title")
	} else {
		insertMaps["quiz_option_title"] = quizOptionTitle
	}

	quizOptionScore := c.FormValue("quiz_option_score")
	if quizOptionScore == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_score", "Missing quiz_option_score")
	} else {
		insertMaps["quiz_option_score"] = quizOptionScore
	}

	err := models.CreateQuizOption(insertMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil tambah option!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func UpdateQuizOption(c echo.Context) error {
	UpdateMaps := make(map[string]string)
	UpdateMaps["rec_modified_by"] = lib.UserIDStr
	UpdateMaps["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	quizQuestionKey := c.FormValue("quiz_question_key")
	if quizQuestionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing question quiz_question_key", "Missing question quiz_question_key")
	} else {
		UpdateMaps["quiz_question_key"] = quizQuestionKey
	}

	quizOptionKey := c.FormValue("quiz_option_key")
	if quizOptionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_key", "Missing quiz_option_key")
	} else {
		// UpdateMaps["quiz_option_key"] = quizOptionKey
	}

	quizOptionLabel := c.FormValue("quiz_option_label")
	if quizOptionLabel == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_label", "Missing quiz_option_label")
	} else {
		UpdateMaps["quiz_option_label"] = quizOptionLabel
	}

	quizOptionTitle := c.FormValue("quiz_option_title")
	if quizOptionTitle == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_title", "Missing quiz_option_title")
	} else {
		UpdateMaps["quiz_option_title"] = quizOptionTitle
	}

	quizOptionScore := c.FormValue("quiz_option_score")
	if quizOptionScore == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_score", "Missing quiz_option_score")
	} else {
		UpdateMaps["quiz_option_score"] = quizOptionScore
	}

	err := models.UpdateQuizOption(quizOptionKey, UpdateMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil ubah option!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func DeleteQuizOption(c echo.Context) error {
	UpdateMaps := make(map[string]string)
	UpdateMaps["rec_status"] = "0"
	UpdateMaps["rec_deleted_by"] = lib.UserIDStr
	UpdateMaps["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	quizOptionKey := c.FormValue("quiz_option_key")
	if quizOptionKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing quiz_option_key", "Missing quiz_option_key")
	} else {
		// UpdateMaps["quiz_option_key"] = quizOptionKey
	}

	err := models.UpdateQuizOption(quizOptionKey, UpdateMaps)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus question!"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
