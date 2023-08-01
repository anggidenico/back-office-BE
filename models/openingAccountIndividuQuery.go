package models

import (
	"mf-bo-api/lib"
	"strconv"
)

func GetOpeningAccountIndividuListQuery(c *[]OaRequestListResponse, backOfficeRole uint64, limit uint64, offset uint64) {
	query := `SELECT t1.oa_request_key, t2.lkp_name AS oa_status, t3.email_address, t3.phone_mobile, 
	t3.date_birth, t3.full_name, t3.idcard_no, t1.oa_entry_start AS oa_date, t4.ulogin_email, t6.branch_name, t5.agent_name 
	FROM oa_request t1
	INNER JOIN gen_lookup t2 ON t1.oa_status = t2.lookup_key
	INNER JOIN oa_personal_data t3 ON t1.oa_request_key = t3.oa_request_key
	INNER JOIN sc_user_login t4 ON t4.user_login_key = t1.user_login_key
	INNER JOIN ms_agent t5 ON t5.agent_key = t1.agent_key 
	INNER JOIN ms_branch t6 ON t6.branch_key = t1.branch_key
	WHERE t1.rec_status = 1 AND t1.request_type = 127 AND t3.rec_status = 1 `
	if backOfficeRole == 11 {
		query += ` AND t1.oa_status = 258`
	}
	if backOfficeRole == 12 {
		query += ` AND t1.oa_status = 259`
	}

	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

}
