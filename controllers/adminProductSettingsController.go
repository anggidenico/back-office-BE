package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

	"github.com/labstack/echo"
)

// func SaveProductSetting2(c echo.Context) (err error) {
// 	m := new(models.MasterProduct)
// 	if err := c.Bind(m); err != nil {
// 		return err
// 	}
// 	if err := c.Validate(m); err != nil {
// 		return err
// 	}

// 	var response lib.Response
// 	response.Status.Code = http.StatusOK
// 	response.Status.MessageServer = "OK"
// 	response.Status.MessageClient = "OK"
// 	response.Data = nil
// 	return c.JSON(http.StatusOK, response)
// }

func SaveMasterProduct(c echo.Context) (err error) {
	insertMsProduct := make(map[string]string)

	productName := c.FormValue("product_name")
	if productName == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_name", "Missing product_name")
	} else {
		insertMsProduct["product_name"] = productName
	}

	productNameAlt := c.FormValue("product_name_alt")
	if productNameAlt == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_name_alt", "Missing product_name_alt")
	} else {
		insertMsProduct["product_name_alt"] = productNameAlt
	}

	productCode := c.FormValue("product_code")
	if productCode == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_code", "Missing product_code")
	} else {
		insertMsProduct["product_code"] = productCode
	}

	currencyKey := c.FormValue("currency_key")
	if currencyKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing currency_key", "Missing currency_key")
	} else {
		insertMsProduct["currency_key"] = currencyKey
	}

	productCategoryKey := c.FormValue("product_category_key")
	if productCategoryKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_category_key", "Missing product_category_key")
	} else {
		insertMsProduct["product_category_key"] = productCategoryKey
	}

	productTypeKey := c.FormValue("product_type_key")
	if productTypeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_type_key", "Missing product_type_key")
	} else {
		insertMsProduct["product_type_key"] = productTypeKey
	}

	fundTypeKey := c.FormValue("fund_type_key")
	if fundTypeKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing fund_type_key", "Missing fund_type_key")
	} else {
		insertMsProduct["fund_type_key"] = fundTypeKey
	}

	fundStructureKey := c.FormValue("fund_structure_key")
	if fundStructureKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing fund_structure_key", "Missing fund_structure_key")
	} else {
		insertMsProduct["fund_structure_key"] = fundStructureKey
	}

	riskProfileKey := c.FormValue("risk_profile_key")
	if riskProfileKey == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing risk_profile_key", "Missing risk_profile_key")
	} else {
		insertMsProduct["risk_profile_key"] = riskProfileKey
	}

	flagSyariah := c.FormValue("flag_syariah")
	if flagSyariah == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_syariah", "Missing flag_syariah")
	} else {
		if flagSyariah == "true" {
			insertMsProduct["flag_syariah"] = "1"
		} else {
			insertMsProduct["flag_syariah"] = "0"
		}
	}

	maxSubFee := c.FormValue("max_sub_fee")
	if maxSubFee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_sub_fee", "Missing max_sub_fee")
	} else {
		insertMsProduct["max_sub_fee"] = maxSubFee
	}

	maxRedFee := c.FormValue("max_red_fee")
	if maxRedFee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_red_fee", "Missing max_red_fee")
	} else {
		insertMsProduct["max_red_fee"] = maxRedFee
	}

	maxSwiFee := c.FormValue("max_swi_fee")
	if maxSwiFee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_swi_fee", "Missing max_swi_fee")
	} else {
		insertMsProduct["max_swi_fee"] = maxSwiFee
	}

	minSubAmount := c.FormValue("min_sub_amount")
	if minSubAmount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_sub_amount", "Missing min_sub_amount")
	} else {
		insertMsProduct["min_sub_amount"] = minSubAmount
	}

	minTopUpAmount := c.FormValue("min_topup_amount")
	if minTopUpAmount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_topup_amount", "Missing min_topup_amount")
	} else {
		insertMsProduct["min_topup_amount"] = minTopUpAmount
	}

	minRedAmount := c.FormValue("min_red_amount")
	if minRedAmount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_amount", "Missing min_red_amount")
	} else {
		insertMsProduct["min_red_amount"] = minRedAmount
	}

	minRedunit := c.FormValue("min_red_unit")
	if minRedunit == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_unit", "Missing min_red_unit")
	} else {
		insertMsProduct["min_red_unit"] = minRedunit
	}

	minUnitAfterRed := c.FormValue("min_unit_after_red")
	if minUnitAfterRed == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_unit_after_red", "Missing min_unit_after_red")
	} else {
		insertMsProduct["min_unit_after_red"] = minUnitAfterRed
	}

	minamountAfterRed := c.FormValue("min_amount_after_red")
	if minamountAfterRed == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_amount_after_red", "Missing min_amount_after_red")
	} else {
		insertMsProduct["min_amount_after_red"] = minamountAfterRed
	}

	managementFee := c.FormValue("management_fee")
	if managementFee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing management_fee", "Missing management_fee")
	} else {
		insertMsProduct["management_fee"] = managementFee
	}

	custodianFee := c.FormValue("custodian_fee")
	if custodianFee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing custodian_fee", "Missing custodian_fee")
	} else {
		insertMsProduct["custodian_fee"] = custodianFee
	}

	overwriteTransactFlag := c.FormValue("overwrite_transact_flag")
	if overwriteTransactFlag == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing overwrite_transact_flag", "Missing overwrite_transact_flag")
	} else {
		insertMsProduct["overwrite_transact_flag"] = overwriteTransactFlag
	}

	overwriteFeeFlag := c.FormValue("overwrite_fee_flag")
	if overwriteFeeFlag == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing overwrite_fee_flag", "Missing overwrite_fee_flag")
	} else {
		insertMsProduct["overwrite_fee_flag"] = overwriteFeeFlag
	}

	otherFeeAmount := c.FormValue("other_fee_amount")
	if otherFeeAmount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing other_fee_amount", "Missing other_fee_amount")
	} else {
		insertMsProduct["other_fee_amount"] = otherFeeAmount
	}

	flagEnabled := c.FormValue("flag_enabled")
	if flagEnabled == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_enabled", "Missing flag_enabled")
	} else {
		insertMsProduct["flag_enabled"] = flagEnabled
	}

	flagSubscription := c.FormValue("flag_subscription")
	if flagSubscription == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_subscription", "Missing flag_subscription")
	} else {
		insertMsProduct["flag_subscription"] = flagSubscription
	}

	flagRedemption := c.FormValue("flag_redemption")
	if flagRedemption == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_redemption", "Missing flag_redemption")
	} else {
		insertMsProduct["flag_redemption"] = flagRedemption
	}

	flagSwitchOut := c.FormValue("flag_switch_out")
	if flagSwitchOut == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_switch_out", "Missing flag_switch_out")
	} else {
		insertMsProduct["flag_switch_out"] = flagSwitchOut
	}

	flagSwitchIn := c.FormValue("flag_switch_in")
	if flagSwitchIn == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing flag_switch_in", "Missing flag_switch_in")
	} else {
		insertMsProduct["flag_switch_in"] = flagSwitchIn
	}

	decNav := c.FormValue("dec_nav")
	if decNav == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing dec_nav", "Missing dec_nav")
	} else {
		insertMsProduct["dec_nav"] = decNav
	}

	decunit := c.FormValue("dec_unit")
	if decunit == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing dec_unit", "Missing dec_unit")
	} else {
		insertMsProduct["dec_unit"] = decNav
	}

	decperformance := c.FormValue("dec_performance")
	if decperformance == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing dec_performance", "Missing dec_performance")
	} else {
		insertMsProduct["dec_performance"] = decperformance
	}

	decamount := c.FormValue("dec_amount")
	if decamount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing dec_amount", "Missing dec_amount")
	} else {
		insertMsProduct["dec_amount"] = decamount
	}

	productProfile := c.FormValue("product_profile")
	insertMsProduct["product_profile"] = productProfile

	investmentObjectives := c.FormValue("investment_objectives")
	insertMsProduct["investment_objectives"] = investmentObjectives

	productPhase := c.FormValue("product_phase")
	if productPhase == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_phase", "Missing product_phase")
	} else {
		insertMsProduct["product_phase"] = productPhase
	}

	navValuationType := c.FormValue("nav_valuation_type")
	if navValuationType == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing nav_valuation_type", "Missing nav_valuation_type")
	} else {
		insertMsProduct["nav_valuation_type"] = navValuationType
	}

	insertMsProduct["rec_status"] = "1"
	err = models.InsertMasterProduct(insertMsProduct)
	if err != nil {
		return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil

	return c.JSON(http.StatusOK, response)
}
