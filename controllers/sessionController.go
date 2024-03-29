package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	_ "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/badoux/checkmail"
	"github.com/denisbrodbeck/machineid"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	ua "github.com/mileusna/useragent"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gopkg.in/gomail.v2"
)

func Register(c echo.Context) error {

	// log.Println(c.Request().Form)
	// log.Println(c.Request().Form.Encode())
	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}
	password := c.FormValue("password")
	if password == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}
	phone := c.FormValue("phone")
	if phone == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}

	token := c.FormValue("token") //player_id

	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		// log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}
	// err = checkmail.ValidateHost(email)
	// if err != nil {
	// 	// log.Error("Email host is not valid")
	// 	return ctx.TextResponse("Email is not valid", fasthttp.StatusBadRequest)
	// }
	// if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
	// 	// log.Error("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
	// 	return ctx.TextResponse("Email is not valid", fasthttp.StatusBadRequest)
	// }
	var user []models.ScUserLogin
	params := make(map[string]string)
	params["ulogin_email"] = email
	params["rec_status"] = "1"
	status, err = models.GetAllScUserLogin(&user, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email " + email)
		return lib.CustomError(status, err.Error(), "Error get email")
	}
	if len(user) > 0 {
		// log.Error("Email " + email + " already registered")
		return lib.CustomError(http.StatusBadRequest, "Email "+email+" already registered", "Data yang kamu masukkan sudah terdaftar.\nSilakan masukkan data lainnya atau hubungi Customer Service - 021 29709696.")
	}
	params = make(map[string]string)
	params["ulogin_mobileno"] = phone
	params["rec_status"] = "1"
	status, err = models.GetAllScUserLogin(&user, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get phone number " + phone)
		return lib.CustomError(status, err.Error(), "Error get phone number")
	}
	if len(user) > 0 {
		// log.Error("Phone number " + phone + " already registered")
		return lib.CustomError(http.StatusBadRequest, "Phone number "+phone+" already registered", "Data yang kamu masukkan sudah terdaftar.\nSilakan masukkan data lainnya atau hubungi Customer Service - 021 29709696.")
	}

	// Validate password
	length, number, upper, special := verifyPassword(password)
	if length == false || number == false || upper == false || special == false {
		// log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])

	// Set expired for token
	date := time.Now().AddDate(0, 0, 1)
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Format(dateLayout)

	// Generate verify key
	verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
	verifyKey := hex.EncodeToString(verifyKeyByte[:])

	// Input to database
	if token != "" {
		params["token_notif"] = token
	}
	params["ulogin_email"] = email
	params["ulogin_name"] = email
	params["ulogin_full_name"] = email
	params["ulogin_must_changepwd"] = "0"
	params["user_category_key"] = "1"
	params["user_dept_key"] = "1"
	params["last_password_changed"] = time.Now().Format(dateLayout)
	params["ulogin_password"] = encryptedPassword
	params["verified_email"] = "0"
	params["verified_mobileno"] = "0"
	params["ulogin_mobileno"] = phone
	params["ulogin_enabled"] = "1"
	params["ulogin_locked"] = "0"
	params["ulogin_failed_count"] = "0"
	params["user_category_key"] = "1"
	params["last_access"] = time.Now().Format(dateLayout)
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["accept_login_tnc"] = "1"
	params["allowed_sharing_login"] = "1"
	params["string_token"] = verifyKey
	params["token_expired"] = expired
	params["rec_status"] = "1"
	params["rec_attribute_id3"] = c.Request().UserAgent()

	status, err = models.CreateScUserLogin(params)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	// Send email
	t := template.New("index-email-activation.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-email-activation.html")
	if err != nil {
		// log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		Url     string
		FileUrl string
	}{Url: config.BaseUrl + "/verifyemail?token=" + verifyKey, FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
		// log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MotionFunds] Verifikasi Email Kamu")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error send email")
	} else {
		// log.Info("Email sent")
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
	// 	// log.Error(err)
	// 	return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error send email")
	// }
	// // log.Info("Email sent")
	var data models.ScUserLoginRegister
	data.UloginEmail = email
	data.UloginMobileno = phone

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	return c.JSON(http.StatusOK, response)
}

func VerifyEmail(c echo.Context) error {
	var err error
	// Get parameter key
	token := c.QueryParam("token")
	if token == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Token tidak ditemukan")
	}
	params := make(map[string]string)
	params["string_token"] = token
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(http.StatusBadRequest, "Error get email", "Gagal mendapatkan data email")
	}
	if len(userLogin) < 1 {
		// log.Error("No matching token " + token)
		return lib.CustomError(http.StatusBadRequest, "Token not found", "Token tidak ditemukan")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)

	// Check if token is expired
	dateLayout := "2006-01-02 15:04:05"
	expired, err := time.Parse(dateLayout, *accountData.TokenExpired)
	if err != nil {
		// log.Error("Error parsing data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error parsing data")
	}
	now := time.Now()
	if now.After(expired) {
		// log.Error("Token is expired")
		return lib.CustomError(http.StatusInternalServerError, "Token is expired", "Token anda sudah kadaluarsa. Silakan kirim ulang email verifikasi.")
	}
	// log.Info("Success verify email")
	// Set expired for otp
	date := time.Now().Add(1 * time.Minute)
	expiredOTP := date.Format(dateLayout)

	// Send otp
	otp, err := sendOTP("0", *accountData.UloginMobileno)
	if err != nil {
		// log.Error(err.Error())
		//return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
	}
	if otp == "" {
		// log.Error("Failed send otp")
		//return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
	} else {
		// log.Info("Success send otp")
	}

	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["otp_number"] = otp
	params["otp_number_expired"] = expiredOTP
	params["verified_email"] = "1"
	params["last_verified_email"] = time.Now().Format(dateLayout)
	params["string_token"] = ""

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func VerifyOtp(c echo.Context) error {
	var err error
	var status int
	// Get parameter key
	otp := c.FormValue("otp")
	if otp == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	params := make(map[string]string)
	params["otp_number"] = otp
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("No matching otp " + otp)
		return lib.CustomError(http.StatusBadRequest, "OTP not found", "OTP not found")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)

	// Check if token is expired
	dateLayout := "2006-01-02 15:04:05"
	expired, err := time.Parse(dateLayout, *accountData.OtpNumberExpired)
	if err != nil {
		// log.Error("Error parsing data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error parsing data")
	}
	now := time.Now()
	if now.After(expired) {
		// log.Error("OTP is expired")
		return lib.CustomError(http.StatusInternalServerError, "OTP is expired", "OTP is expired")
	}
	// log.Info("Success verify OTP")

	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["otp_number"] = ""
	params["ulogin_enabled"] = "1"
	params["verified_mobileno"] = "1"
	params["last_verified_mobileno"] = time.Now().Format(dateLayout)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	// Create session key
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()

	atClaims := jwt.MapClaims{}
	paramsRequest := make(map[string]string)
	paramsRequest["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsRequest["orderBy"] = "oa_request_key"
	paramsRequest["orderType"] = "DESC"
	var request []models.OaRequest
	status, err = models.GetAllOaRequest(&request, config.LimitQuery, 0, true, paramsRequest)
	if err != nil {
		// log.Error(err.Error())
	} else if len(request) > 0 {
		if request[0].Oastatus != nil && *request[0].Oastatus > 0 {
			var lookup models.GenLookup
			status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
			if err != nil {
				// log.Error(err.Error())
			} else {
				if lookup.LkpName != nil && *lookup.LkpName != "" {
					atClaims["oa_status"] = *lookup.LkpName
				}
			}
		}
	}
	if accountData.RoleKey != nil && *accountData.RoleKey > 0 {
		atClaims["role_key"] = *accountData.RoleKey
		paramsRole := make(map[string]string)
		paramsRole["role_key"] = strconv.FormatUint(*accountData.RoleKey, 10)
		var role []models.ScRole
		_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
		if err != nil {
			// log.Error(err.Error())
		} else if len(role) > 0 {
			if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
				atClaims["role_category_key"] = *role[0].RoleCategoryKey
			}
		}
	}
	atClaims["uuid"] = uuidString
	atClaims["email"] = accountData.UloginEmail
	atClaims["exp"] = time.Now().Add(time.Minute * 50).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Secret))
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized, err.Error(), "Login failed")
	}

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, params, true)
	paramsSession["session_id"] = uuidString
	paramsSession["login_date"] = time.Now().Format(dateLayout)
	paramsSession["rec_status"] = "1"
	if err == nil && len(loginSession) > 0 {
		// log.Info("Active session for previous login, overwrite whit new session")
		if len(loginSession) > 1 {

		}
		paramsSession["login_session_key"] = strconv.FormatUint(loginSession[0].LoginSessionKey, 10)

		status, err = models.UpdateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error update session")
			return lib.CustomError(status, "Error update session", "Login failed")
		}
	} else {
		status, err = models.CreateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error create session")
			return lib.CustomError(status, "Error create session", "Login failed")
		}
	}
	// log.Info("Success login")

	var data models.ScLoginSessionInfo
	data.SessionID = token
	// log.Info(data)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	return c.JSON(http.StatusOK, response)
}

