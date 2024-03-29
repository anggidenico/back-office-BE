package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FMNotif struct {
	TransID         string `json:"trans_id" form:"trans_id" query:"trans_id"`
	MerchantCode    string `json:"merchant_code" form:"merchant_code" query:"merchant_code"`
	OrderID         string `json:"order_id" form:"order_id" query:"order_id"`
	Amount          string `json:"amount" form:"amount" query:"amount"`
	PaymentMethod   string `json:"payment_method" form:"payment_method" query:"payment_method"`
	MaskCard        string `json:"mask_card" form:"mask_card" query:"mask_card"`
	VaNumber        string `json:"va_number" form:"va_number" query:"va_number"`
	TimeLimit       string `json:"time_limit" form:"time_limit" query:"time_limit"`
	StatusCode      string `json:"status_code" form:"status_code" query:"status_code"`
	StatusDesc      string `json:"status_desc" form:"status_desc" query:"status_desc"`
	FMRefnum        string `json:"fm_refnum" form:"fm_refnum" query:"fm_refnum"`
	DatetimePayment string `json:"datetime_payment" form:"datetime_payment" query:"datetime_payment"`
	ApprovalCode    string `json:"approval_code" form:"approval_code" query:"approval_code"`
	Signature       string `json:"signature" form:"signature" query:"signature"`
}

func SpinGenerateSignature(trNumber, name string) string {
	str := config.MerchantID + `||` +
		config.Partner + `||` +
		`474e50c41d661e651cf0c094d0551b886e3503d25a78a847854f5dc8e7d034a9` + `||` +
		trNumber + `||` + name
	encryptedByte := sha256.Sum256([]byte(str))
	signature := hex.EncodeToString(encryptedByte[:])
	// log.Info("signature :", signature)
	return signature
}

func Spin(trNumber string, name string, params map[string]string) (int, string, error) {
	spin := make(map[string]map[string]string)
	url := make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/orders"
	spin["CREATE_ORDER"] = url
	url = make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/pay/otp"
	spin["CREATE_OTP"] = url
	url = make(map[string]string)
	url["method"] = "POST"
	url["url"] = "https://staging-paywith.spinpay.id/v1/merchants/pay"
	spin["PAY_ORDER"] = url
	signature := SpinGenerateSignature(trNumber, name)
	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	spinUrl := spin[name]
	req, err := http.NewRequest(spinUrl["method"], spinUrl["url"], payload)
	if err != nil {
		// log.Error("Error1", err.Error())
		return http.StatusBadGateway, "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("auth-merchant", config.MerchantID)
	req.Header.Add("auth-partner", config.Partner)
	req.Header.Add("auth-signature", signature)

	// log.Info(FormatRequest(req))

	res, err := http.DefaultClient.Do(req)
	// log.Info(res.StatusCode)
	// if res.StatusCode != 200 {
	// 	// log.Error("Error : ", res.StatusCode)
	// 	return res.StatusCode, "", err
	// }
	if err != nil {
		// log.Error("Error2", err.Error())
		return http.StatusBadGateway, "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// log.Error("Error3", err.Error())
		return http.StatusBadGateway, "", err
	}
	// log.Info(string(body))
	// var sec map[string]interface{}
	// if err = json.Unmarshal(body, &sec); err != nil {
	// 	// log.Error("Error4", err.Error())
	// 	return err.Error()
	// }

	return http.StatusOK, string(body), nil
}

func GenerateReference(prefix string, id string) string {
	x := 6
	y := len(id)
	z := x - y
	r := prefix + strings.Repeat("0", z) + id + strconv.FormatInt(time.Now().Unix(), 10)
	return r
}

// formatRequest generates ascii representation of a request
func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func FMPostPaymentData(params map[string]string) (int, map[string]string, error) {

	paramLog := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	paramLog["merchant"] = "FLASH MOBILE"
	paramLog["endpoint_name"] = "POST Payment Data"
	paramLog["request_method"] = "POST"
	paramLog["url"] = config.FMUrl
	paramLog["created_date"] = time.Now().Format(dateLayout)
	paramLog["created_by"] = strconv.FormatUint(Profile.UserID, 10)
	paramLog["note"] = "POST Payment Data to Flash Mobile"

	paramsJoin := params["merchant_code"] + params["first_name"] + params["last_name"] + params["email"] +
		params["phone"] + params["order_id"] + params["no_reference"] + params["amount"] + params["currency"] +
		params["item_details"] + params["datetime_request"] + params["payment_method"] + params["time_limit"] +
		params["notif_url"] + params["thanks_url"] + config.SecretKey
	// log.Info(paramsJoin)
	paramsJoinByte := []byte(paramsJoin)
	signatureMD5 := md5.Sum(paramsJoinByte)
	signatureSHA1 := sha1.Sum([]byte(fmt.Sprintf("%x", signatureMD5)))
	params["signature"] = fmt.Sprintf("%x", signatureSHA1)
	// log.Info(params["signature"])
	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	req, err := http.NewRequest("POST", config.FMUrl, payload)
	if err != nil {
		// log.Error("Error1", err.Error())
		return http.StatusBadGateway, nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	paramLog["header"] = FormatRequest(req)
	paramLog["body"] = string(jsonString)

	var sec map[string]interface{}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// log.Error("ERROR POST PAYMENT : ", err.Error())
		return http.StatusBadGateway, nil, err
	}
	// log.Println("RES : ", res)
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		paramLog["response"] = string(body)
		// log.Error("Error2", err)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			// log.Error("Error create log endpoint flash mobile")
			// log.Error(err.Error())
		}

		return http.StatusBadGateway, nil, err
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			paramLog["response"] = err.Error()
			_, err = models.CreateEndpoint3rdPartyLog(paramLog)
			if err != nil {
				// log.Error("Error create log endpoint flash mobile")
				// log.Error(err.Error())
			}
			// log.Error("Error3", err.Error())
			return http.StatusBadGateway, nil, err
		}
		paramLog["response"] = string(body)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			// log.Error("Error create log endpoint flash mobile")
			// log.Error(err.Error())
		}
		if err = json.Unmarshal(body, &sec); err != nil {
			// log.Error("Error4", err.Error())
			return http.StatusBadGateway, nil, err
		}
	}

	response := make(map[string]string)
	response["trans_id"] = sec["trans_id"].(string)
	response["order_id"] = sec["order_id"].(string)
	response["merchant_code"] = sec["merchant_code"].(string)
	response["url"] = sec["frontend_url"].(string)
	response["signature"] = fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(sec["trans_id"].(string)+sec["order_id"].(string)+sec["merchant_code"].(string)+config.SecretKey))))))

	return http.StatusOK, response, nil
}
