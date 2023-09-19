package models

import (
	"fmt"
	"log"
	"math"
	"mf-bo-api/db"
	"regexp"
	"strconv"
)

// type MsCustomerDetail struct {
// 	CustomerKey  uint64  `db:"customer_key"      json:"customer_key"`
// 	Nationality  uint64  `db:"nationality"       json:"nationality"`
// 	Gender       uint64  `db:"gender"            json:"gender"`
// 	IDType       uint64  `db:"id_type"           json:"id_type"`
// 	IDNumber     *string `db:"id_number"         json:"id_number"`
// 	IDHolderName *string `db:"id_holder_name"    json:"id_holder_name"`
// 	FlagEmployee uint8   `db:"flag_employee"     json:"flag_employee"`
// 	FlagGroup    uint8   `db:"flag_group"        json:"flag_group"`
// }

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

type CustomerIndividuListResponse struct {
	OaRequestKey     *uint64 `db:"oa_request_key" json:"oa_request_key"`
	CustomerKey      *uint64 `db:"customer_key" json:"customer_key"`
	BranchName       *string `db:"branch_name" json:"branch_name"`
	CIF              *string `db:"cif" json:"cif"`
	SID              *string `db:"sid" json:"sid"`
	Email            *string `db:"email" json:"email"`
	IdCardNo         *string `db:"idcard_no" json:"idcard_no"`
	FullName         *string `db:"full_name" json:"full_name"`
	MotherMaidenName *string `db:"mother_maiden_name" json:"mother_maiden_name"`
	DateBirth        *string `db:"date_birth" json:"date_birth"`
	PhoneMobile      *string `db:"phone_mobile" json:"phone_mobile"`
	CIFSuspendFlag   *bool   `db:"cif_suspend_flag" json:"cif_suspend_flag"`
}

type CustomerIndividuListDb struct {
	OaRequestKey     *uint64 `db:"oa_request_key" json:"oa_request_key"`
	CustomerKey      *uint64 `db:"customer_key" json:"customer_key"`
	BranchName       *string `db:"branch_name" json:"branch_name"`
	CIF              *string `db:"cif" json:"cif"`
	SID              *string `db:"sid" json:"sid"`
	Email            *string `db:"email" json:"email"`
	IdCardNo         *string `db:"idcard_no" json:"idcard_no"`
	FullName         *string `db:"full_name" json:"full_name"`
	MotherMaidenName *string `db:"mother_maiden_name" json:"mother_maiden_name"`
	DateBirth        *string `db:"date_birth" json:"date_birth"`
	PhoneMobile      *string `db:"phone_mobile" json:"phone_mobile"`
	CIFSuspendFlag   *uint64 `db:"cif_suspend_flag" json:"cif_suspend_flag"`
}

func GetCustomerListWithCondition(params map[string]string, limit uint64, offset uint64) ([]CustomerIndividuListDb, int) {
	query := `SELECT MAX(t1.oa_request_key) AS oa_request_key, t1.customer_key, t4.branch_name, t2.unit_holder_idno AS cif, t2.sid_no AS sid,
	t3.email_address AS email, t3.phone_mobile, t3.idcard_no, t3.full_name, t3.mother_maiden_name, t3.date_birth,
	t2.cif_suspend_flag
	FROM oa_request t1
	INNER JOIN ms_customer t2 ON t2.customer_key = t1.customer_key AND t2.rec_status = 1
	INNER JOIN oa_personal_data t3 ON t3.oa_request_key = t1.oa_request_key
	LEFT JOIN ms_branch t4 ON t4.branch_key = t1.branch_key
	WHERE t1.rec_status = 1  `

	queryCountPage := `SELECT count(*)
	FROM oa_request t1
	INNER JOIN ms_customer t2 ON t2.customer_key = t1.customer_key AND t2.rec_status = 1
	INNER JOIN oa_personal_data t3 ON t3.oa_request_key = t1.oa_request_key
	LEFT JOIN ms_branch t4 ON t4.branch_key = t1.branch_key
	WHERE t1.rec_status = 1 `

	if valueMap, ok := params["branch_key"]; ok {
		query += `AND t4.branch_key = ` + valueMap
		queryCountPage += `AND t4.branch_key = ` + valueMap
	}

	if valueMap, ok := params["cif"]; ok {
		query += `AND t2.unit_holder_idno = ` + valueMap
		queryCountPage += `AND t2.unit_holder_idno = ` + valueMap
	}

	if valueMap, ok := params["idcard_no"]; ok {
		query += `AND t3.idcard_no = ` + valueMap
		queryCountPage += `AND t3.idcard_no = ` + valueMap
	}

	if valueMap, ok := params["full_name"]; ok {
		query += `AND t3.full_name LIKE '%` + valueMap + `%'`
		queryCountPage += `AND t3.idcard_no = ` + valueMap
	}

	if valueMap, ok := params["mother_maiden_name"]; ok {
		query += `AND t3.mother_maiden_name LIKE '%` + valueMap + `%'`
		queryCountPage += `AND t3.idcard_no = ` + valueMap
	}

	if valueMap, ok := params["date_birth"]; ok {
		query += `AND t3.date_birth = ` + valueMap
		queryCountPage += `AND t3.idcard_no = ` + valueMap
	}

	query += ` GROUP BY t2.customer_key ORDER BY t2.customer_key`
	queryCountPage += ` GROUP BY t2.customer_key ORDER BY t2.customer_key`

	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// EXECUTE DATA
	var result []CustomerIndividuListDb
	log.Println(query)
	err := db.Db.Select(&result, query)
	if err != nil {
		log.Println(err.Error())
	}

	// EXECUTE PAGING
	var pagination int
	var count uint64
	err = db.Db.Get(&count, queryCountPage)
	if err != nil {
		log.Println(err.Error())
	}
	if limit > 0 {
		if count < limit {
			pagination = 1
		} else {
			calc := math.Ceil(float64(count) / float64(limit))
			pagination = int(calc)
		}
	}

	return result, pagination
}