func Login(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	password := c.FormValue("password")
	if password == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["rec_status"] = "1"
	params["ulogin_email"] = email
	params["ulogin_name"] = email
	params["user_category_key"] = "1"
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLoginByNameOrEmail(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get Username")
		return lib.CustomError(status, "Email/Username atau Kata Sandi kamu salah", "Email/Username atau Kata Sandi kamu salah")
	}
	if len(userLogin) < 1 {
		// log.Error("Email/Username not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email/Username atau Kata Sandi kamu salah", "Email/Username atau Kata Sandi kamu salah")
	}

	accountData := userLogin[0]
	// log.Info(accountData)

	if *accountData.VerifiedEmail != 1 || accountData.VerifiedMobileno != 1 {
		// log.Error("Email or Mobile number not verified")
		return lib.CustomError(http.StatusUnauthorized, "Email atau Nomor Telepon belum terverifikasi", "Email atau Nomor Telepon belum terverifikasi")
	}

	if accountData.UloginLocked == uint8(1) {
		// log.Error("User is locked")
		countWrongPass := strconv.FormatUint(accountData.UloginFailedCount, 10)
		return lib.CustomError(http.StatusUnauthorized, "Akun kamu terkunci karena salah memasukkan password "+countWrongPass+" kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.", "Akun kamu terkunci karena salah memasukkan password "+countWrongPass+" kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.")
	}

	if accountData.UloginEnabled == uint8(0) {
		// log.Error("User is Disable")
		return lib.CustomError(http.StatusUnauthorized, "Akun anda tidak aktif, Mohon hubungi admin MNCDuit untuk mengaktifkan akun anda kembali.", "Akun anda tidak aktif, Mohon hubungi admin MNCDuit untuk mengaktifkan akun anda kembali.")
	}

	dateLayout := "2006-01-02 15:04:05"

	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != accountData.UloginPassword {
		//update ulogin_failed_count wrong password
		paramsUpdate := make(map[string]string)
		uloginkey := strconv.FormatUint(accountData.UserLoginKey, 10)
		countFalse := accountData.UloginFailedCount + 1
		strCountFalse := strconv.FormatUint(countFalse, 10)
		paramsUpdate["user_login_key"] = uloginkey
		paramsUpdate["ulogin_failed_count"] = strCountFalse

		var scApp models.ScAppConfig
		status, err = models.GetScAppConfigByCode(&scApp, "LOGIN_ATTEMPT")
		if err != nil {
			// log.Error(err.Error())
		}

		countWrong, _ := strconv.ParseUint(*scApp.AppConfigValue, 10, 64)

		if countFalse >= countWrong {
			paramsUpdate["ulogin_locked"] = "1"
			paramsUpdate["locked_date"] = time.Now().Format(dateLayout)
		}

		_, err = models.UpdateScUserLogin(paramsUpdate)
		if err != nil {
			// log.Error(err.Error())
			// log.Error("erroe update ulogin_failed_count wrong password")
		}

		if countFalse >= countWrong {
			// log.Error("Wrong password, user is locked")
			return lib.CustomError(http.StatusUnauthorized, "Akun kamu terkunci karena salah memasukkan password "+*scApp.AppConfigValue+" kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.", "Akun kamu terkunci karena salah memasukkan password "+*scApp.AppConfigValue+" kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.")
		} else {
			// log.Error("Wrong password")
			return lib.CustomError(http.StatusUnauthorized, "Email/Username atau Kata Sandi kamu salah", "Email/Username atau Kata Sandi kamu salah")
		}
	}

	// Create session key
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()

	atClaims := jwt.MapClaims{}
	paramsRequest := make(map[string]string)
	paramsRequest["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsRequest["orderBy"] = "oa_request_key"
	paramsRequest["orderType"] = "DESC"
	var request []models.OaRequest
	status, err = models.GetAllOaRequest(&request, config.LimitQuery, 0, true, paramsRequest)
	if err != nil {
		atClaims["oa_status"] = ""
		// log.Error(err.Error())
	} else if len(request) > 0 {
		if request[0].Oastatus != nil && *request[0].Oastatus > 0 {
			if *request[0].Oastatus == 260 || *request[0].Oastatus == 261 || *request[0].Oastatus == 262 { //sudah approve

				var lookup models.GenLookup
				status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
				if err != nil {
					// log.Error(err.Error())
				} else {
					if lookup.LkpName != nil && *lookup.LkpName != "" {
						atClaims["oa_status"] = *lookup.LkpName
					}
				}
			} else {
				if request[0].OaRequestType != nil && *request[0].OaRequestType != 127 {
					atClaims["oa_status"] = "PENGKINIAN"
				} else {
					var lookup models.GenLookup
					status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
					if err != nil {
						// log.Error(err.Error())
					} else {
						if lookup.LkpName != nil && *lookup.LkpName != "" {
							atClaims["oa_status"] = *lookup.LkpName
						}
					}
				}
			}
		}
	} else {
		atClaims["oa_status"] = ""
	}
	if accountData.RoleKey != nil && *accountData.RoleKey > 0 {
		atClaims["role_key"] = *accountData.RoleKey
		paramsRole := make(map[string]string)
		paramsRole["role_key"] = strconv.FormatUint(*accountData.RoleKey, 10)
		var role []models.ScRole
		_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
		if err != nil {
			// log.Error(err.Error())
		} else if len(role) > 0 {
			if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
				atClaims["role_category_key"] = *role[0].RoleCategoryKey
			}
		}
	}
	atClaims["uuid"] = uuidString
	atClaims["exp"] = time.Now().Add(time.Minute * 50).Unix()
	atClaims["email"] = accountData.UloginEmail
	atClaims["pin"] = accountData.UloginPin
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Secret))
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized, err.Error(), "Login failed")
	}

	// sessionKey := base64.StdEncoding.EncodeToString([]byte(uuidString))
	// expired := date.Add(time.Second * time.Duration(config.SessionExpired)).Format(dateLayout)

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, paramsSession, true)
	paramsSession["session_id"] = uuidString
	paramsSession["login_date"] = time.Now().Format(dateLayout)
	paramsSession["rec_status"] = "1"
	paramsSession["rec_attribute_id3"] = c.Request().UserAgent()
	if err == nil && len(loginSession) > 0 {
		// log.Info("Active session for previous login, overwrite with new session")
		if len(loginSession) > 1 {

		}
		paramsSession["login_session_key"] = strconv.FormatUint(loginSession[0].LoginSessionKey, 10)

		status, err = models.UpdateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error update session")
			return lib.CustomError(status, "Error update session", "Login failed")
		}
	} else {
		status, err = models.CreateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error create session")
			return lib.CustomError(status, "Error create session", "Login failed")
		}
	}

	// update ulogin_failed_count = 0 if success login
	paramsUpdate := make(map[string]string)
	uloginkey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUpdate["user_login_key"] = uloginkey
	paramsUpdate["ulogin_failed_count"] = "0"

	tokenNotif := c.FormValue("token")
	if tokenNotif != "" {
		paramsUpdate["token_notif"] = tokenNotif
		paramsUpdate["last_update_token_notif"] = time.Now().Format(dateLayout)

		//update null token_notif
		_, err = models.SetNullTokenNotif(tokenNotif)
	}

	_, err = models.UpdateScUserLogin(paramsUpdate)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("erroe update ulogin_failed_count = 0 if success login")
	}

	// log.Info("Success login")

	var data models.ScLoginSessionInfo
	data.SessionID = token
	if accountData.UloginMustChangepwd == uint8(1) {
		// data.MustChangePassword = true
	} else {
		// data.MustChangePassword = false
	}
	// log.Info(data)

	//LOG LOGIN
	saveLogLogin(c, paramsSession)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	// log.Info(response)
	return c.JSON(http.StatusOK, response)
}

