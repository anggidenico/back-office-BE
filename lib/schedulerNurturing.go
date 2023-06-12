package lib

import (
	"fmt"
	models "mf-bo-api/models/dashboard"
	"strconv"
	"time"
)

func GetDataUserBelumOA() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON NURTURING 1 : " + time.Now().Format(dateLayout))
	params := make(map[string]string)
	params["id_category"] = "1"
	params["rec_status"] = "1"
	var nurturing []models.NurturingCategoryTime
	_, err := models.GetAllNurturingCategoryTime(&nurturing, params)
	if err == nil {
		if len(nurturing) > 0 {
			for _, nur := range nurturing { //nurturing
				idCat := strconv.FormatUint(nur.IdCategory, 10)
				idTime := strconv.FormatUint(nur.Time, 10)
				var dataNur []models.UserRegistrasiBaruBelumOA
				_, err := models.GetUserRegistrasiBaruBelumOA(&dataNur, idCat, idTime)
				if err == nil {
					if len(dataNur) > 0 {
						for _, data := range dataNur {
							paramsSumm := make(map[string]string)
							paramsSumm["user_login"] = strconv.FormatUint(data.UserLoginKey, 10)
							if data.TokenNotif != nil {
								paramsSumm["token_notif"] = *data.TokenNotif
							}
							paramsSumm["full_name"] = data.FullName
							paramsSumm["email"] = data.UloginEmail
							paramsSumm["id_category"] = strconv.FormatUint(data.IdCategory, 10)
							paramsSumm["id_category_time"] = strconv.FormatUint(data.IdCategoryTime, 10)
							paramsSumm["id_force_type"] = strconv.FormatUint(data.IdForceType, 10)
							paramsSumm["id_message"] = strconv.FormatUint(data.IdMessage, 10)
							paramsSumm["message_title"] = data.MessageTitle
							paramsSumm["message_body"] = data.MessageBody
							paramsSumm["push_notif_status"] = "0"
							paramsSumm["rec_status"] = "1"
							paramsSumm["rec_created_date"] = time.Now().Format(dateLayout)
							paramsSumm["rec_created_by"] = "CRON"
							_, err = models.CreateNurturingNotifSummary(paramsSumm)
							if err != nil {
								fmt.Println("err create notif summary : " + err.Error())
							} else {
								fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
							}
						}
					} else {
						fmt.Println("DATA USER BELUM OA KOSONG - " + idCat + " - " + idTime)
					}
				} else {
					// log.Error(err.Error())
				}
			}
		} else {
			fmt.Println("DATA NURTURING KOSONG")
		}
	} else {
		// log.Error(err.Error())
	}
	fmt.Println("END CRON NURTURING 1 : " + time.Now().Format(dateLayout))
}

func GetDataUserSudahCustomerBelumTransaksi() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON NURTURING 2 : " + time.Now().Format(dateLayout))
	params := make(map[string]string)
	params["id_category"] = "2"
	params["rec_status"] = "1"
	var nurturing []models.NurturingCategoryTime
	_, err := models.GetAllNurturingCategoryTime(&nurturing, params)
	if err == nil {
		if len(nurturing) > 0 {
			for _, nur := range nurturing { //nurturing
				idCat := strconv.FormatUint(nur.IdCategory, 10)
				idTime := strconv.FormatUint(nur.Time, 10)
				var dataNur []models.UserRegistrasiBaruBelumOA
				_, err := models.GetUserSudahCustomerBelumTransaksi(&dataNur, idCat, idTime)
				if err == nil {
					if len(dataNur) > 0 {
						for _, data := range dataNur {
							paramsSumm := make(map[string]string)
							paramsSumm["user_login"] = strconv.FormatUint(data.UserLoginKey, 10)
							if data.TokenNotif != nil {
								paramsSumm["token_notif"] = *data.TokenNotif
							}
							paramsSumm["full_name"] = data.FullName
							paramsSumm["email"] = data.UloginEmail
							paramsSumm["id_category"] = strconv.FormatUint(data.IdCategory, 10)
							paramsSumm["id_category_time"] = strconv.FormatUint(data.IdCategoryTime, 10)
							paramsSumm["id_force_type"] = strconv.FormatUint(data.IdForceType, 10)
							paramsSumm["id_message"] = strconv.FormatUint(data.IdMessage, 10)
							paramsSumm["message_title"] = data.MessageTitle
							paramsSumm["message_body"] = data.MessageBody
							paramsSumm["push_notif_status"] = "0"
							paramsSumm["rec_status"] = "1"
							paramsSumm["rec_created_date"] = time.Now().Format(dateLayout)
							paramsSumm["rec_created_by"] = "CRON"
							_, err = models.CreateNurturingNotifSummary(paramsSumm)
							if err != nil {
								fmt.Println("err create notif summary : " + err.Error())
							} else {
								fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
							}
						}
					} else {
						fmt.Println("DATA USER BELUM OA KOSONG - " + idCat + " - " + idTime)
					}
				} else {
					// log.Error(err.Error())
				}
			}
		} else {
			fmt.Println("DATA NURTURING KOSONG")
		}
	} else {
		// log.Error(err.Error())
	}
	fmt.Println("END CRON NURTURING 2 : " + time.Now().Format(dateLayout))
}

