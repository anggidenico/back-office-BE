package models

import (
	"log"
	"mf-bo-api/db"
	"strconv"

	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	RecPK               *uint64 `db:"rec_pk" json:"rec_pk"`
	RecAction           *string `db:"rec_action" json:"rec_action"`
	ProductKey          *uint64 `db:"product_key"             json:"product_key"`
	ProductID           *uint64 `db:"product_id"              json:"product_id"`
	ProductCode         *string `db:"product_code"            json:"product_code"`
	ProductName         *string `db:"product_name"            json:"product_name"`
	ProductNameAlt      *string `db:"product_name_alt"        json:"product_name_alt"`
	CurrencyKey         *uint64 `db:"currency_key"            json:"currency_key"`
	CurrencyName        *string `db:"currency_name" json:"currency_name"`
	ProductCategoryKey  *uint64 `db:"product_category_key" json:"product_category_key"`
	ProductCategoryName *string `db:"product_category_name" json:"product_category_name"`
	// ProductTypeKey        *uint64          `db:"product_type_key"        json:"product_type_key"`
	ProductTypeName *string `db:"product_type_name"       json:"product_type_name"`
	FundTypeKey     *uint64 `db:"fund_type_key"           json:"fund_type_key"`
	FundTypeName    *string `db:"fund_type_name"           json:"fund_type_name"`
	// FundStructureKey      *uint64          `db:"fund_structure_key"      json:"fund_structure_key"`
	// RiskProfileKey        *uint64          `db:"risk_profile_key"        json:"risk_profile_key"`
	ProductProfile       *string          `db:"product_profile"         json:"product_profile"`
	InvestmentObjectives *string          `db:"investment_objectives"   json:"investment_objectives"`
	ProductPhase         *uint64          `db:"product_phase"           json:"product_phase"`
	NavValuationType     *uint64          `db:"nav_valuation_type"      json:"nav_valuation_type"`
	ProspectusLink       *string          `db:"prospectus_link"         json:"prospectus_link"`
	LaunchDate           *string          `db:"launch_date"             json:"launch_date"`
	InceptionDate        *string          `db:"inception_date"          json:"inception_date"`
	IsinCode             *string          `db:"isin_code"               json:"isin_code"`
	FlagSyariah          *uint8           `db:"flag_syariah"            json:"flag_syariah"`
	MaxSubFee            *decimal.Decimal `db:"max_sub_fee"             json:"max_sub_fee"`
	MaxRedFee            *decimal.Decimal `db:"max_red_fee"             json:"max_red_fee"`
	MaxSwiFee            *decimal.Decimal `db:"max_swi_fee"             json:"max_swi_fee"`
	MinSubAmount         *decimal.Decimal `db:"min_sub_amount"          json:"min_sub_amount"`
	MinTopUpAmount       *decimal.Decimal `db:"min_topup_amount"        json:"min_topup_amount"`
	MinRedAmount         *decimal.Decimal `db:"min_red_amount"          json:"min_red_amount"`
	MinRedUnit           *decimal.Decimal `db:"min_red_unit"            json:"min_red_unit"`
	MinUnitAfterRed      *decimal.Decimal `db:"min_unit_after_red"      json:"min_unit_after_red"`
	MinAmountAfterRed    *decimal.Decimal `db:"min_amount_after_red"    json:"min_amount_after_red"`
	ManagementFee        *decimal.Decimal `db:"management_fee"          json:"management_fee"`
	CustodianFee         *decimal.Decimal `db:"custodian_fee"           json:"custodian_fee"`
	CustodianKey         *uint64          `db:"custodian_key"           json:"custodian_key"`
	// OjkFee                *decimal.Decimal `db:"ojk_fee"                 json:"ojk_fee"`
	// ProductFeeAmount      *decimal.Decimal `db:"product_fee_amount"      json:"product_fee_amount"`
	// OverwriteTransactFlag *uint8           `db:"overwrite_transact_flag" json:"overwrite_transact_flag"`
	// OverwriteFeeFlag      *uint8           `db:"overwrite_fee_flag"      json:"overwrite_fee_flag"`
	// OtherFeeAmount        *decimal.Decimal `db:"other_fee_amount"        json:"other_fee_amount"`
	SettlementPeriod *uint64 `db:"settlement_period"       json:"settlement_period"`
	SinvestFundCode  *string `db:"sinvest_fund_code"       json:"sinvest_fund_code"`
	FlagEnabled      *uint8  `db:"flag_enabled"            json:"flag_enabled"`
	FlagSubscription *uint8  `db:"flag_subscription"       json:"flag_subscription"`
	FlagRedemption   *uint8  `db:"flag_redemption"         json:"flag_redemption"`
	FlagSwitchOut    *uint8  `db:"flag_switch_out"         json:"flag_switch_out"`
	FlagSwitchIn     *uint8  `db:"flag_switch_in"          json:"flag_switch_in"`
	DecUnit          *uint8  `db:"dec_unit" json:"dec_unit"`
	DecAmount        *uint8  `db:"dec_amount" json:"dec_amount"`
	DecNav           *uint8  `db:"dec_nav" json:"dec_nav"`
	DecPerformance   *uint8  `db:"dec_performance" json:"dec_performance"`
	NpwpDateReg      *string `db:"npwp_date_reg" json:"npwp_date_reg"`
	NpwpName         *string `db:"npwp_name" json:"npwp_name"`
	NpwpNumber       *string `db:"npwp_number" json:"npwp_number"`
	PortfolioCode    *string `db:"portfolio_code" json:"portfolio_code"`
	RecCreatedDate   *string `db:"rec_created_date" json:"rec_created_date"`
}

