package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type Portofolio struct {
	Date          string
	Cif           string
	Sid           string
	Name          string
	Address       string
	City          string
	Country       string
	Datas         []ProductPortofolio
	Total         string
	TotalGainLoss string
}

type ProductPortofolio struct {
	ProductName string
	AvgNav      string
	Nav         string
	Unit        string
	CCY         string
	Amount      string
	GainLoss    string
	Kurs        string
	AmountIDR   string
	GainLossIDR string
}

type InstitutionPortofolio struct {
	Date    string
	Cif     string
	Sid     string
	Name    string
	Address string
	City    string
	Country string
	Datas   []InstitutionProductPortofolio
}

type InstitutionProductPortofolio struct {
	TanggalTransaksi string
	ProductName      string
	ProductTujuan    string
	NavDate          string
	TipeTransaksi    string
	JumlahTransaksi  string
	NavValue         string
	Unit             string
	TotalPembelian   string
	Currency         string
}

type TrTransaction struct {
	TransactionKey    uint64           `db:"transaction_key"           json:"transaction_key"`
	ParentKey         *uint64          `db:"parent_key"                json:"parent_key"`
	IDTransaction     *uint64          `db:"id_transaction"            json:"id_transaction"`
	BranchKey         *uint64          `db:"branch_key"                json:"branch_key"`
	AgentKey          *uint64          `db:"agent_key"                 json:"agent_key"`
	CustomerKey       uint64           `db:"customer_key"              json:"customer_key"`
	ProductKey        uint64           `db:"product_key"               json:"product_key"`
	CurrencyKey       uint64           `db:"currency_key"               json:"currency_key"`
	TransStatusKey    uint64           `db:"trans_status_key"          json:"trans_status_key"`
	TransDate         string           `db:"trans_date"                json:"trans_date"`
	TransTypeKey      uint64           `db:"trans_type_key"            json:"trans_type_key"`
	TrxCode           *uint64          `db:"trx_code"                  json:"trx_code"`
	NavDate           string           `db:"nav_date"                  json:"nav_date"`
	EntryMode         *uint64          `db:"entry_mode"                json:"entry_mode"`
	TransCalcMethod   *uint64          `db:"trans_calc_method"         json:"trans_calc_method"`
	TransAmount       decimal.Decimal  `db:"trans_amount"              json:"trans_amount"`
	TransUnit         decimal.Decimal  `db:"trans_unit"                json:"trans_unit"`
	TransUnitPercent  *decimal.Decimal `db:"trans_unit_percent"        json:"trans_unit_percent"`
	FlagRedemtAll     *uint8           `db:"flag_redempt_all"          json:"flag_redempt_all"`
	FlagNewSub        *uint8           `db:"flag_newsub"               json:"flag_newsub"`
	TransFeePercent   decimal.Decimal  `db:"trans_fee_percent"         json:"trans_fee_percent"`
	TransFeeAmount    decimal.Decimal  `db:"trans_fee_amount"          json:"trans_fee_amount"`
	ChargesFeeAmount  decimal.Decimal  `db:"charges_fee_amount"        json:"charges_fee_amount"`
	ServicesFeeAmount decimal.Decimal  `db:"services_fee_amount"       json:"services_fee_amount"`
	StampFeeAmount    decimal.Decimal  `db:"stamp_fee_amount"          json:"stamp_fee_amount"`
	TotalAmount       decimal.Decimal  `db:"total_amount"              json:"total_amount"`
	SettlementDate    *string          `db:"settlement_date"           json:"settlement_date"`
	TransBankAccNo    *string          `db:"trans_bank_accno"          json:"trans_bank_accno"`
	TransBankaccName  *string          `db:"trans_bankacc_name"        json:"trans_bankacc_name"`
	TransBankKey      *uint64          `db:"trans_bank_key"            json:"trans_bank_key"`
	TransRemarks      *string          `db:"trans_remarks"             json:"trans_remarks"`
	TransReferences   *string          `db:"trans_references"          json:"trans_references"`
	PromoCode         *string          `db:"promo_code"                json:"promo_code"`
	SalesCode         *string          `db:"sales_code"                json:"sales_code"`
	RiskWaiver        uint8            `db:"risk_waiver"               json:"risk_waiver"`
	AddtoAutoInvest   *uint8           `db:"addto_auto_invest"         json:"addto_auto_invest"`
	TransSource       *uint64          `db:"trans_source"              json:"trans_source"`
	FileUploadDate    *string          `db:"file_upload_date"          json:"file_upload_date"`
	PaymentMethod     *uint64          `db:"payment_method"            json:"payment_method"`
	Check1Date        *string          `db:"check1_date"               json:"check1_date"`
	Check1Flag        *uint8           `db:"check1_flag"               json:"check1_flag"`
	Check1References  *string          `db:"check1_references"         json:"check1_references"`
	Check1Notes       *string          `db:"check1_notes"              json:"check1_notes"`
	Check2Date        *string          `db:"check2_date"               json:"check2_date"`
	Check2Flag        *uint8           `db:"check2_flag"               json:"check2_flag"`
	Check2References  *string          `db:"check2_references"         json:"check2_references"`
	Check2Notes       *string          `db:"check2_notes"              json:"check2_notes"`
	TrxRiskLevel      *uint64          `db:"trx_risk_level"            json:"trx_risk_level"`
	ProceedDate       *string          `db:"proceed_date"              json:"proceed_date"`
	ProceedAmount     *decimal.Decimal `db:"proceed_amount"            json:"proceed_amount"`
	SentDate          *string          `db:"sent_date"                 json:"sent_date"`
	SentReferences    *string          `db:"sent_references"           json:"sent_references"`
	ConfirmedDate     *string          `db:"confirmed_date"            json:"confirmed_date"`
	PostedDate        *string          `db:"posted_date"               json:"posted_date"`
	PostedUnits       *decimal.Decimal `db:"posted_units"              json:"posted_units"`
	AcaKey            *uint64          `db:"aca_key"                   json:"aca_key"`
	SettledDate       *string          `db:"settled_date"              json:"settled_date"`
	BatchKey          *uint64          `db:"batch_key"                 json:"batch_key"`
	RecOrder          *uint64          `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8            `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string          `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string          `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string          `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string          `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string          `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string          `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8           `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64          `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string          `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string          `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string          `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string          `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string          `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string          `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string          `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrTransactionList struct {
	TransactionKey    uint64                        `json:"transaction_key"`
	ProductName       string                        `json:"product_name"`
	TransStatus       string                        `json:"trans_status"`
	TransDate         string                        `json:"trans_date"`
	TransType         string                        `json:"trans_type"`
	NavDate           string                        `json:"nav_date"`
	NavValue          decimal.Decimal               `json:"nav_value"`
	TransAmount       decimal.Decimal               `json:"trans_amount,omitempty"`
	TransUnit         decimal.Decimal               `json:"trans_unit,omitempty"`
	TotalAmount       decimal.Decimal               `json:"total_amount"`
	TransFeePercent   decimal.Decimal               `json:"trans_fee_percent"`
	TransFeeAmount    decimal.Decimal               `json:"trans_fee_amount"`
	ChargesFeeAmount  decimal.Decimal               `json:"charges_fee_amount"`
	ServicesFeeAmount decimal.Decimal               `json:"services_fee_amount"`
	StampFeeAmount    decimal.Decimal               `json:"stamp_fee_amount"`
	Uploaded          bool                          `json:"uploaded"`
	DateUploaded      *string                       `json:"date_uploaded"`
	ProdBankAcc       *TransactionPoductBankAccount `json:"prod_bank_acc,omitempty"`
	BankName          *string                       `json:"bank_name"`
	BankAccNo         *string                       `json:"bank_accno"`
	BankAccName       *string                       `json:"bankacc_name"`
	PaymentMethod     *string                       `json:"payment_method"`
	ProductOut        *string                       `json:"product_name_out"`
	ProductIn         *string                       `json:"product_name_in"`
	Currency          *MsCurrencyInfo               `json:"currency"`
}

type AdminTrTransactionList struct {
	TransactionKey   uint64          `json:"transaction_key"`
	BranchName       string          `json:"branch_name"`
	AgentName        string          `json:"agent_name"`
	CustomerName     string          `json:"customer_name"`
	ProductName      string          `json:"product_name"`
	TransStatus      string          `json:"trans_status"`
	TransDate        string          `json:"trans_date"`
	TransType        string          `json:"trans_type"`
	NavDate          string          `json:"nav_date"`
	TransAmount      decimal.Decimal `json:"trans_amount"`
	TransUnit        decimal.Decimal `json:"trans_unit"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	TransBankName    string          `json:"trans_bank_name"`
	TransBankAccNo   *string         `json:"trans_bank_accno"`
	TransBankaccName *string         `json:"trans_bankacc_name"`
	ProductOut       *string         `json:"product_name_out"`
	ProductIn        *string         `json:"product_name_in"`
	PaymentMethod    *string         `json:"payment_method,omitempty"`
	PaymentChannel   *string         `json:"payment_channel,omitempty"`
	TransSource      *string         `json:"trans_source,omitempty"`
	RecImage1        *string         `json:"rec_image1"`
}

