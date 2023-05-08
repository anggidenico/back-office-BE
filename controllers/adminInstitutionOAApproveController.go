package controllers

import (
	"bytes"
	"html/template"
	"mf-bo-api/config"
	"mf-bo-api/db"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func CSApproveOAInstitution(c echo.Context) error {
	errorAuth := initAuthCs()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error

	params := make(map[string]string)

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey == "" {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(oaRequestKey, 10, 64)
	if err != nil || n == 0 {
		log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	var oareq models.OaRequest
	_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
	if err != nil {
		log.Error("OA Request not found.")
		return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
	}

	if strconv.FormatUint(lib.Profile.UserID, 10) == *oareq.RecCreatedBy {
		log.Error("User not autorized.")
		return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
	}

	if *oareq.Oastatus != uint64(lib.OA_ENTRIED) {
		log.Error("User not autorized.")
		return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
	}

	params["oa_request_key"] = oaRequestKey

	oastatus := c.FormValue("oa_status") //259 = approve --------- 444 = reject
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != strconv.FormatUint(uint64(lib.CS_APPROVED), 10)) && (oastatus != strconv.FormatUint(uint64(lib.DRAFT), 10)) {
			log.Error("Wrong input for parameter: oa_status must 444/258")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		params["oa_status"] = oastatus
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	check1notes := c.FormValue("notes")
	params["check1_notes"] = check1notes

	if oastatus != strconv.FormatUint(uint64(lib.CS_APPROVED), 10) { //jika reject
		if check1notes == "" {
			log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["check1_flag"] = "0"
	} else {
		params["check1_flag"] = "1"
	}

	params["check1_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check1_references"] = strKey
	params["rec_modified_by"] = strKey

	_, err = models.UpdateOaRequest(params)
	if err != nil {
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	if oastatus == strconv.FormatUint(uint64(lib.CS_APPROVED), 10) { //jika approve
		//sent email to all KYC
		var institutionData models.OaInstitutionData
		_, err = models.GetOaInstitutionData(&institutionData, oaRequestKey, "oa_request_key")
		if err == nil {
			SentEmailInstitusiOaPengkinianToBackOfficeSales(oareq, institutionData, "12", false)
		}
	} else {
		//Gak ngirim Email apa2 (confirmed)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func KYCAPproveOAInstitution(c echo.Context) error {
	errorAuth := initAuthKyc()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int

	params := make(map[string]string)

	oaRequestKey := c.FormValue("oa_request_key")
	if oaRequestKey == "" {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(oaRequestKey, 10, 64)
	if err != nil || n == 0 {
		log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	var oareq models.OaRequest
	_, err = models.GetOaRequestInstitution(&oareq, oaRequestKey, "")
	if err != nil {
		log.Error("OA Request not found.")
		return lib.CustomError(http.StatusBadRequest, "OA Request not found.", "OA Request not found.")
	}

	var institutionData models.OaInstitutionData
	_, err = models.GetOaInstitutionData(&institutionData, oaRequestKey, "oa_request_key")
	if err != nil {
		log.Error("OA Institution data not found.")
		return lib.CustomError(http.StatusBadRequest, "OA Institution data not found.", "OA Institution data not found.")
	}

	if *oareq.Oastatus != uint64(lib.CS_APPROVED) {
		log.Error("User not autorized.")
		return lib.CustomError(http.StatusBadRequest, "User not autorized.", "User not autorized.")
	}

	params["oa_request_key"] = oaRequestKey

	oastatus := c.FormValue("oa_status") //260 = approve --------- 444 = reject
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != strconv.FormatUint(uint64(lib.KYC_APPROVED), 10)) && (oastatus != strconv.FormatUint(uint64(lib.DRAFT), 10)) {
			log.Error("Wrong input for parameter: oa_status must 444/259")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		params["oa_status"] = oastatus
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	check2notes := c.FormValue("notes")
	params["check2_notes"] = check2notes

	if oastatus != strconv.FormatUint(uint64(lib.KYC_APPROVED), 10) { //jika reject
		if check2notes == "" {
			log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["check2_flag"] = "0"
	} else {
		params["check2_flag"] = "1"
	}

	oarisklevel := c.FormValue("oa_risk_level")
	if oarisklevel == "" {
		log.Error("Missing required parameter: oa_risk_level")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: oa_risk_level", "Missing required parameter: oa_risk_level")
	}
	n, err = strconv.ParseUint(oarisklevel, 10, 64)
	if err == nil && n > 0 {
		params["oa_risk_level"] = oarisklevel
	} else {
		log.Error("Wrong input for parameter: oa_risk_level")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_risk_level", "Wrong input for parameter: oa_risk_level")
	}

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	params["check2_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strKey
	params["check2_references"] = strKey

	tx, err := db.Db.Begin()

	//update oa request
	_, err = models.UpdateOaRequest(params)
	if err != nil {
		tx.Rollback()
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}
	log.Info("Success update approved KYC Transaction")

	var oadata models.OaInstitutionData
	_, err = models.GetOaInstitutionData(&oadata, oaRequestKey, "oa_request_key")
	if err != nil {
		tx.Rollback()
		log.Error("Error Institution Data not Found")
		return lib.CustomError(status, err.Error(), "Institution data not found")
	}

	if oastatus == strconv.FormatUint(uint64(lib.KYC_APPROVED), 10) { //APPROVE
		if *oareq.OaRequestType == uint64(127) { // new
			//generate customer
			paramsCustomer := make(map[string]string)
			paramsCustomer["id_customer"] = "0"

			year, month, _ := time.Now().Date()

			var customer models.MsCustomer
			tahun := strconv.FormatUint(uint64(year), 10)
			bulan := strconv.FormatUint(uint64(month), 10)
			if len(bulan) == 1 {
				bulan = "0" + bulan
			}
			unitHolderLike := tahun + bulan
			status, err = models.GetLastUnitHolder(&customer, unitHolderLike)
			if err != nil {
				paramsCustomer["unit_holder_idno"] = unitHolderLike + "000001"
			} else {
				dgt := customer.UnitHolderIDno[len(customer.UnitHolderIDno)-6:]
				productKeyCek, _ := strconv.ParseUint(dgt, 10, 64)
				hasil := strconv.FormatUint(productKeyCek+1, 10)
				if len(hasil) == 1 {
					hasil = "00000" + hasil
				} else if len(hasil) == 2 {
					hasil = "0000" + hasil
				} else if len(hasil) == 3 {
					hasil = "000" + hasil
				} else if len(hasil) == 4 {
					hasil = "00" + hasil
				} else if len(hasil) == 5 {
					hasil = "0" + hasil
				}

				resultData := unitHolderLike + hasil
				paramsCustomer["unit_holder_idno"] = resultData
			}

			paramsCustomer["full_name"] = *oadata.FullName
			paramsCustomer["investor_type"] = "264"
			paramsCustomer["customer_category"] = "265"
			paramsCustomer["cif_suspend_flag"] = "0"

			if oareq.BranchKey == nil {
				paramsCustomer["openacc_branch_key"] = "1"
			} else {
				paramsCustomer["openacc_branch_key"] = strconv.FormatUint(*oareq.BranchKey, 10)
			}

			if oareq.AgentKey == nil {
				paramsCustomer["openacc_agent_key"] = "1"
			} else {
				paramsCustomer["openacc_agent_key"] = strconv.FormatUint(*oareq.AgentKey, 10)
			}

			paramsCustomer["openacc_date"] = time.Now().Format(dateLayout)
			paramsCustomer["flag_employee"] = "0"
			paramsCustomer["flag_group"] = "0"
			paramsCustomer["merging_flag"] = "0"
			paramsCustomer["rec_status"] = "1"
			paramsCustomer["rec_created_date"] = time.Now().Format(dateLayout)
			paramsCustomer["rec_created_by"] = strKey

			sliceName := strings.Fields(*oadata.FullName)

			if len(sliceName) > 0 {
				if len(sliceName) == 1 {
					paramsCustomer["first_name"] = sliceName[0]
					paramsCustomer["last_name"] = sliceName[0]
				}
				if len(sliceName) == 2 {
					paramsCustomer["first_name"] = sliceName[0]
					paramsCustomer["last_name"] = sliceName[1]
				}
				if len(sliceName) > 2 {
					ln := len(sliceName)
					paramsCustomer["first_name"] = sliceName[0]
					paramsCustomer["middle_name"] = sliceName[1]
					lastName := strings.Join(sliceName[2:ln], " ")
					paramsCustomer["last_name"] = lastName
				}
			}

			strNationality := strconv.FormatUint(*oadata.Nationality, 10)
			if strNationality == "97" {
				paramsCustomer["fatca_status"] = "278"
			} else if strNationality == "225" {
				paramsCustomer["fatca_status"] = "279"
			} else {
				paramsCustomer["fatca_status"] = "280"
			}
			if oadata.TinNumber != nil {
				paramsCustomer["tin_number"] = *oadata.TinNumber
			}

			status, err, cusID := models.CreateMsCustomer(paramsCustomer)
			if err != nil {
				tx.Rollback()
				log.Error("Error create customer")
				return lib.CustomError(status, err.Error(), "failed input data")
			}
			request, err := strconv.ParseUint(cusID, 10, 64)
			if request == 0 {
				tx.Rollback()
				log.Error("Failed create customer")
				return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
			}

			paramOaUpdate := make(map[string]string)
			paramOaUpdate["customer_key"] = cusID
			paramOaUpdate["oa_request_key"] = oaRequestKey

			_, err = models.UpdateOaRequest(paramOaUpdate)
			if err != nil {
				tx.Rollback()
				log.Error("Error update oa request")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}
			//create all ms_customer_bank_account by oa_req_key
			var accBank []models.OaRequestByField
			status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
			if err != nil {
				log.Error(err.Error())
			}
			if len(accBank) > 0 {
				var bindVarMsBank []interface{}
				for _, value := range accBank {
					var row []string
					row = append(row, cusID)                                        //customer_key
					row = append(row, strconv.FormatUint(value.BankAccountKey, 10)) //bank_account_key
					row = append(row, strconv.FormatUint(value.FlagPriority, 10))   //flag_priority
					row = append(row, *value.AccountHolderName)                     //bank_account_name
					row = append(row, "1")                                          //rec_status
					row = append(row, time.Now().Format(dateLayout))                //rec_created_date
					row = append(row, strconv.FormatUint(lib.Profile.UserID, 10))   //rec_created_by
					bindVarMsBank = append(bindVarMsBank, row)
				}
				_, err = models.CreateMultipleMsCustomerBankkAccount(bindVarMsBank)
				if err != nil {
					tx.Rollback()
					log.Error("Failed create promo product: " + err.Error())
					return lib.CustomError(status, err.Error(), "failed input data")
				}
			}
		} else { // pengkinian
			//delete all ms_customer_bank_account by customer
			deleteParam := make(map[string]string)
			deleteParam["rec_status"] = "0"
			deleteParam["rec_modified_date"] = time.Now().Format(dateLayout)
			deleteParam["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			deleteParam["rec_deleted_date"] = time.Now().Format(dateLayout)
			deleteParam["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			_, err = models.UpdateDataByField(deleteParam, "customer_key", strconv.FormatUint(*oareq.CustomerKey, 10))
			if err != nil {
				log.Error("Error delete all ms_customer_bank_account")
			}
			//create all ms_customer_bank_account by oa_req_key
			var accBank []models.OaRequestByField
			status, err = models.GetOaRequestBankByField(&accBank, "oa_request_key", strconv.FormatUint(oareq.OaRequestKey, 10))
			if err != nil {
				log.Error(err.Error())
			}
			if len(accBank) > 0 {
				var bindVarMsBank []interface{}
				for _, value := range accBank {
					var row []string
					row = append(row, strconv.FormatUint(*oareq.CustomerKey, 10))   //customer_key
					row = append(row, strconv.FormatUint(value.BankAccountKey, 10)) //bank_account_key
					row = append(row, strconv.FormatUint(value.FlagPriority, 10))   //flag_priority
					row = append(row, *value.AccountHolderName)                     //bank_account_name
					row = append(row, "1")                                          //rec_status
					row = append(row, time.Now().Format(dateLayout))                //rec_created_date
					row = append(row, strconv.FormatUint(lib.Profile.UserID, 10))   //rec_created_by
					bindVarMsBank = append(bindVarMsBank, row)
				}
				_, err = models.CreateMultipleMsCustomerBankkAccount(bindVarMsBank)
				if err != nil {
					tx.Rollback()
					log.Error("Failed create promo product: " + err.Error())
					return lib.CustomError(status, err.Error(), "failed input data")
				}
			}
		}
		//kirim email ke FundAdmin
		SentEmailInstitusiOaPengkinianToBackOfficeSales(oareq, institutionData, "13", false)

		if oareq.PipelineKey != nil {
			//kirim email Sukses PIC cc Sales --> jika dari web (pipeline)
			SentEmailOaApprovePicInstitutionCcSales(oareq, institutionData)
		}
	} else { //reject
		if oareq.PipelineKey != nil {
			//kirim email Reject PIC cc Sales --> jika dari web (pipeline)
			SentEmailOaRejectPicInstitutionCcSales(oareq, institutionData, check2notes)
		} else {
			//kirim email ke sales
			SentEmailOaRejectToSales(oareq, institutionData, check2notes)
		}
	}
	tx.Commit()

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func SentEmailInstitusiOaPengkinianToBackOfficeSales(
	oaRequest models.OaRequest,
	institutionData models.OaInstitutionData,
	roleKey string,
	isSentSales bool) {

	var err error
	var mailTemp, subject string
	mailParam := make(map[string]string)
	if roleKey == "11" {
		mailParam["BackOfficeGroup"] = "Customer Service"
	} else if roleKey == "12" {
		mailParam["BackOfficeGroup"] = "Compliance"
	} else if roleKey == "13" {
		mailParam["BackOfficeGroup"] = "FundAdmin"
	}

	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	if institutionData.DeedDate != nil {
		date, _ := time.Parse(layout, *institutionData.DeedDate)
		mailParam["TanggalPendirian"] = date.Format(dateLayout)
	} else {
		mailParam["TanggalPendirian"] = "-"
	}
	if institutionData.FullName != nil {
		mailParam["NamaInstitusi"] = *institutionData.FullName
	} else {
		mailParam["NamaInstitusi"] = "-"
	}
	if institutionData.TinNumber != nil {
		mailParam["Npwp"] = *institutionData.TinNumber
	} else {
		mailParam["Npwp"] = "-"
	}
	if institutionData.EmailAddress != nil {
		mailParam["Email"] = *institutionData.EmailAddress
	} else {
		mailParam["Email"] = "-"
	}
	if institutionData.PhoneNo != nil {
		mailParam["NoTelp"] = *institutionData.PhoneNo
	} else {
		mailParam["NoTelp"] = "-"
	}
	if institutionData.PhoneNo != nil {
		mailParam["NoHp"] = *institutionData.MobileNo
	} else {
		mailParam["NoHp"] = "-"
	}
	mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

	if *oaRequest.OaRequestType == uint64(127) { // oa new
		mailParam["JenisPengajuanData"] = "melengkapi Pembukaan Rekening Reksa Dana melalui MotionFunds"

		subject = "[MotionFunds] Mohon Verifikasi Pembukaan Rekening Reksa Dana"
	} else { // pengkinian
		if *oaRequest.OaRequestType == uint64(128) {
			mailParam["JenisPengajuanData"] = "melakukan Pengkinian Profile Risiko"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Profile Risiko"
		} else {
			mailParam["JenisPengajuanData"] = "melakukan Pengkinian Personal Data"
			subject = "[MotionFunds] Mohon Verifikasi Pengkinian Personal Data"
		}
	}
	mailTemp = "email-oa-institution-to-cs-kyc-fundadmin-sales.html"

	paramsScLogin := make(map[string]string)
	paramsScLogin["role_key"] = roleKey
	paramsScLogin["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
	if err != nil {
		log.Error("User BO tidak ditemukan")
		log.Error(err)
	} else {
		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
		} else {
			for _, scLogin := range userLogin {
				if scLogin.UserLoginKey != lib.Profile.UserID {
					strUserCat := strconv.FormatUint(scLogin.UserCategoryKey, 10)
					if (strUserCat == "2") || (strUserCat == "3") {
						var tpl bytes.Buffer
						if err := t.Execute(&tpl, mailParam); err != nil {
							log.Error("Failed send mail: " + err.Error())
						} else {
							result := tpl.String()

							mailer := gomail.NewMessage()
							mailer.SetHeader("From", config.EmailFrom)
							mailer.SetHeader("To", scLogin.UloginEmail)
							mailer.SetHeader("Subject", subject)
							mailer.SetBody("text/html", result)

							err = lib.SendEmail(mailer)
							if err != nil {
								log.Error("Failed send mail to: " + scLogin.UloginEmail)
								log.Error("Failed send mail: " + err.Error())
							} else {
								log.Println("Sukses kirim email")
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
							// 	log.Error("Failed send mail to: " + scLogin.UloginEmail)
							// 	log.Error("Failed send mail: " + err.Error())
							// }
						}
					}
				}
			}
		}
	}

	if isSentSales { //kirim email ke sales
		mailParam["BackOfficeGroup"] = "Sales"
		var agentKey string
		if oaRequest.AgentKey == nil {
			agentKey = "1"
		} else {
			agentKey = strconv.FormatUint(*oaRequest.AgentKey, 10)
		}

		var agent models.MsAgent
		_, err = models.GetMsAgent(&agent, agentKey)
		if err != nil {
			log.Error("Agent not found")
		} else {
			if agent.AgentEmail != nil {
				t := template.New(mailTemp)

				t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
				if err != nil {
					log.Error("Failed send mail to sales: " + err.Error())
				} else {
					var tpl bytes.Buffer
					if err := t.Execute(&tpl, mailParam); err != nil {
						log.Error("Failed send mail to sales: " + err.Error())
					} else {
						result := tpl.String()

						mailer := gomail.NewMessage()
						mailer.SetHeader("From", config.EmailFrom)
						mailer.SetHeader("To", *agent.AgentEmail)
						mailer.SetHeader("Subject", subject)
						mailer.SetBody("text/html", result)

						err = lib.SendEmail(mailer)
						if err != nil {
							log.Error("Failed send mail to sales : " + *agent.AgentEmail)
							log.Error("Failed send mail: " + err.Error())
						} else {
							log.Println("Sukses kirim email")
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
						// 	log.Error("Failed send mail to sales : " + *agent.AgentEmail)
						// 	log.Error("Failed send mail: " + err.Error())
						// }
					}
				}
			} else {
				log.Error("Sales tidak punya email")
			}
		}
	}
}

func SentEmailOaApprovePicInstitutionCcSales(
	oaRequest models.OaRequest,
	institutionData models.OaInstitutionData) {

	var err error
	if oaRequest.PipelineKey != nil {
		var mailTemp, subject string
		mailParam := make(map[string]string)
		mailParam["FileUrl"] = config.ImageUrl + "/images/mail"

		if *oaRequest.OaRequestType == uint64(127) { // oa new
			mailTemp = "email-sukses-verifikasi-kyc-oa-institusi.html"
			subject = "[MotionFunds] Pembukaan Rekening Anda telah Disetujui"
		} else { // pengkinian
			mailTemp = "email-sukses-verifikasi-kyc-oa-institusi-pengkinian.html"
			if *oaRequest.OaRequestType == uint64(128) {
				mailParam["JenisPengkinian"] = "Profile Risiko"
				subject = "[MotionFunds] Pengkinian Data Institusi Profil Risiko Anda telah Berhasil"
			} else {
				mailParam["JenisPengajuanData"] = "Institution"
				subject = "[MotionFunds] Pengkinian Data Institusi Anda telah Berhasil"
			}
		}

		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
		} else {
			var pipeline models.OaCustomerPipeline
			_, err = models.GetOaCustomerPipeline(&pipeline, strconv.FormatUint(*oaRequest.PipelineKey, 10), "pipeline_key")
			if err == nil {
				var tpl bytes.Buffer
				if err := t.Execute(&tpl, mailParam); err != nil {
					log.Error("Failed send mail: " + err.Error())
				} else {
					if pipeline.PicEmailAddress != nil {
						result := tpl.String()

						mailer := gomail.NewMessage()
						mailer.SetHeader("From", config.EmailFrom)
						mailer.SetHeader("To", *pipeline.PicEmailAddress)

						var agentKey string
						if oaRequest.AgentKey == nil {
							agentKey = "1"
						} else {
							agentKey = strconv.FormatUint(*oaRequest.AgentKey, 10)
						}
						var agent models.MsAgent
						_, err = models.GetMsAgent(&agent, agentKey)
						if err != nil {
							log.Error("Agent not found")
						} else {
							if agent.AgentEmail != nil {
								mailer.SetAddressHeader("Cc", *agent.AgentEmail, agent.AgentName)
							}
						}

						mailer.SetHeader("Subject", subject)
						mailer.SetBody("text/html", result)

						err = lib.SendEmail(mailer)
						if err != nil {
							log.Error("Failed send mail pipeline to: " + *pipeline.PicEmailAddress)
							log.Error("Failed send mail pipeline : " + err.Error())
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
						// 	log.Error("Failed send mail pipeline to: " + *pipeline.PicEmailAddress)
						// 	log.Error("Failed send mail pipeline : " + err.Error())
						// }
					} else {
						log.Error("Pipeline tidak punya email")
					}
				}
			}
		}
	}
}

func SentEmailOaRejectPicInstitutionCcSales(
	oaRequest models.OaRequest,
	institutionData models.OaInstitutionData,
	notes string) {

	var err error
	if oaRequest.PipelineKey != nil {
		var pipeline models.OaCustomerPipeline
		_, err = models.GetOaCustomerPipeline(&pipeline, strconv.FormatUint(*oaRequest.PipelineKey, 10), "pipeline_key")
		if err == nil {
			layout := "2006-01-02 15:04:05"
			dateLayout := "02 Jan 2006"
			var mailTemp, subject string
			mailParam := make(map[string]string)
			mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
			mailParam["NamaPerusahaan"] = *institutionData.FullName
			mailParam["AlamatPerusahaan"] = "-"
			if institutionData.CorrespondenceKey != nil {
				var postalAddress models.OaPostalAddress
				_, err = models.GetOaPostalAddress(&postalAddress, strconv.FormatUint(*institutionData.CorrespondenceKey, 10))
				if err == nil && postalAddress.AddressLine1 != nil {
					mailParam["AlamatPerusahaan"] = *postalAddress.AddressLine1
				}
			}
			mailParam["NPWP"] = *institutionData.TinNumber
			mailParam["NamaPIC"] = *pipeline.FullName

			date, _ := time.Parse(layout, oaRequest.OaEntryEnd)
			mailParam["TanggalPendaftaran"] = date.Format(dateLayout)
			mailParam["Status"] = "Rejected"
			mailParam["Keterangan"] = notes

			if *oaRequest.OaRequestType == uint64(127) { // oa new
				mailParam["Jenis"] = "pendaftaran rekening"
				subject = "[MotionFunds] Pembukaan New Opening Account " + *institutionData.FullName + " Belum dapat Diproses"
			} else { // pengkinian
				if *oaRequest.OaRequestType == uint64(128) {
					mailParam["Jenis"] = "pengkinian Profil Risiko"
					subject = "[MotionFunds] Pengkinian Data Institusi Profil Risiko Anda Belum dapat Diproses"
				} else {
					mailParam["Jenis"] = "pengkinian Data Institusi"
					subject = "[MotionFunds] Pengkinian Data Institusi Anda Belum dapat Diproses"
				}
			}
			mailTemp = "email-oa-institution-rejected-kyc.html"

			t := template.New(mailTemp)

			t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
			if err != nil {
				log.Error("Failed send mail: " + err.Error())
			} else {
				var tpl bytes.Buffer
				if err := t.Execute(&tpl, mailParam); err != nil {
					log.Error("Failed send mail: " + err.Error())
				} else {
					if pipeline.PicEmailAddress != nil {
						result := tpl.String()

						mailer := gomail.NewMessage()
						mailer.SetHeader("From", config.EmailFrom)
						mailer.SetHeader("To", *pipeline.PicEmailAddress)

						var agentKey string
						if oaRequest.AgentKey == nil {
							agentKey = "1"
						} else {
							agentKey = strconv.FormatUint(*oaRequest.AgentKey, 10)
						}
						var agent models.MsAgent
						_, err = models.GetMsAgent(&agent, agentKey)
						if err != nil {
							log.Error("Agent not found")
						} else {
							if agent.AgentEmail != nil {
								mailer.SetAddressHeader("Cc", *agent.AgentEmail, agent.AgentName)
							}
						}

						mailer.SetHeader("Subject", subject)
						mailer.SetBody("text/html", result)

						err = lib.SendEmail(mailer)
						if err != nil {
							log.Error("Failed send mail pipeline to: " + *pipeline.PicEmailAddress)
							log.Error("Failed send mail pipeline : " + err.Error())
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
						// 	log.Error("Failed send mail pipeline to: " + *pipeline.PicEmailAddress)
						// 	log.Error("Failed send mail pipeline : " + err.Error())
						// }
					} else {
						log.Error("Pipeline tidak punya email")
					}
				}
			}
		}
	}
}

func SentEmailOaRejectToSales(
	oaRequest models.OaRequest,
	institutionData models.OaInstitutionData,
	notes string) {

	var err error
	var agentKey string
	if oaRequest.AgentKey == nil {
		agentKey = "1"
	} else {
		agentKey = strconv.FormatUint(*oaRequest.AgentKey, 10)
	}

	var agent models.MsAgent
	_, err = models.GetMsAgent(&agent, agentKey)
	if err == nil && agent.AgentEmail != nil {
		layout := "2006-01-02 15:04:05"
		dateLayout := "02 Jan 2006"
		var mailTemp, subject string
		mailParam := make(map[string]string)
		mailParam["FileUrl"] = config.ImageUrl + "/images/mail"
		mailParam["NamaPerusahaan"] = *institutionData.FullName
		mailParam["AlamatPerusahaan"] = "-"
		if institutionData.CorrespondenceKey != nil {
			var postalAddress models.OaPostalAddress
			_, err = models.GetOaPostalAddress(&postalAddress, strconv.FormatUint(*institutionData.CorrespondenceKey, 10))
			if err == nil && postalAddress.AddressLine1 != nil {
				mailParam["AlamatPerusahaan"] = *postalAddress.AddressLine1
			}
		}
		mailParam["NPWP"] = *institutionData.TinNumber
		mailParam["NamaPIC"] = "-"

		date, _ := time.Parse(layout, oaRequest.OaEntryEnd)
		mailParam["TanggalPendaftaran"] = date.Format(dateLayout)
		mailParam["Status"] = "Rejected"
		mailParam["Keterangan"] = notes

		if *oaRequest.OaRequestType == uint64(127) { // oa new
			mailParam["Jenis"] = "pendaftaran rekening"
			subject = "[MotionFunds] Pembukaan New Opening Account " + *institutionData.FullName + " Belum dapat Diproses"
		} else { // pengkinian
			if *oaRequest.OaRequestType == uint64(128) {
				mailParam["Jenis"] = "pengkinian Profil Risiko"
				subject = "[MotionFunds] Pengkinian Data Institusi Profil Risiko Belum dapat Diproses"
			} else {
				mailParam["Jenis"] = "pengkinian Data Institusi"
				subject = "[MotionFunds] Pengkinian Data Institusi Belum dapat Diproses"
			}
		}
		mailTemp = "email-oa-institution-rejected-to-sales.html"

		t := template.New(mailTemp)

		t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
		} else {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, mailParam); err != nil {
				log.Error("Failed send mail: " + err.Error())
			} else {
				result := tpl.String()

				mailer := gomail.NewMessage()
				mailer.SetHeader("From", config.EmailFrom)
				mailer.SetHeader("To", *agent.AgentEmail)
				mailer.SetHeader("Subject", subject)
				mailer.SetBody("text/html", result)

				err = lib.SendEmail(mailer)
				if err != nil {
					log.Error("Failed send mail sales to: " + *agent.AgentEmail)
					log.Error("Failed send mail sales : " + err.Error())
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
				// 	log.Error("Failed send mail sales to: " + *agent.AgentEmail)
				// 	log.Error("Failed send mail sales : " + err.Error())
				// }
			}
		}
	} else {
		log.Error("Agent not found")
	}
}
