package models

import (
	"fmt"
	"log"
	"mf-bo-api/db"
	"reflect"
)

func ProductApprovalAction(params map[string]string) error {
	tx, err := db.Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	recBy := params["rec_by"]
	recDate := params["rec_date"]

	var ProdReqData ProductRequest

	qGetData := `SELECT rec_pk, rec_action, product_key, product_id, product_code, product_name, product_name_alt,currency_key, product_category_key, fund_type_key, product_profile, investment_objectives, product_phase, nav_valuation_type, prospectus_link, launch_date, inception_date, isin_code, flag_syariah, max_sub_fee, max_red_fee, max_swi_fee, min_sub_amount, min_topup_amount, min_red_amount, min_red_amount, min_red_unit, min_unit_after_red, min_amount_after_red, management_fee, custodian_fee, custodian_key, settlement_period, sinvest_fund_code, flag_enabled, flag_subscription, flag_redemption, flag_redemption, flag_switch_out, flag_switch_in, dec_unit, dec_amount, dec_nav, dec_performance, npwp_date_reg, npwp_name, npwp_number, portfolio_code, rec_created_date
	FROM ms_product_request WHERE rec_status = 1 AND rec_pk = ` + params["rec_pk"]

	row := tx.QueryRow(qGetData)
	err = row.Scan(&ProdReqData.RecPK, &ProdReqData.RecAction, &ProdReqData.ProductKey, &ProdReqData.ProductID, &ProdReqData.ProductCode, &ProdReqData.ProductName, &ProdReqData.ProductNameAlt, &ProdReqData.CurrencyKey, &ProdReqData.ProductCategoryKey, &ProdReqData.FundTypeKey, &ProdReqData.ProductProfile, &ProdReqData.InvestmentObjectives, &ProdReqData.ProductPhase, &ProdReqData.NavValuationType, &ProdReqData.ProspectusLink, &ProdReqData.LaunchDate, &ProdReqData.InceptionDate, &ProdReqData.IsinCode, &ProdReqData.FlagSyariah, &ProdReqData.MaxSubFee, &ProdReqData.MaxRedFee, &ProdReqData.MaxSwiFee, &ProdReqData.MinSubAmount, &ProdReqData.MinTopUpAmount, &ProdReqData.MinRedAmount, &ProdReqData.MinRedAmount, &ProdReqData.MinRedUnit, &ProdReqData.MinUnitAfterRed, &ProdReqData.MinAmountAfterRed, &ProdReqData.ManagementFee, &ProdReqData.CustodianFee, &ProdReqData.CustodianKey, &ProdReqData.SettlementPeriod, &ProdReqData.SinvestFundCode, &ProdReqData.FlagEnabled, &ProdReqData.FlagSubscription, &ProdReqData.FlagRedemption, &ProdReqData.FlagRedemption, &ProdReqData.FlagSwitchOut, &ProdReqData.FlagSwitchIn, &ProdReqData.DecUnit, &ProdReqData.DecAmount, &ProdReqData.DecNav, &ProdReqData.DecPerformance, &ProdReqData.NpwpDateReg, &ProdReqData.NpwpName, &ProdReqData.NpwpNumber, &ProdReqData.PortfolioCode, &ProdReqData.RecCreatedDate)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// UpdProductRequest := `UPDATE ms_product_request SET rec_approval_status = ` + params["approval"] + ` AND rec_approved_date = '` + params["rec_date"] + ` AND rec_approved_by = '` + params["rec_by"] + ` AND rec_attribute_id1 = '` + params["reason"] + `' WHERE rec_pk = ` + params["rec_pk"]

	// _, err = tx.Exec(UpdProductRequest)
	// if err != nil {
	// 	tx.Rollback()
	// 	log.Println(err.Error())
	// 	return err
	// }

	if *ProdReqData.RecAction == "CREATE" {
		insertMsProduct := make(map[string]string)
		var reflectValue = reflect.ValueOf(ProdReqData)
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
		var reflectType = reflectValue.Type()
		for i := 0; i < reflectValue.NumField(); i++ {
			columnName := reflectType.Field(i).Tag.Get("db")
			columnValue := fmt.Sprintf("%v", reflectValue.Field(i).Interface())
			insertMsProduct[columnName] = columnValue
			insertMsProduct["rec_status"] = "1"
			insertMsProduct["rec_created_by"] = recBy
			insertMsProduct["rec_created_date"] = recDate
		}
		var fields, values string
		var bindvars []interface{}
		for key, value := range insertMsProduct {
			fields += key + ", "
			values += ` "` + value + `", `
			bindvars = append(bindvars, value)
		}
		fields = fields[:(len(fields) - 2)]
		values = values[:(len(values) - 2)]
		query := `INSERT INTO ms_product (` + fields + `) VALUES(` + values + `)`

		log.Println(query)

	}

	if *ProdReqData.RecAction == "UPDATE" {

	}

	if *ProdReqData.RecAction == "DELETE" {

	}

	tx.Commit()
	return nil
}