type AdminTrTransactionInquiryList struct {
	TransactionKey uint64          `json:"transaction_key"`
	CustomerKey    uint64          `json:"customer_key,omitempty"`
	ProductKey     uint64          `json:"product_key,omitempty"`
	BranchName     string          `json:"branch_name"`
	AgentName      string          `json:"agent_name"`
	CustomerName   string          `json:"customer_name"`
	ProductName    string          `json:"product_name"`
	TransStatus    string          `json:"trans_status"`
	TransDate      string          `json:"trans_date"`
	TransType      string          `json:"trans_type"`
	NavDate        string          `json:"nav_date"`
	TransAmount    decimal.Decimal `json:"trans_amount"`
	TransUnit      decimal.Decimal `json:"trans_unit"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	PaymentMethod  *string         `json:"payment_method,omitempty"`
	PaymentChannel *string         `json:"payment_channel,omitempty"`
	TransSource    *string         `json:"trans_source,omitempty"`
	RecImage1      *string         `json:"rec_image1"`
}

type CountData struct {
	CountData int `db:"count_data"             json:"count_data"`
}

type NavValue struct {
	NavValue *decimal.Decimal `db:"nav_value"              json:"nav_value"`
}

type AdminTransactionDetail struct {
	TransactionKey              uint64                               `json:"transaction_key"`
	Branch                      *BranchTrans                         `json:"branch"`
	Agent                       *AgentTrans                          `json:"agent"`
	Customer                    CustomerTrans                        `json:"customer"`
	Product                     ProductTrans                         `json:"product"`
	TransStatus                 TransStatus                          `json:"trans_status"`
	TransDate                   string                               `json:"trans_date"`
	TransType                   TransType                            `json:"trans_type"`
	TrxCode                     *LookupTrans                         `json:"trx_code"`
	NavDate                     string                               `json:"nav_date"`
	EntryMode                   *LookupTrans                         `json:"entry_mode"`
	TransAmount                 decimal.Decimal                      `json:"trans_amount"`
	TransUnit                   decimal.Decimal                      `json:"trans_unit"`
	TransUnitPercent            *decimal.Decimal                     `json:"trans_unit_percent"`
	FlagRedemtAll               bool                                 `json:"flag_redempt_all"`
	FlagNewSub                  bool                                 `json:"flag_newsub"`
	TransFeePercent             decimal.Decimal                      `json:"trans_fee_percent"`
	TransFeeAmount              decimal.Decimal                      `json:"trans_fee_amount"`
	ChargesFeeAmount            decimal.Decimal                      `json:"charges_fee_amount"`
	ServicesFeeAmount           decimal.Decimal                      `json:"services_fee_amount"`
	StampFeeAmount              decimal.Decimal                      `json:"stamp_fee_amount"`
	TotalAmount                 decimal.Decimal                      `json:"total_amount"`
	SettlementDate              *string                              `json:"settlement_date"`
	TransBankAccNo              *string                              `json:"trans_bank_accno"`
	TransBankaccName            *string                              `json:"trans_bankacc_name"`
	TransBank                   *TransBank                           `json:"trans_bank"`
	TransRemarks                *string                              `json:"trans_remarks"`
	TransReferences             *string                              `json:"trans_references"`
	PromoCode                   *string                              `json:"promo_code"`
	SalesCode                   *string                              `json:"sales_code"`
	RiskWaiver                  bool                                 `json:"risk_waiver"`
	FileUploadDate              *string                              `json:"file_upload_date"`
	UrlUpload                   []*string                            `json:"url_upload_date"`
	PaymentMethod               *LookupTrans                         `json:"payment_method"`
	TransactionSettlement       *[]TransactionSettlement             `json:"transaction_settlement"`
	TrxRiskLevel                *LookupTrans                         `json:"trx_risk_level"`
	ProceedDate                 *string                              `json:"proceed_date"`
	ProceedAmount               *decimal.Decimal                     `json:"proceed_amount"`
	SentDate                    *string                              `json:"sent_date"`
	SentReferences              *string                              `json:"sent_references"`
	ConfirmedDate               *string                              `json:"confirmed_date"`
	PostedDate                  *string                              `json:"posted_date"`
	PostedUnits                 *decimal.Decimal                     `json:"posted_units"`
	Aca                         *AcaTrans                            `json:"aca"`
	SettledDate                 *string                              `json:"settled_date"`
	RecImage1                   *string                              `json:"rec_image1"`
	RecCreatedDate              *string                              `json:"rec_created_date"`
	RecCreatedBy                *string                              `json:"rec_created_by"`
	TransactionConfirmation     *TransactionConfirmation             `json:"transaction_confirmation"`
	ProductBankAccount          *MsProductBankAccountTransactionInfo `json:"product_bank_account"`
	CustomerBankAccount         *MsCustomerBankAccountInfo           `json:"customer_bank_account"`
	IsEnableUnposting           bool                                 `json:"is_enable_unposting"`
	MessageEnableUnposting      string                               `json:"message_enable_unposting"`
	TransactionConfirmationInfo *TrTransactionConfirmationInfo       `json:"transaction_confirmation_info"`
	Promo                       *TrPromoData                         `json:"promo"`
	FundTypeKey                 *uint64                              `json:"fund_type_key,omitempty"`
	TransCalcMethod             *uint64                              `json:"trans_calc_method,omitempty"`
	NavDateReal                 string                               `json:"nav_date_real,omitempty"`
	TransSource                 *string                              `json:"trans_source,omitempty"`
}

type DownloadFormatExcelList struct {
	IDTransaction   uint64           `json:"id_transaction"`
	IDCategory      string           `json:"id_category"`
	ProductName     string           `json:"product_name"`
	FullName        string           `json:"full_name"`
	NavDate         string           `json:"nav_date"`
	TransactionDate string           `json:"transaction_date"`
	Units           decimal.Decimal  `json:"units"`
	NetAmount       decimal.Decimal  `json:"net_amount"`
	NavValue        *decimal.Decimal `json:"nav_value"`
	ApproveUnits    decimal.Decimal  `json:"approve_units"`
	ApproveAmount   decimal.Decimal  `json:"approve_amount"`
	Keterangan      string           `json:"keterangan"`
	Result          string           `json:"result"`
}

type BranchTrans struct {
	BranchKey  uint64 `json:"branch_key"`
	BranchCode string `json:"branch_code"`
	BranchName string `json:"branch_name"`
}

type AgentTrans struct {
	AgentKey  uint64 `json:"agent_key"`
	AgentCode string `json:"agent_code"`
	AgentName string `json:"agent_name"`
}

type CustomerTrans struct {
	CustomerKey    uint64  `json:"customer_key"`
	FullName       string  `json:"full_name"`
	SidNo          *string `json:"sid_no"`
	UnitHolderIDno string  `json:"unit_holder_idno"`
}

type ProductTrans struct {
	ProductKey  uint64 `json:"product_key"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
}

type TransStatus struct {
	TransStatusKey    uint64  `json:"trans_status_key"`
	StatusCode        *string `json:"status_code"`
	StatusDescription *string `json:"status_description"`
}