func GetDataUserHanyaSubs1Kali() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON NURTURING 3 : " + time.Now().Format(dateLayout))
	params := make(map[string]string)
	params["id_category"] = "3"
	params["rec_status"] = "1"
	var nurturing []models.NurturingCategoryTime
	_, err := models.GetAllNurturingCategoryTime(&nurturing, params)
	if err == nil {
		if len(nurturing) > 0 {
			for _, nur := range nurturing { //nurturing
				forLoop := 1
				if nur.IsLooping != nil && *nur.IsLooping == 1 {
					forLoop = 3
				}

				firstCatTime := strconv.FormatUint(nur.Time, 10)

				for i := 1; i <= forLoop; i++ {
					idCat := strconv.FormatUint(nur.IdCategory, 10)
					idTime := strconv.FormatUint(nur.Time*uint64(i), 10)
					forceType := strconv.FormatUint(uint64(i), 10)
					var dataNur []models.UserRegistrasiBaruBelumOA
					_, err := models.GetUserHanyaSubs1kali(&dataNur, idCat, idTime, firstCatTime, forceType)
					if err == nil {
						if len(dataNur) > 0 {
							for _, data := range dataNur {
								paramsSumm := make(map[string]string)
								paramsSumm["user_login"] = strconv.FormatUint(data.UserLoginKey, 10)
								if data.TokenNotif != nil {
									paramsSumm["token_notif"] = *data.TokenNotif
								}
								paramsSumm["full_name"] = data.FullName
								paramsSumm["email"] = data.UloginEmail
								paramsSumm["id_category"] = strconv.FormatUint(data.IdCategory, 10)
								paramsSumm["id_category_time"] = strconv.FormatUint(data.IdCategoryTime, 10)
								paramsSumm["id_force_type"] = strconv.FormatUint(data.IdForceType, 10)
								paramsSumm["id_message"] = strconv.FormatUint(data.IdMessage, 10)
								paramsSumm["message_title"] = data.MessageTitle
								paramsSumm["message_body"] = data.MessageBody
								paramsSumm["push_notif_status"] = "0"
								paramsSumm["rec_status"] = "1"
								paramsSumm["rec_created_date"] = time.Now().Format(dateLayout)
								paramsSumm["rec_created_by"] = "CRON"
								_, err = models.CreateNurturingNotifSummary(paramsSumm)
								if err != nil {
									fmt.Println("err create notif summary : " + err.Error())
								} else {
									fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
								}
							}
						} else {
							fmt.Println("DATA USER BELUM OA KOSONG - " + idCat + " - " + idTime)
						}
					} else {
						// log.Error(err.Error())
					}
				}
			}
		} else {
			fmt.Println("DATA NURTURING KOSONG")
		}
	} else {
		// log.Error(err.Error())
	}
	fmt.Println("END CRON NURTURING 3 : " + time.Now().Format(dateLayout))
}
