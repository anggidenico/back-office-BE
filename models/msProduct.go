package models

import (
	"database/sql"
	"log"
	"mf-bo-api/config"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type MsProductList struct {
	ProductKey                uint64                 `json:"product_key"`
	ProductID                 uint64                 `json:"product_id"`
	ProductCode               string                 `json:"product_code"`
	ProductName               string                 `json:"product_name"`
	ProductNameAlt            string                 `json:"product_name_alt"`
	MinSubAmount              decimal.Decimal        `json:"min_sub_amount"`
	RecImage1                 string                 `json:"rec_image1"`
	FundType                  *MsFundTypeInfo        `json:"fund_type,omitempty"`
	NavPerformance            *FfsNavPerformanceInfo `json:"nav_performance,omitempty"`
	Nav                       *TrNavInfo             `json:"nav,omitempty"`
	RiskProfile               *MsRiskProfileInfo     `json:"risk_profile,omitempty"`
	IsAllowRedemption         bool                   `json:"is_allow_redemption"`
	IsAllowSwitchin           bool                   `json:"is_allow_switchin"`
	IsAllowProductDestination bool                   `json:"is_allow_product_destination"`
	Currency                  *MsCurrencyInfo        `json:"currency"`
}

type MsProductListDropdown struct {
	ProductKey  uint64 `json:"product_key"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
}

type MsProductInfo struct {
	ProductKey     uint64 `json:"product_key"`
	ProductCode    string `json:"product_code"`
	ProductName    string `json:"product_name"`
	ProductNameAlt string `json:"product_name_alt"`
}

type MsProductData struct {
	ProductKey       uint64                     `json:"product_key"`
	ProductID        uint64                     `json:"product_id"`
	ProductCode      string                     `json:"product_code"`
	ProductName      string                     `json:"product_name"`
	ProductNameAlt   string                     `json:"product_name_alt"`
	MinSubAmount     decimal.Decimal            `json:"min_sub_amount"`
	MinRedAmount     decimal.Decimal            `json:"min_red_amount"`
	MinRedUnit       decimal.Decimal            `json:"min_red_unit"`
	MinUnitAfterRed  decimal.Decimal            `json:"min_unit_after_red"`
	ProspectusLink   string                     `json:"prospectus_link"`
	FundFactSheet    string                     `json:"ffs_link"`
	RecImage1        string                     `json:"rec_image1"`
	FlagSubscription bool                       `json:"flag_subscription"`
	FlagRedemption   bool                       `json:"flag_redemption"`
	FlagSwitchOut    bool                       `json:"flag_switch_out"`
	FlagSwitchIn     bool                       `json:"flag_switch_in"`
	FeeService       string                     `json:"fee_service"`
	FeeTransfer      string                     `json:"fee_transfer"`
	InvestValue      string                     `json:"invest_value"`
	RedSuspend       bool                       `json:"red_suspend"`
	SubSuspend       bool                       `json:"sub_suspend"`
	BalanceUnit      decimal.Decimal            `json:"balance_unit"`
	IsNew            bool                       `json:"is_new"`
	TncIsNew         string                     `json:"tnc_is_new"`
	Currency         *MsCurrencyInfo            `json:"currency"`
	BankAcc          []MsProductBankAccountInfo `json:"bank_account"`
	ProductFee       []MsProductFeeInfo         `json:"product_fee"`
	NavPerformance   *FfsNavPerformanceInfo     `json:"nav_performance,omitempty"`
	Nav              *TrNavInfo                 `json:"nav,omitempty"`
	CustodianBank    *MsCustodianBankInfo       `json:"custodian_bank,omitempty"`
	RiskProfile      *MsRiskProfileInfo         `json:"risk_profile,omitempty"`
}

type MsProduct struct {
	ProductKey            uint64           `db:"product_key"             json:"product_key"`
	ProductID             *uint64          `db:"product_id"              json:"product_id"`
	ProductCode           string           `db:"product_code"            json:"product_code"`
	ProductName           string           `db:"product_name"            json:"product_name"`
	ProductNameAlt        string           `db:"product_name_alt"        json:"product_name_alt"`
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
	FlagSyariah           uint8            `db:"flag_syariah"            json:"flag_syariah"`
	MaxSubFee             decimal.Decimal  `db:"max_sub_fee"             json:"max_sub_fee"`
	MaxRedFee             decimal.Decimal  `db:"max_red_fee"             json:"max_red_fee"`
	MaxSwiFee             decimal.Decimal  `db:"max_swi_fee"             json:"max_swi_fee"`
	MinSubAmount          decimal.Decimal  `db:"min_sub_amount"          json:"min_sub_amount"`
	MinTopUpAmount        *decimal.Decimal `db:"min_topup_amount"        json:"min_topup_amount"`
	MinRedAmount          decimal.Decimal  `db:"min_red_amount"          json:"min_red_amount"`
	MinRedUnit            decimal.Decimal  `db:"min_red_unit"            json:"min_red_unit"`
	MinUnitAfterRed       decimal.Decimal  `db:"min_unit_after_red"      json:"min_unit_after_red"`
	MinAmountAfterRed     *decimal.Decimal `db:"min_amount_after_red"    json:"min_amount_after_red"`
	ManagementFee         decimal.Decimal  `db:"management_fee"          json:"management_fee"`
	CustodianFee          decimal.Decimal  `db:"custodian_fee"           json:"custodian_fee"`
	CustodianKey          *uint64          `db:"custodian_key"           json:"custodian_key"`
	OjkFee                *decimal.Decimal `db:"ojk_fee"                 json:"ojk_fee"`
	ProductFeeAmount      *decimal.Decimal `db:"product_fee_amount"      json:"product_fee_amount"`
	OverwriteTransactFlag uint8            `db:"overwrite_transact_flag" json:"overwrite_transact_flag"`
	OverwriteFeeFlag      uint8            `db:"overwrite_fee_flag"      json:"overwrite_fee_flag"`
	OtherFeeAmount        decimal.Decimal  `db:"other_fee_amount"        json:"other_fee_amount"`
	SettlementPeriod      *uint64          `db:"settlement_period"       json:"settlement_period"`
	SinvestFundCode       *string          `db:"sinvest_fund_code"       json:"sinvest_fund_code"`
	FlagEnabled           uint8            `db:"flag_enabled"            json:"flag_enabled"`
	FlagSubscription      uint8            `db:"flag_subscription"       json:"flag_subscription"`
	FlagRedemption        uint8            `db:"flag_redemption"         json:"flag_redemption"`
	FlagSwitchOut         uint8            `db:"flag_switch_out"         json:"flag_switch_out"`
	FlagSwitchIn          uint8            `db:"flag_switch_in"          json:"flag_switch_in"`
	DecUnit               uint8            `db:"dec_unit" json:"dec_unit"`
	DecAmount             uint8            `db:"dec_amount" json:"dec_amount"`
	DecNav                uint8            `db:"dec_nav" json:"dec_nav"`
	DecPerformance        uint8            `db:"dec_performance" json:"dec_performance"`
	NpwpDateReg           *string          `db:"npwp_date_reg" json:"npwp_date_reg"`
	NpwpName              *string          `db:"npwp_name" json:"npwp_name"`
	NpwpNumber            *string          `db:"npwp_number" json:"npwp_number"`
	PortfolioCode         *string          `db:"portfolio_code" json:"portfolio_code"`
	RecOrder              *uint64          `db:"rec_order"               json:"rec_order"`
	RecStatus             uint8            `db:"rec_status"              json:"rec_status"`
	RecCreatedDate        *string          `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy          *string          `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate       *string          `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy         *string          `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1             *string          `db:"rec_image1"              json:"rec_image1"`
	RecImage2             *string          `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus     *uint8           `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage      *uint64          `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate       *string          `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy         *string          `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate        *string          `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy          *string          `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1       *string          `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2       *string          `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3       *string          `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type AdminMsProductList struct {
	ProductKey          uint64  `json:"product_key"`
	ProductCode         string  `json:"product_code"`
	ProductName         string  `json:"product_name"`
	ProductNameAlt      string  `json:"product_name_alt"`
	CurrencyName        *string `json:"currency_name"`
	ProductCategoryName *string `json:"product_category_name"`
	ProductTypeName     *string `json:"product_type_name"`
	RiskProfileName     *string `json:"risk_profile_name"`
	LaunchDate          *string `json:"launch_date"`
	InceptionDate       *string `json:"inception_date"`
	IsinCode            *string `json:"isin_code"`
	Syariah             bool    `json:"syariah"`
	CustodianFullName   *string `json:"custodian_full_name"`
	SinvestFundCode     *string `json:"sinvest_fund_code"`
	Enabled             bool    `json:"enabled"`
	Subscription        bool    `json:"subscription"`
	Redemption          bool    `json:"redemption"`
	SwitchOut           bool    `json:"switch_out"`
	SwitchIn            bool    `json:"switch_in"`
	StatusUpdate        bool    `json:"status_update"`
}

func CreateMsProduct(params map[string]string) (int, error) {
	query := "INSERT INTO ms_product"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

type AdminMsProductDetail struct {
	ProductKey            uint64                   `json:"product_key"`
	ProductCode           string                   `json:"product_code"`
	ProductName           string                   `json:"product_name"`
	ProductNameAlt        string                   `json:"product_name_alt"`
	Currency              *MsCurrencyInfo          `json:"currency"`
	ProductCategory       *MsProductCategoryInfo   `json:"product_category"`
	ProductType           *MsProductTypeInfo       `json:"product_type"`
	FundType              *MsFundTypeInfo          `json:"fund_type"`
	FundStructure         *MsFundStructureInfo     `json:"fund_structure"`
	RiskProfile           *MsRiskProfileInfoAdmin  `json:"risk_profile"`
	ProductProfile        *string                  `json:"product_profile"`
	InvestmentObjectives  *string                  `json:"investment_objectives"`
	ProductPhase          *LookupTrans             `json:"product_phase"`
	NavValuationType      *LookupTrans             `json:"nav_valuation_type"`
	ProspectusLink        *string                  `json:"prospectus_link"`
	LaunchDate            *string                  `json:"launch_date"`
	InceptionDate         *string                  `json:"inception_date"`
	IsinCode              *string                  `json:"isin_code"`
	FlagSyariah           bool                     `json:"flag_syariah"`
	MaxSubFee             decimal.Decimal          `json:"max_sub_fee"`
	MaxRedFee             decimal.Decimal          `json:"max_red_fee"`
	MaxSwiFee             decimal.Decimal          `json:"max_swi_fee"`
	MinSubAmount          decimal.Decimal          `json:"min_sub_amount"`
	MinRedAmount          decimal.Decimal          `json:"min_red_amount"`
	MinRedUnit            decimal.Decimal          `json:"min_red_unit"`
	MinUnitAfterRed       decimal.Decimal          `json:"min_unit_after_red"`
	ManagementFee         decimal.Decimal          `json:"management_fee"`
	CustodianFee          decimal.Decimal          `json:"custodian_fee"`
	Custodian             *MsCustodianBankInfoList `json:"custodian"`
	OjkFee                *decimal.Decimal         `json:"ojk_fee"`
	ProductFeeAmount      *decimal.Decimal         `json:"product_fee_amount"`
	OverwriteTransactFlag bool                     `json:"overwrite_transact_flag"`
	OverwriteFeeFlag      bool                     `json:"overwrite_fee_flag"`
	OtherFeeAmount        decimal.Decimal          `json:"other_fee_amount"`
	SettlementPeriod      *uint64                  `json:"settlement_period"`
	SinvestFundCode       *string                  `json:"sinvest_fund_code"`
	FlagEnabled           bool                     `json:"flag_enabled"`
	FlagSubscription      bool                     `json:"flag_subscription"`
	FlagRedemption        bool                     `json:"flag_redemption"`
	FlagSwitchOut         bool                     `json:"flag_switch_out"`
	FlagSwitchIn          bool                     `json:"flag_switch_in"`
	MinTopupAmount        *decimal.Decimal         `json:"min_topup_amount"`
	MinAmountAfterRed     *decimal.Decimal         `json:"min_amount_after_red"`
	DecNav                uint64                   `json:"dec_nav"`
	DecPerformance        uint64                   `json:"dec_performance"`
	DecUnit               uint64                   `json:"dec_unit"`
	DecAmount             uint64                   `json:"dec_amount"`
	NpwpName              *string                  `json:"npwp_name"`
	NpwpNumber            *string                  `json:"npwp_number"`
	NpwpDateReg           *string                  `json:"npwp_date_reg"`
	PortfolioCode         *string                  `json:"portfolio_code"`
}

type ProductSubscription struct {
	ProductKey     uint64          `db:"product_key" json:"product_key"`
	FundTypeName   string          `db:"fund_type_name"        json:"fund_type_name"`
	ProductName    string          `db:"product_name"          json:"product_name"`
	NavDate        string          `db:"nav_date"              json:"nav_date"`
	NavValue       decimal.Decimal `db:"nav_value"             json:"nav_value"`
	ProductImage   *string         `db:"product_image"         json:"product_image"`
	MinSubAmount   decimal.Decimal `db:"min_sub_amount"        json:"min_sub_amount"`
	MinRedAmount   decimal.Decimal `db:"min_red_amount"        json:"min_red_amount"`
	MinRedUnit     decimal.Decimal `db:"min_red_unit"          json:"min_red_unit"`
	ProspectusLink string          `db:"prospectus_link"       json:"prospectus_link"`
	FfsLink        *string         `db:"ffs_link"              json:"ffs_link"`
	RiskName       string          `db:"risk_name"             json:"risk_name"`
	CurrencyKey    uint64          `db:"currency_key"          json:"currency_key"`
	Symbol         string          `db:"currency_symbol"       json:"symbol"`
	CurrencyCode   string          `db:"currency_code"         json:"currency_code"`
	CurrencyName   string          `db:"currency_name"         json:"currency_name"`
	FlagShowOntnc  *uint64         `db:"flag_show_ontnc"       json:"flag_show_ontnc"`
	FeeAnnotation  *string         `db:"fee_annotation"        json:"fee_annotation"`
	FeeValue       decimal.Decimal `db:"fee_value"             json:"fee_value"`
}

type ProductSubscriptionFundType struct {
	ProductKey  uint64 `db:"product_key"             json:"product_key"`
	ProductName string `db:"product_name"            json:"product_name"`
}

func GetAllMsProduct(c *[]MsProduct, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_product.* FROM 
			  ms_product `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
func GetAutoinvestProduct(c *[]MsProduct, customerKey string) (int, error) {
	query := `SELECT
			p.* 
			FROM ms_product AS p 
			LEFT JOIN tr_account AS ta ON ta.product_key = p.product_key 
			AND ta.customer_key = ` + customerKey +
		` AND ta.rec_status = 1
			WHERE p.flag_enabled = 1 AND p.flag_subscription = 1 
			AND p.rec_status = 1 
			AND p.flag_enabled = 1 
			AND (ta.sub_suspend_flag IS NULL OR ta.sub_suspend_flag = 0) `

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsProduct(c *MsProduct, key string) (int, error) {
	query := `SELECT ms_product.* FROM ms_product WHERE ms_product.product_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductIn(c *[]MsProduct, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_product.* FROM 
				ms_product WHERE 
				ms_product.rec_status = 1 `
	query := query2 + " AND ms_product." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("========= QUERY GET PRODUCT ========= >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllMsProductWithLike(c *[]MsProduct, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_product.* FROM 
			  ms_product `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, "ms_product."+fieldLike+" like '%"+valueLike+"%'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ProductStatusUpdate(productKey string) bool {
	query := `SELECT count(*) FROM ms_product_request 
	WHERE rec_status = 1 AND rec_approval_status IS NULL 
	AND product_key = ` + productKey
	// log.Println(query)
	var count uint64
	err := db.Db.Get(&count, query)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	var result bool
	if count > 0 { // kalau ada request gantung maka false
		result = false
	} else {
		result = true
	}
	return result
}

func AdminGetCountMsProductWithLike(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT
			  count(ms_product.product_key) as count_data 
			  FROM ms_product `
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, "ms_product."+fieldLike+" like '%"+valueLike+"%'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsProduct(params map[string]string) (int, error) {
	query := "UPDATE ms_product SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "product_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE product_key = " + params["product_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func AdminGetValidateUniqueDataInsertUpdate(c *CountData, paramsOr map[string]string, paramsAnd map[string]string, updateKey *string) (int, error) {
	query := `SELECT
			  count(ms_product.product_key) as count_data 
			  FROM ms_product `
	var orWhereClause []string
	var andWhereClause []string
	var condition string

	for fieldOr, valueOr := range paramsOr {
		orWhereClause = append(orWhereClause, "ms_product."+fieldOr+" = '"+valueOr+"'")
	}

	for fieldAnd, valueAnd := range paramsAnd {
		andWhereClause = append(andWhereClause, "ms_product."+fieldAnd+" like '"+valueAnd+"'")
	}

	// Combile where Or clause
	if len(orWhereClause) > 0 {
		condition += " WHERE ("
		for index, where := range orWhereClause {
			condition += where
			if (len(orWhereClause) - 1) > index {
				condition += " OR "
			} else {
				condition += ") "
			}
		}
	}

	// Combile where And clause
	if len(andWhereClause) > 0 {
		if len(orWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		for index, where := range andWhereClause {
			condition += where
			if (len(andWhereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if updateKey != nil {
		if len(orWhereClause) > 0 {
			condition += " AND "
		} else if len(andWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		condition += " ms_product.product_key != '" + *updateKey + "'"
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetProductSubscription(c *[]ProductSubscriptionFundType, fundtypeKey string) (int, error) {
	query := `SELECT 
				p.product_key as product_key,
				p.product_name_alt as product_name 
			FROM ms_product AS p 
			WHERE p.rec_status = 1 AND p.flag_subscription = 1 AND p.fund_type_key = '` + fundtypeKey + `'
			ORDER BY p.rec_order ASC`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetProductSubscriptionByProductKey(c *ProductSubscription, productKey string) (int, error) {
	query := `SELECT 
				p.product_key as product_key,
				f.fund_type_name as fund_type_name,
				p.product_name_alt as product_name,
				DATE_FORMAT(nav.nav_date, '%d %M %Y') AS nav_date, 
				nav.nav_value as nav_value,
				(CASE
					WHEN p.rec_image1 IS NULL THEN CONCAT('` + config.BaseUrl + `', '/images/product/default.png')
					ELSE CONCAT('` + config.BaseUrl + `', '/images/product/', p.rec_image1)
				END) AS product_image,
				p.min_sub_amount as min_sub_amount, 
				p.min_red_amount as min_red_amount,
				p.min_red_unit as min_red_unit,
				p.prospectus_link as prospectus_link,
				pub.ffs_link as ffs_link,
				pr.risk_name as risk_name,
				cur.currency_key AS currency_key,
				cur.code AS currency_code,
				cur.symbol AS currency_symbol,
				cur.name AS currency_name,
				mpf.flag_show_ontnc AS flag_show_ontnc, 
				mpf.fee_annotation AS fee_annotation, 
				mpfi.fee_value AS fee_value 
			FROM ms_product AS p 
			INNER JOIN ms_fund_type AS f ON p.fund_type_key = f.fund_type_key 
			INNER JOIN (
				SELECT MAX(nav_date) AS nav_date, product_key 
				FROM tr_nav
				WHERE rec_status = 1
				AND publish_mode = 236
				GROUP BY product_key
			) b ON (b.product_key = p.product_key)
			INNER JOIN tr_nav AS nav ON nav.product_key = p.product_key AND nav.nav_date = b.nav_date
			INNER JOIN (
				SELECT product_key, MAX(nav_date) AS nav_date
				FROM ffs_nav_performance
				WHERE rec_status = 1
				GROUP BY product_key
			) c ON (c.product_key = p.product_key)
			INNER JOIN ms_risk_profile AS pr ON p.risk_profile_key = pr.risk_profile_key 
			LEFT JOIN ms_currency AS cur ON cur.currency_key = p.currency_key 
			LEFT JOIN ms_product_fee AS mpf ON mpf.product_key = p.product_key AND mpf.fee_type = '183' AND mpf.rec_status = 1 
			LEFT JOIN ms_product_fee_item AS mpfi ON mpfi.product_fee_key = mpf.fee_key AND mpfi.rec_status = 1
			LEFT JOIN ffs_publish AS pub ON pub.product_key = p.product_key 
			WHERE p.rec_status = 1 AND p.flag_subscription = 1 
			AND p.product_key = '` + productKey + `'
			ORDER BY f.rec_order ASC`

	// Main query
	// log.Println("========== QUERY CEK PRODUCT SUBSCRIPTION =========", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type ProductRedemption struct {
	ProductKey      uint64          `db:"product_key"             json:"product_key"`
	RiskProfileKey  uint64          `db:"risk_profile_key"        json:"risk_profile_key"`
	FundTypeName    string          `db:"fund_type_name"          json:"fund_type_name"`
	ProductName     string          `db:"product_name"            json:"product_name"`
	NavDate         string          `db:"nav_date"                json:"nav_date"`
	FlagRedemption  string          `db:"flag_redemption"         json:"flag_redemption"`
	FlagSwitchOut   string          `db:"flag_switch_out"         json:"flag_switch_out"`
	NavValue        decimal.Decimal `db:"nav_value"               json:"nav_value"`
	MinRedAmount    decimal.Decimal `db:"min_red_amount"          json:"min_red_amount"`
	MinRedUnit      decimal.Decimal `db:"min_red_unit"            json:"min_red_unit"`
	MinUnitAfterRed decimal.Decimal `db:"min_unit_after_red"      json:"min_unit_after_red"`
	RiskName        string          `db:"risk_name"               json:"risk_name"`
	AcaKey          uint64          `db:"aca_key"                 json:"aca_key"`
	Unit            decimal.Decimal `db:"unit"                    json:"unit,omitempty"`
	NilaiInvestasi  decimal.Decimal `db:"nilai_investasi"         json:"nilai_investasi,omitempty"`
	SalesName       *string         `db:"sales_name"              json:"sales_name,omitempty"`
}

func AdminGetProductRedemption(c *[]ProductRedemption, customerKey string) (int, error) {
	query := `SELECT 
				p.product_key AS product_key,
				p.risk_profile_key AS risk_profile_key,
				f.fund_type_name AS fund_type_name,
				p.product_name_alt AS product_name,
				DATE_FORMAT(nav.nav_date, '%d %M %Y') AS nav_date, 
				p.flag_redemption,
				p.flag_switch_out,
				nav.nav_value AS nav_value,
				p.min_red_amount AS min_red_amount,
				p.min_red_unit AS min_red_unit,
				p.min_unit_after_red AS min_unit_after_red,
				pr.risk_name AS risk_name,
				t.aca_key AS aca_key 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key 
			INNER JOIN ms_fund_type AS f ON p.fund_type_key = f.fund_type_key 
			INNER JOIN (
				SELECT MAX(nav_date) AS nav_date, product_key 
				FROM tr_nav
				WHERE rec_status = 1
				AND publish_mode = 236
				GROUP BY product_key
			) b ON (b.product_key = p.product_key)
			INNER JOIN tr_nav AS nav ON nav.product_key = p.product_key AND nav.nav_date = b.nav_date
			INNER JOIN ms_risk_profile AS pr ON p.risk_profile_key = pr.risk_profile_key
			WHERE p.rec_status = 1 AND f.show_home = 1 AND
			t.rec_status = 1 AND t.customer_key = "` + customerKey + `" AND trans_status_key = 9
			AND t.trans_type_key IN (1,4) 
			GROUP BY p.product_key ORDER BY f.rec_order ASC`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type ProductHaveBalanceSwitchIn struct {
	ProductKey     uint64          `db:"product_key"             json:"product_key"`
	RiskProfileKey uint64          `db:"risk_profile_key"        json:"risk_profile_key"`
	FundTypeName   string          `db:"fund_type_name"          json:"fund_type_name"`
	ProductName    string          `db:"product_name"            json:"product_name"`
	NavDate        string          `db:"nav_date"                json:"nav_date"`
	FlagSwitchIn   string          `db:"flag_switch_in"          json:"flag_switch_in"`
	NavValue       decimal.Decimal `db:"nav_value"               json:"nav_value"`
	MinSubAmount   decimal.Decimal `db:"min_sub_amount"          json:"min_sub_amount"`
	RiskName       string          `db:"risk_name"               json:"risk_name"`
	AcaKey         *uint64         `db:"aca_key"                 json:"aca_key"`
	Unit           decimal.Decimal `db:"unit"                    json:"unit,omitempty"`
	NilaiInvestasi decimal.Decimal `db:"nilai_investasi"         json:"nilai_investasi,omitempty"`
}

func AdminGetProductHaveBalanceSwitchIn(c *[]ProductHaveBalanceSwitchIn, customerKey string, productSwOutKey string) (int, error) {
	query := `SELECT 
				p.product_key AS product_key,
				p.risk_profile_key AS risk_profile_key,
				f.fund_type_name AS fund_type_name,
				p.product_name_alt AS product_name,
				DATE_FORMAT(nav.nav_date, '%d %M %Y') AS nav_date, 
				p.flag_switch_in,
				nav.nav_value AS nav_value,
				p.min_sub_amount AS min_sub_amount,
				pr.risk_name AS risk_name,
				t.aca_key AS aca_key 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key 
			INNER JOIN ms_fund_type AS f ON p.fund_type_key = f.fund_type_key 
			INNER JOIN (
				SELECT MAX(nav_date) AS nav_date, product_key 
				FROM tr_nav
				WHERE rec_status = 1
				AND publish_mode = 236
				GROUP BY product_key
			) b ON (b.product_key = p.product_key)
			INNER JOIN tr_nav AS nav ON nav.product_key = p.product_key AND nav.nav_date = b.nav_date
			INNER JOIN ms_risk_profile AS pr ON p.risk_profile_key = pr.risk_profile_key
			WHERE p.rec_status = 1 AND f.show_home = 1 AND
			t.rec_status = 1 AND t.customer_key = "` + customerKey + `" AND trans_status_key = 9
			AND t.trans_type_key IN (1,4) AND p.product_key != "` + productSwOutKey + `"
			GROUP BY p.product_key ORDER BY f.rec_order ASC`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetProductNotInBalanceSwitchIn(c *[]ProductHaveBalanceSwitchIn, value []string) (int, error) {
	inQuery := strings.Join(value, ",")
	query := `SELECT 
				p.product_key AS product_key,
				p.risk_profile_key AS risk_profile_key,
				f.fund_type_name AS fund_type_name,
				p.product_name_alt AS product_name,
				DATE_FORMAT(nav.nav_date, '%d %M %Y') AS nav_date, 
				p.flag_switch_in,
				nav.nav_value AS nav_value,
				p.min_sub_amount AS min_sub_amount,
				pr.risk_name AS risk_name,
				NULL AS aca_key,
				0 AS unit,
				0 AS nilai_investasi 
			FROM ms_product AS p
			INNER JOIN ms_fund_type AS f ON p.fund_type_key = f.fund_type_key 
			INNER JOIN (
				SELECT MAX(nav_date) AS nav_date, product_key 
				FROM tr_nav
				WHERE rec_status = 1
				AND publish_mode = 236
				GROUP BY product_key
			) b ON (b.product_key = p.product_key)
			INNER JOIN tr_nav AS nav ON nav.product_key = p.product_key AND nav.nav_date = b.nav_date
			INNER JOIN ms_risk_profile AS pr ON p.risk_profile_key = pr.risk_profile_key
			WHERE p.rec_status = 1 AND f.show_home = 1 AND p.product_key NOT IN (` + inQuery + `)
			GROUP BY p.product_key ORDER BY p.product_key ASC`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
