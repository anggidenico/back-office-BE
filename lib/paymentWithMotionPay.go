package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func CreateOrder(transactionKey string,
	phoneNumber string,
	amount decimal.Decimal,
	transAmountPayment decimal.Decimal,
	transFeeAmountPayment decimal.Decimal,
	productName string,
) (int, map[string]string, error) {
	params := make(map[string]interface{})
	params["reference_code"] = transactionKey
	params["amount"] = amount
	params["description"] = "Order By MNC Asset Managament"

	var dataitem []interface{}

	paramsProduct := make(map[string]interface{})
	paramsProduct["name"] = productName
	paramsProduct["quantity"] = 1
	paramsProduct["amount"] = transAmountPayment.Add(transFeeAmountPayment)

	dataitem = append(dataitem, paramsProduct)

	params["item"] = dataitem

	refCode := transactionKey
	apiName := CREATE_ORDER
	path := PATH_CREATE_ORDER
	requestMethod := "POST"

	status, res, err := requestPaymentMotionPay(
		refCode,
		apiName,
		path,
		requestMethod,
		params,
	)

	response := make(map[string]string)
	var orderId string
	var dataBody map[string]interface{}
	err = json.Unmarshal([]byte(res), &dataBody)
	response["message_action"] = dataBody["message_action"].(string)
	response["message_desc"] = dataBody["message_desc"].(string)
	if err == nil && status == http.StatusOK {
		log.Println(err)
		messageData := dataBody["message_data"].(map[string]interface{})
		orderId = messageData["order_id"].(string)
		status, err, resss := CreateOtp(orderId, phoneNumber)
		var dataBodyCrOtp map[string]interface{}
		err = json.Unmarshal([]byte(resss), &dataBodyCrOtp)
		response["message_action"] = dataBodyCrOtp["message_action"].(string)
		response["message_desc"] = dataBodyCrOtp["message_desc"].(string)
		if err != nil || status != http.StatusOK {
			return status, response, err
		}
	} else {
		return status, response, err
	}

	response["order_id"] = orderId
	response["phone"] = phoneNumber

	return http.StatusOK, response, nil
}

func CreateOtp(orderId string, phoneNumber string) (int, error, string) {
	params := make(map[string]interface{})
	params["order_id"] = orderId
	params["phone"] = phoneNumber

	apiName := CREATE_OTP
	path := PATH_CREATE_OTP
	requestMethod := "POST"

	status, res, err := requestPaymentMotionPay(
		orderId,
		apiName,
		path,
		requestMethod,
		params,
	)
	if status == http.StatusOK {
		return status, err, res
	}

	return http.StatusOK, nil, res
}

func PayOrderMotionPay(orderId string, phoneNumber string, authCode string) (int, string, error) {
	params := make(map[string]interface{})
	params["order_id"] = orderId
	params["phone"] = phoneNumber
	params["auth_code"] = authCode

	apiName := PAY_ORDER
	path := PATH_PAY_ORDER
	requestMethod := "POST"

	status, body, err := requestPaymentMotionPay(
		orderId,
		apiName,
		path,
		requestMethod,
		params,
	)
	if status != http.StatusOK {
		return status, body, err

	}

	return http.StatusOK, body, nil
}

func requestPaymentMotionPay(
	refCodeOrderId string,
	apiName string,
	path string,
	requestMethod string,
	params map[string]interface{}) (int, string, error) {

	paramLog := make(map[string]string)

	url := config.SANDBOX_MP_PAYMENT + path
	dateLayout := "2006-01-02 15:04:05"
	paramLog["merchant"] = "MOTION PAY - PAYMENT"
	paramLog["endpoint_name"] = apiName
	paramLog["request_method"] = requestMethod
	paramLog["url"] = url
	paramLog["created_date"] = time.Now().Format(dateLayout)
	paramLog["created_by"] = strconv.FormatUint(Profile.UserID, 10)
	paramLog["note"] = "PAYMENT WITH MOTION PAY " + apiName

	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	req, err := http.NewRequest(requestMethod, url, payload)
	if err != nil {
		log.Error("Error1", err.Error())
		return http.StatusBadGateway, "", err
	}
	req.Header.Add("auth-merchant", config.MERCHANT_ID_MP_PAYMENT)
	req.Header.Add("auth-partner", config.PARTNER_ID_MP_PAYMENT)
	req.Header.Add("auth-signature", generateSignatureMotionPayment(apiName, refCodeOrderId))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	paramLog["header"] = FormatRequest(req)
	paramLog["body"] = string(jsonString)

	res, err := http.DefaultClient.Do(req)
	log.Info(res.StatusCode)
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)

	if res.StatusCode != http.StatusOK {
		bodyRes, err := ioutil.ReadAll(res.Body)
		paramLog["response"] = string(bodyRes)
		log.Error("Error2", err)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			log.Error("Error create log endpoint motion pay")
			log.Error(err.Error())
		}

		return res.StatusCode, string(bodyRes), err
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			paramLog["response"] = err.Error()
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

func generateSignatureMotionPayment(apiName string, refCodeOrderId string) string {
	token := config.TOKEN_MP_PAYMENT
	var scApp models.ScAppConfig
	_, err := models.GetScAppConfigByCode(&scApp, "TOKEN_MOTION_PAY_PAYMENT")
	if err == nil {
		token = *scApp.AppConfigValue
	}
	sig := config.MERCHANT_ID_MP_PAYMENT + "||" +
		config.PARTNER_ID_MP_PAYMENT + "||" +
		token + "||" +
		refCodeOrderId + "||" +
		apiName

	encryptedPasswordByte := sha256.Sum256([]byte(sig))
	signature := hex.EncodeToString(encryptedPasswordByte[:])

	return signature

}
