package models

import (
	"log"
	"mf-bo-api/db"
)

type OaRequestBankDetails struct {
	BankAccountKey  *uint64 `db:"bank_account_key" json:"bank_account_key"`
	BankKey         *uint64 `db:"bank_key" json:"bank_key"`
	BankValue       *string `db:"bank_value" json:"bank_value"`
	BankAccountNo   *string `db:"bank_account_no" json:"bank_account_no"`
	BankAccountName *string `db:"bank_account_name" json:"bank_account_name"`
	BankBranchName  *string `db:"bank_branch_name" json:"bank_branch_name"`
	FlagPriority    *uint64 `db:"flag_priority" json:"flag_priority"`
}

type PengkinianPersonalDataResponse struct {
	OaRequestKey           uint64                  `db:"oa_request_key" json:"oa_request_key"`
	OaRequestType          string                  `db:"oa_request_type" json:"oa_request_type"`
	OaRiskLevel            *string                 `db:"oa_risk_level" json:"oa_risk_level"`
	OaEntryStart           *string                 `db:"oa_entry_start" json:"oa_entry_start"`
	OaEntryEnd             *string                 `db:"oa_entry_end" json:"oa_entry_end"`
	OaStatus               *string                 `db:"oa_status" json:"oa_status"`
	EmailAddress           *string                 `db:"email_address" json:"email_address"`
	PhoneMobile            *string                 `db:"phone_mobile" json:"phone_mobile"`
	PhoneHome              *string                 `db:"phone_home" json:"phone_home"`
	PlaceBirth             *string                 `db:"place_birth" json:"place_birth"`
	DateBirth              *string                 `db:"date_birth" json:"date_birth"`
	FullName               *string                 `db:"full_name" json:"full_name"`
	Nationality            *string                 `db:"nationality" json:"nationality"`
	IdCardType             *string                 `db:"idcard_type" json:"idcard_type"`
	Gender                 *string                 `db:"gender" json:"gender"`
	Religion               *string                 `db:"religion" json:"religion"`
	Education              *string                 `db:"education" json:"education"`
	MaritalStatus          *string                 `db:"marital_status" json:"marital_status"`
	PepStatus              *string                 `db:"pep_status" json:"pep_status"`
	PepName                *string                 `db:"pep_name" json:"pep_name"`
	PepPosition            *string                 `db:"pep_position" json:"pep_position"`
	SalesCode              *string                 `db:"sales_code" json:"sales_code"`
	PicKtp                 *string                 `db:"pic_ktp" json:"pic_ktp"`
	PicSelfieKtp           *string                 `db:"pic_selfie_ktp" json:"pic_selfie_ktp"`
	OccupJob               *string                 `db:"occup_job" json:"occup_job"`
	OccupCompany           *string                 `db:"occup_company" json:"occup_company"`
	OccupPosition          *string                 `db:"occup_position" json:"occup_position"`
	OccupBusinessFields    *string                 `db:"occup_business_fields" json:"occup_business_fields"`
	AnnualIncome           *string                 `db:"annual_income" json:"annual_income"`
	SourceOfFund           *string                 `db:"sourceof_fund" json:"sourceof_fund"`
	InvesmentObjectives    *string                 `db:"invesment_objectives" json:"invesment_objectives"`
	MotherMaidenName       *string                 `db:"mother_maiden_name" json:"mother_maiden_name"`
	BeneficialRelation     *string                 `db:"beneficial_relation" json:"beneficial_relation"`
	BeneficialFullName     *string                 `db:"beneficial_full_name" json:"beneficial_full_name"`
	RelationFullName       *string                 `db:"relation_full_name" json:"relation_full_name"`
	IdCardAddress          *string                 `db:"idcard_address" json:"idcard_address"`
	IdCardProvince         *string                 `db:"idcard_province" json:"idcard_province"`
	IdCardCity             *string                 `db:"idcard_city" json:"idcard_city"`
	IdCardPostalCode       *string                 `db:"idcard_postal_code" json:"idcard_postal_code"`
	DomicileAddress        *string                 `db:"domicile_address" json:"domicile_address"`
	DomicileProvince       *string                 `db:"domicile_province" json:"domicile_province"`
	DomicileCity           *string                 `db:"domicile_city" json:"domicile_city"`
	DomicilePostalCode     *string                 `db:"domicile_postal_code" json:"domicile_postal_code"`
	OccupAddress           *string                 `db:"occup_address" json:"occup_address"`
	RelationType           *string                 `db:"relation_type" json:"relation_type"`
	RelationOccupation     *string                 `db:"relation_occupation" json:"relation_occupation"`
	RelationBusinessFields *string                 `db:"relation_business_fields" json:"relation_business_fields"`
	EmergencyRelation      *string                 `db:"emergency_relation" json:"emergency_relation"`
	EmergencyFullName      *string                 `db:"emergency_full_name" json:"emergency_full_name"`
	EmergencyPhoneNo       *string                 `db:"emergency_phone_no" json:"emergency_phone_no"`
	SiteReferrer           *string                 `db:"site_referrer" json:"site_referrer"`
	Agent                  *string                 `db:"agent" json:"agent"`
	Branch                 *string                 `db:"branch" json:"branch"`
	BankAccount            *[]OaRequestBankDetails `db:"bank_account_request" json:"bank_account_request"`
	// IdCardProvinceAlter    *string                 `db:"idcard_province_alter" json:"idcard_province_alter"`
	// IdCardCityAlter        *string                 `db:"idcard_city_alter" json:"idcard_city_alter"`
	// DomicileProvinceAlter  *string                 `db:"domicile_province_alter" json:"domicile_province_alter"`
	// DomicileCityAlter      *string                 `db:"domicile_city_alter" json:"domicile_city_alter"`

}

