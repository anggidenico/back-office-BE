package models

import (
	"log"
	"math"
	"mf-bo-api/db"
	"strconv"
)

func GetOARequestIndividuListQuery(c *[]OaRequestListResponse, requestType uint64, backOfficeRole uint64, limit uint64, offset uint64) int {
	query := `SELECT t1.oa_request_key, t2.lkp_name AS oa_status, t3.email_address, t3.phone_mobile, 
	t3.date_birth, t3.full_name, t3.idcard_no, t1.oa_entry_start AS oa_date, 
	t4.ulogin_email AS email_address,t4.ulogin_name AS created_by, t5.agent_name AS agent, 
	t6.branch_name AS branch, t4.customer_key
	FROM oa_request t1
	INNER JOIN gen_lookup t2 ON t1.oa_status = t2.lookup_key
	INNER JOIN oa_personal_data t3 ON t1.oa_request_key = t3.oa_request_key
	INNER JOIN sc_user_login t4 ON t4.user_login_key = t1.user_login_key
	INNER JOIN ms_agent t5 ON t5.agent_key = t1.agent_key 
	INNER JOIN ms_branch t6 ON t6.branch_key = t1.branch_key
	WHERE t1.rec_status = 1 AND t3.rec_status = 1 AND t1.oa_request_type = ` + strconv.FormatUint(requestType, 10)

	queryPage := `SELECT count(*) FROM oa_request t1
	INNER JOIN gen_lookup t2 ON t1.oa_status = t2.lookup_key
	INNER JOIN oa_personal_data t3 ON t1.oa_request_key = t3.oa_request_key
	INNER JOIN sc_user_login t4 ON t4.user_login_key = t1.user_login_key
	INNER JOIN ms_agent t5 ON t5.agent_key = t1.agent_key 
	INNER JOIN ms_branch t6 ON t6.branch_key = t1.branch_key
	WHERE t1.rec_status = 1 AND t3.rec_status = 1 AND t1.oa_request_type = ` + strconv.FormatUint(requestType, 10)

	if backOfficeRole == 11 {
		query += ` AND t1.oa_status = 258`
		queryPage += ` AND t1.oa_status = 258`
	}
	if backOfficeRole == 12 {
		query += ` AND t1.oa_status = 259`
		queryPage += ` AND t1.oa_status = 259`
	}

	log.Println(limit)
	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// EXECUTE DATANYA
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err.Error())
	}

	// EXECUTE PAGING
	var pagination int
	var count uint64
	log.Println(queryPage)
	err = db.Db.Get(&count, queryPage)
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

	return pagination
}