type TransType struct {
	TransTypeKey    uint64  `json:"trans_type_key"`
	TypeCode        *string `json:"type_code"`
	TypeDescription *string `json:"type_description"`
}

type TransBank struct {
	BankKey  uint64 `json:"bank_key"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type LookupTrans struct {
	LookupKey   uint64  `json:"lookup_key"`
	LkpGroupKey uint64  `json:"lkp_group_key"`
	LkpCode     *string `json:"lkp_code"`
	LkpName     *string `json:"lkp_name"`
}
type AcaTrans struct {
	AcaKey    uint64 `json:"aca_key"`
	AccKey    uint64 `json:"acc_key"`
	AgentKey  uint64 `json:"agent_key"`
	AgentCode string `json:"agent_code"`
	AgentName string `json:"agent_name"`
}
type TransactionConfirmation struct {
	TcKey           uint64          `json:"tc_key"`
	ConfirmDate     string          `json:"confirm_date"`
	ConfirmedAmount decimal.Decimal `json:"confirmed_amount"`
	ConfirmedUnit   decimal.Decimal `json:"confirmed_unit"`
}

type ParamBatchTrTransaction struct {
	ProductCode    string  `db:"product_code"     json:"product_code"`
	TypeCode       string  `db:"type_code"        json:"type_code"`
	Bulan          string  `db:"bulan"            json:"bulan"`
	Tahun          string  `db:"tahun"            json:"tahun"`
	NavDate        string  `db:"nav_date"         json:"nav_date"`
	ProductKey     uint64  `db:"product_key"      json:"product_key"`
	TransTypeKey   uint64  `db:"trans_type_key"   json:"trans_type_key"`
	TransactionKey string  `db:"transaction_key"  json:"transaction_key"`
	TransDate      string  `db:"trans_date"       json:"trans_date"`
	Batch          *uint64 `db:"batch"            json:"batch"`
}

type ProductCheckAllowRedmSwtching struct {
	ProductKey uint64 `db:"product_key"             json:"product_key"`
}

type TransactionCustomerHistory struct {
	ProductName string  `db:"product_name"      json:"product_name"`
	FullName    string  `db:"full_name"         json:"full_name"`
	Cif         *string `db:"cif"               json:"cif"`
	Sid         *string `db:"sid"               json:"sid"`
	ParamDetail *string `db:"param_detail"      json:"param_detail"`
}

type TransactionConsumenProduct struct {
	TransactionKey  uint64          `db:"transaction_key"     json:"transaction_key"`
	TransTypeKey    uint64          `db:"trans_type_key"      json:"trans_type_key"`
	NavDate         string          `db:"nav_date"            json:"nav_date"`
	TypeDescription string          `db:"type_description"    json:"type_description"`
	NavValue        decimal.Decimal `db:"nav_value"            json:"nav_value"`
	Unit            decimal.Decimal `db:"unit"                json:"unit"`
	GrossAmount     decimal.Decimal `db:"gross_amount"        json:"gross_amount"`
	FeeAmount       decimal.Decimal `db:"fee_amount"          json:"fee_amount"`
	NetAmount       decimal.Decimal `db:"net_amount"          json:"net_amount"`
}

type DataDetailTransaksiCustomerProduct struct {
	DataTransaksi     DetailHeaderTransaksiCustomer `json:"data_transaksi"`
	DetailTransaksi   *[]TransactionConsumenProduct `json:"detail_transaksi"`
	CountSubscription decimal.Decimal               `json:"total_subscription"`
	CountRedemption   decimal.Decimal               `json:"total_redemption"`
	CountNetSub       decimal.Decimal               `json:"total_netsub"`
}

type DetailHeaderTransaksiCustomer struct {
	UnitHolder  string `db:"unit_holder"          json:"unit_holder"`
	FullName    string `db:"full_name"            json:"full_name"`
	Sid         string `db:"sid"                  json:"sid"`
	IfuaNo      string `db:"ifua_no"              json:"ifua_no"`
	NavDateFrom string `db:"nav_date_from"        json:"nav_date_from"`
	NavDateTo   string `db:"nav_date_to"          json:"nav_date_to"`
	ProductName string `db:"product_name"         json:"product_name"`
}

func AdminGetAllTrTransaction(c *[]TrTransaction, limit uint64, offset uint64, nolimit bool,
	params map[string]string, valueIn []string, fieldIn string, isAll bool, userID string) (int, error) {
	query := `SELECT
              t.*
			  FROM tr_transaction as t
			  inner join ms_customer as c on c.customer_key = t.customer_key
			  WHERE t.rec_status = 1 AND t.trans_status_key != 3`

	if isAll == false {
		query += " AND t.trans_type_key != 3"
	}
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			if field == "c.openacc_branch_key" {
				whereClause = append(whereClause, field+" = '"+value+"'")
			} else {
				whereClause = append(whereClause, "t."+field+" = '"+value+"'")
			}
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"

		_, cekCsApprove := search(valueIn, "2")
		if cekCsApprove {
			condition += " AND t.rec_created_by != " + userID
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
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetCountTrTransaction(c *CountData, params map[string]string, valueIn []string, fieldIn string, userID string) (int, error) {
	query := `SELECT
              count(t.transaction_key) as count_data
			  FROM tr_transaction as t
			  inner join ms_customer as c on c.customer_key = t.customer_key
			  WHERE t.trans_type_key != 3 `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			if field == "c.openacc_branch_key" {
				whereClause = append(whereClause, field+" = '"+value+"'")
			} else {
				whereClause = append(whereClause, "t."+field+" = '"+value+"'")
			}
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"

		_, cekCsApprove := search(valueIn, "2")
		if cekCsApprove {
			condition += " AND t.rec_created_by != " + userID
		}
	}

	query += condition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTrTransaction(c *[]TrTransaction, params map[string]string) (int, error) {
	query := `SELECT
              tr_transaction.* FROM 
			  tr_transaction`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
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

	// Main query
	log.Info("========= QUERY GET ALL TRANSACTION ======== >>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTrTransactionCount(c *CountData, params map[string]string) (int, error) {

	query := `SELECT count(tr_transaction.transaction_key) as count_data  
			  FROM tr_transaction `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
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
	/*
		var orderBy string
		var orderType string
		if orderBy, present = params["orderBy"]; present == true {
			condition += " ORDER BY " + orderBy
			if orderType, present = params["orderType"]; present == true {
				condition += " " + orderType
			}
		}
	*/
	query += condition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTransactionJoinSettlement(c *[]TrTransaction, customerKey string, paymentMethodKey string, productKey string) (int, error) {
	query := `SELECT t.* FROM tr_transaction t
			INNER JOIN tr_transaction_settlement s
			ON t.transaction_key = s.transaction_key 
			WHERE t.customer_key =  ` + customerKey +
		` AND t.product_key = ` + productKey +
		` AND s.settle_payment_method = ` + paymentMethodKey +
		` AND s.settled_status = 243` +
		` AND t.rec_status = 1`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateTrTransaction(params map[string]string) (int, error) {
	query := "UPDATE tr_transaction SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "transaction_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE transaction_key = " + params["transaction_key"]
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func SetTransactionInactive(customerKey string, paymentMethodKey string, productKey string) (int, error) {
	query := `UPDATE tr_transaction t
				INNER JOIN tr_transaction_settlement s
				ON t.transaction_key = s.transaction_key 
				SET t.rec_status = 0, s.rec_status = 0
				WHERE t.customer_key =  ` + customerKey +
		` AND t.product_key = ` + productKey +
		` AND s.settle_payment_method = ` + paymentMethodKey +
		` AND s.settled_status = 243`
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CreateTrTransaction(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_transaction"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		if value == "NULL" {
			var s *string
			bindvars = append(bindvars, s)
		} else {
			bindvars = append(bindvars, value)
		}

	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetTrTransaction(c *TrTransaction, key string) (int, error) {
	query := `SELECT tr_transaction.* FROM tr_transaction WHERE tr_transaction.rec_status = "1" AND tr_transaction.transaction_key = ` + key
	log.Println("========== QUERY GetTrTransaction ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionByField(c *TrTransaction, field string, value string) (int, error) {
	query := `SELECT tr_transaction.* FROM tr_transaction WHERE tr_transaction.rec_status = 1 AND tr_transaction.` + field + ` = '` + value + `'`
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateTrTransactionByKeyIn(params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := "UPDATE tr_transaction SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	inQuery := strings.Join(valueIn, ",")
	query += " WHERE tr_transaction." + fieldIn + " IN(" + inQuery + ")"

	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetTrTransactionIn(c *[]TrTransaction, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_transaction.* FROM 
				tr_transaction `
	query := query2 + " WHERE tr_transaction.rec_status = 1 AND tr_transaction." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTransactionByParamAndValueIn(c *[]TrTransaction, limit uint64, offset uint64,
	nolimit bool, params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := `SELECT
              tr_transaction.*
			  FROM tr_transaction`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
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

	if len(valueIn) > 0 {
		if len(whereClause) < 1 {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " WHERE tr_transaction." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND tr_transaction." + fieldIn + " IN(" + inQuery + ")"
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
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionDateRange(c *[]TrTransaction, params map[string]string, start string, end string) (int, error) {
	query := `SELECT
              tr_transaction.* FROM 
			  tr_transaction`
	query += " WHERE tr_transaction.trans_date >= " + start + " AND tr_transaction.trans_date <= " + end
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
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

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionOnProcess(c *[]TrTransaction, customerKey string) (int, error) {
	query := `SELECT 
				*
			FROM tr_transaction 
			WHERE trans_status_key IN (2,4,5,6,7,8) AND rec_status = 1 AND customer_key = ` + customerKey + `
			AND DATE_FORMAT(fn_AddDate(nav_date, 1),'%Y-%m-%d') >= DATE_FORMAT(NOW(),'%Y-%m-%d') ORDER BY trans_date DESC`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ParamBatchTrTransactionByKey(c *ParamBatchTrTransaction, transactionKey string) (int, error) {
	query := `SELECT
				p.product_code AS product_code,
				tt.type_code AS type_code,
				MONTH(t.nav_date) AS bulan,
				YEAR(t.nav_date) AS tahun,
				t.nav_date AS nav_date,
				p.product_key AS product_key,
				t.trans_type_key AS trans_type_key,
				t.transaction_key AS transaction_key,
				t.trans_date AS trans_date,
			    (SELECT batch_number AS bat FROM tr_transaction_batch ORDER BY batch_number DESC LIMIT 1) AS batch 
			FROM tr_transaction AS t
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key
			WHERE t.transaction_key = ` + transactionKey + ` LIMIT 1`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckTrTransactionLastProductCustomer(c *TrTransaction, customerKey string, productKey string, transKey string) (int, error) {
	query := `SELECT
				tr_transaction.* FROM 
				tr_transaction `
	query += " WHERE tr_transaction.rec_status = 1"
	query += " AND trans_status_key = 9"
	query += " AND customer_key = " + customerKey
	query += " AND product_key = " + productKey
	query += " AND transaction_key > " + transKey + " LIMIT 1"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckProductAllowRedmOrSwitching(c *[]ProductCheckAllowRedmSwtching, customerKey string, productKeyIn []string) (int, error) {

	inQuery := strings.Join(productKeyIn, ",")

	query := `SELECT product_key 
				FROM tr_transaction`
	query += " WHERE rec_status = 1"
	query += " AND trans_type_key IN (2,3)"
	query += " AND trans_status_key NOT IN (3,9)"
	query += " AND customer_key = " + customerKey
	query += " AND product_key IN(" + inQuery + ")"
	query += " GROUP BY product_key"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetTransactionCustomerHistory(c *[]TransactionCustomerHistory, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	dateFrom := ""
	dateTo := ""

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	query := `SELECT 
				p.product_name_alt AS product_name,
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				TO_BASE64(CONCAT(c.customer_key, ",", p.product_key, ",", "` + dateFrom + `", ",", "` + dateTo + `")) AS param_detail 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key
			WHERE t.trans_status_key = 9 AND t.rec_status = 1 AND tc.rec_status = 1`

	var present bool
	var condition string
	var conditionOrder string

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " GROUP BY t.customer_key, t.product_key"

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		conditionOrder += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			conditionOrder += " " + orderType
		}
	}
	query += conditionOrder

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountTransactionCustomerHistory(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT 
				count(t.transaction_key) AS count_data 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key
			WHERE t.trans_status_key = 9 AND t.rec_status = 1 AND tc.rec_status = 1`

	var whereClause []string
	var condition string
	dateFrom := ""
	dateTo := ""

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " GROUP BY t.customer_key, t.product_key"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetTransactionConsumenProduct(c *[]TransactionConsumenProduct, params map[string]string, paramsLike map[string]string, dateFrom string, dateTo string) (int, error) {
	query := `SELECT 
				t.transaction_key as transaction_key,
				t.trans_type_key as trans_type_key,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date, 
				ty.type_description as type_description,
				nav.nav_value as nav_value,
				tc.confirmed_unit as unit,
				(CASE
					WHEN t.total_amount IS NULL OR t.total_amount = 0 THEN tc.confirmed_amount
					ELSE t.total_amount
				END) AS gross_amount,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee_amount, 
				(CASE
					WHEN t.total_amount IS NULL OR t.total_amount = 0 THEN tc.confirmed_amount - (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount)
					ELSE (t.total_amount - (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount))
				END) AS net_amount 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON t.transaction_key = tc.transaction_key
			INNER JOIN tr_nav AS nav ON t.nav_date = nav.nav_date AND t.product_key = nav.product_key 
			INNER JOIN tr_transaction_type AS ty ON ty.trans_type_key = t.trans_type_key
			WHERE t.rec_status = 1 AND tc.rec_status = 1 AND t.trans_status_key = 9 `

	query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"

	var present bool
	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
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

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetDetailHeaderTransaksiCustomer(c *DetailHeaderTransaksiCustomer, dateFrom string, dateTo string, params map[string]string) (int, error) {
	query := `SELECT
				c.unit_holder_idno AS unit_holder,
				c.full_name AS full_name,
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN a.ifua_no IS NULL THEN ""
					ELSE a.ifua_no
				END) AS ifua_no, 
				DATE_FORMAT("` + dateFrom + `", '%d %M %Y') AS nav_date_from,
				DATE_FORMAT("` + dateTo + `", '%d %M %Y') AS nav_date_to, 
				CONCAT(p.product_code, " - ", p.product_name_alt) AS product_name  
			FROM tr_account AS a 
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			INNER JOIN ms_product AS p ON p.product_key = a.product_key
			INNER JOIN sc_user_login AS l ON l.customer_key = c.customer_key 
			INNER JOIN sc_user_dept AS d ON d.user_dept_key = l.user_dept_key 
			WHERE a.rec_status = 1 AND c.rec_status = 1 `

	var present bool
	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
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

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllTrTransactionPosting(c *[]TrTransaction, params map[string]string, valueIn []string, fieldIn string, isAll bool) (int, error) {
	query := `SELECT
              t.*
			  FROM tr_transaction as t
			  WHERE t.rec_status = 1 AND t.trans_status_key != 3`

	if isAll == false {
		query += " AND t.trans_type_key != 3"
	}
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "t."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"
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

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminLastAvgNav(c *NavValue, productKey string, customerKey string, date string) (int, error) {
	query := `SELECT 
				tc.avg_nav AS nav_value
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key AND tc.rec_status = '1'
			WHERE t.product_key = '` + productKey + `' 
			AND t.customer_key = '` + customerKey + `' 
			AND t.nav_date <= '` + date + `' 
			AND t.trans_status_key = '9'
			ORDER BY tc.tc_key DESC LIMIT 1`

	// Main query
	log.Println("========== QUERY GET LAST AVERAGE NAV ==========", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type DetailTransactionDataSentEmail struct {
	TransactionKey        uint64          `db:"transaction_key"           json:"transaction_key"`
	NavDate               string          `db:"nav_date"                  json:"nav_date"`
	NavValue              decimal.Decimal `db:"nav_value"                 json:"nav_value"`
	CustomerKey           uint64          `db:"customer_key"              json:"customer_key,omitempty"`
	InvestorType          uint64          `db:"investor_type"             json:"investor_type,omitempty"`
	TransTypeKey          uint64          `db:"trans_type_key"            json:"trans_type_key"`
	FullName              string          `db:"full_name"                 json:"full_name"`
	Cif                   *string         `db:"cif"                       json:"cif"`
	TransDate             string          `db:"trans_date"                json:"trans_date"`
	TransTime             string          `db:"trans_time"                json:"trans_time"`
	ProductName           string          `db:"product_name"              json:"product_name"`
	CurrencySymbol        string          `db:"currency_symbol"           json:"currency_symbol"`
	EntryMode             *uint64         `db:"entry_mode"                json:"entry_mode"`
	TransAmount           decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	TransUnit             decimal.Decimal `db:"trans_unit"                json:"trans_unit"`
	Fee                   decimal.Decimal `db:"fee"                       json:"fee"`
	PaymentMethod         *uint64         `db:"payment_method"            json:"payment_method"`
	PaymentMethodName     *string         `db:"payment_method_name"       json:"payment_method_name"`
	RekBankCustodian      *string         `db:"rek_bank_custodian"        json:"rek_bank_custodian"`
	NoRekBankCustomer     *string         `db:"no_rek_bank_customer"      json:"no_rek_bank_customer"`
	NameRekBankCustomer   *string         `db:"name_rek_bank_customer"    json:"name_rek_bank_customer"`
	CabangRekBankCustomer *string         `db:"cabang_rek_bank_customer"  json:"cabang_rek_bank_customer"`
	BankRekBankCustomer   *string         `db:"bank_rek_bank_customer"    json:"bank_rek_bank_customer"`
	Sales                 *string         `db:"sales"                     json:"sales"`
	SalesEmail            *string         `db:"sales_email"               json:"sales_email"`
	BuktiTransafer        *string         `db:"bukti_transafer"           json:"bukti_transafer"`
	ProductTujuan         *string         `db:"product_tujuan"            json:"product_tujuan"`
	UserLoginKey          string          `db:"user_login_key"            json:"user_login_key"`
	UloginEmail           string          `db:"ulogin_email"              json:"ulogin_email"`
	FlagNewSub            *uint8          `db:"flag_newsub"               json:"flag_newsub"`
}

func AdminDetailTransactionDataSentEmail(c *DetailTransactionDataSentEmail, tansactionKey string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date,
				COALESCE(nv.nav_value, 0) AS nav_value,
				t.customer_key,
				c.investor_type,
				t.flag_newsub,
				t.trans_type_key, 
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				DATE_FORMAT(t.trans_date, '%d %M %Y') AS trans_date,
				CONCAT(DATE_FORMAT(t.trans_date, '%H:%i'), " WIB") AS trans_time,
				p.product_name_alt AS product_name,
				cu.symbol as currency_symbol,
				t.entry_mode,
				t.trans_amount,
				t.trans_unit,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				t.payment_method,
				mp.lkp_name AS payment_method_name,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE CONCAT(ba.account_no, " - ", ba.account_holder_name)
				END) AS rek_bank_custodian,
				t.rec_image1 AS bukti_transafer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.account_no
				END) AS no_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.account_holder_name
				END) AS name_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.branch_name
				END) AS cabang_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE b.bank_name
				END) AS bank_rek_bank_customer,
				CONCAT(a.agent_code, " - ", a.agent_name) AS sales,
				a.agent_email AS sales_email,
				t.rec_image1 AS bukti_transafer,
				p_t.product_name_alt AS product_tujuan,
				ul.user_login_key,  
				ul.ulogin_email 
			FROM tr_transaction AS t 
			LEFT JOIN tr_nav AS nv ON nv.product_key = t.product_key AND nv.nav_date = t.nav_date AND nv.rec_status = "1" AND nv.nav_status = "231" 
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key
			LEFT JOIN ms_agent AS a ON a.agent_key = c.openacc_agent_key
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			LEFT join ms_currency as cu on cu.currency_key = p.currency_key 
			LEFT JOIN ms_currency AS cur ON cur.currency_key = p.currency_key
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN tr_transaction_bank_account AS tbank ON tbank.transaction_key = t.transaction_key
			LEFT JOIN ms_product_bank_account AS pbank ON pbank.prod_bankacc_key = tbank.prod_bankacc_key
			LEFT JOIN ms_customer_bank_account AS cbank ON cbank.cust_bankacc_key = tbank.cust_bankacc_key
			LEFT JOIN ms_bank_account AS ba ON ba.bank_account_key = pbank.bank_account_key 
			LEFT JOIN ms_bank_account AS ba_c ON ba_c.bank_account_key = cbank.bank_account_key 
			LEFT JOIN ms_bank AS b ON b.bank_key = ba_c.bank_key 
			LEFT JOIN tr_transaction AS t_ch ON t_ch.parent_key = t.transaction_key
			LEFT JOIN ms_product AS p_t ON p_t.product_key = t_ch.product_key 
			INNER JOIN sc_user_login as ul on ul.customer_key = t.customer_key 
			WHERE t.rec_status = 1 AND t.transaction_key = '` + tansactionKey + `'`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func TransactionProductCustomerVA(c *TrTransaction, productKey string, customerKey string) (int, error) {
	query := `SELECT
				t.*
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_settlement AS ts ON ts.transaction_key = t.transaction_key
			WHERE t.rec_status = "1" AND ts.rec_status = "1" 
			AND t.payment_method = "287" 
			AND ts.settled_status = "243" 
			AND t.product_key = "` + productKey + `" 
			AND t.customer_key = "` + customerKey + `" order by t.transaction_key DESC LIMIT 1`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type DetailTransactionVaBelumBayar struct {
	TransactionKey    uint64          `db:"transaction_key"           json:"transaction_key"`
	SettlementKey     uint64          `db:"settlement_key"            json:"settlement_key,omitempty"`
	TotalPembayaran   decimal.Decimal `db:"total_pembayaran"          json:"total_pembayaran"`
	NoVa              string          `db:"no_va"                     json:"no_va"`
	TransTypeKey      uint64          `db:"trans_type_key"            json:"trans_type_key"`
	FullName          string          `db:"full_name"                 json:"full_name"`
	Cif               *string         `db:"cif"                       json:"cif"`
	TransDate         string          `db:"trans_date"                json:"trans_date"`
	TransTime         string          `db:"trans_time"                json:"trans_time"`
	LastDayPay        string          `db:"last_day_pay"              json:"last_day_pay"`
	LastDatePay       string          `db:"last_date_pay"             json:"last_date_pay"`
	LastTimePay       string          `db:"last_time_pay"             json:"last_time_pay"`
	ProductName       string          `db:"product_name"              json:"product_name"`
	CurrencySymbol    string          `db:"currency_symbol"           json:"currency_symbol"`
	EntryMode         *uint64         `db:"entry_mode"                json:"entry_mode"`
	TransAmount       decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	Fee               decimal.Decimal `db:"fee"                       json:"fee"`
	PaymentMethod     *uint64         `db:"payment_method"            json:"payment_method"`
	PaymentMethodName *string         `db:"payment_method_name"       json:"payment_method_name"`
	Sales             *string         `db:"sales"                     json:"sales"`
	SalesEmail        *string         `db:"sales_email"               json:"sales_email"`
	UserLoginKey      string          `db:"user_login_key"            json:"user_login_key"`
	UloginEmail       string          `db:"ulogin_email"              json:"ulogin_email"`
	FlagNewSub        *uint8          `db:"flag_newsub"               json:"flag_newsub"`
	TokenNotif        *string         `db:"token_notif"               json:"token_notif,omitempty"`
	PchannelKey       *uint64         `db:"pchannel_key"              json:"pchannel_key,omitempty"`
}

func AdminDetailTransactionVaBelumBayar(c *DetailTransactionVaBelumBayar, tansactionKey string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				ts.settlement_key,
				ts.settle_nominal AS total_pembayaran,
				ts.client_subaccount_no AS no_va,
				t.flag_newsub,
				t.trans_type_key, 
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				DATE_FORMAT(t.trans_date, '%d %M %Y') AS trans_date,
				CONCAT(DATE_FORMAT(t.trans_date, '%H:%i'), " WIB") AS trans_time,
				(CASE DAYOFWEEK(ts.expired_date)
					WHEN 1 THEN 'Minggu'
					WHEN 2 THEN 'Senin'
					WHEN 3 THEN 'Selasa'
					WHEN 4 THEN 'Rabu'
					WHEN 5 THEN 'Kamis'
					WHEN 6 THEN 'Jumat'
					WHEN 7 THEN 'Sabtu'
				END) AS last_day_pay,
				CONCAT(
					DAY(ts.expired_date),
				" ",
					CASE MONTH(ts.expired_date) 
					WHEN 1 THEN 'Januari' 
					WHEN 2 THEN 'Februari' 
					WHEN 3 THEN 'Maret' 
					WHEN 4 THEN 'April' 
					WHEN 5 THEN 'Mei' 
					WHEN 6 THEN 'Juni' 
					WHEN 7 THEN 'Juli' 
					WHEN 8 THEN 'Agustus' 
					WHEN 9 THEN 'September'
					WHEN 10 THEN 'Oktober' 
					WHEN 11 THEN 'November' 
					WHEN 12 THEN 'Desember' 
				END, 
				" ",
				YEAR(ts.expired_date)
				) AS last_date_pay,
				CONCAT(
				DATE_FORMAT(ts.expired_date, '%H:%i'), 
				" WIB"
				) AS last_time_pay,
				p.product_name_alt AS product_name,
				cu.symbol AS currency_symbol,
				t.entry_mode,
				t.trans_amount,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				t.payment_method,
				mp.lkp_name AS payment_method_name,
				CONCAT(a.agent_code, " - ", a.agent_name) AS sales,
				a.agent_email AS sales_email,
				ul.user_login_key,  
				ul.ulogin_email 
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_settlement AS ts ON ts.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key
			LEFT JOIN ms_agent AS a ON a.agent_key = c.openacc_agent_key
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			LEFT JOIN ms_currency AS cu ON cu.currency_key = p.currency_key 
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN tr_transaction AS t_ch ON t_ch.parent_key = t.transaction_key
			INNER JOIN sc_user_login AS ul ON ul.customer_key = t.customer_key 
			WHERE t.rec_status = 1 AND ts.rec_status = 1 
			AND t.transaction_key = '` + tansactionKey + `'`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type TrTransactionExpired struct {
	TransactionKey uint64 `db:"transaction_key"           json:"transaction_key"`
	TransStatusKey uint64 `db:"trans_status_key"          json:"trans_status_key"`
	TransDate      string `db:"trans_date"                json:"trans_date"`
	TransTypeKey   uint64 `db:"trans_type_key"            json:"trans_type_key"`
	SettlementKey  uint64 `db:"settlement_key"            json:"settlement_key"`
	SettledStatus  uint64 `db:"settled_status"            json:"settled_status"`
}

func AdminGetTransactionExpired(c *[]TrTransactionExpired) (int, error) {
	query := `SELECT 
				t.transaction_key,
				t.trans_date,
				t.trans_status_key,
				t.trans_type_key,
				ts.settlement_key,
				ts.settled_status 
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_settlement AS ts ON ts.transaction_key = t.transaction_key 
			WHERE t.rec_status = 1 AND ts.rec_status AND ts.settled_status = 243 
			AND t.trans_status_key = 2 AND t.trans_type_key = 1
			AND ts.expired_date IS NOT NULL AND ts.expired_date <= NOW()`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetTransactionBeforeExpired(c *[]DetailTransactionVaBelumBayar, timeBefore string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				ts.settlement_key,
				ts.settle_nominal AS total_pembayaran,
				ts.client_subaccount_no AS no_va,
				t.flag_newsub,
				t.trans_type_key, 
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				DATE_FORMAT(t.trans_date, '%d %M %Y') AS trans_date,
				CONCAT(DATE_FORMAT(t.trans_date, '%H:%i'), " WIB") AS trans_time,
				(CASE DAYOFWEEK(ts.expired_date)
					WHEN 1 THEN 'Minggu'
					WHEN 2 THEN 'Senin'
					WHEN 3 THEN 'Selasa'
					WHEN 4 THEN 'Rabu'
					WHEN 5 THEN 'Kamis'
					WHEN 6 THEN 'Jumat'
					WHEN 7 THEN 'Sabtu'
				END) AS last_day_pay,
				CONCAT(
					DAY(ts.expired_date),
				" ",
					CASE MONTH(ts.expired_date) 
					WHEN 1 THEN 'Januari' 
					WHEN 2 THEN 'Februari' 
					WHEN 3 THEN 'Maret' 
					WHEN 4 THEN 'April' 
					WHEN 5 THEN 'Mei' 
					WHEN 6 THEN 'Juni' 
					WHEN 7 THEN 'Juli' 
					WHEN 8 THEN 'Agustus' 
					WHEN 9 THEN 'September'
					WHEN 10 THEN 'Oktober' 
					WHEN 11 THEN 'November' 
					WHEN 12 THEN 'Desember' 
				END, 
				" ",
				YEAR(ts.expired_date)
				) AS last_date_pay,
				CONCAT(
				DATE_FORMAT(ts.expired_date, '%H:%i'), 
				" WIB"
				) AS last_time_pay,
				p.product_name_alt AS product_name,
				cu.symbol AS currency_symbol,
				t.entry_mode,
				t.trans_amount,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				t.payment_method,
				mp.lkp_name AS payment_method_name,
				CONCAT(a.agent_code, " - ", a.agent_name) AS sales,
				a.agent_email AS sales_email,
				ul.user_login_key,  
				ul.ulogin_email, 
				ts.settle_payment_method AS pchannel_key 
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_settlement AS ts ON ts.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key
			LEFT JOIN ms_agent AS a ON a.agent_key = c.openacc_agent_key
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			LEFT JOIN ms_currency AS cu ON cu.currency_key = p.currency_key 
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN tr_transaction AS t_ch ON t_ch.parent_key = t.transaction_key
			INNER JOIN sc_user_login AS ul ON ul.customer_key = t.customer_key 
			WHERE t.rec_status = 1 AND ts.rec_status = 1 
			AND t.trans_status_key = 2 
			AND ts.settled_status = 243 
			AND t.trans_type_key = 1 
			AND ts.rec_attribute_id1 IS NULL 
			AND ts.expired_date IS NOT NULL 
			AND ts.expired_date > NOW() 
			AND ((ts.expired_date - NOW()) <= ` + timeBefore + `)`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func search(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type ListTransactionInstitution struct {
	TransactionKey         uint64           `db:"transaction_key"                 json:"transaction_key"`
	TransTypeKey           uint64           `db:"trans_type_key"                  json:"trans_type_key"`
	TransType              string           `db:"trans_type"                      json:"trans_type"`
	ProductName            string           `db:"product_name"                    json:"product_name"`
	EntryMode              *uint64          `db:"entry_mode"                      json:"entry_mode"`
	TransAmount            decimal.Decimal  `db:"trans_amount"                    json:"trans_amount"`
	TransUnit              *decimal.Decimal `db:"trans_unit"                      json:"trans_unit"`
	TransactionDate        string           `db:"transaction_date"                json:"transaction_date"`
	TransactionTime        string           `db:"transaction_time"                json:"transaction_time"`
	BuktiTransfer          *string          `db:"bukti_transfer"                  json:"bukti_transfer"`
	RecApprovalStage       *string          `db:"rec_approval_stage"              json:"rec_approval_stage"`
	RecApprovalStatus      *string          `db:"rec_approval_status"             json:"rec_approval_status"`
	FlagNewsub             *uint64          `db:"flag_newsub"                     json:"flag_newsub"`
	LastStage              *string          `db:"last_stage"                      json:"last_stage"`
	PaymentMethod          *string          `db:"payment_method"                  json:"payment_method"`
	StatusTransaction      *string          `db:"status_transaction"              json:"status_transaction"`
	ProductBankName        *string          `db:"product_bank_name"               json:"product_bank_name"`
	ProductBankNoAccount   *string          `db:"product_bank_no_account"         json:"product_bank_no_account"`
	ProductBankNameAccount *string          `db:"product_bank_name_account"       json:"product_bank_name_account"`
	CustBankName           *string          `db:"cust_bank_name"                  json:"cust_bank_name"`
	CustBankNoAccount      *string          `db:"cust_bank_no_account"            json:"cust_bank_no_account"`
	CustBankNameAccount    *string          `db:"cust_bank_name_account"          json:"cust_bank_name_account"`
	Remarks                *string          `db:"remarks"                         json:"remarks"`
	ProductTujuan          *string          `db:"product_tujuan"                  json:"product_tujuan"`
	Currency               *string          `db:"currency"                        json:"currency"`
}

func GetTransactionInstitution(c *[]ListTransactionInstitution, trStatusKeyIn []string, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	statusInQuery := strings.Join(trStatusKeyIn, ",")

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query := `SELECT 
				t.transaction_key AS transaction_key, 
				t.trans_type_key AS trans_type_key, 
				(CASE 
					WHEN t.trans_type_key = 1 AND t.flag_newsub = 0 THEN "Top Up" 
					WHEN t.trans_type_key = 3 THEN "Switching" 
					ELSE tt.type_description 
				END) AS trans_type, 
				p.product_name_alt AS product_name, 
				t.entry_mode AS entry_mode, 
				t.trans_amount AS trans_amount, 
				t.trans_unit AS trans_unit, 
				DATE_FORMAT(t.trans_date, '%d %M %Y') AS transaction_date, 
				CONCAT(DATE_FORMAT(t.trans_date, '%H:%i'), " WIB") AS transaction_time, 
				t.rec_image1 AS bukti_transfer, 
				t.rec_approval_stage AS rec_approval_stage, 
				t.rec_approval_status AS rec_approval_status, 
				t.flag_newsub AS flag_newsub, 
				s.stage_code AS last_stage, 
				mp.lkp_name AS payment_method, 
				gl.lkp_name AS status_transaction, 
				bproduct.bank_name AS product_bank_name, 
				bankaccproduct.account_no AS product_bank_no_account, 
				bankaccproduct.account_holder_name AS product_bank_name_account , 
				bcust.bank_name AS cust_bank_name, 
				bankacccust.account_no AS cust_bank_no_account, 
				bankacccust.account_holder_name AS cust_bank_name_account, 
				t.trans_remarks AS remarks, 
				parentproduct.product_name_alt AS product_tujuan, 
				cur.symbol AS currency 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key 
			LEFT JOIN tr_transaction_bank_account AS tba ON tba.transaction_key = t.transaction_key 
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN ms_product_bank_account AS bankproduct ON bankproduct.prod_bankacc_key = tba.prod_bankacc_key 
			LEFT JOIN ms_bank_account AS bankaccproduct ON bankaccproduct.bank_account_key = bankproduct.bank_account_key 
			LEFT JOIN ms_bank AS bproduct ON bproduct.bank_key = bankaccproduct.bank_key 
			LEFT JOIN ms_customer_bank_account AS bankcust ON bankcust.cust_bankacc_key = tba.cust_bankacc_key  
			LEFT JOIN ms_bank_account AS bankacccust ON bankacccust.bank_account_key = bankcust.bank_account_key 
			LEFT JOIN ms_bank AS bcust ON bcust.bank_key = bankacccust.bank_key 
			LEFT JOIN tr_transaction AS tp ON tp.parent_key = t.transaction_key 
			LEFT JOIN ms_product AS parentproduct ON parentproduct.product_key = tp.product_key 
			LEFT JOIN wf_stage AS s ON s.wf_stage_key = t.rec_approval_stage 
			LEFT JOIN gen_lookup AS gl ON gl.lookup_key = t.rec_approval_status 
			LEFT JOIN ms_currency AS cur ON cur.currency_key = p.currency_key 
			WHERE t.trans_type_key IN (1, 2, 3) 
			AND t.rec_status = 1 
			AND t.trans_status_key IN (` + statusInQuery + `) ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCountTransactionInstitution(c *CountData, trStatusKeyIn []string, params map[string]string) (int, error) {
	var whereClause []string
	var condition string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	statusInQuery := strings.Join(trStatusKeyIn, ",")

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query := `SELECT 
				count(t.transaction_key) AS count_data 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key 
			WHERE t.trans_type_key IN (1, 2, 3) 
			AND t.rec_status = 1 
			AND t.trans_status_key IN (` + statusInQuery + `) ` + condition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type ListMutasiTransactionInstitution struct {
	TransactionKey        uint64           `db:"transaction_key"                 json:"transaction_key"`
	TransTypeKey          uint64           `db:"trans_type_key"                  json:"trans_type_key"`
	TransType             string           `db:"trans_type"                      json:"trans_type"`
	ProductName           string           `db:"product_name"                    json:"product_name"`
	EntryMode             string           `db:"entry_mode"                      json:"entry_mode"`
	TransAmount           decimal.Decimal  `db:"trans_amount"                    json:"trans_amount"`
	TotalAmount           decimal.Decimal  `db:"total_amount"                    json:"total_amount"`
	ConfirmedAmount       decimal.Decimal  `db:"confirmed_amount"                json:"confirmed_amount"`
	ConfirmedUnit         decimal.Decimal  `db:"confirmed_unit"                  json:"confirmed_unit"`
	TransactionDate       string           `db:"transaction_date"                json:"transaction_date"`
	NavDate               string           `db:"nav_date"                        json:"nav_date"`
	NavValue              decimal.Decimal  `db:"nav_value"                       json:"nav_value"`
	BuktiTransfer         *string          `db:"bukti_transfer"                  json:"bukti_transfer"`
	FlagNewsub            *uint64          `db:"flag_newsub"                     json:"flag_newsub"`
	PaymentMethod         *string          `db:"payment_method"                  json:"payment_method"`
	Remarks               *string          `db:"remarks"                         json:"remarks"`
	ProductTujuan         *string          `db:"product_tujuan"                  json:"product_tujuan"`
	ConfirmedAmountParent *decimal.Decimal `db:"confirmed_amount_parent"         json:"confirmed_amount_parent"`
	ConfirmedUnitParent   *decimal.Decimal `db:"confirmed_unit_parent"           json:"confirmed_unit_parent"`
	NavValueParent        *decimal.Decimal `db:"nav_value_parent"                json:"nav_value_parent"`
	Currency              *string          `db:"currency"                        json:"currency"`
}

func GetMutasiTransactionInstitution(c *[]ListMutasiTransactionInstitution, productKey string, dateFrom string, dateTo string, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	if productKey != "" {
		condition += " AND (tp.product_key = '" + productKey + "' or t.product_key = '" + productKey + "') "
	}

	if (dateFrom != "") && (dateTo != "") {
		condition += " AND (t.trans_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "') "
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query := `SELECT 
				t.transaction_key AS transaction_key, 
				t.trans_type_key AS trans_type_key, 
				(CASE 
					WHEN t.trans_type_key = 1 AND t.flag_newsub = 0 THEN "Top Up" 
					WHEN t.trans_type_key = 3 THEN "Switching" 
					ELSE tt.type_description 
				END) AS trans_type, 
				p.product_name_alt AS product_name, 
				(CASE 
					WHEN t.entry_mode IS NOT NULL AND t.entry_mode = 140 THEN "AMOUNT" 
					ELSE "UNIT"
				END) AS entry_mode,  
				t.trans_amount AS trans_amount,
				t.total_amount AS total_amount,
				tc.confirmed_amount AS confirmed_amount, 
				tc.confirmed_unit AS confirmed_unit, 
				DATE_FORMAT(t.trans_date, '%d-%m-%Y') AS transaction_date, 
				DATE_FORMAT(t.nav_date, '%d-%m-%Y') AS nav_date, 
				tn.nav_value AS nav_value,
				t.rec_image1 AS bukti_transfer, 
				t.flag_newsub AS flag_newsub, 
				mp.lkp_name AS payment_method, 
				t.trans_remarks AS remarks, 
				parentproduct.product_name_alt AS product_tujuan, 
				tcp.confirmed_amount AS confirmed_amount_parent, 
				tcp.confirmed_unit AS confirmed_unit_parent, 
				tnp.nav_value AS nav_value_parent,
				cur.symbol AS currency 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key 
			INNER JOIN tr_nav AS tn ON tn.product_key = t.product_key AND tn.rec_status = 1 AND tn.nav_date = t.nav_date AND tn.nav_status = "231" 
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN tr_transaction AS tp ON tp.parent_key = t.transaction_key 
			LEFT JOIN tr_transaction_confirmation AS tcp ON tcp.transaction_key = tp.transaction_key 
			LEFT JOIN tr_nav AS tnp ON tnp.product_key = tp.product_key AND tnp.rec_status = 1 AND tnp.nav_date = tp.nav_date AND tnp.nav_status = "231" 
			LEFT JOIN ms_product AS parentproduct ON parentproduct.product_key = tp.product_key 
			LEFT JOIN ms_currency AS cur ON cur.currency_key = p.currency_key 
			WHERE t.trans_type_key IN (1, 2, 3) 
			AND t.rec_status = 1 AND tc.rec_status = 1  
			AND t.trans_status_key IN(9) ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCountMutasiTransactionInstitution(c *CountData, productKey string, dateFrom string, dateTo string, params map[string]string) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	if productKey != "" {
		condition += " AND (tp.product_key = '" + productKey + "' or t.product_key = '" + productKey + "') "
	}

	if (dateFrom != "") && (dateTo != "") {
		condition += " AND (t.trans_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "') "
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query := `SELECT 
				count(t.transaction_key) AS count_data 
			FROM tr_transaction AS t 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key 
			INNER JOIN tr_nav AS tn ON tn.product_key = t.product_key AND tn.rec_status = 1 AND tn.nav_date = t.nav_date AND tn.nav_status = "231" 
			LEFT JOIN tr_transaction AS tp ON tp.parent_key = t.transaction_key 
			WHERE t.trans_type_key IN (1, 2, 3) 
			AND t.rec_status = 1 AND tc.rec_status = 1 
			AND t.trans_status_key IN(9) ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	query += orderCondition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type AdminTransactionCorrection struct {
	TransactionKey    uint64          `db:"transaction_key"           json:"transaction_key"`
	CustomerName      string          `db:"customer_name"             json:"customer_name"`
	ProductName       string          `db:"product_name"              json:"product_name"`
	BranchName        string          `db:"branch_name"               json:"branch_name"`
	AgentName         string          `db:"agent_name"                json:"agent_name"`
	StatusCode        string          `db:"status_code"               json:"status_code"`
	TransDate         string          `db:"trans_date"                json:"trans_date"`
	NavDate           string          `db:"nav_date"                  json:"nav_date"`
	TypeDescription   string          `db:"type_description"          json:"type_description"`
	TransTypeKey      uint64          `db:"trans_type_key"            json:"trans_type_key"`
	TransAmount       decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	TransUnit         decimal.Decimal `db:"trans_unit"                json:"trans_unit"`
	TotalAmount       decimal.Decimal `db:"total_amount"              json:"total_amount"`
	PaymentMethod     *string         `db:"payment_method"            json:"payment_method"`
	TransSource       *string         `db:"trans_source"              json:"trans_source"`
	ProductNameTujuan *string         `db:"product_name_tujuan"       json:"product_name_tujuan"`
}

func AdminGetListTransactionCorrection(c *[]AdminTransactionCorrection, limit uint64, offset uint64, params map[string]string, nolimit bool, productKey string, userId string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				c.full_name AS customer_name,
				p.product_name_alt AS product_name,
				b.branch_name,
				a.agent_name,
				tts.status_code,
				DATE_FORMAT(t.trans_date, '%d %M %Y %H:%m') AS trans_date,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date,
				t.trans_type_key,
				(CASE 
					WHEN t.trans_type_key IN (1,2) THEN ttt.type_description 
					ELSE "Switching"
				END) AS type_description,
				t.trans_amount,
				t.trans_unit,
				t.total_amount,
				pm.lkp_name AS payment_method,
				ts.lkp_name AS trans_source,
				(CASE 
					WHEN par_pro.product_key IS NOT NULL THEN par_pro.product_name_alt 
					ELSE NULL 
				END) product_name_tujuan 
			FROM tr_transaction AS t 
			LEFT JOIN tr_transaction AS t_par ON t_par.parent_key = t.transaction_key  
			LEFT JOIN ms_product AS par_pro ON par_pro.product_key = t_par.product_key 
			INNER JOIN tr_transaction_status AS tts ON tts.trans_status_key = t.trans_status_key 
			INNER JOIN ms_branch AS b ON b.branch_key = t.branch_key 
			INNER JOIN ms_agent AS a ON a.agent_key = t.agent_key 
			INNER JOIN tr_transaction_type AS ttt ON ttt.trans_type_key = t.trans_type_key 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			LEFT JOIN gen_lookup AS pm ON pm.lookup_key = t.payment_method 
			LEFT JOIN gen_lookup AS ts ON ts.lookup_key = t.trans_source  
			LEFT JOIN sc_user_login AS ul ON t.rec_created_by = ul.user_login_key 
			WHERE t.trans_type_key IN (1,2,3) AND t.rec_status = 1 AND t.trans_status_key = 1`

	var present bool
	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if productKey != "" {
		condition += " AND (t.product_key = " + productKey + " OR t_par.product_key = " + productKey + ") "
	}

	if userId != "" {
		condition += " AND (ul.user_login_key IS NULL OR ul.user_category_key = '1' OR ul.user_login_key = " + userId + ")"
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

	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetListTransactionCorrection(c *CountData, params map[string]string, productKey string, userId string) (int, error) {
	query := `SELECT 
				count(t.transaction_key) as count_data 
			FROM tr_transaction AS t 
			LEFT JOIN tr_transaction AS t_par ON t_par.parent_key = t.transaction_key  
			LEFT JOIN ms_product AS par_pro ON par_pro.product_key = t_par.product_key 
			INNER JOIN tr_transaction_status AS tts ON tts.trans_status_key = t.trans_status_key 
			INNER JOIN ms_branch AS b ON b.branch_key = t.branch_key 
			INNER JOIN ms_agent AS a ON a.agent_key = t.agent_key 
			INNER JOIN tr_transaction_type AS ttt ON ttt.trans_type_key = t.trans_type_key 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key 
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			LEFT JOIN gen_lookup AS pm ON pm.lookup_key = t.payment_method 
			LEFT JOIN gen_lookup AS ts ON ts.lookup_key = t.trans_source  
			LEFT JOIN sc_user_login AS ul ON t.rec_created_by = ul.user_login_key 
			WHERE t.trans_type_key IN (1,2,3) AND t.rec_status = 1 AND t.trans_status_key = 1`

	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if productKey != "" {
		condition += " AND (t.product_key = " + productKey + " OR t_par.product_key = " + productKey + ") "
	}

	if userId != "" {
		condition += " AND (ul.user_login_key IS NULL OR ul.user_category_key = '1' OR ul.user_login_key = " + userId + ")"
	}

	query += condition

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckProductAllowRedmOrSwitchingInUpdate(c *[]ProductCheckAllowRedmSwtching, customerKey string, productKeyIn []string, transactionKey string) (int, error) {

	inQuery := strings.Join(productKeyIn, ",")

	query := `SELECT product_key 
				FROM tr_transaction`
	query += " WHERE rec_status = 1"
	query += " AND trans_type_key IN (2,3)"
	query += " AND trans_status_key NOT IN (3,9)"
	query += " AND transaction_key != " + transactionKey
	query += " AND customer_key = " + customerKey
	query += " AND product_key IN(" + inQuery + ")"
	query += " GROUP BY product_key"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
