package lib

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"strconv"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func GenerateUserInstitution() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON GENERATE USER INSTITUTION")
	var err error

	//cek user
	var userIns []models.OaInstitutionUserGenerateLogin
	_, err = models.GetOaInstitutionUserGenerateLogin(&userIns)
	if err == nil {
		if len(userIns) > 0 {
			for _, us := range userIns { //user notif
				var countData models.OaRequestCountData
				var pagination int
				_, err = models.GetCountUserActive(&countData, strconv.FormatUint(us.Oarequestkey, 10))
				if err != nil {
					pagination = 1
				} else {
					pagination = int(countData.CountData) + 1
				}

				username := *us.ShortName + strconv.FormatUint(uint64(pagination), 10)

				var count models.CountData
				_, err := models.ValidateUniqueData(&count, "ulogin_name", username, nil)
				if err != nil {
					fmt.Println("error get validate unique ulogin_name")
					continue
				}
				if int(count.CountData) > int(0) {
					username = *us.ShortName + "." + RandStringBytesMaskImprSrc(2) + strconv.FormatUint(uint64(pagination), 10)
				}

				date := time.Now().AddDate(0, 0, 10)
				expired := date.Format(dateLayout)

				//generate user login
				verifyKeyByte := sha256.Sum256([]byte(*us.EmailAddress + "_" + expired))
				verifyKey := hex.EncodeToString(verifyKeyByte[:])
				params := make(map[string]string)
				params["user_category_key"] = "1"
				params["user_dept_key"] = "1"
				params["ulogin_name"] = username
				params["ulogin_full_name"] = *us.FullName
				params["ulogin_password"] = ""
				params["ulogin_must_changepwd"] = "1"
				params["ulogin_email"] = *us.EmailAddress
				params["verified_email"] = "0"
				params["string_token"] = verifyKey
				params["ulogin_mobileno"] = *us.PhoneNumber
				params["verified_mobileno"] = "0"
				params["ulogin_enabled"] = "0"
				params["ulogin_locked"] = "0"
				params["ulogin_failed_count"] = "0"
				params["accept_login_tnc"] = "1"
				params["must_change_pin"] = "1"
				params["allowed_sharing_login"] = "1"
				params["customer_key"] = strconv.FormatUint(*us.CustomerKey, 10)
				params["role_key"] = strconv.FormatUint(*us.RoleKey, 10)
				params["rec_order"] = "0"
				params["rec_status"] = "1"
				params["rec_created_date"] = time.Now().Format(dateLayout)
				params["rec_created_by"] = "CRON"
				params["rec_modified_date"] = time.Now().Format(dateLayout)
				params["rec_modified_by"] = "CRON"

				//update oa_institution_user
				_, err, userID := models.CreateScUserLoginWithReturnPK(params)
				if err == nil {
					paramsUser := make(map[string]string)
					paramsUser["insti_user_key"] = strconv.FormatUint(us.InstiUserKey, 10)
					paramsUser["ulogin_created_date"] = time.Now().Format(dateLayout)
					paramsUser["user_login_key"] = userID
					paramsUser["rec_modified_date"] = time.Now().Format(dateLayout)
					paramsUser["rec_modified_by"] = "CRON"

					err = sendEmail(*us.EmailAddress, verifyKey)
					if err == nil {
						paramsUser["rec_attribute_id1"] = "1"
					} else {
						paramsUser["rec_attribute_id1"] = "0"
						paramsUser["rec_attribute_id2"] = err.Error()
					}

					_, err := models.UpdateOaInstitutionUser(paramsUser)
					if err != nil {
						fmt.Println("sukses update institution user")
					} else {
						fmt.Println("error update institution user")
					}
				} else {
					log.Println("fffffffff")
					fmt.Println("error save user_login")
				}
			}
		} else {
			fmt.Println("user institution tidak ada")
		}
	} else {
		fmt.Println("err get user institution : " + err.Error())
	}

	fmt.Println("======END CRON GENERATE USER INSTITUTION============")
}

func ResendEmailFailedGenerateUserInstitution() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON GENERATE USER INSTITUTION RESEND EMAIL")
	var err error

	//cek user
	var userIns []models.OaInstitutionUserGenerateLogin
	_, err = models.GetOaInstitutionUserGenerateLoginFailedSendEmail(&userIns)
	if err == nil {
		if len(userIns) > 0 {
			for _, us := range userIns { //user notif
				//kirim email
				paramsUser := make(map[string]string)
				paramsUser["insti_user_key"] = strconv.FormatUint(us.InstiUserKey, 10)
				paramsUser["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsUser["rec_modified_by"] = "CRON"

				err = sendEmail(*us.EmailAddress, *us.StringToken)
				if err == nil {
					paramsUser["rec_attribute_id1"] = "1"
				} else {
					paramsUser["rec_attribute_id1"] = "0"
					paramsUser["rec_attribute_id2"] = err.Error()
				}

				_, err := models.UpdateOaInstitutionUser(paramsUser)
				if err == nil {
					fmt.Println("sukses update institution user")
				} else {
					fmt.Println("error update institution user")
				}
			}
		} else {
			fmt.Println("user institution tidak ada")
		}
	} else {
		fmt.Println("err get user institution : " + err.Error())
	}

	fmt.Println("======END CRON GENERATE USER INSTITUTION RESEND EMAIL============")
}

func sendEmail(email string, token string) error {
	fmt.Println("======SEND EMAIL TO " + email + "============")
	var err error
	t := template.New("index-email-activation-web.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-email-activation-web.html")
	if err != nil {
		log.Println(err)
		fmt.Println("======FAILED SEND EMAIL TO " + email + "============")
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		Url     string
		FileUrl string
	}{Url: config.BaseUrl + "/verify-email-user?token=" + token, FileUrl: config.FileUrl + "/images/mail"}); err != nil {
		log.Println(err)
		fmt.Println("======FAILED SEND EMAIL TO " + email + "============")
		return err
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MotionFunds] Verifikasi Email Kamu")
	mailer.SetBody("text/html", result)

	err = SendEmail(mailer)
	if err != nil {
		log.Error(err)
		fmt.Println("======FAILED SEND EMAIL TO " + email + "============")
		return err
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
	// 	log.Error(err)
	// 	fmt.Println("======FAILED SEND EMAIL TO " + email + "============")
	// 	return err
	// }
	fmt.Println("======SUKSES SEND EMAIL TO " + email + "============")
	return nil
}
