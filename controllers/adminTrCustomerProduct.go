package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func AdminGetProductListCanRedeem(c echo.Context) error {
	CustomerKey := c.QueryParam("customer_key")
	if CustomerKey == "" {
		return lib.CustomError(http.StatusBadRequest, "missing customer_key", "missing customer_key")
	}

	FundTypeKey := c.QueryParam("fund_type_key")

	rData := make([]interface{}, 0)

	var CustProdList []models.CustomerProductModel
	_, _ = models.GetCustomerProductList(&CustProdList, CustomerKey, FundTypeKey)

	if len(CustProdList) > 0 {
		for _, CustProd := range CustProdList {
			Customer_Key_Uint, _ := strconv.ParseUint(CustomerKey, 10, 64)
			var getbalance models.CustomerBalanceModel
			_, _ = models.GetCustomerBalance(&getbalance, Customer_Key_Uint, CustProd.ProductKey, lib.TIMENOW_TIMESTAMPFORMAT)

			maps1 := make(map[string]interface{})
			if getbalance.BalanceUnit.Cmp(decimal.Zero) > 0 {
				maps1["product_key"] = CustProd.ProductKey
				maps1["product_name"] = CustProd.ProductName
				rData = append(rData, maps1)
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = rData

	return c.JSON(http.StatusOK, response)
}
