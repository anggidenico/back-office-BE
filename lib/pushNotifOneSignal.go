package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mf-bo-api/config"
	"mf-bo-api/models"
	"net/http"
	"strings"

	"github.com/tbalthazar/onesignal-go"
)

func CreateNotifCustomerFromApp(heading string, content string, category string) {
	if Profile.TokenNotif == nil {
		// log.Println("token kosong")
	} else {
		// log.Println("token : " + *Profile.TokenNotif)
		playerID := &Profile.TokenNotif
		CreateNotificationHelper(**playerID, heading, content, category)
	}
}

func CreateNotifCustomerFromAdminByCustomerId(customerId string, heading string, content string, category string) {
	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByCustomerKey(&userData, customerId)
	if err == nil {
		if userData.TokenNotif == nil {
			// log.Println("token kosong")
		} else {
			// log.Println("token : " + *userData.TokenNotif)
			playerID := &userData.TokenNotif
			CreateNotificationHelper(**playerID, heading, content, category)
		}
	}
}

func CreateNotifCustomerFromAdminByUserLoginKey(userLoginKey string, heading string, content string, category string) {
	var userData models.ScUserLogin
	_, err := models.GetScUserKey(&userData, userLoginKey)
	if err == nil {
		if userData.TokenNotif == nil {
			// log.Println("token kosong")
		} else {
			// log.Println("token : " + *userData.TokenNotif)
			playerID := &userData.TokenNotif
			CreateNotificationHelper(**playerID, heading, content, category)
		}
	}
}

func CreateNotificationHelper(playerID string, heading string, content string, category string) *onesignal.NotificationCreateResponse {
	// log.Println("playerID : " + playerID)
	// log.Println("Heading : " + heading)
	// log.Println("Content : " + content)
	client := onesignal.NewClient(nil)
	client.AppKey = config.OneSignalAppKey

	DataNotif := make(map[string]interface{})
	DataNotif["category"] = category
	notificationReq := &onesignal.NotificationRequest{
		AppID:            config.OneSignalAppID,
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		Data:             DataNotif,
		SmallIcon:        "ic_stat_onesignal_default",
		LargeIcon:        config.ImageUrl + "/images/mail/icon_mncduit.png",
		IncludePlayerIDs: []string{playerID},
	}
	createRes, _, err := client.Notifications.Create(notificationReq)
	if err != nil {
		// log.Println("OneSignal Message")
		fmt.Println(err)
	} else {
		return createRes
	}
	return createRes
}

func BlastAllNotificationHelper(playerIDs []string, heading string, content string, data map[string]interface{}) *onesignal.NotificationCreateResponse {
	// log.Println("Heading : " + heading)
	// log.Println("Content : " + content)
	client := onesignal.NewClient(nil)
	client.AppKey = config.OneSignalAppKey

	notificationReq := &onesignal.NotificationRequest{
		AppID:            config.OneSignalAppID,
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		Data:             data,
		SmallIcon:        "ic_stat_onesignal_default",
		LargeIcon:        config.ImageUrl + "/images/mail/icon_mncduit.png",
		IncludePlayerIDs: playerIDs,
	}
	createRes, _, err := client.Notifications.Create(notificationReq)
	if err != nil {
		// log.Println("OneSignal Message")
		fmt.Println(err)
	} else {
		return createRes
	}
	return createRes
}

func CreateNotifOneSignal(params map[string]string) error { // INI YANG DIPAKE UNTUK HIT ONESIGNAL
	curlParam := make(map[string]interface{})

	var TokenNotif string
	if valueMap, ok := params["token_notif"]; ok {
		TokenNotif = valueMap
	} else {
		return fmt.Errorf("missing token_notif")
	}

	var PhoneNumber string
	if valueMap, ok := params["phone_number"]; ok {
		PhoneNumber = valueMap
	} else {
		return fmt.Errorf("missing phone_number")
	}

	var Description string
	if valueMap, ok := params["description"]; ok {
		Description = valueMap
	} else {
		return fmt.Errorf("missing description")
	}

	if TokenNotif != "" {
		include_player_id := [...]string{TokenNotif}
		curlParam["include_player_ids"] = include_player_id
	}
	external_user_id := [...]string{PhoneNumber}
	included_segments := [...]string{"Active User"}
	contentsMap := make(map[string]string)
	contentsMap["en"] = Description

	curlParam["app_id"] = config.OneSignalAppID
	curlParam["include_external_user_ids"] = external_user_id
	curlParam["chanel_for_external_user_ids"] = "push"
	curlParam["included_segments"] = included_segments
	curlParam["contents"] = contentsMap

	jsonString, err := json.Marshal(curlParam)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("HIT ONESIGNAL WITH PARAMETER:", string(jsonString))

	payload := strings.NewReader(string(jsonString))
	url := "https://onesignal.com/api/v1/notifications"
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+config.OneSignalAppKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
