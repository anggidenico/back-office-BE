package controllers

import (
	"database/sql"
	"errors"
	"mf-bo-api/db"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func GetEndpointscController(c echo.Context) error {
	var endpoint []models.ScEndpointt
	status, err := models.GetEndpointscModels(&endpoint)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = endpoint
	return c.JSON(http.StatusOK, response)
}

func GetEndpointDetailController(c echo.Context) error {
	endpointKey := c.Param("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing endpoint key", "Missing endpoint key")
	}
	var detailendpoint models.ScEndpointDetail
	status, err := models.GetDetailEndpointModels(&detailendpoint, endpointKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.CustomError(http.StatusNotFound, "endpoint_key not found", "endpoint_key not found")
		}
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = detailendpoint
	return c.JSON(http.StatusOK, response)
}

func CreateEndpointController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	endpointCode := c.FormValue("endpoint_code")
	if endpointCode != "" {
		if len(endpointCode) > 150 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_code should be <= 150 characters", "endpoint_code should be <= 150 characters")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_code can not be blank", "endpoint_code can not be blank")
	}
	endpointCategoryKey := c.FormValue("endpoint_category_key")
	if endpointCategoryKey != "" {
		if len(endpointCategoryKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_category_key should be <= 11 characters", "endpoint_category_key should be <= 11 characters")
		}
		categoryKey, err := strconv.Atoi(endpointCategoryKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "endpoint_category_key must be a number", "endpoint_category_key must be a number")
		}
		// Periksa apakah nilai categoryKey sudah ada di tabel sc_endpoint_category
		query := "SELECT COUNT(*) FROM sc_endpoint_category WHERE endpoint_category_key = ?"
		var count int
		err = db.Db.QueryRow(query, categoryKey).Scan(&count)
		if err != nil {
			// Tangani kesalahan saat menjalankan query
			return lib.CustomError(http.StatusBadRequest, "error in query", "error in query")
		}
		if count == 0 {
			// Kembalikan kesalahan jika nilai categoryKey tidak ditemukan di tabel sc_endpoint_category
			return lib.CustomError(http.StatusBadRequest, "Invalid endpoint_category_key", "Invalid endpoint_category_key")
		}

	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_category_key can not be blank", "endpoint_category_key can not be blank")
	}

	endpointName := c.FormValue("endpoint_name")
	if endpointName != "" {
		if len(endpointName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_name should be <= 150 characters", "endpoint_name should be <= 150 characters")
		}
	}
	endpointVerb := c.FormValue("endpoint_verb")
	if endpointVerb != "" {
		if len(endpointVerb) > 20 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_verb should be <= 20 characters", "endpoint_verb should be <= 20 characters")
		}
		endpointVerb = strings.ToUpper(endpointVerb)
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_verb can not be blank", "endpoint_verb can not be blank")
	}
	endpointUrl := c.FormValue("endpoint_uri")
	if endpointUrl != "" {
		if len(endpointUrl) > 255 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_name should be <= 255 characters", "endpoint_name should be <= 255 characters")
		}
	}
	endpointRoute := c.FormValue("endpoint_route")
	if endpointRoute != "" {
		if len(endpointRoute) > 100 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_route should be <= 100 characters", "endpoint_route should be <= 100 characters")
		}
	}
	endpointController := c.FormValue("endpoint_controller")
	if endpointController != "" {
		if len(endpointController) > 50 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_controller should be <= 50 characters", "endpoint_route should be <= 50 characters")
		}
	}
	endpointVersion := c.FormValue("endpoint_version")
	if endpointVersion != "" {
		if len(endpointController) > 11 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_version must be <= 11 characters", "endpoint_version must be <= 11 characters")
		}
		_, err := strconv.Atoi(endpointVersion)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "endpoint_version must be number", "endpoint_version must be number")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_version can not be blank", "endpoint_version can not be blank")
	}
	menuKey := c.FormValue("menu_key")
	if menuKey != "" {
		if len(menuKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "menu_key must be <= 11 characters", "menu_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(menuKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "menu_key must be number", "menu_key must be number")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	}
	endpointDesc := c.FormValue("endpoint_desc")
	if endpointDesc != "" {
		if len(endpointDesc) > 255 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_desc should be <= 255 characters", "endpoint_desc should be <= 255 characters")
		}
	}
	params["endpoint_code"] = endpointCode
	params["endpoint_category_key"] = endpointCategoryKey
	params["endpoint_name"] = endpointName
	params["endpoint_verb"] = endpointVerb
	params["endpoint_uri"] = endpointUrl
	params["endpoint_route"] = endpointRoute
	params["endpoint_controller"] = endpointController
	params["endpoint_desc"] = endpointDesc
	params["rec_status"] = "1"
	params["menu_key"] = menuKey
	params["endpoint_version"] = "1"
	params["privileges_key"] = "READ"

	status, err = models.CreateEndpointSc(params)
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

func UpdateEndpointController(c echo.Context) error {
	var err error
	params := make(map[string]string)

	endpointKey := c.FormValue("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	}
	endpointCode := c.FormValue("endpoint_code")
	if endpointCode != "" {
		if len(endpointCode) > 150 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_code should be <= 150 characters", "endpoint_code should be <= 150 characters")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_code can not be blank", "endpoint_code can not be blank")
	}
	endpointCategoryKey := c.FormValue("endpoint_category_key")
	if endpointCategoryKey != "" {
		if len(endpointCategoryKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_category_key should be <= 11 characters", "endpoint_category_key should be <= 11 characters")
		}
		categoryKey, err := strconv.Atoi(endpointCategoryKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "endpoint_category_key must be a number", "endpoint_category_key must be a number")
		}
		// Periksa apakah nilai categoryKey sudah ada di tabel sc_endpoint_category
		query := "SELECT COUNT(*) FROM sc_endpoint_category WHERE endpoint_category_key = ?"
		var count int
		err = db.Db.QueryRow(query, categoryKey).Scan(&count)
		if err != nil {
			// Tangani kesalahan saat menjalankan query
			return lib.CustomError(http.StatusBadRequest, "error in query", "error in query")
		}
		if count == 0 {
			// Kembalikan kesalahan jika nilai categoryKey tidak ditemukan di tabel sc_endpoint_category
			return lib.CustomError(http.StatusBadRequest, "Invalid endpoint_category_key", "Invalid endpoint_category_key")
		}

	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_category_key can not be blank", "endpoint_category_key can not be blank")
	}

	endpointName := c.FormValue("endpoint_name")
	if endpointName != "" {
		if len(endpointName) > 150 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_name should be <= 150 characters", "endpoint_name should be <= 150 characters")
		}
	}

	endpointVerb := c.FormValue("endpoint_verb")
	if endpointVerb != "" {
		if len(endpointVerb) > 20 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_verb should be <= 20 characters", "endpoint_verb should be <= 20 characters")
		}
		endpointVerb = strings.ToUpper(endpointVerb)
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_verb can not be blank", "endpoint_verb can not be blank")
	}
	endpointUrl := c.FormValue("endpoint_uri")
	if endpointUrl != "" {
		if len(endpointUrl) > 255 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_name should be <= 255 characters", "endpoint_name should be <= 255 characters")
		}
	}
	endpointRoute := c.FormValue("endpoint_route")
	if endpointRoute != "" {
		if len(endpointRoute) > 100 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_route should be <= 100 characters", "endpoint_route should be <= 100 characters")
		}
	}
	endpointController := c.FormValue("endpoint_controller")
	if endpointController != "" {
		if len(endpointController) > 50 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_controller should be <= 50 characters", "endpoint_route should be <= 50 characters")
		}
	}
	endpointVersion := c.FormValue("endpoint_version")
	if endpointVersion != "" {
		if len(endpointController) > 11 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_version must be <= 11 characters", "endpoint_version must be <= 11 characters")
		}
		_, err := strconv.Atoi(endpointVersion)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "endpoint_version must be number", "endpoint_version must be number")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "endpoint_version can not be blank", "endpoint_version can not be blank")
	}
	menuKey := c.FormValue("menu_key")
	if menuKey != "" {
		if len(menuKey) > 11 {
			return lib.CustomError(http.StatusBadRequest, "menu_key must be <= 11 characters", "menu_key must be <= 11 characters")
		}
		_, err := strconv.Atoi(menuKey)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, "menu_key must be number", "menu_key must be number")
		}
	} else {
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	}
	endpointDesc := c.FormValue("endpoint_desc")
	if endpointDesc != "" {
		if len(endpointDesc) > 255 {
			return lib.CustomError(http.StatusBadRequest, "endpoint_desc should be <= 255 characters", "endpoint_desc should be <= 255 characters")
		}
	}
	params["endpoint_code"] = endpointCode
	params["endpoint_category_key"] = endpointCategoryKey
	params["endpoint_name"] = endpointName
	params["endpoint_verb"] = endpointVerb
	params["endpoint_uri"] = endpointUrl
	params["endpoint_route"] = endpointRoute
	params["endpoint_controller"] = endpointController
	params["endpoint_desc"] = endpointDesc
	params["rec_status"] = "1"
	params["menu_key"] = menuKey
	params["endpoint_version"] = "1"
	params["privileges_key"] = "READ"

	status, err = models.UpdateEndpointSc(endpointKey, params)
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

func DeleteEndpointController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	endpointKey := c.FormValue("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing endpointKey", "Missing endpointKey")
	}
	params["endpoint_key"] = endpointKey

	status, err := models.UpdateEndpointSc(endpointKey, params)
	if err != nil {
		return lib.CustomError(status, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Endpoint!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetEndpointCategoryController(c echo.Context) error {
	var endpoint []models.EndpointCategory
	status, err := models.GetEndpointCategoryModels(&endpoint)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = endpoint
	return c.JSON(http.StatusOK, response)
}

func GetScMenuController(c echo.Context) error {
	var endpoint []models.ScMenuu
	status, err := models.GetScMenuModels(&endpoint)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = endpoint
	return c.JSON(http.StatusOK, response)
}