func GetProductRequestList() (result []ProductRequest) {
	query := `SELECT rec_pk, rec_action, product_key, product_name, product_code, rec_created_date
	FROM ms_product_request WHERE rec_status = 1`
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

type ProductUpdateDetails struct {
	Existing ProductRequest `json:"existing"`
	Updates  ProductRequest `json:"updates"`
}

func GetProductRequestDetail(RecPK string) ProductUpdateDetails {
	var result ProductUpdateDetails
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	// GET UPDATES DATA
	query := `SELECT rec_pk,rec_action, product_key, product_id, product_code, product_name, product_name_alt,
	currency_key, product_category_key, fund_type_key, product_profile, investment_objectives, product_phase, nav_valuation_type, prospectus_link, launch_date, inception_date, isin_code, flag_syariah, max_sub_fee, max_red_fee, max_swi_fee, min_sub_amount, min_topup_amount, min_red_amount, min_red_amount, min_red_unit, min_unit_after_red, min_amount_after_red, management_fee, custodian_fee, custodian_key, settlement_period, sinvest_fund_code, flag_enabled, flag_subscription, flag_redemption, flag_redemption, flag_switch_out, flag_switch_in, dec_unit, dec_amount, dec_nav, dec_performance, npwp_date_reg, npwp_name, npwp_number, portfolio_code, rec_created_date
	FROM ms_product_request WHERE rec_status = 1 AND rec_pk = ` + RecPK

	row := tx.QueryRow(query)
	err = row.Scan(&result.Updates.RecPK, &result.Updates.RecAction, &result.Updates.ProductKey, &result.Updates.ProductID, &result.Updates.ProductCode, &result.Updates.ProductName, &result.Updates.ProductNameAlt, &result.Updates.CurrencyKey, &result.Updates.ProductCategoryKey, &result.Updates.FundTypeKey, &result.Updates.ProductProfile, &result.Updates.InvestmentObjectives, &result.Updates.ProductPhase, &result.Updates.NavValuationType, &result.Updates.ProspectusLink, &result.Updates.LaunchDate, &result.Updates.InceptionDate, &result.Updates.IsinCode, &result.Updates.FlagSyariah, &result.Updates.MaxSubFee, &result.Updates.MaxRedFee, &result.Updates.MaxSwiFee, &result.Updates.MinSubAmount, &result.Updates.MinTopUpAmount, &result.Updates.MinRedAmount, &result.Updates.MinRedAmount, &result.Updates.MinRedUnit, &result.Updates.MinUnitAfterRed, &result.Updates.MinAmountAfterRed, &result.Updates.ManagementFee, &result.Updates.CustodianFee, &result.Updates.CustodianKey, &result.Updates.SettlementPeriod, &result.Updates.SinvestFundCode, &result.Updates.FlagEnabled, &result.Updates.FlagSubscription, &result.Updates.FlagRedemption, &result.Updates.FlagRedemption, &result.Updates.FlagSwitchOut, &result.Updates.FlagSwitchIn, &result.Updates.DecUnit, &result.Updates.DecAmount, &result.Updates.DecNav, &result.Updates.DecPerformance, &result.Updates.NpwpDateReg, &result.Updates.NpwpName, &result.Updates.NpwpNumber, &result.Updates.PortfolioCode, &result.Updates.RecCreatedDate)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	// GET EXISTING DATA
	query2 := `SELECT product_key, product_id, product_code, product_name, product_name_alt,
	currency_key, product_category_key, fund_type_key, product_profile, investment_objectives, product_phase, nav_valuation_type, prospectus_link, launch_date, inception_date, isin_code, flag_syariah, max_sub_fee, max_red_fee, max_swi_fee, min_sub_amount, min_topup_amount, min_red_amount, min_red_amount, min_red_unit, min_unit_after_red, min_amount_after_red, management_fee, custodian_fee, custodian_key, settlement_period, sinvest_fund_code, flag_enabled, flag_subscription, flag_redemption,flag_redemption, flag_switch_out, flag_switch_in, dec_unit, dec_amount, dec_nav, dec_performance, npwp_date_reg, npwp_name, npwp_number, portfolio_code
	FROM ms_product_request WHERE rec_status = 1 AND product_key = ` + strconv.FormatUint(result.Updates.ProductKey, 10)
	row2 := tx.QueryRow(query2)
	err = row2.Scan(&result.Existing.ProductKey, &result.Existing.ProductID, &result.Existing.ProductCode,
		&result.Existing.ProductName, &result.Existing.ProductNameAlt, &result.Existing.CurrencyKey,
		&result.Existing.ProductCategoryKey, &result.Existing.FundTypeKey, &result.Existing.ProductProfile,
		&result.Existing.InvestmentObjectives, &result.Existing.ProductPhase, &result.Existing.NavValuationType,
		&result.Existing.ProspectusLink, &result.Existing.LaunchDate, &result.Existing.InceptionDate,
		&result.Existing.IsinCode, &result.Existing.FlagSyariah, &result.Existing.MaxSubFee,
		&result.Existing.MaxRedFee, &result.Existing.MaxSwiFee, &result.Existing.MinSubAmount,
		&result.Existing.MinTopUpAmount, &result.Existing.MinRedAmount, &result.Existing.MinRedAmount, &result.Existing.MinRedUnit, &result.Existing.MinUnitAfterRed, &result.Existing.MinAmountAfterRed, &result.Existing.ManagementFee, &result.Existing.CustodianFee, &result.Existing.CustodianKey, &result.Existing.SettlementPeriod, &result.Existing.SinvestFundCode, &result.Existing.FlagEnabled, &result.Existing.FlagSubscription, &result.Existing.FlagRedemption, &result.Existing.FlagRedemption, &result.Existing.FlagSwitchOut, &result.Existing.FlagSwitchIn, &result.Existing.DecUnit, &result.Existing.DecAmount, &result.Existing.DecNav, &result.Existing.DecPerformance, &result.Existing.NpwpDateReg, &result.Existing.NpwpName, &result.Existing.NpwpNumber, &result.Existing.PortfolioCode)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err.Error())
	}

	return result
}

func CreateProductRequest(params map[string]string) error {
	query := "INSERT INTO ms_product_request"
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += ` "` + value + `", `
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	query += "(" + fields + ") VALUES(" + values + ")"

	// log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
