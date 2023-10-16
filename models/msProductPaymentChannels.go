package models

import (
	"log"
	"mf-bo-api/db"
)

type ProductPaymentChannels struct {
	ProdChannelKey uint64 `db:"prod_channel_key" json:"prod_channel_key"`
	ProductKey     uint64 `db:"product_key" json:"product_key"`
	PChannelKey    uint64 `db:"pchannel_key" json:"pchannel_key"`
	PChannelName   string `db:"pchannel_name" json:"pchannel_name"`
	CotSettlement  string `db:"cot_settlement" json:"cot_settlement"`
	CotTransaction string `db:"cot_transaction" json:"cot_transaction"`
}

func GetProductPaymentChannels(productKey string) (result []ProductPaymentChannels) {
	query := `SELECT t1.prod_channel_key, t1.product_key, t1.pchannel_key, 
	t2.pchannel_name, t1.cot_settlement, t1.cot_transaction
	FROM ms_product_channel t1
	INNER JOIN ms_payment_channel t2 ON (t1.pchannel_key = t2.pchannel_key AND t2.rec_status = 1)
	WHERE t1.rec_status = 1 AND t1.product_key = ` + productKey

	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	return
}
