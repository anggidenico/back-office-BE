package models

import (
	"log"
	"mf-bo-api/db"

	"github.com/shopspring/decimal"
)

type MasterProduct struct {
	ProductKey            *uint64          `validate:"required" db:"product_key"    json:"product_key"      gorm:"column:product_key;primary_key;autoIncrement:true"`                 /*  */
	ProductId             *string          `validate:"" db:"product_id"    json:"product_id" gorm:"column:product_id;type:int(11)"`                                                   /*  */
	ProductCode           *string          `validate:"required" db:"product_code"    json:"product_code" gorm:"column:product_code;type:varchar(30)"`                                 /*  */
	ProductName           *string          `validate:"required" db:"product_name"    json:"product_name" gorm:"column:product_name;type:varchar(150)"`                                /*  */
	ProductNameAlt        *string          `validate:"required" db:"product_name_alt"    json:"product_name_alt" gorm:"column:product_name_alt;type:varchar(150)"`                    /*  */
	CurrencyKey           *uint64          `validate:"number" db:"currency_key"    json:"currency_key" gorm:"column:currency_key;type:int(11)"`                                       /*  */
	ProductCategoryKey    *uint64          `validate:"number" db:"product_category_key"    json:"product_category_key" gorm:"column:product_category_key;type:int(11)"`               /*  */
	ProductTypeKey        *uint64          `validate:"number" db:"product_type_key"    json:"product_type_key" gorm:"column:product_type_key;type:int(11)"`                           /*  */
	FundTypeKey           *uint64          `validate:"number" db:"fund_type_key"    json:"fund_type_key" gorm:"column:fund_type_key;type:int(11)"`                                    /*  */
	FundStructureKey      *uint64          `validate:"number" db:"fund_structure_key"    json:"fund_structure_key" gorm:"column:fund_structure_key;type:int(11)"`                     /*  */
	RiskProfileKey        *uint64          `validate:"number" db:"risk_profile_key"    json:"risk_profile_key" gorm:"column:risk_profile_key;type:int(11)"`                           /*  */
	ProductProfile        *string          `validate:"" db:"product_profile"    json:"product_profile" gorm:"column:product_profile;type:text"`                                       /*  */
	InvestmentObjectives  *string          `validate:"" db:"investment_objectives"    json:"investment_objectives" gorm:"column:investment_objectives;type:text"`                     /*  */
	ProductPhase          *string          `validate:"" db:"product_phase"    json:"product_phase" gorm:"column:product_phase;type:int(11)"`                                          /*  */
	NavValuationType      *string          `validate:"" db:"nav_valuation_type"    json:"nav_valuation_type" gorm:"column:nav_valuation_type;type:int(11)"`                           /*  */
	ProspectusLink        *string          `validate:"" db:"prospectus_link"    json:"prospectus_link" gorm:"column:prospectus_link;type:varchar(255)"`                               /*  */
	LaunchDate            *string          `validate:"" db:"launch_date"    json:"launch_date" gorm:"column:launch_date;type:datetime"`                                               /*  */
	InceptionDate         *string          `validate:"" db:"inception_date"    json:"inception_date" gorm:"column:inception_date;type:datetime"`                                      /*  */
	IsinCode              *string          `validate:"" db:"isin_code"    json:"isin_code" gorm:"column:isin_code;type:varchar(50)"`                                                  /*  */
	FlagSyariah           *uint64          `validate:"required" db:"flag_syariah"    json:"flag_syariah" gorm:"column:flag_syariah;type:tinyint(1)"`                                  /*  */
	MaxSubFee             *string          `validate:"required" db:"max_sub_fee"    json:"max_sub_fee" gorm:"column:max_sub_fee;type:decimal(9,4)"`                                   /*  */
	MaxRedFee             *string          `validate:"required" db:"max_red_fee"    json:"max_red_fee" gorm:"column:max_red_fee;type:decimal(9,4)"`                                   /*  */
	MaxSwiFee             *string          `validate:"required" db:"max_swi_fee"    json:"max_swi_fee" gorm:"column:max_swi_fee;type:decimal(9,4)"`                                   /*  */
	MinSubAmount          *decimal.Decimal `validate:"required" db:"min_sub_amount"    json:"min_sub_amount" gorm:"column:min_sub_amount;type:decimal(18,4)"`                         /*  */
	MinTopupAmount        *decimal.Decimal `validate:"" db:"min_topup_amount"    json:"min_topup_amount" gorm:"column:min_topup_amount;type:decimal(18,4)"`                           /*  */
	MinRedAmount          *decimal.Decimal `validate:"required" db:"min_red_amount"    json:"min_red_amount" gorm:"column:min_red_amount;type:decimal(18,4)"`                         /*  */
	MinRedUnit            *decimal.Decimal `validate:"required" db:"min_red_unit"    json:"min_red_unit" gorm:"column:min_red_unit;type:decimal(18,4)"`                               /*  */
	MinUnitAfterRed       *decimal.Decimal `validate:"" db:"min_unit_after_red"    json:"min_unit_after_red" gorm:"column:min_unit_after_red;type:decimal(18,4)"`                     /*  */
	MinAmountAfterRed     *decimal.Decimal `validate:"" db:"min_amount_after_red"    json:"min_amount_after_red" gorm:"column:min_amount_after_red;type:decimal(18,4)"`               /*  */
	ManagementFee         *decimal.Decimal `validate:"" db:"management_fee"    json:"management_fee" gorm:"column:management_fee;type:decimal(9,4)"`                                  /*  */
	CustodianFee          *decimal.Decimal `validate:"" db:"custodian_fee"    json:"custodian_fee" gorm:"column:custodian_fee;type:decimal(9,4)"`                                     /*  */
	CustodianKey          *uint64          `validate:"number" db:"custodian_key"    json:"custodian_key" gorm:"column:custodian_key;type:int(11)"`                                    /*  */
	OjkFee                *string          `validate:"" db:"ojk_fee"    json:"ojk_fee" gorm:"column:ojk_fee;type:decimal(9,4)"`                                                       /*  */
	ProductFeeAmount      *string          `validate:"" db:"product_fee_amount"    json:"product_fee_amount" gorm:"column:product_fee_amount;type:decimal(18,2)"`                     /*  */
	OverwriteTransactFlag *string          `validate:"required" db:"overwrite_transact_flag"    json:"overwrite_transact_flag" gorm:"column:overwrite_transact_flag;type:tinyint(4)"` /*  */
	OverwriteFeeFlag      *string          `validate:"required" db:"overwrite_fee_flag"    json:"overwrite_fee_flag" gorm:"column:overwrite_fee_flag;type:tinyint(4)"`                /*  */
	OtherFeeAmount        *decimal.Decimal `validate:"required" db:"other_fee_amount"    json:"other_fee_amount" gorm:"column:other_fee_amount;type:decimal(18,2)"`                   /*  */
	SettlementPeriod      *uint64          `validate:"" db:"settlement_period"    json:"settlement_period" gorm:"column:settlement_period;type:int(11)"`                              /*  */
	SinvestFundCode       *string          `validate:"" db:"sinvest_fund_code"    json:"sinvest_fund_code" gorm:"column:sinvest_fund_code;type:varchar(30)"`                          /*  */
	FlagEnabled           *uint64          `validate:"required" db:"flag_enabled"    json:"flag_enabled" gorm:"column:flag_enabled;type:tinyint(1)"`                                  /*  */
	FlagSubscription      *uint64          `validate:"required" db:"flag_subscription"    json:"flag_subscription" gorm:"column:flag_subscription;type:tinyint(1)"`                   /*  */
	FlagRedemption        *uint64          `validate:"required" db:"flag_redemption"    json:"flag_redemption" gorm:"column:flag_redemption;type:tinyint(1)"`                         /*  */
	FlagSwitchOut         *uint64          `validate:"required" db:"flag_switch_out"    json:"flag_switch_out" gorm:"column:flag_switch_out;type:tinyint(1)"`                         /*  */
	FlagSwitchIn          *uint64          `validate:"required" db:"flag_switch_in"    json:"flag_switch_in" gorm:"column:flag_switch_in;type:tinyint(1)"`                            /*  */
	DecNav                *uint64          `validate:"required" db:"dec_nav"    json:"dec_nav" gorm:"column:dec_nav;type:tinyint(4)"`                                                 /*  */
	DecPerformance        *uint64          `validate:"required" db:"dec_performance"    json:"dec_performance" gorm:"column:dec_performance;type:tinyint(4)"`                         /*  */
	DecUnit               *uint64          `validate:"required" db:"dec_unit"    json:"dec_unit" gorm:"column:dec_unit;type:tinyint(4)"`                                              /*  */
	DecAmount             *uint64          `validate:"required" db:"dec_amount"    json:"dec_amount" gorm:"column:dec_amount;type:tinyint(4)"`                                        /*  */
	NpwpNumber            *string          `validate:"" db:"npwp_number"    json:"npwp_number" gorm:"column:npwp_number;type:varchar(30)"`                                            /*  */
	NpwpDateReg           *string          `validate:"" db:"npwp_date_reg"    json:"npwp_date_reg" gorm:"column:npwp_date_reg;type:date"`                                             /*  */
	NpwpName              *string          `validate:"" db:"npwp_name"    json:"npwp_name" gorm:"column:npwp_name;type:varchar(150)"`                                                 /*  */
	PortfolioCode         *string          `validate:"" db:"portfolio_code"    json:"portfolio_code" gorm:"column:portfolio_code;type:varchar(50)"`                                   /*  */
}

func InsertMasterProduct(params map[string]string) error {
	query := "INSERT INTO ms_product"
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

	log.Println(query)
	_, err := db.Db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
