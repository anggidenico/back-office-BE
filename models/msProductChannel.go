package models

import (
	"mf-bo-api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsProductChannel struct {
	ProdChannelKey    uint64  `db:"prod_channel_key"       json:"prod_channel_key"`
	ProductKey        uint64  `db:"product_key"            json:"product_key"`
	PchannelKey       uint64  `db:"pchannel_key"           json:"pchannel_key"`
	PrefixCode        *string `db:"prefix_code"            json:"prefix_code"`
	PrefixCodeAlt     *string `db:"prefix_code_alt"        json:"prefix_code_alt"`
	ExtProdCode       *string `db:"ext_prod_code"          json:"ext_prod_code"`
	ExtProdName       *string `db:"ext_prod_name"          json:"ext_prod_name"`
	BankAccountKey    *uint8  `db:"bank_account_key"       json:"bank_account_key"`
	CotTransaction    string  `db:"cot_transaction"        json:"cot_transaction"`
	CotSettlement     string  `db:"cot_settlement"         json:"cot_settlement"`
	RecOrder          *uint64 `db:"rec_order"              json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"             json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"             json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func GetAllMsProductChannel(c *[]MsProductChannel, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_channel.* FROM 
			  ms_product_channel `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_channel."+field+" = '"+value+"'")
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
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsProductChannel(c *MsProductChannel, key string) (int, error) {
	query := `SELECT ms_product_channel.* FROM ms_product_channel WHERE ms_product_channel.prod_channel_key = ` + key
	log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductChannelIn(c *[]MsBank, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
	ms_product_channel.* FROM 
	ms_product_channel `
	query := query2 + " WHERE ms_product_channel." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
