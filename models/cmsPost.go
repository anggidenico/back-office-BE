package models

import (
	"database/sql"
	"mf-bo-api/config"
	"mf-bo-api/db"
	"net/http"
	"strconv"
	"strings"
)

type CmsPost struct {
	PostKey              uint64  `db:"post_key"                  json:"post_key"`
	PostSubtypeKey       uint64  `db:"post_subtype_key"          json:"post_subtype_key"`
	PostTitle            string  `db:"post_title"                json:"post_title"`
	PostSubTitle         *string `db:"post_sub_title"            json:"post_sub_title"`
	PostContent          *string `db:"post_content"              json:"post_content"`
	PostContentAuthor    *string `db:"post_content_author"       json:"post_content_author"`
	PostContentSources   *string `db:"post_content_sources"      json:"post_content_sources"`
	PostPublishStart     string  `db:"post_publish_start"        json:"post_publish_start"`
	PostPublishThru      string  `db:"post_publish_thru"         json:"post_publish_thru"`
	PostPageAllowed      uint8   `db:"post_page_allowed"         json:"post_page_allowed"`
	PostCommentAllowed   uint8   `db:"post_comment_allowed"      json:"post_comment_allowed"`
	PostCommentDisplayed uint8   `db:"post_comment_displayed"    json:"post_comment_displayed"`
	PostFilesAllowed     uint8   `db:"post_files_allowed"        json:"post_files_allowed"`
	PostVideoAllowed     uint8   `db:"post_video_allowed"        json:"post_video_allowed"`
	PostVideoUrl         *string `db:"post_video_url"            json:"post_video_url"`
	PostPinned           uint8   `db:"post_pinned"               json:"post_pinned"`
	PostOwnerKey         *uint64 `db:"post_owner_key"            json:"post_owner_key"`
	RecOrder             *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus            uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate       *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy         *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate      *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy        *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1            *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2            *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus    *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage     *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate      *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy        *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate       *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy         *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1      *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2      *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3      *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type CmsPostData struct {
	PostKey              uint64             `json:"post_key"`
	PostSubtype          CmsPostSubtypeInfo `json:"post_subtype"`
	PostTitle            string             `json:"post_title"`
	PostSubTitle         string             `json:"post_sub_title"`
	PostContent          string             `json:"post_content"`
	PostContentAuthor    string             `json:"post_content_author"`
	PostContentSources   string             `json:"post_content_sources"`
	PostPublishStart     string             `json:"post_publish_start"`
	PostPublishThru      string             `json:"post_publish_thru"`
	PostPageAllowed      bool               `json:"post_page_allowed"`
	PostCommentAllowed   bool               `json:"post_comment_allowed"`
	PostCommentDisplayed bool               `json:"post_comment_displayed"`
	PostFilesAllowed     bool               `json:"post_files_allowed"`
	PostVideoAllowed     bool               `json:"post_video_allowed"`
	PostVideoUrl         string             `json:"post_video_url"`
	PostPinned           bool               `json:"post_pinned"`
	RecImage1            string             `json:"rec_image1"`
	RecImage2            string             `json:"rec_image2"`
}

type CmsPostList struct {
	PostKey            uint64             `json:"post_key"`
	PostSubtype        CmsPostSubtypeInfo `json:"post_subtype"`
	PostTitle          string             `json:"post_title"`
	PostSubTitle       string             `json:"post_sub_title"`
	PostContentAuthor  string             `json:"post_content_author"`
	PostContentSources string             `json:"post_content_sources"`
	PostPublishStart   string             `json:"post_publish_start"`
	PostPublishThru    string             `json:"post_publish_thru"`
	PostPinned         bool               `json:"post_pinned"`
	RecImage1          string             `json:"rec_image1"`
}

type AdminCmsPostList struct {
	PostKey            uint64  `json:"post_key"`
	PostTypeKey        uint64  `json:"post_type_key"`
	PostTypeCode       string  `json:"post_type_code"`
	PostTypeName       *string `json:"post_type_name"`
	PostTypeDesc       *string `json:"post_type_desc"`
	PostSubtypeKey     uint64  `json:"post_subtype_key"`
	PostSubtypeCode    string  `json:"post_subtype_code"`
	PostSubtypeName    *string `json:"post_subtype_name"`
	PostTitle          string  `json:"post_title"`
	PostSubTitle       *string `json:"post_sub_title"`
	PostContentAuthor  *string `json:"post_content_author"`
	PostContentSources *string `json:"post_content_sources"`
	PostPublishStart   string  `json:"post_publish_start"`
	PostPublishThru    string  `json:"post_publish_thru"`
	PostPinned         bool    `json:"post_pinned"`
	RecImage1          string  `json:"rec_image1"`
}

type PublicGuidance struct {
	PostSubtypeKey  uint64  `db:"post_subtype_key"  json:"post_subtype_key"`
	PostSubtypeCode string  `db:"post_subtype_code" json:"post_subtype_code"`
	PostSubtypeName *string `db:"post_subtype_name" json:"post_subtype_name"`
	PostTitle       string  `db:"post_title"        json:"post_title"`
	PostSubTitle    *string `db:"post_sub_title"    json:"post_sub_title"`
	PostContent     string  `db:"post_content"      json:"post_content"`
	RecImage1       string  `db:"rec_image1"        json:"rec_image1"`
}

type CmsPostCount struct {
	CountData int `db:"count_data"             json:"count_data"`
}

func GetAllCmsPost(c *[]CmsPost, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              cms_post.* FROM 
			  cms_post WHERE 
			  cms_post.post_publish_start <= NOW() AND 
			  cms_post.post_publish_thru > NOW() AND 
			  cms_post.rec_status = 1 `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "cms_post."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
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
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCmsPostIn(c *[]CmsPost, value []string, field string, params map[string]string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				cms_post.* FROM 
				cms_post WHERE 
				cms_post.post_publish_start <= NOW() AND 
				cms_post.post_publish_thru > NOW() AND 
				cms_post.rec_status = 1 `

	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
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

	query := query2 + " AND cms_post." + field + " IN(" + inQuery + ")"

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
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCmsPost(c *CmsPost, key string) (int, error) {
	query := `SELECT cms_post.* FROM cms_post WHERE cms_post.rec_status = 1 AND cms_post.post_key = ` + key
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAdminCmsPostListIn(c *[]CmsPost, limit uint64, offset uint64, nolimit bool, params map[string]string, ids []string) (int, error) {
	var present bool
	var whereClause []string
	var condition string

	query := `SELECT
				cms_post.* FROM 
				cms_post WHERE 
				cms_post.rec_status = 1 `

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "cms_post."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if len(ids) > 0 {
		inQuery := strings.Join(ids, ",")
		condition += " AND cms_post.post_subtype_key IN(" + inQuery + ")"
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
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCountCmsPost(c *CmsPostCount, params map[string]string, ids []string) (int, error) {
	var whereClause []string
	var condition string

	query := `SELECT
              count(cms_post.post_key) as count_data
			  FROM cms_post WHERE 
			  cms_post.rec_status = 1 `

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "cms_post."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if len(ids) > 0 {
		inQuery := strings.Join(ids, ",")
		condition += " AND cms_post.post_subtype_key IN(" + inQuery + ")"
	}

	query += condition

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreatePost(params map[string]string) (int, error) {
	query := "INSERT INTO cms_post"
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
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateCmsPost(params map[string]string) (int, error) {
	query := "UPDATE cms_post SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "post_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE post_key = " + params["post_key"]
	// log.Println("==========  ==========>>>", query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetPublicGuidance(c *[]PublicGuidance) (int, error) {
	query := `SELECT 
	t2.post_subtype_key, 
	t2.post_subtype_code, 
	t2.post_subtype_name, 
	t3.post_title, 
	t3.post_sub_title, 
	t3.post_content, 
	CONCAT("` + config.ImageUrl + `", "/images/post/guidance/", t3.rec_image1) AS rec_image1
	FROM cms_post_type t1
	INNER JOIN cms_post_subtype t2 ON (t1.post_type_key=t2.post_type_key)
	INNER JOIN cms_post t3 ON (t3.post_subtype_key = t2.post_subtype_key AND t3.rec_status = 1)
	where t1.rec_status = 1
	AND t2.rec_status = 1
	AND t1.post_type_code='GUIDANCE'
	AND t2.post_subtype_code IN ('PUBLIC_GUIDANCE','PUBLIC_ADDRESS')`

	// Main query
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
