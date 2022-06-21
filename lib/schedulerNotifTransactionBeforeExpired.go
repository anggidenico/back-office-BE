package lib

import (
	"bytes"
	"fmt"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"strconv"
	"text/template"
	"time"

	"github.com/leekchan/accounting"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func NotifTransactionBeforeExpired() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON NOTIF TRANSAKSI BEFORE EXPIRED")
	var err error

	timeBefore := "20000" //2 jam
	var trans []models.DetailTransactionVaBelumBayar
	_, err = models.AdminGetTransactionBeforeExpired(&trans, timeBefore)
	if err == nil {
		if len(trans) > 0 {
			//insert message
			var bindVarMessage []interface{}
			var playerIds []string
			var settleIds []string
			var transVaMandiri []string

			heading := "Segera Lakukan Pembayaran"
			content := "Segera Lakukan Pembayaran Kamu agar Transaksi Kamu dapat segera kami proses."
			for _, cus := range trans { //customer
				if cus.PchannelKey != nil {
					if *cus.PchannelKey == uint64(9) { //FM_VA_MANDIRI
						if _, ok := Find(transVaMandiri, strconv.FormatUint(cus.TransactionKey, 10)); !ok {
							transVaMandiri = append(transVaMandiri, strconv.FormatUint(cus.TransactionKey, 10))
						}
					}
				}
				if _, ok := Find(settleIds, strconv.FormatUint(cus.SettlementKey, 10)); !ok {
					settleIds = append(settleIds, strconv.FormatUint(cus.SettlementKey, 10))
				}
				var row []string
				row = append(row, "247")                         //umessage_type
				row = append(row, "4")                           //notif_hdr_key
				row = append(row, cus.UserLoginKey)              //umessage_recipient_key
				row = append(row, time.Now().Format(dateLayout)) //umessage_receipt_date
				row = append(row, "0")                           //flag_read
				row = append(row, "1")                           //flag_sent
				row = append(row, heading)                       //umessage_subject
				row = append(row, content)                       //umessage_body
				row = append(row, "249")                         //umessage_category
				row = append(row, "0")                           //flag_archieved
				row = append(row, time.Now().Format(dateLayout)) //archieved_date
				row = append(row, "1")                           //rec_status
				row = append(row, time.Now().Format(dateLayout)) //rec_created_date
				row = append(row, "CRON")                        //rec_created_by
				bindVarMessage = append(bindVarMessage, row)

				if cus.TokenNotif != nil {
					if _, ok := Find(playerIds, *cus.TokenNotif); !ok {
						playerIds = append(playerIds, *cus.TokenNotif)
					}
				}
			}

			//create message
			_, err = models.CreateMultipleUserMessageFromUserNotif(bindVarMessage)
			if err != nil {
				fmt.Println("err create multiple user message : " + err.Error())
			} else {
				fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
			}

			//push notif
			if len(playerIds) > 0 {
				DataNotif := make(map[string]interface{})
				DataNotif["category"] = "MESSAGE"
				BlastAllNotificationHelper(playerIds, heading, content, DataNotif)
			}

			//update sudah pernah kirim notif
			if len(settleIds) > 0 {
				fmt.Println("jml data : ")
				fmt.Println(len(settleIds))
				paramsSett := make(map[string]string)
				paramsSett["rec_attribute_id1"] = "1"
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

			//send email VA mandiri payment
			if len(transVaMandiri) > 0 {

				for _, trID := range transVaMandiri {
					var err error
					var transaction models.DetailTransactionVaBelumBayar
					_, err = models.AdminDetailTransactionVaBelumBayar(&transaction, trID)
					if err != nil {
						log.Error("Failed get transaction: " + err.Error())
					}
					sentEmailCustomerTransaksiVaMandiriSebelumExpired(transaction)
				}

			}
		} else {
			fmt.Println("transaksi hampir expired tidak ada")
		}
	} else {
		fmt.Println("err get transaction : " + err.Error())
	}

	fmt.Println("======END CRON NOTIF TRANSAKSI BEFORE EXPIRED============")
}

func sentEmailCustomerTransaksiVaMandiriSebelumExpired(transaction models.DetailTransactionVaBelumBayar) {
	var mailTemp, subject string
	mailParam := make(map[string]string)

	mailParam["FileUrl"] = config.FileUrl + "/images/mail"
	mailParam["Name"] = transaction.FullName
	mailParam["Cif"] = *transaction.Cif
	mailParam["Date"] = transaction.TransDate
	mailParam["Time"] = transaction.TransTime
	mailParam["PaymentMethod"] = "Mandiri Virtual Account"
	mailParam["Day"] = transaction.LastDayPay
	mailParam["LastDate"] = transaction.LastDatePay
	mailParam["LastTime"] = transaction.LastTimePay
	mailParam["NoVA"] = transaction.NoVa
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}

	subject = "[MotionFunds] Kami Masih Menunggu Pembayaran Subscription melalui " + mailParam["PaymentMethod"] + " (" + transaction.NoVa + ")"
	mailTemp = "index-subscription-va-uncomplete-expired.html"

	mailParam["ProductName"] = transaction.ProductName
	mailParam["Symbol"] = transaction.CurrencySymbol + " "
	mailParam["TotalPembayaran"] = ac0.FormatMoneyDecimal(transaction.TotalPembayaran.Truncate(0))
	mailParam["Amount"] = ac0.FormatMoneyDecimal(transaction.TransAmount.Truncate(0))
	mailParam["Fee"] = ac0.FormatMoneyDecimal(transaction.Fee.Truncate(0))

	t := template.New(mailTemp)

	t, _ = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mailParam); err != nil {
		log.Error("Failed send mail: " + err.Error())
	} else {
		result := tpl.String()

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", transaction.UloginEmail)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", result)

		err = SendEmail(mailer)
		if err != nil {
			log.Error("Failed send mail transaksi sebelum expired to: " + transaction.UloginEmail)
			log.Error("Failed send mail: " + err.Error())
		} else {
			log.Println("Sukses email transaksi sebelum expired : " + transaction.UloginEmail)
		}

		// dialer := gomail.NewDialer(
		// 	config.EmailSMTPHost,
		// 	int(config.EmailSMTPPort),
		// 	config.EmailFrom,
		// 	config.EmailFromPassword,
		// )
		// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// err = dialer.DialAndSend(mailer)
		// if err != nil {
		// 	log.Error("Failed send mail transaksi sebelum expired to: " + transaction.UloginEmail)
		// 	log.Error("Failed send mail: " + err.Error())
		// } else {
		// 	log.Println("Sukses email transaksi sebelum expired : " + transaction.UloginEmail)
		// }
	}
}
