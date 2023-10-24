package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func SaveProductUpdateRequest(c echo.Context) error {
	var err error
	params := make(map[string]string)

	// VALIDASI PARAMETER

	product_key := c.FormValue("product_key")
	if product_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_key")
	}
	params["product_key"] = product_key

	product_code := c.FormValue("product_code")
	if product_code == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_code")
	}
	params["product_code"] = product_code

	product_name := c.FormValue("product_name")
	if product_name == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_name")
	}
	params["product_name"] = product_name

	product_name_alt := c.FormValue("product_name_alt")
	if product_name_alt == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_name_alt")
	}
	params["product_name_alt"] = product_name_alt

	currency_key := c.FormValue("currency_key")
	if currency_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing currency_key")
	}
	params["currency_key"] = currency_key

	product_category_key := c.FormValue("product_category_key")
	if product_category_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_category_key")
	}
	params["product_category_key"] = product_category_key

	fund_type_key := c.FormValue("fund_type_key")
	if fund_type_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing fund_type_key")
	}
	params["fund_type_key"] = fund_type_key

	product_profile := c.FormValue("product_profile")
	if product_profile == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_profile")
	}
	params["product_profile"] = product_profile

	investment_objectives := c.FormValue("investment_objectives")
	if investment_objectives == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing investment_objectives")
	}
	params["investment_objectives"] = investment_objectives

	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_action"] = "REQUEST_UPDATE"
	err = models.CreateProductUpdateRequest(params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}

	return nil
}
