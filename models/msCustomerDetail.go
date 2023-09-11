package models

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

	// query := `SELECT client_code FROM ms_customer ORDER BY client_code DESC LIMIT 1`

	// var maxClientCode *string
	// err := db.Db.Get(&maxClientCode, query)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// if maxClientCode != nil {
	// 	intVar, _ := strconv.ParseInt(*maxClientCode, 10, 0)
	// 	// fmt.Println("integernya: ", intVar)
	// 	incr := intVar + 1
	// 	// fmt.Println("tambah 1 jadi: ", incr)
	// 	digit6 := fmt.Sprintf("%06d", incr)
	// 	// fmt.Println("dibuat 6 digit jadi: ", digit6)
	// 	NewClientCode = digit6

	// 	country := "London"
	// 	firstLetter := country[0:1]
	// 	asciiCode := int(firstLetter)
	// }

	return NewClientCode
}
