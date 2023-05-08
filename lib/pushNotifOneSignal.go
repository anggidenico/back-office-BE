package lib

import (
	"fmt"
	"log"
	"mf-bo-api/config"
	"mf-bo-api/models"

	"github.com/tbalthazar/onesignal-go"
)

func CreateNotifCustomerFromApp(heading string, content string, category string) {
	if Profile.TokenNotif == nil {
		log.Println("token kosong")
	} else {
		log.Println("token : " + *Profile.TokenNotif)
		playerID := &Profile.TokenNotif
		CreateNotificationHelper(**playerID, heading, content, category)
	}
}

func CreateNotifCustomerFromAdminByCustomerId(customerId string, heading string, content string, category string) {
	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByCustomerKey(&userData, customerId)
	if err == nil {
		if userData.TokenNotif == nil {
			log.Println("token kosong")
		} else {
			log.Println("token : " + *userData.TokenNotif)
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
			log.Println("token kosong")
		} else {
			log.Println("token : " + *userData.TokenNotif)
			playerID := &userData.TokenNotif
			CreateNotificationHelper(**playerID, heading, content, category)
		}
	}
}

func CreateNotificationHelper(playerID string, heading string, content string, category string) *onesignal.NotificationCreateResponse {
	log.Println("playerID : " + playerID)
	log.Println("Heading : " + heading)
	log.Println("Content : " + content)
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
		log.Println("OneSignal Message")
		fmt.Println(err)
	} else {
		return createRes
	}
	return createRes
}

func BlastAllNotificationHelper(playerIDs []string, heading string, content string, data map[string]interface{}) *onesignal.NotificationCreateResponse {
	log.Println("Heading : " + heading)
	log.Println("Content : " + content)
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
		log.Println("OneSignal Message")
		fmt.Println(err)
	} else {
		return createRes
	}
	return createRes
}
