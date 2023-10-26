package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func ProductCreateRequest(c echo.Context) error {
	var err error
	params := make(map[string]string)

	// VALIDASI PARAMETER

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

	product_phase := c.FormValue("product_phase")
	if product_phase == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_phase")
	}
	params["product_phase"] = product_phase

	nav_valuation_type := c.FormValue("nav_valuation_type")
	if nav_valuation_type == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing nav_valuation_type")
	}
	params["nav_valuation_type"] = nav_valuation_type

	prospectus_link := c.FormValue("prospectus_link")
	if prospectus_link == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing prospectus_link")
	}
	params["prospectus_link"] = prospectus_link

	launch_date := c.FormValue("launch_date")
	if launch_date == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing launch_date")
	}
	params["launch_date"] = launch_date

	inception_date := c.FormValue("inception_date")
	if inception_date == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing inception_date")
	}
	params["inception_date"] = inception_date

	isin_code := c.FormValue("isin_code")
	if isin_code == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing isin_code")
	}
	params["isin_code"] = isin_code

	max_sub_fee := c.FormValue("max_sub_fee")
	if max_sub_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_sub_fee")
	}
	params["max_sub_fee"] = max_sub_fee

	max_red_fee := c.FormValue("max_red_fee")
	if max_red_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_red_fee")
	}
	params["max_red_fee"] = max_red_fee

	max_swi_fee := c.FormValue("max_swi_fee")
	if max_swi_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_swi_fee")
	}
	params["max_swi_fee"] = max_swi_fee

	min_sub_amount := c.FormValue("min_sub_amount")
	if min_sub_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_sub_amount")
	}
	params["min_sub_amount"] = min_sub_amount

	min_topup_amount := c.FormValue("min_topup_amount")
	if min_topup_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_topup_amount")
	}
	params["min_topup_amount"] = min_topup_amount

	min_red_amount := c.FormValue("min_red_amount")
	if min_red_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_amount")
	}
	params["min_red_amount"] = min_red_amount

	min_red_unit := c.FormValue("min_red_unit")
	if min_red_unit == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_unit")
	}
	params["min_red_unit"] = min_red_unit

	min_unit_after_red := c.FormValue("min_unit_after_red")
	if min_unit_after_red == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_unit_after_red")
	}
	params["min_unit_after_red"] = min_unit_after_red

	min_amount_after_red := c.FormValue("min_amount_after_red")
	if min_amount_after_red == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_amount_after_red")
	}
	params["min_amount_after_red"] = min_amount_after_red

	management_fee := c.FormValue("management_fee")
	if management_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing management_fee")
	}
	params["management_fee"] = management_fee

	custodian_fee := c.FormValue("custodian_fee")
	if custodian_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing custodian_fee")
	}
	params["custodian_fee"] = custodian_fee

	custodian_key := c.FormValue("custodian_key")
	if custodian_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing custodian_key")
	}
	params["custodian_key"] = custodian_key

	settlement_period := c.FormValue("settlement_period")
	if settlement_period == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing settlement_period")
	}
	params["settlement_period"] = settlement_period

	sinvest_fund_code := c.FormValue("sinvest_fund_code")
	if sinvest_fund_code == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sinvest_fund_code")
	}
	params["sinvest_fund_code"] = sinvest_fund_code

	riskProfileKey := models.GetRiskProfileByFundType(fund_type_key)
	params["risk_profile_key"] = strconv.FormatUint(riskProfileKey, 10)

	digit_dec := "4"

	params["dec_nav"] = digit_dec
	params["dec_performance"] = digit_dec
	params["dec_unit"] = digit_dec
	params["dec_amount"] = digit_dec

	params["flag_enabled"] = "1"
	params["flag_subscription"] = "1"
	params["flag_redemption"] = "1"
	params["flag_switch_out"] = "1"
	params["flag_switch_in"] = "1"

	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_action"] = "CREATE"
	err = models.CreateProductRequest(params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ProductUpdateRequest(c echo.Context) error {
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

	product_phase := c.FormValue("product_phase")
	if product_phase == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing product_phase")
	}
	params["product_phase"] = product_phase

	nav_valuation_type := c.FormValue("nav_valuation_type")
	if nav_valuation_type == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing nav_valuation_type")
	}
	params["nav_valuation_type"] = nav_valuation_type

	prospectus_link := c.FormValue("prospectus_link")
	if prospectus_link == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing prospectus_link")
	}
	params["prospectus_link"] = prospectus_link

	launch_date := c.FormValue("launch_date")
	if launch_date == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing launch_date")
	}
	params["launch_date"] = launch_date

	inception_date := c.FormValue("inception_date")
	if inception_date == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing inception_date")
	}
	params["inception_date"] = inception_date

	isin_code := c.FormValue("isin_code")
	if isin_code == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing isin_code")
	}
	params["isin_code"] = isin_code

	max_sub_fee := c.FormValue("max_sub_fee")
	if max_sub_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_sub_fee")
	}
	params["max_sub_fee"] = max_sub_fee

	max_red_fee := c.FormValue("max_red_fee")
	if max_red_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_red_fee")
	}
	params["max_red_fee"] = max_red_fee

	max_swi_fee := c.FormValue("max_swi_fee")
	if max_swi_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing max_swi_fee")
	}
	params["max_swi_fee"] = max_swi_fee

	min_sub_amount := c.FormValue("min_sub_amount")
	if min_sub_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_sub_amount")
	}
	params["min_sub_amount"] = min_sub_amount

	min_topup_amount := c.FormValue("min_topup_amount")
	if min_topup_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_topup_amount")
	}
	params["min_topup_amount"] = min_topup_amount

	min_red_amount := c.FormValue("min_red_amount")
	if min_red_amount == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_amount")
	}
	params["min_red_amount"] = min_red_amount

	min_red_unit := c.FormValue("min_red_unit")
	if min_red_unit == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_red_unit")
	}
	params["min_red_unit"] = min_red_unit

	min_unit_after_red := c.FormValue("min_unit_after_red")
	if min_unit_after_red == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_unit_after_red")
	}
	params["min_unit_after_red"] = min_unit_after_red

	min_amount_after_red := c.FormValue("min_amount_after_red")
	if min_amount_after_red == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing min_amount_after_red")
	}
	params["min_amount_after_red"] = min_amount_after_red

	management_fee := c.FormValue("management_fee")
	if management_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing management_fee")
	}
	params["management_fee"] = management_fee

	custodian_fee := c.FormValue("custodian_fee")
	if custodian_fee == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing custodian_fee")
	}
	params["custodian_fee"] = custodian_fee

	custodian_key := c.FormValue("custodian_key")
	if custodian_key == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing custodian_key")
	}
	params["custodian_key"] = custodian_key

	settlement_period := c.FormValue("settlement_period")
	if settlement_period == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing settlement_period")
	}
	params["settlement_period"] = settlement_period

	sinvest_fund_code := c.FormValue("sinvest_fund_code")
	if sinvest_fund_code == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing sinvest_fund_code")
	}
	params["sinvest_fund_code"] = sinvest_fund_code

	dec_nav := c.FormValue("dec_nav")
	if dec_nav != "" {
		params["dec_nav"] = dec_nav
	}

	dec_performance := c.FormValue("dec_performance")
	if dec_performance != "" {
		params["dec_performance"] = dec_performance
	}

	dec_unit := c.FormValue("dec_unit")
	if dec_unit != "" {
		params["dec_unit"] = dec_unit
	}

	dec_amount := c.FormValue("dec_amount")
	if dec_amount != "" {
		params["dec_amount"] = dec_amount
	}

	flag_syariah := c.FormValue("flag_syariah")
	if flag_syariah != "" {
		if flag_syariah == "true" {
			params["flag_syariah"] = "1"
		} else {
			params["flag_syariah"] = "0"
		}
	}

	flag_enabled := c.FormValue("flag_enabled")
	if flag_enabled != "" {
		if flag_enabled == "true" {
			params["flag_enabled"] = "1"
		} else {
			params["flag_enabled"] = "0"
		}
	}

	flag_subscription := c.FormValue("flag_subscription")
	if flag_subscription != "" {
		if flag_subscription == "true" {
			params["flag_subscription"] = "1"
		} else {
			params["flag_subscription"] = "0"
		}
	}

	flag_redemption := c.FormValue("flag_redemption")
	if flag_redemption != "" {
		if flag_redemption == "true" {
			params["flag_redemption"] = "1"
		} else {
			params["flag_redemption"] = "0"
		}
	}

	flag_switch_out := c.FormValue("flag_switch_out")
	if flag_switch_out != "" {
		if flag_switch_out == "true" {
			params["flag_switch_out"] = "1"
		} else {
			params["flag_switch_out"] = "0"
		}
	}

	flag_switch_in := c.FormValue("flag_switch_in")
	if flag_switch_in != "" {
		if flag_switch_in == "true" {
			params["flag_switch_in"] = "1"
		} else {
			params["flag_switch_in"] = "0"
		}
	}

	riskProfileKey := models.GetRiskProfileByFundType(fund_type_key)
	params["risk_profile_key"] = strconv.FormatUint(riskProfileKey, 10)

	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(lib.TIMESTAMPFORMAT)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_action"] = "UPDATE"
	err = models.CreateProductRequest(params)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), err.Error())
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ProductApprovalList(c echo.Context) error {

	result := models.GetProductRequestList()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result
	return c.JSON(http.StatusOK, response)
}
