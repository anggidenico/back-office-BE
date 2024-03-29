package lib

import (
	"fmt"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	ua "github.com/mileusna/useragent"
)

type CProfile struct {
	UserID          uint64  `json:"user_id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	PhoneNumber     string  `json:"phone_number"`
	RoleKey         uint64  `json:"role_key"`
	RoleCategoryKey uint64  `json:"role_category_key"`
	RecImage1       string  `json:"rec_image1"`
	CustomerKey     *uint64 `json:"customer_key"`
	UserCategoryKey uint64  `json:"user_category_key"`
	RolePrivileges  *uint64 `json:"role_privileges"`
	BranchKey       *uint64 `json:"branch_key"`
	TokenNotif      *string `json:"token_notif"`
}

var Profile CProfile
var UserIDStr string

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tokenString string
		request := c.Request()
		authorization := request.Header["Authorization"]
		if authorization != nil {
			if strings.HasPrefix(authorization[0], "Bearer ") == true {
				tokenString = authorization[0][7:]
				// log.Info(tokenString)
			}
		}
		token, err := VerifyToken(tokenString)
		if err != nil {
			// log.Error(err.Error())
			return CustomError(http.StatusForbidden, err.Error(), "Authentication failed : cannot verified user")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		// log.Info(claims)
		if ok && token.Valid {
			accessUuid, ok := claims["uuid"].(string)
			if !ok {
				// log.Error("Cannot get uuid")
				return CustomError(http.StatusForbidden, "Cannot get uuid", "Authentication failed : cannot verified user")
			}
			params := make(map[string]string)
			params["session_id"] = accessUuid
			var loginSession []models.ScLoginSession
			_, err := models.GetAllScLoginSession(&loginSession, config.LimitQuery, 0, params, true)
			if err != nil {
				// log.Error("Error get email")
				return CustomError(http.StatusForbidden, "Forbidden", "you have to login first")
			}
			if len(loginSession) < 1 {
				// log.Error("No matching token " + tokenString)
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}

			paramsUser := make(map[string]string)
			paramsUser["user_login_key"] = strconv.FormatUint(loginSession[0].UserLoginKey, 10)
			paramsUser["rec_status"] = "1"
			var userLogin []models.ScUserLogin
			_, err = models.GetAllScUserLogin(&userLogin, config.LimitQuery, 0, paramsUser, true)
			if err != nil {
				// log.Error("Error get email")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}
			if len(userLogin) < 1 {
				// log.Error("No user login")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}

			user := userLogin[0]

			//check allowed endpoint
			isCheckAllowedEndpoint := false //kondisi check / tidak check
			if isCheckAllowedEndpoint {
				strRoleKey := ""
				if user.RoleKey != nil && *user.RoleKey > 0 {
					strRoleKey = strconv.FormatUint(*user.RoleKey, 10)
				}
				if strRoleKey != "1" { //selain user mobile
					var countData models.CountData
					_, err = models.CheckAllowedEndpoint(&countData, strRoleKey, c.Path())
					if err != nil {
						// log.Error("Error Check Allowed Endpoint")
						return CustomError(http.StatusForbidden, "Forbidden", "Error Check Allowed Endpoint")
					}
					if int(countData.CountData) < 1 {
						// log.Error("Action Not Allowed")
						return CustomError(http.StatusForbidden, "Forbidden", "Action Not Allowed")
					}
				}
			}

			if user.RoleKey != nil && *user.RoleKey > 0 {
				Profile.RoleKey = *user.RoleKey
				paramsRole := make(map[string]string)
				paramsRole["role_key"] = strconv.FormatUint(*user.RoleKey, 10)
				var role []models.ScRole
				_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
				if err != nil {
					// log.Error(err.Error())
				} else if len(role) > 0 {
					if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
						Profile.RoleCategoryKey = *role[0].RoleCategoryKey
					}
				}

				if user.UserDeptKey != nil {
					var dept models.ScUserDept
					strDept := strconv.FormatUint(*user.UserDeptKey, 10)
					_, err = models.GetScUserDept(&dept, strDept)
					if err != nil {
						// log.Error(err.Error())
					} else {
						Profile.RolePrivileges = dept.RolePrivileges
						Profile.BranchKey = dept.BranchKey
					}
				}
			}

			Profile.UserID = user.UserLoginKey
			Profile.Name = user.UloginFullName
			Profile.Email = user.UloginEmail
			Profile.PhoneNumber = *user.UloginMobileno
			Profile.CustomerKey = user.CustomerKey
			Profile.UserCategoryKey = user.UserCategoryKey
			Profile.TokenNotif = user.TokenNotif
			if user.RecImage1 != nil && *user.RecImage1 != "" {
				Profile.RecImage1 = config.ImageUrl + "/images/user/" + strconv.FormatUint(user.UserLoginKey, 10) + "/profile/" + *user.RecImage1
			} else {
				Profile.RecImage1 = config.BaseUrl + "/user/default.png"
			}
			UserIDStr = strconv.Itoa(int(Profile.UserID))
			//log endpoint
			dateLayout := "2006-01-02 15:04:05"
			paramLog := make(map[string]string)
			paramLog["path"] = c.Path()
			paramLog["url"] = c.Request().URL.String()
			paramLog["user_login_key"] = strconv.FormatUint(Profile.UserID, 10)
			paramLog["created_date"] = time.Now().Format(dateLayout)
			paramLog["created_by"] = strconv.FormatUint(Profile.UserID, 10)
			id, err := machineid.ID()
			if err == nil {
				paramLog["device_id"] = id
			}
			ua := ua.Parse(c.Request().UserAgent())
			paramLog["device_ip"] = c.RealIP()
			paramLog["device_os"] = ua.OS
			paramLog["device_name"] = ua.Name
			paramLog["notes"] = c.Request().UserAgent()

			// // log.Println("-----------------------")
			// json_map := make(map[string]interface{})
			// // log.Println(c.Request())
			// // log.Println(json.NewDecoder(c.Request().Body))
			// err = json.NewDecoder(c.Request().Body).Decode(&json_map)
			// if err != nil {
			// 	// log.Println("HAHAHAHAHAHAHA")
			// 	// log.Println(err)
			// 	return err
			// } else {
			// 	// log.Println(json_map)
			// 	//json_map has the JSON Payload decoded into a map
			// 	// cb_type := json_map["type"]
			// 	// challenge := json_map["challenge"]
			// }

			// // log.Println(request)
			// // log.Println(request.PostForm.Encode())
			// // log.Println(c.Echo().Binder)
			// // log.Println("-----------------------")
			// paramLog["data"] = c.Request().Body.Close().Error()

			// // log.Println("BODY: " + c.Request().Header)
			_, err = models.CreateEndpointAuditTrail(paramLog)
			if err != nil {
				// log.Error("Failed Log Audit Trail: " + err.Error())
			}

		} else {
			// log.Error("Invalid token")
			return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
		}
		return next(c)
	}
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// log.Error("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
