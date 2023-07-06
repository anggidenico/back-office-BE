package controllers

import (
	"database/sql"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

func AdminGetProductListCanRedeem(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true

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

				cek_in_proses := models.CekInProcessByUnit(CustProd.ProductKey, Customer_Key_Uint)
				avls := getbalance.BalanceUnit.Sub(cek_in_proses)
				var getNav models.NavModel
				status, err := models.GetLastNAV(&getNav, CustProd.ProductKey)
				if err != nil {
					if err == sql.ErrNoRows {
						return lib.CustomError(status, err.Error(), "Failed get nav", "Failed get nav")
					}
				}

				maps1["product_key"] = CustProd.ProductKey
				maps1["product_name"] = CustProd.ProductName
				// maps1["nav_price_per_unit"] = getNav.NavValue
				// maps1["balance"] = avls.Mul(getNav.NavValue)
				maps1["available_units"] = avls
				maps1["min_redeem_amount"] = CustProd.MinRedAmount
				maps1["risk_profile_name"] = CustProd.RiskProfileName
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