func ResendVerification(c echo.Context) error {
	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}

	params := make(map[string]string)
	params["ulogin_email"] = email
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(status, err.Error(), "Failed get email")
	}
	if len(userLogin) < 1 {
		// log.Error("No matching email " + email)
		return lib.CustomError(http.StatusBadRequest, "Email not registered", "Email not registered")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)

	dateLayout := "2006-01-02 15:04:05"
	if accountData.VerifiedEmail != nil && *accountData.VerifiedEmail == 1 {
		date := time.Now().Add(1 * time.Minute)
		expiredOTP := date.Format(dateLayout)

		// Send otp
		otp, err := sendOTP("0", *accountData.UloginMobileno)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
		}
		if otp == "" {
			// log.Error("Failed send otp")
			return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
		}

		params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
		params["otp_number"] = otp
		params["otp_number_expired"] = expiredOTP
		params["verified_email"] = "1"
		params["last_verified_email"] = time.Now().Format(dateLayout)
		params["string_token"] = ""

		_, err = models.UpdateScUserLogin(params)
		if err != nil {
			// log.Error("Error update user data")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}

		// log.Info("Success send otp")
	} else {
		// Set expired for token
		date := time.Now().AddDate(0, 0, 1)
		expired := date.Format(dateLayout)

		// Generate verify key
		verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
		verifyKey := hex.EncodeToString(verifyKeyByte[:])

		// Update token
		params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
		params["string_token"] = verifyKey
		params["token_expired"] = expired

		status, err = models.UpdateScUserLogin(params)
		if err != nil {
			// log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update token")
		}

		// Send email
		t := template.New("index-email-activation.html")

		var err error
		t, err = t.ParseFiles(config.BasePath + "/mail/index-email-activation.html")
		if err != nil {
			// log.Println(err)
		}

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, struct {
			Url     string
			FileUrl string
		}{Url: config.BaseUrl + "/verifyemail?token=" + verifyKey, FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
			// log.Println(err)
		}

		result := tpl.String()

		mailer := gomail.NewMessage()
		// mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", email)
		mailer.SetHeader("Subject", "[MotionFunds] Verifikasi Email Kamu")
		mailer.SetBody("text/html", result)

		err = lib.SendEmail(mailer)
		if err != nil {
			// log.Error(err)
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
		} else {
			// log.Info("Email sent")
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
		// 	// log.Error(err)
		// 	return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
		// }
		// // log.Info("Email sent")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
func ForgotPassword(c echo.Context) error {
	var err error
	var status int

	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}
	params := make(map[string]string)
	params["ulogin_email"] = email
	params["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(status, err.Error(), "Failed get email")
	}
	if len(userLogin) < 1 {
		// log.Error("No matching email " + email)
		return lib.CustomError(http.StatusBadRequest, "Pastikan email yang anda gunakan terdaftar, cek kembali email yang anda masukkan", "Pastikan email yang anda gunakan terdaftar, cek kembali email yang anda masukkan")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)

	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "=+%*/!@#$?"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	encryptedPasswordByte := sha256.Sum256([]byte(str))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	dateLayout := "2006-01-02 15:04:05"
	params = make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["string_token"] = str
	params["last_password_changed"] = time.Now().Format(dateLayout)
	params["ulogin_password"] = encryptedPassword
	params["ulogin_locked"] = "0"
	params["ulogin_failed_count"] = "0"
	params["ulogin_must_changepwd"] = "1"

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	name := ""

	var fullName models.FullNameData
	_, err = models.GetFullName(&fullName, strconv.FormatUint(accountData.UserLoginKey, 10))
	if err != nil {
		name = accountData.UloginFullName
	} else {
		name = fullName.FullName
	}

	// Send email
	t := template.New("index-forget-password.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-forget-password.html")
	if err != nil {
		// log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		Password string
		FileUrl  string
		Name     string
	}{Password: str, FileUrl: config.ImageUrl + "/images/mail", Name: name}); err != nil {
		// log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MotionFunds] Lupa Kata Sandi")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	} else {
		// log.Info("Email sent")
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
	// 	// log.Error(err)
	// 	return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	// }
	// // log.Info("Email sent")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func GetUserLogin(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	var oaRequestDB []models.OaRequest
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["orderBy"] = "oa_request_key"
	params["orderType"] = "DESC"
	_, err = models.GetAllOaRequest(&oaRequestDB, 0, 0, true, params)
	if err != nil {
		// log.Error(err.Error())
	}

	var oaRequestActive models.OaRequest
	var requestKey string
	if len(oaRequestDB) > 0 {
		oaRequestActive = oaRequestDB[0]
		requestKey = strconv.FormatUint(oaRequestDB[0].OaRequestKey, 10)
		if len(oaRequestDB) > 1 {
			for _, oareq := range oaRequestDB {
				if *oareq.Oastatus == 260 || *oareq.Oastatus == 261 || *oareq.Oastatus == 262 { //sudah approve
					oaRequestActive = oareq
					requestKey = strconv.FormatUint(oareq.OaRequestKey, 10)
					break
				}
			}
		}
	}

	var personalDataDB models.OaPersonalData
	var riskProfileDB models.OaRiskProfile
	var responseData models.UserProfile
	if requestKey != "" {
		_, err = models.GetOaPersonalData(&personalDataDB, requestKey, "oa_request_key")
		if err != nil {
			// log.Error(err.Error())
		}
		_, err = models.GetOaRiskProfile(&riskProfileDB, requestKey, "oa_request_key")
		if err != nil {
			// log.Error(err.Error())
		}

		var riskDB models.MsRiskProfile
		if riskProfileDB.RiskProfileKey > 0 {
			_, err = models.GetMsRiskProfile(&riskDB, strconv.FormatUint(riskProfileDB.RiskProfileKey, 10))
			if err != nil {
				// log.Error(err.Error())
			}
		}

		var customerDB models.MsCustomer
		if lib.Profile.CustomerKey != nil && *lib.Profile.CustomerKey > 0 {
			_, err = models.GetMsCustomer(&customerDB, strconv.FormatUint(*lib.Profile.CustomerKey, 10))
			if err != nil {
				// log.Error(err.Error())
			}
		}

		responseData.FullName = personalDataDB.FullName
		responseData.CIF = customerDB.UnitHolderIDno
		if customerDB.SidNo != nil {
			responseData.SID = *customerDB.SidNo
		}
		if customerDB.CifSuspendFlag == 0 {
			responseData.CifSuspendFlag = false
		} else {
			responseData.CifSuspendFlag = true
		}
		responseData.Email = lib.Profile.Email
		responseData.PhoneNumber = lib.Profile.PhoneNumber
		responseData.RiskProfile.RiskProfileKey = riskDB.RiskProfileKey
		responseData.RiskProfile.RiskCode = riskDB.RiskCode
		responseData.RiskProfile.RiskName = riskDB.RiskName
		responseData.RiskProfile.RiskDesc = riskDB.RiskDesc
		if riskProfileDB.ScoreResult != nil {
			responseData.RiskProfile.Score = *riskProfileDB.ScoreResult
		}
		responseData.RecImage1 = lib.Profile.RecImage1

		bankAccountKey := ""
		if oaRequestActive.CustomerKey != nil {
			var cusBank []models.MsCustomerBankAccount
			paramsCusBank := make(map[string]string)
			paramsCusBank["customer_key"] = strconv.FormatUint(*oaRequestActive.CustomerKey, 10)
			paramsCusBank["rec_status"] = "1"
			paramsCusBank["flag_priority"] = "1"
			_, err = models.GetAllMsCustomerBankAccount(&cusBank, paramsCusBank)
			if err == nil {
				bankAccountKey = strconv.FormatUint(cusBank[0].BankAccountKey, 10)
				if len(cusBank) > 1 {
					for _, bk := range cusBank {
						if bk.FlagPriority == uint8(1) {
							bankAccountKey = strconv.FormatUint(bk.BankAccountKey, 10)
							break
						}
					}
				}
			}
		} else {
			var cusBank []models.OaRequestByField
			_, err = models.GetOaRequestBankByField(&cusBank, "oa_request_key", requestKey)
			if err == nil && len(cusBank) > 0 {
				bankAccountKey = strconv.FormatUint(cusBank[0].BankAccountKey, 10)
				if len(cusBank) > 1 {
					for _, bk := range cusBank {
						if bk.FlagPriority == uint64(1) {
							bankAccountKey = strconv.FormatUint(bk.BankAccountKey, 10)
							break
						}
					}
				}
			}
		}

		if bankAccountKey != "" {
			var bankAccountDB models.MsBankAccount
			if personalDataDB.BankAccountKey != nil && *personalDataDB.BankAccountKey > 0 {
				_, err = models.GetBankAccount(&bankAccountDB, strconv.FormatUint(*personalDataDB.BankAccountKey, 10))
				if err != nil {
					// log.Error(err.Error())
				}
			}

			var bankDB models.MsBank
			if bankAccountDB.BankKey > 0 {
				_, err = models.GetMsBank(&bankDB, strconv.FormatUint(bankAccountDB.BankKey, 10))
				if err != nil {
					// log.Error(err.Error())
				}
			}
			responseData.BankAcc.BankName = bankDB.BankName
			responseData.BankAcc.AccountNo = bankAccountDB.AccountNo
			responseData.BankAcc.AccountHolderName = bankAccountDB.AccountHolderName
			responseData.BankAcc.BranchName = bankAccountDB.BranchName
		}
	} else {
		responseData.FullName = lib.Profile.Email
		// responseData.CIF = customerDB.UnitHolderIDno
		// responseData.SID = *customerDB.SidNo
		responseData.CifSuspendFlag = false
		responseData.Email = lib.Profile.Email
		responseData.PhoneNumber = lib.Profile.PhoneNumber
		// responseData.RiskProfile.RiskProfileKey = riskDB.RiskProfileKey
		// responseData.RiskProfile.RiskCode = riskDB.RiskCode
		// responseData.RiskProfile.RiskName = riskDB.RiskName
		// responseData.RiskProfile.RiskDesc = riskDB.RiskDesc
		// responseData.RiskProfile.Score = *riskProfileDB.ScoreResult
		responseData.RecImage1 = lib.Profile.RecImage1
		// responseData.BankAcc.BankName = bankDB.BankName
		// responseData.BankAcc.AccountNo = bankAccountDB.AccountNo
		// responseData.BankAcc.AccountHolderName = bankAccountDB.AccountHolderName
		// responseData.BankAcc.BranchName = bankAccountDB.BranchName

	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func UploadProfilePic(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)
	filePath := config.ImageUrl + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/profile"
	err = os.MkdirAll(filePath, 0755)
	if err != nil {
		// log.Error(err.Error())
	} else {

		var file *multipart.FileHeader
		file, err = c.FormFile("pic")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: pic")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				// log.Println("Generate filename:", filename)
				var trans []models.TrTransaction
				getParams := make(map[string]string)
				getParams["rec_image1"] = filename + extension
				_, err = os.Stat(filePath + "/" + filename + extension)
				if err != nil {
					if os.IsNotExist(err) {
						_, err = models.GetAllTrTransaction(&trans, getParams)
						if (err == nil && len(trans) < 1) || err != nil {
							break
						}
					}
				}
				// // log.Infoln("========== get parameter filename dan extension ==========")
				// log.Infoln(getParams)
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, filePath+"/"+filename+extension)
			// log.Println("UPLOAD PATH: ", filePath+"/"+filename+extension)
			if err != nil {
				// log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
		}
	}
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ChangePassword(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	recentPassword := c.FormValue("recent_password")
	if recentPassword == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPassword1 := c.FormValue("new_password1")
	if newPassword1 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPassword2 := c.FormValue("new_password2")
	if newPassword2 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["ulogin_email"] = lib.Profile.Email
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(status, "Error get email", "Error get email")
	}
	if len(userLogin) < 1 {
		// log.Error("Email not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email not registered", "Email not registered")
	}

	accountData := userLogin[0]
	// log.Info(accountData)

	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(recentPassword))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != accountData.UloginPassword {
		// log.Error("Wrong password")
		return lib.CustomError(http.StatusUnauthorized, "Wrong password", "Wrong password")
	}

	if newPassword1 != newPassword2 {
		// log.Error("Password doesnt match")
		return lib.CustomError(http.StatusBadRequest, "Password doesnt match", "Password doesnt match")
	}
	// Validate password
	length, number, upper, special := verifyPassword(newPassword1)
	if length == false || number == false || upper == false || special == false {
		// log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	// Encrypt password
	encryptedPasswordByteNew := sha256.Sum256([]byte(newPassword1))
	encryptedPasswordNew := hex.EncodeToString(encryptedPasswordByteNew[:])

	dateLayout := "2006-01-02 15:04:05"
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["ulogin_password"] = encryptedPasswordNew
	params["last_password_changed"] = time.Now().Format(dateLayout)
	params["ulogin_must_changepwd"] = "0"
	params["ulogin_failed_count"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//insert table sc_user_message
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	strUserLoginKey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Perubahan Kata Sandi Berhasil"
	body := "Kata sandi kamu berhasil berubah. Apabila kamu tidak merasa melakukan perubahan kata sandi, mohon segera menghubungi customer service kami."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strKey

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

	//kirim email
	t := template.New("index-sukses-ubah-password.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-sukses-ubah-password.html")
	if err != nil {
		// log.Println(err)
	}

	var customer models.MsCustomer

	fullnameuser := accountData.UloginFullName

	if accountData.CustomerKey != nil {
		strCustomerKey := strconv.FormatUint(*accountData.CustomerKey, 10)
		status, err = models.GetMsCustomer(&customer, strCustomerKey)
		if err == nil {
			fullnameuser = customer.FullName
		}
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl,
		struct {
			Name    string
			FileUrl string
		}{
			Name:    fullnameuser,
			FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
		// log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", accountData.UloginEmail)
	mailer.SetHeader("Subject", "[MotionFunds] Berhasil Merubah Kata Sandi")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Error(err)
	} else {
		// log.Info("Email sent")
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
	// 	// log.Error(err)
	// }
	// // log.Info("Email sent")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func TransactionCutOffTime(c echo.Context) error {
	var err error
	var str_message string
	var cfg_value string

	var appConfig models.ScAppConfig
	_, err = models.GetScAppConfigByCode(&appConfig, "TRX_CUTOFF_TIME")
	if err != nil {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(http.StatusBadRequest, str_message, "Fail to get Config TRX_CUTOFF_TIME")
	}

	cfg_value = *appConfig.AppConfigValue

	dateLayout := "2006-01-02 15:04:05"
	curr_date_time := time.Now()
	str_date_part := curr_date_time.Format("2006-01-02")
	str_time_part := strings.Split(cfg_value, ":")

	date_time_string := str_date_part + " " + str_time_part[0] + ":" + str_time_part[1] + ":00"
	curr_date_time, err = time.Parse(dateLayout, date_time_string)
	if err != nil {
		str_message = err.Error()
		// log.Error(str_message)
		return lib.CustomError(http.StatusBadRequest, str_message, "Fail to parse time")

	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = curr_date_time.Format(dateLayout)

	return c.JSON(http.StatusOK, response)

}

func CurrentTime(c echo.Context) error {

	dateLayout := "2006-01-02 15:04:05"

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = time.Now().Format(dateLayout)
	return c.JSON(http.StatusOK, response)
}

func verifyPassword(s string) (length, number, upper, special bool) {
	var letter bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c):
			letter = true
		default:
			//return false, false, false, false
		}
	}
	length = letter && len(s) >= 8

	return
}

func sendOTP(gateway, phone string) (string, error) {

	paramLog := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	paramLog["merchant"] = "CITCALL"
	paramLog["endpoint_name"] = "Asynccall"
	paramLog["request_method"] = "POST"
	paramLog["url"] = config.CitcallUrl
	paramLog["created_date"] = time.Now().Format(dateLayout)
	paramLog["created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramLog["note"] = "send OTP with citcall"

	curlParam := make(map[string]string)
	curlParam["retry"] = gateway
	curlParam["msisdn"] = phone
	jsonString, err := json.Marshal(curlParam)
	payload := strings.NewReader(string(jsonString))
	req, err := http.NewRequest("POST", config.CitcallUrl, payload)
	if err != nil {
		// log.Error("Error1", err.Error())
		return "", err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Apikey 7f837aea98ceea9efcd33ca1d435c9cf")

	paramLog["header"] = lib.FormatRequest(req)
	paramLog["body"] = string(jsonString)

	var sec map[string]interface{}
	res, err := http.DefaultClient.Do(req)
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		paramLog["response"] = string(body)
		// log.Error("Error2", err)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			// log.Error("Error create log endpoint citcall")
			// log.Error(err.Error())
		}

		return "", err
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			paramLog["response"] = err.Error()
			_, err = models.CreateEndpoint3rdPartyLog(paramLog)
			if err != nil {
				// log.Error("Error create log endpoint flash citcall")
				// log.Error(err.Error())
			}
			// log.Error("Error3", err.Error())
			return "", err
		}
		paramLog["response"] = string(body)
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			// log.Error("Error create log endpoint citcall")
			// log.Error(err.Error())
		}
		if err = json.Unmarshal(body, &sec); err != nil {
			// log.Error("Error4", err.Error())
			return "", err
		}
	}
	var otp string
	if sec["rc"].(float64) == 0 {
		token := sec["token"].(string)
		otp = token[len(token)-4:]
	}
	return otp, nil
}