type PengkinianPersonalDataModels struct {
	UserLoginKey           uint64  `db:"user_login_key" json:"user_login_key"`
	OaRequestKey           uint64  `db:"oa_request_key" json:"oa_request_key"`
	OaRequestType          string  `db:"oa_request_type" json:"oa_request_type"`
	OaRiskLevel            *string `db:"oa_risk_level" json:"oa_risk_level"`
	OaEntryStart           *string `db:"oa_entry_start" json:"oa_entry_start"`
	OaEntryEnd             *string `db:"oa_entry_end" json:"oa_entry_end"`
	OaStatus               *string `db:"oa_status" json:"oa_status"`
	EmailAddress           *string `db:"email_address" json:"email_address"`
	PhoneMobile            *string `db:"phone_mobile" json:"phone_mobile"`
	PhoneHome              *string `db:"phone_home" json:"phone_home"`
	PlaceBirth             *string `db:"place_birth" json:"place_birth"`
	DateBirth              *string `db:"date_birth" json:"date_birth"`
	FullName               *string `db:"full_name" json:"full_name"`
	Nationality            *string `db:"nationality" json:"nationality"`
	IdCardType             *string `db:"idcard_type" json:"idcard_type"`
	Gender                 *string `db:"gender" json:"gender"`
	Religion               *string `db:"religion" json:"religion"`
	Education              *string `db:"education" json:"education"`
	MaritalStatus          *string `db:"marital_status" json:"marital_status"`
	PepStatus              *string `db:"pep_status" json:"pep_status"`
	PepName                *string `db:"pep_name" json:"pep_name"`
	PepPosition            *string `db:"pep_position" json:"pep_position"`
	SalesCode              *string `db:"sales_code" json:"sales_code"`
	PicKtp                 *string `db:"pic_ktp" json:"pic_ktp"`
	PicSelfieKtp           *string `db:"pic_selfie_ktp" json:"pic_selfie_ktp"`
	OccupJob               *string `db:"occup_job" json:"occup_job"`
	OccupCompany           *string `db:"occup_company" json:"occup_company"`
	OccupPosition          *string `db:"occup_position" json:"occup_position"`
	OccupBusinessFields    *string `db:"occup_business_fields" json:"occup_business_fields"`
	AnnualIncome           *string `db:"annual_income" json:"annual_income"`
	SourceOfFund           *string `db:"sourceof_fund" json:"sourceof_fund"`
	InvesmentObjectives    *string `db:"invesment_objectives" json:"invesment_objectives"`
	MotherMaidenName       *string `db:"mother_maiden_name" json:"mother_maiden_name"`
	BeneficialRelation     *string `db:"beneficial_relation" json:"beneficial_relation"`
	BeneficialFullName     *string `db:"beneficial_full_name" json:"beneficial_full_name"`
	RelationFullName       *string `db:"relation_full_name" json:"relation_full_name"`
	IdCardAddress          *string `db:"idcard_address" json:"idcard_address"`
	IdCardProvinceAlter    *string `db:"idcard_province_alter" json:"idcard_province_alter"`
	IdCardCityAlter        *string `db:"idcard_city_alter" json:"idcard_city_alter"`
	IdCardProvince         *string `db:"idcard_province" json:"idcard_province"`
	IdCardCity             *string `db:"idcard_city" json:"idcard_city"`
	IdCardPostalCode       *string `db:"idcard_postal_code" json:"idcard_postal_code"`
	DomicileAddress        *string `db:"domicile_address" json:"domicile_address"`
	DomicileProvinceAlter  *string `db:"domicile_province_alter" json:"domicile_province_alter"`
	DomicileCityAlter      *string `db:"domicile_city_alter" json:"domicile_city_alter"`
	DomicileProvince       *string `db:"domicile_province" json:"domicile_province"`
	DomicileCity           *string `db:"domicile_city" json:"domicile_city"`
	DomicilePostalCode     *string `db:"domicile_postal_code" json:"domicile_postal_code"`
	OccupAddress           *string `db:"occup_address" json:"occup_address"`
	RelationType           *string `db:"relation_type" json:"relation_type"`
	RelationOccupation     *string `db:"relation_occupation" json:"relation_occupation"`
	RelationBusinessFields *string `db:"relation_business_fields" json:"relation_business_fields"`
	EmergencyRelation      *string `db:"emergency_relation" json:"emergency_relation"`
	EmergencyFullName      *string `db:"emergency_full_name" json:"emergency_full_name"`
	EmergencyPhoneNo       *string `db:"emergency_phone_no" json:"emergency_phone_no"`
	SiteReferrer           *string `db:"site_referrer" json:"site_referrer"`
	Agent                  *string `db:"agent" json:"agent"`
	Branch                 *string `db:"branch" json:"branch"`
}

