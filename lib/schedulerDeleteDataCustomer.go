package lib

import (
	"fmt"
	"mf-bo-api/models"
	"time"
)

func DeleteUserBelumOASamaSekali() {
	fmt.Println("START CRON DELETE USER BELUM OA SAMA SEKALI")
	var err error

	_, err = models.UpdateDeleteUserNotOA("30")
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err.Error())
	} else {
		fmt.Println("SUKSES")
	}

	fmt.Println("======END CRON DELETE USER BELUM OA SAMA SEKALI============")
}

func DeleteDataOAAndUserBelumSelesaiOA(oaStatus []string, oaType []string, isDeleteUser bool) {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON DELETE OA DATA DAN ATAU USER")
	var err error

	//cek user
	var userIds []string
	var oaIds []string
	var oaDelete []models.OaDataDelete
	_, err = models.GetOaDataDelete(&oaDelete, oaStatus, oaType)
	if err == nil {
		if len(oaDelete) > 0 {
			for _, od := range oaDelete {
				userIds = append(userIds, od.UserLoginKey)
				oaIds = append(oaIds, od.OaRequestKey)
			}
		}
	}

	//delete update user login
	if isDeleteUser {
		if len(userIds) > 0 {
			params := make(map[string]string)
			params["rec_status"] = "0"
			params["rec_modified_by"] = "CRON DELETE USER"
			params["rec_modified_date"] = time.Now().Format(dateLayout)
			_, err = models.UpdateScUserLoginByKeyIn(params, userIds, "user_login_key")
			if err != nil {
				// log.Error("ERROR DELETE USER : " + err.Error())
			} else {
				fmt.Println("DELETE USER DONE. Jumlah Data : ")
				fmt.Println(len(userIds))
			}
		} else {
			fmt.Println("NO DATA USER_LOGIN_KEY")
		}
	}

	//delete update oa_request & oa_personal_data
	if len(oaIds) > 0 {
		paramsOa := make(map[string]string)
		paramsOa["rec_status"] = "0"
		paramsOa["rec_modified_by"] = "CRON DELETE USER"
		paramsOa["rec_modified_date"] = time.Now().Format(dateLayout)
		_, err = models.UpdateOaRequestByKeyIn(paramsOa, oaIds, "oa_request_key")
		if err != nil {
			// log.Error("Error update delete oa request")
		}

		_, err = models.UpdateOaPersonalDataByKeyIn(paramsOa, oaIds, "oa_request_key")
		if err != nil {
			// log.Error("Error update delete oa personal data")
		}
	} else {
		fmt.Println("NO DATA OA_REQUEST IDS")
	}

	fmt.Println("======END CRON DELETE OA DATA DAN ATAU USER============")
}
