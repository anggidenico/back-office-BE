package models

import (
	"fmt"
	"log"
	"mf-bo-api/db"
	"regexp"
	"strconv"
)

type MsCustomerDetail struct {
	CustomerKey  uint64  `db:"customer_key"      json:"customer_key"`
	Nationality  uint64  `db:"nationality"       json:"nationality"`
	Gender       uint64  `db:"gender"            json:"gender"`
	IDType       uint64  `db:"id_type"           json:"id_type"`
	IDNumber     *string `db:"id_number"         json:"id_number"`
	IDHolderName *string `db:"id_holder_name"    json:"id_holder_name"`
	FlagEmployee uint8   `db:"flag_employee"     json:"flag_employee"`
	FlagGroup    uint8   `db:"flag_group"        json:"flag_group"`
}

func NewClientCode() string {
	var NewClientCode string

	query := `SELECT client_code FROM ms_customer ORDER BY client_code DESC LIMIT 1`

	var maxClientCode *string
	err := db.Db.Get(&maxClientCode, query)
	if err != nil {
		log.Println(err.Error())
	}

	if maxClientCode != nil {
		stringz := *maxClientCode
		firstLetter := stringz[0:1]
		rn := []rune(firstLetter)
		asciiCode := rn[0]
		var intVar int64
		var incr int64

		if asciiCode >= 65 && asciiCode <= 90 {
			re := regexp.MustCompile("[0-9]+")
			str1 := re.FindAllString(*maxClientCode, -1)
			intVar, _ = strconv.ParseInt(str1[0], 10, 0)
		} else {
			intVar, _ = strconv.ParseInt(*maxClientCode, 10, 0)
		}

		if intVar == 999999 {
			// fmt.Println("lewat 1")

			if asciiCode >= 65 && asciiCode <= 90 {
				// fmt.Println("lewat 2")

				asciiNum := asciiCode + 1
				firstLetter = string(asciiNum)
				*maxClientCode = firstLetter + "0000001"
			}

			intVar, _ = strconv.ParseInt(*maxClientCode, 10, 0)
			incr = intVar + 1
			NewClientCode = firstLetter + fmt.Sprintf("%06d", incr)

		} else {
			// fmt.Println("lewat 3")

			incr = intVar + 1
			if asciiCode >= 65 && asciiCode <= 90 {
				NewClientCode = firstLetter + fmt.Sprintf("%06d", incr)
			} else {
				NewClientCode = fmt.Sprintf("%06d", incr)
			}

		}
	}

	return NewClientCode
}
