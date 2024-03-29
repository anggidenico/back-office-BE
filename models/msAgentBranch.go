package models

import (
	"mf-bo-api/db"
	"net/http"
	"strconv"
)

type MsAgentBranch struct {
	AgentBranchKey    uint64  `db:"agent_branch_key"     json:"agent_branch_key"`
	AgentKey          uint64  `db:"agent_key"            json:"agent_key"`
	BranchKey         uint64  `db:"branch_key"           json:"branch_key"`
	EffDate           *string `db:"eff_date"             json:"eff_date"`
	Remarks           *string `db:"remarks"              json:"remarks"`
	RecOrder          *uint64 `db:"rec_order"            json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"           json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"           json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}

type MsAgentLastBranch struct {
	AgentKey  uint64 `db:"agent_key"            json:"agent_key"`
	AgentCode string `db:"agent_code"            json:"agent_code"`
	Agentname string `db:"agent_name"           json:"agent_name"`
}

func GetMsAgentLastBranch(c *[]MsAgentLastBranch, branchKey string) (int, error) {
	query := `SELECT 
				a.agent_key AS agent_key,
				a.agent_code AS agent_code,
				CONCAT(a.agent_code, " - ", a.agent_name) AS agent_name 
			FROM ms_agent_branch AS mab 
			INNER JOIN ms_agent AS a ON mab.agent_key = a.agent_key
			JOIN (
				SELECT MAX(eff_date) eff_date, agent_key 
				FROM ms_agent_branch 
				GROUP BY agent_key 
				ORDER BY eff_date DESC
				) t2 ON mab.agent_key = t2.agent_key 
			WHERE mab.rec_status = 1 AND a.rec_status = 1 
			AND t2.eff_date = mab.eff_date AND mab.branch_key = '` + branchKey + `'  
			ORDER BY mab.agent_branch_key ASC`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAllMsAgentBranch(c *[]MsAgentBranch, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_agent_branch.* FROM 
			  ms_agent_branch`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_agent_branch."+field+" = '"+value+"'")
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

func UpdateDeleteBranchAgent(params map[string]string, field string, val string) (int, error) {
	query := "UPDATE ms_agent_branch SET "
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " = '" + val + "' AND rec_status = '1'"
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

func CreateMsAgentBranch(params map[string]string) (int, error) {
	query := "INSERT INTO ms_agent_branch"
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
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetLastBranchAgent(c *MsAgentBranch, agentKey string) (int, error) {
	query := `SELECT * 
			FROM ms_agent_branch
			WHERE agent_key = "` + agentKey + `" AND rec_status = 1
			AND eff_date <= NOW() ORDER BY eff_date DESC LIMIT 1`
	// log.Println("==========  ==========>>>", query)
	err := db.Db.Get(c, query)
	if err != nil {
		// log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
