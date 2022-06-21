package lib

import (
	"fmt"
	"mf-bo-api/models"
	"strconv"
	"time"
)

func ExpiredTransaction() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON EXPIRED TRANSACTION" + time.Now().Format(dateLayout))
	var err error

	var trans []models.TrTransactionExpired
	_, err = models.AdminGetTransactionExpired(&trans)

	if err == nil && len(trans) > 0 {
		var settleIds []string
		for _, tr := range trans { //TRANS
			if _, ok := Find(settleIds, strconv.FormatUint(tr.SettlementKey, 10)); !ok {
				settleIds = append(settleIds, strconv.FormatUint(tr.SettlementKey, 10))
			}
		}
		if len(settleIds) > 0 {
			fmt.Println("jml data : ")
			fmt.Println(len(settleIds))
			paramsSett := make(map[string]string)
			paramsSett["settled_status"] = "301"
			paramsSett["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsSett["rec_modified_by"] = "CRON"
			_, err = models.UpdateTrTransactionSettlementExpired(paramsSett, settleIds)
			if err != nil {
				fmt.Println("Error update tr_transaction_settlement")
				fmt.Println(err.Error())
			} else {
				fmt.Println("SUKSES UPDATE DATA EXPIRED")
			}
		}
	} else {
		fmt.Println("NO DATA TRANSACTION EXPIRED")
	}

	fmt.Println("======END CRON EXPIRED TRANSACTION============")
}