func GetPersonalDataOnlyQuery(c *PengkinianPersonalDataModels, oa_request_key string) error {

	query := `SELECT t1.user_login_key,
	t1.oa_request_key, ortype.lkp_name AS oa_request_type, 
	orlv.lkp_name AS oa_risk_level,t1.oa_entry_start, t1.oa_entry_end, 
	orst.lkp_name AS oa_status, t2.email_address, t2.phone_mobile, t2.phone_home,
	t2.place_birth, t2.date_birth, t2.full_name, msco.country_name AS nationality, 
	idtype.lkp_name AS idcard_type, gend.lkp_name AS gender, rel.lkp_name AS religion, 
	edu.lkp_name AS education, mar.lkp_name AS marital_status, pep.lkp_name AS pep_status, 
	t2.pep_name, t2.pep_position, t1.sales_code, t2.pic_ktp, t2.pic_selfie_ktp, 
	jobz.lkp_name AS occup_job, t2.occup_company, posit.lkp_name AS occup_position,
	bfield.lkp_name AS occup_business_fields, aincm.lkp_name AS annual_income,
	isrce.lkp_name AS sourceof_fund, ivobj.lkp_name AS invesment_objectives, 
	t2.mother_maiden_name, bnr.lkp_name AS beneficial_relation, t2.beneficial_full_name, 
	t2.relation_full_name, ktp.address_line1 AS idcard_address, 
	ktp.address_line2 AS idcard_province_alter, ktp.address_line3 AS idcard_city_alter, 
	ktprov.city_name AS idcard_province, ktpct.city_name AS idcard_city, 
	ktp.postal_code AS idcard_postal_code, doms.address_line1 AS domicile_address, 
	doms.address_line2 AS domicile_province_alter, doms.address_line3 AS domicile_city_alter, 
	dcp.city_name AS domicile_province, dct.city_name AS domicile_city, 
	doms.postal_code AS domicile_postal_code, cAddr.address_line1 AS occup_address,
	rlt.lkp_name AS relation_type, job.lkp_name AS relation_occupation, 
	bfi.lkp_name AS relation_business_fields, emr.lkp_name AS emergency_relation,
	t2.emergency_full_name, t2.emergency_phone_no, rfr.lkp_name AS site_referrer, t5.agent_name AS agent, 
	t6.branch_name AS branch

	FROM oa_request t1
	LEFT JOIN oa_personal_data t2 ON t2.oa_request_key = t1.oa_request_key
	LEFT JOIN oa_postal_address ktp ON ktp.postal_address_key = t2.idcard_address_key
	LEFT JOIN oa_postal_address doms ON doms.postal_address_key = t2.domicile_address_key
	LEFT JOIN oa_postal_address cAddr ON cAddr.postal_address_key = t2.occup_address_key
	LEFT JOIN ms_country msco ON msco.country_key = t2.nationality 
	LEFT JOIN ms_city ktpct ON ktpct.city_key = ktp.kabupaten_key
	LEFT JOIN ms_city ktprov ON ktprov.city_key = ktp.province_key
	LEFT JOIN ms_city dct ON dct.city_key = doms.kabupaten_key
	LEFT JOIN ms_city dcp ON dcp.city_key = doms.province_key
	LEFT JOIN ms_agent t5 ON t5.agent_key = t1.agent_key 
	LEFT JOIN ms_branch t6 ON t6.branch_key = t1.branch_key
	LEFT JOIN gen_lookup idtype ON idtype.lookup_key = t2.idcard_type
	LEFT JOIN gen_lookup jobz ON jobz.lookup_key = t2.occup_job
	LEFT JOIN gen_lookup posit ON posit.lookup_key = t2.occup_position
	LEFT JOIN gen_lookup bfield ON bfield.lookup_key = t2.occup_business_fields
	LEFT JOIN gen_lookup aincm ON aincm.lookup_key = t2.annual_income
	LEFT JOIN gen_lookup isrce ON isrce.lookup_key = t2.sourceof_fund
	LEFT JOIN gen_lookup ivobj ON ivobj.lookup_key = t2.invesment_objectives
	LEFT JOIN gen_lookup rfr ON rfr.lookup_key = t1.site_referer 
	LEFT JOIN gen_lookup rlt ON rlt.lookup_key = t2.relation_type
	LEFT JOIN gen_lookup job ON job.lookup_key = t2.relation_occupation
	LEFT JOIN gen_lookup bfi ON bfi.lookup_key = t2.relation_business_fields
	LEFT JOIN gen_lookup emr ON emr.lookup_key = t2.emergency_relation
	LEFT JOIN gen_lookup gend ON gend.lookup_key = t2.gender
	LEFT JOIN gen_lookup mar ON mar.lookup_key = t2.marital_status
	LEFT JOIN gen_lookup rel ON rel.lookup_key = t2.religion
	LEFT JOIN gen_lookup edu ON edu.lookup_key = t2.education
	LEFT JOIN gen_lookup ortype ON ortype.lookup_key = t1.oa_request_type
	LEFT JOIN gen_lookup orlv ON orlv.lookup_key = t1.oa_risk_level
	LEFT JOIN gen_lookup orst ON orst.lookup_key = t1.oa_status
	LEFT JOIN gen_lookup bnr ON bnr.lookup_key = t2.beneficial_relation
	LEFT JOIN gen_lookup pep ON pep.lookup_key = t2.pep_status
	
	WHERE t1.rec_status = 1  AND t2.rec_status = 1 AND t1.oa_request_key = ` + oa_request_key

	// EXECUTE DATANYA
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetPengkinianBankAccount(c *[]OaRequestBankDetails, oaRequestKey string) error {
	query := `SELECT 
	a2.bank_key,
	a2.bank_account_key, 
	a2.account_no as bank_account_no, 
	a2.account_holder_name as bank_account_name, 
	a3.bank_name as bank_value, 
	a2.branch_name as bank_branch_name,
	a1.flag_priority
	FROM oa_request_bank_account a1
	INNER JOIN ms_bank_account a2 ON a1.bank_account_key = a2.bank_account_key
	INNER JOIN ms_bank a3 ON a2.bank_key = a3.bank_key 
	WHERE a1.rec_status = 1 AND a1.oa_request_key = ` + oaRequestKey

	log.Println(query)

	err := db.Db.Select(c, query)
	if err != nil {
		// log.Println(err)
		return err
	}

	return nil
}