func LoginBo(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	password := c.FormValue("password")
	if password == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["rec_status"] = "1"
	params["ulogin_email"] = email
	params["ulogin_name"] = email
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLoginByNameOrEmail(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get Username")
		return lib.CustomError(status, "Error get Email/Username", "Error get Email/Username")
	}
	if len(userLogin) < 1 {
		// log.Error("Email/Username not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email/Username not registered", "Email/Username not registered")
	}

	accountData := userLogin[0]

	//check user USR_BOHO / USR_BOBRANCH
	if (accountData.UserCategoryKey != 2) && (accountData.UserCategoryKey != 3) {
		// log.Error("Email/Username not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email/Username not registered", "Email/Username not registered")
	}

	// log.Info(accountData)

	if *accountData.VerifiedEmail != 1 || accountData.VerifiedMobileno != 1 {
		// log.Error("Email or Mobile number not verified")
		return lib.CustomError(http.StatusUnauthorized, "Email or Mobile number not verified", "Email or Mobile number not verified")
	}

	if accountData.UloginLocked == uint8(1) {
		// log.Error("User is locked")
		return lib.CustomError(http.StatusUnauthorized, "Akun kamu terkunci karena salah memasukkan password 3 kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.", "Akun kamu terkunci karena salah memasukkan password 3 kali berturut-turut. Silakan menunggu 1 jam lagi untuk login atau hubungi Customer Service untuk informasi lebih lanjut.")
	}

	if accountData.UloginEnabled == uint8(0) {
		// log.Error("User is Disable")
		return lib.CustomError(http.StatusUnauthorized, "Akun kamu tidak aktif. Silakan menghubungi Customer Service untuk informasi lebih lanjut.", "Akun kamu tidak aktif. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
	}

	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != accountData.UloginPassword {
		//update ulogin_failed_count wrong password
		paramsUpdate := make(map[string]string)
		uloginkey := strconv.FormatUint(accountData.UserLoginKey, 10)
		countFalse := accountData.UloginFailedCount + 1
		strCountFalse := strconv.FormatUint(countFalse, 10)
		paramsUpdate["user_login_key"] = uloginkey
		paramsUpdate["ulogin_failed_count"] = strCountFalse

		var scApp models.ScAppConfig
		status, err = models.GetScAppConfigByCode(&scApp, "LOGIN_ATTEMPT")
		if err != nil {
			// log.Error(err.Error())
		}

		countWrong, _ := strconv.ParseUint(*scApp.AppConfigValue, 10, 64)

		if countFalse >= countWrong {
			paramsUpdate["ulogin_locked"] = "1"
		}

		_, err = models.UpdateScUserLogin(paramsUpdate)
		if err != nil {
			// log.Error(err.Error())
			// log.Error("erroe update ulogin_failed_count wrong password")
		}

		if countFalse >= countWrong {
			// log.Error("Wrong password, user is locked")
			return lib.CustomError(http.StatusUnauthorized, "Akun kamu terkunci karena salah memasukkan password 3 kali berturut-turut. Silakan menghubungi Customer Service untuk informasi lebih lanjut.", "Akun kamu terkunci karena salah memasukkan password 3 kali berturut-turut. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
		} else {
			// log.Error("Wrong password")
			return lib.CustomError(http.StatusUnauthorized, "Password yang kamu masukkan salah", "Password yang kamu masukkan salah")
		}
	}

	// Create session key
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()

	atClaims := jwt.MapClaims{}
	// paramsRequest := make(map[string]string)
	// paramsRequest["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	// paramsRequest["orderBy"] = "oa_request_key"
	// paramsRequest["orderType"] = "DESC"
	// var request []models.OaRequest
	// status, err = models.GetAllOaRequest(&request, config.LimitQuery, 0, true, paramsRequest)
	// if err != nil {
	// 	// log.Error(err.Error())
	// } else if len(request) > 0 {
	// 	if request[0].Oastatus != nil && *request[0].Oastatus > 0 {
	// 		var lookup models.GenLookup
	// 		status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
	// 		if err != nil {
	// 			// log.Error(err.Error())
	// 		} else {
	// 			if lookup.LkpName != nil && *lookup.LkpName != "" {
	// 				atClaims["oa_status"] = *lookup.LkpName
	// 			}
	// 		}
	// 	}
	// }

	if accountData.RoleKey != nil && *accountData.RoleKey > 0 {
		atClaims["role_key"] = *accountData.RoleKey
		atClaims["user_category_key"] = accountData.UserCategoryKey
		atClaims["user_dept_key"] = accountData.UserDeptKey
		var dept models.ScUserDept
		strDept := strconv.FormatUint(*accountData.UserDeptKey, 10)
		_, err = models.GetScUserDept(&dept, strDept)
		if err == nil {
			atClaims["user_branch"] = dept.BranchKey
		} else {
			atClaims["user_branch"] = 1
		}
		paramsRole := make(map[string]string)
		paramsRole["role_key"] = strconv.FormatUint(*accountData.RoleKey, 10)
		var role []models.ScRole
		_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
		if err != nil {
			// log.Error(err.Error())
		} else if len(role) > 0 {
			if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
				atClaims["role_category_key"] = *role[0].RoleCategoryKey
			}
		}
	}
	atClaims["uuid"] = uuidString
	atClaims["exp"] = time.Now().Add(time.Minute * 50).Unix()
	atClaims["email"] = accountData.UloginEmail
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Secret))
	if err != nil {
		// log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized, err.Error(), "Login failed")
	}

	// sessionKey := base64.StdEncoding.EncodeToString([]byte(uuidString))
	dateLayout := "2006-01-02 15:04:05"
	// expired := date.Add(time.Second * time.Duration(config.SessionExpired)).Format(dateLayout)

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, paramsSession, true)
	paramsSession["session_id"] = uuidString
	paramsSession["login_date"] = time.Now().Format(dateLayout)
	paramsSession["rec_status"] = "1"
	paramsSession["rec_attribute_id3"] = c.Request().UserAgent()
	if err == nil && len(loginSession) > 0 {
		// log.Info("Active session for previous login, overwrite with new session")
		if len(loginSession) > 1 {

		}
		paramsSession["login_session_key"] = strconv.FormatUint(loginSession[0].LoginSessionKey, 10)

		status, err = models.UpdateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error update session")
			return lib.CustomError(status, "Error update session", "Login failed")
		}
	} else {
		status, err = models.CreateScLoginSession(paramsSession)
		if err != nil {
			// log.Error("Error create session")
			return lib.CustomError(status, "Error create session", "Login failed")
		}
	}

	// update ulogin_failed_count = 0 if success login
	paramsUpdate := make(map[string]string)
	uloginkey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUpdate["user_login_key"] = uloginkey
	paramsUpdate["ulogin_failed_count"] = "0"
	_, err = models.UpdateScUserLogin(paramsUpdate)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("erroe update ulogin_failed_count = 0 if success login")
	}

	// log.Info("Success login")

	var data models.ScLoginSessionInfo
	data.SessionID = token
	if accountData.UloginMustChangepwd == uint8(1) {
		// data.MustChangePassword = true
	} else {
		// data.MustChangePassword = false
	}
	// log.Info(data)

	//LOG LOGIN
	saveLogLogin(c, paramsSession)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	// log.Info(response)
	return c.JSON(http.StatusOK, response)
}

