package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func GenerateTokenMPLinking() {
	fmt.Println("start cron generate token linking motion pay")
	var scApp models.ScAppConfig
	_, err := models.GetScAppConfigTokenMotionPay(&scApp, "TOKEN_MOTION_PAY_LINKING")
	if err != nil {
		fmt.Println("error get data config")
		log.Error(err.Error())
	} else {
		merchantId := config.MERCHANT_ID
		partnerId := config.PARTNER_ID
		secretKey := config.SECRET_ID
		apiName := "TOKEN_MP_LINKING"
		urlPath := config.SandBox + PATH_GENERATE_TOKEN_LINKING
		reqMethod := "POST"
		status, res, err := requestTokenMotionPay(
			merchantId,
			partnerId,
			secretKey,
			apiName,
			urlPath,
			reqMethod,
		)
		if err == nil && status == http.StatusOK {
			var dataBody map[string]interface{}
			err := json.Unmarshal([]byte(res), &dataBody)
			if err == nil {
				messageData := dataBody["message_data"].(map[string]interface{})
				dateLayout := "2006-01-02 15:04:05"
				paramsConfig := make(map[string]string)
				paramsConfig["app_config_value"] = messageData["token"].(string)
				paramsConfig["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsConfig["rec_modified_by"] = "CRON"
				_, err = models.UpdateMsCustomerByConfigCode(paramsConfig, "TOKEN_MOTION_PAY_LINKING")
				if err != nil {
					log.Error("Error update App Config")
					log.Error(err.Error())
				}
			}
		} else {
			log.Error(status, " error get data")
		}
	}
	fmt.Println("end cron generate token linking motion pay")
}

func GenerateTokenMPPayment() {
	fmt.Println("start cron generate token payment motion pay")
	var scApp models.ScAppConfig
	_, err := models.GetScAppConfigTokenMotionPay(&scApp, "TOKEN_MOTION_PAY_PAYMENT")
	if err != nil {
		fmt.Println("error get data config")
		log.Error(err.Error())
	} else {
		merchantId := config.MERCHANT_ID_MP_PAYMENT
		partnerId := config.PARTNER_ID_MP_PAYMENT
		secretKey := config.SECRET_ID_MP_PAYMENT
		apiName := "TOKEN_MP_PAYMENT"
		urlPath := config.SANDBOX_MP_PAYMENT + PATH_GENERATE_TOKEN_PAYMENT
		reqMethod := "POST"
		status, res, err := requestTokenMotionPay(
			merchantId,
			partnerId,
			secretKey,
			apiName,
			urlPath,
			reqMethod,
		)
		if err == nil && status == http.StatusOK {
			var dataBody map[string]interface{}
			err := json.Unmarshal([]byte(res), &dataBody)
			if err == nil {
				messageData := dataBody["message_data"].(map[string]interface{})
				dateLayout := "2006-01-02 15:04:05"
				paramsConfig := make(map[string]string)
				paramsConfig["app_config_value"] = messageData["token"].(string)
				paramsConfig["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsConfig["rec_modified_by"] = "CRON"
				_, err = models.UpdateMsCustomerByConfigCode(paramsConfig, "TOKEN_MOTION_PAY_PAYMENT")
				if err != nil {
					log.Error("Error update App Config")
					log.Error(err.Error())
				}
			}
		} else {
			log.Error(status, " error get data")
		}
	}
	fmt.Println("end cron generate token payment motion pay")
}

func requestTokenMotionPay(
	merchantId string,
	partnerId string,
	secretKey string,
	apiName string,
	urlPath string,
	reqMethod string) (int, string, error) {

	params := make(map[string]interface{})
	params["merchant_id"] = merchantId
	params["partner_id"] = partnerId
	params["secret_key"] = secretKey

	paramLog := make(map[string]string)

	dateLayout := "2006-01-02 15:04:05"
	paramLog["merchant"] = "MOTION PAY"
	paramLog["endpoint_name"] = apiName
	paramLog["request_method"] = "POST"
	paramLog["url"] = urlPath
	paramLog["created_date"] = time.Now().Format(dateLayout)
	paramLog["created_by"] = "CRON"
	paramLog["note"] = "GET TOKEN MOTION PAY " + apiName

	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	req, err := http.NewRequest(reqMethod, urlPath, payload)
	if err != nil {
		log.Error("Error1", err.Error())
	}

	req.Header.Add("content-type", "application/json")

	paramLog["header"] = FormatRequest(req)
	paramLog["body"] = string(jsonString)

	res, err := http.DefaultClient.Do(req)
	log.Info(res.StatusCode)
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		paramLog["response"] = string(body)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			log.Error("Error create log endpoint motion pay")
			log.Error(err.Error())
		}
		return res.StatusCode, string(body), err
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			paramLog["response"] = string(body)
			_, err = models.CreateEndpoint3rdPartyLog(paramLog)
			if err != nil {
				log.Error("Error create log endpoint motion pay")
				log.Error(err.Error())
			}
			log.Error("Error3", err.Error())
			return http.StatusBadGateway, string(body), err
		}
		paramLog["response"] = string(body)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			log.Error("Error create log endpoint motion pay")
			log.Error(err.Error())
		}

		return http.StatusOK, string(body), nil
	}
}
