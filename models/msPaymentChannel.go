package models

import (
	"mf-bo-api/config"
	"mf-bo-api/db"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type MsPaymentChannel struct {
	PchannelKey         uint64           `db:"pchannel_key"             json:"pchannel_key"`
	PchannelCode        *string          `db:"pchannel_code"            json:"pchannel_code"`
	PchannelName        *string          `db:"pchannel_name"            json:"pchannel_name"`
	SettleChannel       uint64           `db:"settle_channel"           json:"settle_channel"`
	SettlePaymentMethod uint64           `db:"settle_payment_method"    json:"settle_payment_method"`
	MinNominalTrx       *decimal.Decimal `db:"min_nominal_trx"          json:"min_nominal_trx"`
	ValueType           uint64           `db:"value_type"               json:"value_type"`
	FeeValue            decimal.Decimal  `db:"fee_value"                json:"fee_value"`
	HasMinMax           uint8            `db:"has_min_max"              json:"has_min_max"`
	FeeMinValue         *decimal.Decimal `db:"fee_min_value"            json:"fee_min_value"`
	FeeMaxValue         *decimal.Decimal `db:"fee_max_value"            json:"fee_max_value"`
	FixedAmountFee      *decimal.Decimal `db:"fixed_amount_fee"         json:"fixed_amount_fee"`
	FixedDmrFee         *decimal.Decimal `db:"fixed_dmr_fee" json:"fixed_dmr_fee"`
	PgTnc               *string          `db:"pg_tnc"                   json:"pg_tnc"`
	PgRemarks           *string          `db:"pg_remarks"               json:"pg_remarks"`
	PaymentLoginUrl     *string          `db:"payment_login_url"        json:"payment_login_url"`
	PaymentEntryUrl     *string          `db:"payment_entry_url"        json:"payment_entry_url"`
	PaymentErrorUrl     *string          `db:"payment_error_url"        json:"payment_error_url"`
	PaymentSuccessUrl   *string          `db:"payment_success_url"      json:"payment_success_url"`
	PgPrefix            *string          `db:"pg_prefix"                json:"pg_prefix"`
	PicName             *string          `db:"pic_name"                 json:"pic_name"`
	PicPhoneNo          *string          `db:"pic_phone_no"             json:"pic_phone_no"`
	PicEmailAddress     *string          `db:"pic_email_address"        json:"pic_email_address"`
	RecOrder            *uint64          `db:"rec_order"                json:"rec_order"`
	RecStatus           uint8            `db:"rec_status"               json:"rec_status"`
	RecCreatedDate      *string          `db:"rec_created_date"         json:"rec_created_date"`
	RecCreatedBy        *string          `db:"rec_created_by"           json:"rec_created_by"`
	RecModifiedDate     *string          `db:"rec_modified_date"        json:"rec_modified_date"`
	RecModifiedBy       *string          `db:"rec_modified_by"          json:"rec_modified_by"`
	RecImage1           *string          `db:"rec_image1"               json:"rec_image1"`
	RecImage2           *string          `db:"rec_image2"               json:"rec_image2"`
	RecApprovalStatus   *uint8           `db:"rec_approval_status"      json:"rec_approval_status"`
	RecApprovalStage    *uint64          `db:"rec_approval_stage"       json:"rec_approval_stage"`
	RecApprovedDate     *string          `db:"rec_approved_date"        json:"rec_approved_date"`
	RecApprovedBy       *string          `db:"rec_approved_by"          json:"rec_approved_by"`
	RecDeletedDate      *string          `db:"rec_deleted_date"         json:"rec_deleted_date"`
	RecDeletedBy        *string          `db:"rec_deleted_by"           json:"rec_deleted_by"`
	RecAttributeID1     *string          `db:"rec_attribute_id1"        json:"rec_attribute_id1"`
	RecAttributeID2     *string          `db:"rec_attribute_id2"        json:"rec_attribute_id2"`
	RecAttributeID3     *string          `db:"rec_attribute_id3"        json:"rec_attribute_id3"`
}

type SubscribePaymentChannel struct {
	PchannelKey    uint64           `db:"pchannel_key"             json:"pchannel_key"`
	PchannelCode   *string          `db:"pchannel_code"            json:"pchannel_code"`
	PchannelName   *string          `db:"pchannel_name"            json:"pchannel_name"`
	MinNominalTrx  *decimal.Decimal `db:"min_nominal_trx"          json:"min_nominal_trx"`
	ValueType      uint64           `db:"value_type"               json:"value_type"`
	FeeValue       decimal.Decimal  `db:"fee_value"                json:"fee_value"`
	HasMinMax      uint8            `db:"has_min_max"              json:"has_min_max"`
	FeeMinValue    *decimal.Decimal `db:"fee_min_value"            json:"fee_min_value"`
	FeeMaxValue    *decimal.Decimal `db:"fee_max_value"            json:"fee_max_value"`
	FixedAmountFee *decimal.Decimal `db:"fixed_amount_fee"         json:"fixed_amount_fee"`
	PgTNC          *string          `db:"pg_tnc"                   json:"pg_tnc"`
	Channel        *string          `db:"channel"                  json:"channel"`
	Method         *string          `db:"method"                   json:"method"`
	Logo           string           `db:"logo"                     json:"logo"`
	InUse          bool             `db:"in_use"                   json:"in_use"`
}

type MsPaymentChannelDropdown struct {
	PchannelKey  uint64  `db:"pchannel_key"             json:"pchannel_key"`
	PchannelCode *string `db:"pchannel_code"            json:"pchannel_code"`
	PchannelName *string `db:"pchannel_name"            json:"pchannel_name"`
}

func GetAllMsPaymentChannel(c *[]MsPaymentChannel, params map[string]string) (int, error) {
	query := `SELECT
              ms_payment_channel.* FROM 
			  ms_payment_channel `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_payment_channel."+field+" = '"+value+"'")
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
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsPaymentChannel(c *MsPaymentChannel, key string) (int, error) {
	query := `SELECT ms_payment_channel.* FROM ms_payment_channel WHERE ms_payment_channel.pchannel_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsPaymentChannelIn(c *[]MsPaymentChannel, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
	ms_payment_channel.* FROM 
	ms_payment_channel `
	query := query2 + " WHERE ms_payment_channel." + field + " IN(" + inQuery + ")"

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetPaymentChannelByCusomerKey(c *[]SubscribePaymentChannel, product string, customer string) (int, error) {
	query := `SELECT 
			c.pchannel_key,
			c.pchannel_name,
			c.pchannel_code,
			c.min_nominal_trx,
			c.value_type,
			c.fee_value,
			c.has_min_max,
			c.fee_min_value,
			c.fee_max_value,
			c.fixed_amount_fee,
			c.pg_tnc,
			lc.lkp_name AS channel,
			lm.lkp_name AS method,
			(CASE
				WHEN c.rec_image1 IS NULL THEN ""
				ELSE CONCAT("` + config.ImageUrl + `", "/images/payment/", c.rec_image1)
			END) AS logo, 
			(CASE
				WHEN ts.settle_payment_method IS NULL OR c.pchannel_key != 9 THEN "false"
				ELSE "true"
			END) AS in_use FROM 
			ms_payment_channel AS c 
			INNER JOIN  ms_product_channel AS p ON c.pchannel_key = p.pchannel_key 
			LEFT JOIN (SELECT s.settle_payment_method FROM tr_transaction_settlement AS s
			INNER JOIN tr_transaction AS t ON s.transaction_key = t.transaction_key 
			WHERE t.customer_key = ` + customer + ` AND s.settled_status = 243 GROUP BY s.settle_payment_method) AS ts ON ts.settle_payment_method = c.pchannel_key
			LEFT JOIN gen_lookup AS lc ON lc.lookup_key = c.settle_channel
			LEFT JOIN gen_lookup AS lm ON lm.lookup_key = c.settle_payment_method
			WHERE p.product_key = ` + product

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