func ChangeForgotPassword(c echo.Context) error {

	var err error

	token := c.FormValue("token")
	if token == "" {
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Token tidak ditemukan")
	}
	params := make(map[string]string)
	params["string_token"] = token
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(http.StatusBadRequest, "Error get email", "Gagal mendapatkan data email")
	}
	if len(userLogin) < 1 {
		// log.Error("No matching token " + token)
		return lib.CustomError(http.StatusBadRequest, "Token not found", "Token tidak ditemukan")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)
	// Check parameters
	newPassword1 := c.FormValue("new_password1")
	if newPassword1 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPassword2 := c.FormValue("new_password2")
	if newPassword2 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	if newPassword1 != newPassword2 {
		// log.Error("Password doesnt match")
		return lib.CustomError(http.StatusBadRequest, "Password doesnt match", "Password doesnt match")
	}
	// Validate password
	length, number, upper, special := verifyPassword(newPassword1)
	if length == false || number == false || upper == false || special == false {
		// log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(newPassword1))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])

	dateLayout := "2006-01-02 15:04:05"
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["ulogin_password"] = encryptedPassword
	params["last_password_changed"] = time.Now().Format(dateLayout)
	params["ulogin_must_changepwd"] = "0"
	params["string_token"] = ""
	params["ulogin_failed_count"] = "0"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(accountData.UserLoginKey, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//insert table sc_user_message
	strKey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	strUserLoginKey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	paramsUserMessage["umessage_subject"] = "Perubahan Kata Sandi Berhasil"
	paramsUserMessage["umessage_body"] = "Kata sandi kamu berhasil berubah. Apabila kamu tidak merasa melakukan perubahan kata sandi, mohon segera menghubungi customer service kami."
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strKey

	_, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	}

	//kirim email
	t := template.New("index-sukses-ubah-password.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-sukses-ubah-password.html")
	if err != nil {
		// log.Println(err)
	}

	var customer models.MsCustomer

	fullnameuser := accountData.UloginFullName

	if accountData.CustomerKey != nil {
		strCustomerKey := strconv.FormatUint(*accountData.CustomerKey, 10)
		_, err = models.GetMsCustomer(&customer, strCustomerKey)
		if err == nil {
			fullnameuser = customer.FullName
		}
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl,
		struct {
			Name    string
			FileUrl string
		}{
			Name:    fullnameuser,
			FileUrl: config.ImageUrl + "/images/mail"}); err != nil {
		// log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", accountData.UloginEmail)
	mailer.SetHeader("Subject", "[MotionFunds] Berhasil Merubah Kata Sandi")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Error(err)
	} else {
		// log.Info("Email sent")
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
	// 	// log.Error(err)
	// } else {
	// 	// log.Info("Email sent")
	// }

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func CreatePin(c echo.Context) error {

	var err error
	// var status int
	// Check parameters
	pin1 := c.FormValue("pin1")
	if pin1 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	pin2 := c.FormValue("pin2")
	if pin2 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	// params["ulogin_key"] = lib.Profile.UserID
	// var userLogin []models.ScUserLogin
	// status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	// if err != nil {
	// 	// log.Error("Error get email")
	// 	return lib.CustomError(status, "Error get email", "Error get email")
	// }
	// if len(userLogin) < 1 {
	// 	// log.Error("Email not registered")
	// 	return lib.CustomError(http.StatusUnauthorized, "Email not registered", "Email not registered")
	// }

	// accountData := userLogin[0]
	// // log.Info(accountData)

	if pin1 != pin2 {
		// log.Error("Pin doesnt match")
		return lib.CustomError(http.StatusBadRequest, "Pin doesnt match", "Pin doesnt match")
	}

	// Encrypt pin
	// encryptedPasswordByte := sha256.Sum256([]byte(pin1))
	// encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	encryptedPassword := base64.StdEncoding.EncodeToString([]byte(pin1))

	dateLayout := "2006-01-02 15:04:05"
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["ulogin_pin"] = encryptedPassword
	params["must_change_pin"] = "0"
	params["last_changed_pin"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//insert table sc_user_message
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	strUserLoginKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Pembuatan pin Berhasil"
	body := "Pin Berhasil dibuat. Apabila kamu tidak merasa melakukan pembuatan pin, mohon segera menghubungi customer service kami."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strKey

	_, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ChangePin(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	recentPin := c.FormValue("recent_pin")
	if recentPin == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPin1 := c.FormValue("new_pin1")
	if newPin1 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPin2 := c.FormValue("new_pin2")
	if newPin2 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(status, "Error get email", "Error get email")
	}
	if len(userLogin) < 1 {
		// log.Error("Email not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email not registered", "Email not registered")
	}

	accountData := userLogin[0]
	// log.Info(accountData)

	// Check valid password
	// encryptedPasswordByte := sha256.Sum256([]byte(recentPin))
	// encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	encryptedPassword := base64.StdEncoding.EncodeToString([]byte(recentPin))
	if accountData.UloginPin != nil && encryptedPassword != *accountData.UloginPin {
		// log.Error("Wrong password")
		return lib.CustomError(http.StatusUnauthorized, "Wrong pin", "Wrong pin")
	}

	if newPin1 != newPin2 {
		// log.Error("Password doesnt match")
		return lib.CustomError(http.StatusBadRequest, "Pin doesnt match", "Pin doesnt match")
	}

	// Encrypt password
	// encryptedPasswordByte = sha256.Sum256([]byte(newPin1))
	// encryptedPassword = hex.EncodeToString(encryptedPasswordByte[:])
	encryptedPassword = base64.StdEncoding.EncodeToString([]byte(newPin1))

	dateLayout := "2006-01-02 15:04:05"
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["ulogin_pin"] = encryptedPassword
	params["last_changed_pin"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["must_change_pin"] = "0"

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//insert table sc_user_message
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	strUserLoginKey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Perubahan Pin Berhasil"
	body := "Pin kamu berhasil berubah. Apabila kamu tidak merasa melakukan perubahan pin, mohon segera menghubungi customer service kami."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strKey

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ForgotPin(c echo.Context) error {
	var err error
	var status int

	// Check valid email
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get user")
		return lib.CustomError(status, "Error get user", "Error get user")
	}
	if len(userLogin) < 1 {
		// log.Error("Email not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email not registered", "Email not registered")
	}

	accountData := userLogin[0]
	// log.Info("Found account with email " + accountData.UloginEmail)

	pin := ""
	length := 6
	for i := 0; i < length; i++ {
		pin += strconv.FormatUint(uint64(rand.Intn(9)), 10)
	}

	encryptedPassword := base64.StdEncoding.EncodeToString([]byte(pin))

	params = make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["ulogin_pin"] = encryptedPassword
	params["must_change_pin"] = "1"
	dateLayout := "2006-01-02 15:04:05"
	params["last_changed_pin"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	// // log.Println(rand.Intn(9))

	// Send email
	t := template.New("index-forget-pin.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-forget-pin.html")
	if err != nil {
		// log.Println(err)
	}

	name := ""

	var fullName models.FullNameData
	_, err = models.GetFullName(&fullName, strconv.FormatUint(lib.Profile.UserID, 10))
	if err != nil {
		name = accountData.UloginFullName
	} else {
		name = fullName.FullName
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		PIN     string
		FileUrl string
		Name    string
	}{PIN: pin, FileUrl: config.ImageUrl + "/images/mail", Name: name}); err != nil {
		// log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	// mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MotionFunds] Lupa PIN")
	mailer.SetBody("text/html", result)

	err = lib.SendEmail(mailer)
	if err != nil {
		// log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	} else {
		// log.Info("Email sent")
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
	// 	// log.Error(err)
	// 	return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	// }
	// // log.Info("Email sent")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func ChangeForgotPin(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	token := c.FormValue("token")
	if token == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPin1 := c.FormValue("new_pin1")
	if newPin1 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}
	newPin2 := c.FormValue("new_pin2")
	if newPin2 == "" {
		// log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter", "Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		// log.Error("Error get email")
		return lib.CustomError(status, "Error get email", "Error get email")
	}
	if len(userLogin) < 1 {
		// log.Error("Email not registered")
		return lib.CustomError(http.StatusUnauthorized, "Email not registered", "Email not registered")
	}

	accountData := userLogin[0]
	// log.Info(accountData)

	if newPin1 != newPin2 {
		// log.Error("PIN doesnt match")
		return lib.CustomError(http.StatusBadRequest, "PIN doesnt match", "PIN doesnt match")
	}

	// Encrypt password
	// encryptedPasswordByte := sha256.Sum256([]byte(newPin1))
	// encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	encryptedPassword := base64.StdEncoding.EncodeToString([]byte(newPin1))

	dateLayout := "2006-01-02 15:04:05"
	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["ulogin_pin"] = encryptedPassword
	params["must_change_pin"] = "0"
	params["last_changed_pin"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//insert table sc_user_message
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	strUserLoginKey := strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	subject := "Perubahan Pin Berhasil"
	body := "Pin kamu berhasil berubah. Apabila kamu tidak merasa melakukan perubahan pin, mohon segera menghubungi customer service kami."
	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strKey

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		// log.Error(err.Error())
		// log.Error("Error create user message")
	}
	lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body, "TRANSACTION")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func saveLogLogin(c echo.Context, session map[string]string) {
	var err error
	params := make(map[string]string)
	params["session_id"] = session["session_id"]
	params["login_date"] = session["login_date"]
	params["user_login_key"] = session["user_login_key"]
	params["workstation_ipaddress"] = c.RealIP()
	params["client_agent"] = session["rec_attribute_id3"]
	params["rec_status"] = "1"
	id, err := machineid.ID()
	if err == nil {
		params["terminal_name"] = id
	}
	ua := ua.Parse(c.Request().UserAgent())
	params["workstation_name"] = ua.OS
	params["device_model"] = ua.Name
	// _, err = models.CreateScLoginLog(params)
	// if err != nil {
	// 	// log.Error("Error create log loginsession")
	// }
}

func AccountDeletionIndividu(c echo.Context) error {
	var err error

	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	paramsDeactive := make(map[string]string)

	deactivatereason := c.FormValue("deactivate_reason")
	if deactivatereason != "" {
		n, err := strconv.ParseUint(deactivatereason, 10, 64)
		if err != nil || n == 0 {
			// log.Error("Wrong input for parameter: deactivate_reason")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: deactivate_reason", "Wrong input for parameter: deactivate_reason")
		}

		if len(deactivatereason) > 11 {
			// log.Error("Wrong input for parameter: deactivate_reason too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: deactivate_reason too long, max 11 character", "Missing required parameter: deactivate_reason too long, max 11 character")
		}
		paramsDeactive["deactivate_reason"] = deactivatereason
	} else {
		// log.Error("Missing required parameter: deactivate_reason")
		return lib.CustomError(http.StatusBadRequest, "deactivate_reason can not be blank", "deactivate_reason can not be blank")
	}

	deactivatenotes := c.FormValue("deactivate_notes")
	if deactivatenotes != "" {
		if len(deactivatenotes) > 250 {
			// log.Error("Wrong input for parameter: deactivate_notes too long")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: deactivate_notes too long, max 250 character", "Missing required parameter: deactivate_notes too long, max 250 character")
		}
		paramsDeactive["deactivate_notes"] = deactivatenotes
	}

	dateLayout := "2006-01-02 15:04:05"

	paramsDeactive["deactivate_status"] = "449"
	paramsDeactive["deactivate_date"] = time.Now().Format(dateLayout)
	paramsDeactive["user_login_key"] = strIDUserLogin
	paramsDeactive["rec_status"] = "1"
	paramsDeactive["rec_order"] = "0"
	paramsDeactive["rec_created_date"] = time.Now().Format(dateLayout)
	paramsDeactive["rec_created_by"] = strIDUserLogin
	_, err = models.CreateScLoginDeactivation(paramsDeactive)
	if err != nil {
		// log.Error("Error create sc_login_deactivation")
	}

	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strIDUserLogin
	paramsSession["logout_date"] = time.Now().Format(dateLayout)
	paramsSession["login_session_key"] = ""

	_, err = models.UpdateScLoginSession(paramsSession)
	if err != nil {
		// log.Error("Error update session in logout")
	}

	params := make(map[string]string)
	params["user_login_key"] = strIDUserLogin
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strIDUserLogin
	params["rec_attribute_id1"] = "ACCOUNT_DELETION"

	_, err = models.UpdateScUserLogin(params)
	if err != nil {
		// log.Error("Error account deletion user data")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
