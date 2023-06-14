package models

import (
	"log"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"
)

type MsFile struct {
	FileKey         *uint64 `db:"file_key"                json:"file_key"`
	RefFkKey        *uint64 `db:"ref_fk_key"              json:"ref_fk_key"`
	RefFkDomain     *string `db:"ref_fk_domain"           json:"ref_fk_domain"`
	FileName        *string `db:"file_name"               json:"file_name"`
	FileExt         *string `db:"file_ext"                json:"file_ext"`
	BlobMode        *uint8  `db:"blob_mode"               json:"blob_mode"`
	FilePath        *string `db:"file_path"               json:"file_path"`
	FileUrl         *string `db:"file_url"                json:"file_url"`
	FileNotes       *string `db:"file_notes"              json:"file_notes"`
	FileObj         *uint64 `db:"file_obj"                json:"properties"`
	RecStatus       *uint64 `db:"rec_status"                json:"rec_status"`
	RecCreatedDate  *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy    *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy   *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecAttributeId1 *string `db:"rec_attribute_id1" json:"rec_attribute_id1"`
	RecAttributeId2 *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeId3 *string `db:"rec_attribute_id3"                json:"rec_attribute_id3"`
}

type MsFileDetail struct {
	FileKey        uint64  `db:"file_key"                json:"file_key"`
	RefFkKey       uint64  `db:"ref_fk_key"              json:"ref_fk_key"`
	FileName       string  `db:"file_name"               json:"file_name"`
	Path           string  `db:"path"                    json:"path"`
	FileExt        string  `db:"file_ext"                json:"file_ext"`
	FileNotes      *string `db:"file_notes"              json:"file_notes"`
	RecCreatedDate *string `db:"rec_created_date"        json:"rec_created_date"`
}

type CustomerDocumentDetail struct {
	Customer CustomerIndividuStatusSuspend `json:"customer"`
	Document []MsFileDetail                `json:"document"`
}

func GetAllMsFile(c *[]MsFile, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_file.* FROM 
			  ms_file`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_file."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsFile(params map[string]string) (int, error) {
	query := "UPDATE ms_file SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "file_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE file_key = " + params["file_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func CreateMsFile(params map[string]string) (int, error, string) {
	query := "INSERT INTO ms_file"
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
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	ret, err := tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetALlDetailMsFile(c *[]MsFileDetail, params map[string]string) (int, error) {
	query := `SELECT 
				ms_file.file_key,
				ms_file.ref_fk_key,
				ms_file.file_name,
				ms_file.file_ext,
				ms_file.file_notes,
				DATE_FORMAT(ms_file.rec_created_date, '%d %M %Y %H:%i') AS rec_created_date 
			FROM ms_file AS ms_file`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_file."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsFileWithIn(params map[string]string, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query := "UPDATE ms_file SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " IN(" + inQuery + ")"
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

type MsFileModels struct {
	FileKey uint64 `db:"file_key" json:"file_key"`
	// RefFkKey    uint64 `db:"ref_fk_key" json:"ref_fk_key"`
	// RefFkDomain string `db:"ref_fk_domain" json:"ref_fk_domain"`
	FilePath string `db:"file_path" json:"file_path"`
}

func GetMsFileDataWithCondition(c *[]MsFileModels, params map[string]string) (int, error) {

	query := `SELECT t1.file_key, t1.file_path FROM ms_file t1 WHERE t1.rec_status = 1`
	if valueMap, ok := params["ref_fk_key"]; ok {
		query += ` AND t1.ref_fk_key = ` + valueMap
	}
	if valueMap, ok := params["ref_fk_domain"]; ok {
		query += ` AND t1.ref_fk_domain = "` + valueMap + `"`
	}
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
