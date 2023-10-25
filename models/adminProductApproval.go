package models

import (
	"log"
	"mf-bo-api/db"

	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	RecPK                 *uint64          `db:"rec_pk" json:"rec_pk"`
	RecAction             *string          `db:"rec_action" json:"rec_action"`
	ProductKey            uint64           `db:"product_key"             json:"product_key"`
	ProductID             *uint64          `db:"product_id"              json:"product_id"`
	ProductCode           *string          `db:"product_code"            json:"product_code"`
	ProductName           *string          `db:"product_name"            json:"product_name"`
	ProductNameAlt        *string          `db:"product_name_alt"        json:"product_name_alt"`
	CurrencyKey           *uint64          `db:"currency_key"            json:"currency_key"`
	ProductCategoryKey    *uint64          `db:"product_category_key"    json:"product_category_key"`
	ProductTypeKey        *uint64          `db:"product_type_key"        json:"product_type_key"`
	FundTypeKey           *uint64          `db:"fund_type_key"           json:"fund_type_key"`
	FundStructureKey      *uint64          `db:"fund_structure_key"      json:"fund_structure_key"`
	RiskProfileKey        *uint64          `db:"risk_profile_key"        json:"risk_profile_key"`
	ProductProfile        *string          `db:"product_profile"         json:"product_profile"`
	InvestmentObjectives  *string          `db:"investment_objectives"   json:"investment_objectives"`
	ProductPhase          *uint64          `db:"product_phase"           json:"product_phase"`
	NavValuationType      *uint64          `db:"nav_valuation_type"      json:"nav_valuation_type"`
	ProspectusLink        *string          `db:"prospectus_link"         json:"prospectus_link"`
	LaunchDate            *string          `db:"launch_date"             json:"launch_date"`
	InceptionDate         *string          `db:"inception_date"          json:"inception_date"`
	IsinCode              *string          `db:"isin_code"               json:"isin_code"`
	FlagSyariah           *uint8           `db:"flag_syariah"            json:"flag_syariah"`
	MaxSubFee             *decimal.Decimal `db:"max_sub_fee"             json:"max_sub_fee"`
	MaxRedFee             *decimal.Decimal `db:"max_red_fee"             json:"max_red_fee"`
	MaxSwiFee             *decimal.Decimal `db:"max_swi_fee"             json:"max_swi_fee"`
	MinSubAmount          *decimal.Decimal `db:"min_sub_amount"          json:"min_sub_amount"`
	MinTopUpAmount        *decimal.Decimal `db:"min_topup_amount"        json:"min_topup_amount"`
	MinRedAmount          *decimal.Decimal `db:"min_red_amount"          json:"min_red_amount"`
	MinRedUnit            *decimal.Decimal `db:"min_red_unit"            json:"min_red_unit"`
	MinUnitAfterRed       *decimal.Decimal `db:"min_unit_after_red"      json:"min_unit_after_red"`
	MinAmountAfterRed     *decimal.Decimal `db:"min_amount_after_red"    json:"min_amount_after_red"`
	ManagementFee         *decimal.Decimal `db:"management_fee"          json:"management_fee"`
	CustodianFee          *decimal.Decimal `db:"custodian_fee"           json:"custodian_fee"`
	CustodianKey          *uint64          `db:"custodian_key"           json:"custodian_key"`
	OjkFee                *decimal.Decimal `db:"ojk_fee"                 json:"ojk_fee"`
	ProductFeeAmount      *decimal.Decimal `db:"product_fee_amount"      json:"product_fee_amount"`
	OverwriteTransactFlag *uint8           `db:"overwrite_transact_flag" json:"overwrite_transact_flag"`
	OverwriteFeeFlag      *uint8           `db:"overwrite_fee_flag"      json:"overwrite_fee_flag"`
	OtherFeeAmount        *decimal.Decimal `db:"other_fee_amount"        json:"other_fee_amount"`
	SettlementPeriod      *uint64          `db:"settlement_period"       json:"settlement_period"`
	SinvestFundCode       *string          `db:"sinvest_fund_code"       json:"sinvest_fund_code"`
	FlagEnabled           *uint8           `db:"flag_enabled"            json:"flag_enabled"`
	FlagSubscription      *uint8           `db:"flag_subscription"       json:"flag_subscription"`
	FlagRedemption        *uint8           `db:"flag_redemption"         json:"flag_redemption"`
	FlagSwitchOut         *uint8           `db:"flag_switch_out"         json:"flag_switch_out"`
	FlagSwitchIn          *uint8           `db:"flag_switch_in"          json:"flag_switch_in"`
	DecUnit               *uint8           `db:"dec_unit" json:"dec_unit"`
	DecAmount             *uint8           `db:"dec_amount" json:"dec_amount"`
	DecNav                *uint8           `db:"dec_nav" json:"dec_nav"`
	DecPerformance        *uint8           `db:"dec_performance" json:"dec_performance"`
	NpwpDateReg           *string          `db:"npwp_date_reg" json:"npwp_date_reg"`
	NpwpName              *string          `db:"npwp_name" json:"npwp_name"`
	NpwpNumber            *string          `db:"npwp_number" json:"npwp_number"`
	PortfolioCode         *string          `db:"portfolio_code" json:"portfolio_code"`
}

func GetProductRequestList() (result []ProductRequest) {
	query := `SELECT rec_pk, rec_action, product_key, product_name, fund_type_key
	FROM ms_product_request WHERE rec_status = 1`
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return

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
