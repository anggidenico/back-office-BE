package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func GetEndpointscController(c echo.Context) error {

	result := models.GetEndpointscModels()
	log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func GetEndpointDetailController(c echo.Context) error {
	endpointKey := c.Param("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing endpoint key", "Missing endpoint key")
	}
	result := models.GetDetailEndpointModels(endpointKey)
	// log.Println("Not Found")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}

func CreateEndpointController(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["rec_created_by"] = lib.UserIDStr
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)

	endpointCode := c.FormValue("endpoint_code")
	if endpointCode == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_code can not be blank", "endpoint_code can not be blank")
	}
	endpointCategoryKey := c.FormValue("endpoint_category_key")
	if endpointCategoryKey == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_category_key can not be blank", "endpoint_category_key can not be blank")
	}
	endpointName := c.FormValue("endpoint_name")
	if endpointName == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_name can not be blank", "endpoint_code can not be blank")
	}
	endpointVerb := c.FormValue("endpoint_verb")
	if endpointVerb == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_verb can not be blank", "endpoint_verb can not be blank")
	}
	endpointUrl := c.FormValue("endpoint_uri")
	if endpointUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_uri can not be blank", "endpoint_uri can not be blank")
	}
	menuKey := c.FormValue("menu_key")
	if menuKey == "" {
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	}
	privilegesKey := c.FormValue("privileges_key")
	if privilegesKey == "" {
		return lib.CustomError(http.StatusBadRequest, "privileges_key can not be blank", "privileges_key can not be blank")
	}
	endpointDesc := c.FormValue("endpoint_desc")
	if endpointDesc == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_desc can not be blank", "endpoint_desc can not be blank")
	}
	params["endpoint_code"] = endpointCode
	params["endpoint_category_key"] = endpointCategoryKey
	params["endpoint_name"] = endpointName
	params["endpoint_verb"] = endpointVerb
	params["endpoint_uri"] = endpointUrl
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

	endpointCategoryKey := c.FormValue("endpoint_category_key")
	if endpointCategoryKey == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_category_key can not be blank", "endpoint_category_key can not be blank")
	}

	endpointName := c.FormValue("endpoint_name")
	if endpointName == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_name can not be blank", "endpoint_name can not be blank")
	}

	endpointVerb := c.FormValue("endpoint_verb")
	if endpointVerb == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_verb can not be blank", "endpoint_verb can not be blank")
	}
	endpointCode := c.FormValue("endpoint_code")
	if endpointCode == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_code can not be blank", "endpoint_code can not be blank")
	}

	endpointUrl := c.FormValue("endpoint_uri")
	if endpointUrl == "" {
		return lib.CustomError(http.StatusBadRequest, "endpoint_uri can not be blank", "endpoint_uri can not be blank")
	}

	menuKey := c.FormValue("menu_key")
	if menuKey == "" {
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	}
	params["endpoint_key"] = endpointKey
	params["endpoint_category_key"] = endpointCategoryKey
	params["endpoint_name"] = endpointName
	params["endpoint_verb"] = endpointVerb
	params["endpoint_uri"] = endpointUrl
	params["menu_key"] = menuKey
	params["rec_modified_by"] = lib.UserIDStr
	params["rec_modified_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_status"] = "1"
	params["endpoint_version"] = "1"
	params["privileges_key"] = "READ"

	err = models.UpdateEndpointSc(endpointKey, params)
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

func DeleteEndpointController(c echo.Context) error {
	params := make(map[string]string)
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_deleted_by"] = lib.UserIDStr

	endpointKey := c.FormValue("endpoint_key")
	if endpointKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing endpointKey", "Missing endpointKey")
	}

	err := models.DeleteEndpoint(endpointKey, params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "Berhasil hapus Endpoint!"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
