package models

import (
	"database/sql"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OaPostalAddress struct {
	PostalAddressKey  uint64  `db:"postal_address_key"         json:"postal_address_key"`
	AddressType       *string `db:"address_type"               json:"address_type"`
	ProvinceKey       *uint64 `db:"province_key" json:"province"`
	KabupatenKey      *uint64 `db:"kabupaten_key"              json:"kabupaten_key"`
	KecamatanKey      *uint64 `db:"kecamatan_key"              json:"kecamatan_key"`
	AddressLine1      *string `db:"address_line1"              json:"address_line1"`
	AddressLine2      *string `db:"address_line2"              json:"address_line2"`
	AddressLine3      *string `db:"address_line3"              json:"address_line3"`
	PostalCode        *string `db:"postal_code"                json:"postal_code"`
	GeolocName        *string `db:"geoloc_name"                json:"geoloc_name"`
	GeolocLongitude   *string `db:"geoloc_longitude"           json:"geoloc_longitude"`
	GeolocLatitude    *string `db:"geoloc_latitude"            json:"geoloc_latitude"`
	RecOrder          *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

type AddressDetail struct {
	PostalAddressKey uint64  `db:"postal_address_key"       json:"postal_address_key"`
	AddressType      *uint64 `db:"address_type"             json:"address_type"`
	AddressTypeName  *string `db:"address_type_name"        json:"address_type_name"`
	ProvinsiKey      *uint64 `db:"provinsi_key"             json:"provinsi_key"`
	ProvinsiName     *string `db:"provinsi_name"            json:"provinsi_name"`
	KabupatenKey     *uint64 `db:"kabupaten_key"            json:"kabupaten_key"`
	KabupatenName    *string `db:"kabupaten_name"           json:"kabupaten_name"`
	KecamatanKey     *uint64 `db:"kecamatan_key"            json:"kecamatan_key"`
	KecamatanName    *string `db:"kecamatan_name"           json:"kecamatan_name"`
	AddressLine1     *string `db:"address_line1"            json:"address_line1"`
	PostalCode       *string `db:"postal_code"              json:"postal_code"`
}

func CreateOaPostalAddress(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_postal_address"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func UpdateOaPostalAddress(params map[string]string) (int, error) {
	query := "UPDATE oa_postal_address SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "postal_address_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE postal_address_key = " + params["postal_address_key"]
	log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func GetOaPostalAddressIn(c *[]OaPostalAddress, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				oa_postal_address.* FROM 
				oa_postal_address `
	query := query2 + " WHERE oa_postal_address." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaPostalAddress(c *OaPostalAddress, key string) (int, error) {
	query := `SELECT oa_postal_address.* FROM oa_postal_address WHERE oa_postal_address.rec_status = 1 AND oa_postal_address.postal_address_key = ` + key
	// log.Println("==========  ==========>>>",query)
	log.Println("===== QUERY OA POSTAL ADDRESS ===== >>", query)

	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetOaPostalAddressDetailIn(c *[]AddressDetail, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query := `SELECT 
				p.postal_address_key AS postal_address_key,
				p.address_type AS address_type,
				ty.lkp_name AS address_type_name,
				prov.city_key AS provinsi_key,
				prov.city_name AS provinsi_name,
				p.kabupaten_key AS kabupaten_key,
				kab.city_name AS kabupaten_name,
				p.kecamatan_key AS kecamatan_key,
				kec.city_name AS kecamatan_name,
				p.address_line1 AS address_line1,
				p.postal_code AS postal_code 
			FROM oa_postal_address AS p
			LEFT JOIN gen_lookup AS ty ON ty.lookup_key = p.address_type
			LEFT JOIN ms_city AS kab ON kab.city_key = p.kabupaten_key 
			LEFT JOIN ms_city AS prov ON prov.city_key = kab.parent_key 
			LEFT JOIN ms_city AS kec ON kec.city_key = p.kecamatan_key 
			WHERE p.rec_status = 1 AND p.` + field + ` IN(` + inQuery + `)`

	// Main query
	log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
