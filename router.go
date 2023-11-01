package main

import (
	"log"
	"mf-bo-api/config"
	"mf-bo-api/controllers"
	"mf-bo-api/db"
	"mf-bo-api/lib"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func router() *echo.Echo {

	db.Db = db.DBInit()
	db.DbDashboard = db.DBDashboardInit()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.CORSAllowOrigin,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Use(printUrlMiddleware)
	// e.Use(middleware.Logger())

	admin := e.Group("/admin")

	admin.Use(lib.AuthenticationMiddleware)

	e.POST("/loginbo", controllers.LoginBo).Name = "LoginBo"

	admin.GET("/lookup", controllers.GetGenLookup).Name = "GetGenLookup"
	// admin.GET("/city/:field/:key", controllers.GetMsCityList).Name = "GetMsCityList"
	admin.GET("/country", controllers.GetMsCountryList).Name = "GetMsCountryList"
	admin.GET("/bank", controllers.GetMsBankList).Name = "GetMsBankList"

	//Admin OA Request
	admin.GET("/oarequestlist", controllers.GetOaRequestList).Name = "GetOaRequestList"
	admin.GET("/oarequestdata/:key", controllers.GetOaRequestData).Name = "GetOaRequestData"
	admin.POST("/updatestatusapproval/cs", controllers.UpdateStatusApprovalCS).Name = "UpdateStatusApprovalCS"
	admin.POST("/updatestatusapproval/compliance", controllers.UpdateStatusApprovalCompliance).Name = "UpdateStatusApprovalCompliance"
	admin.GET("/oarequestlist/dotransaction", controllers.GetOaRequestListDoTransaction).Name = "GetOaRequestListDoTransaction"
	admin.GET("/downloadformatsinvest", controllers.DownloadOaRequestFormatSinvest).Name = "DownloadOaRequestFormatSinvest"
	admin.POST("/uploadformatsinvest", controllers.UploadOaRequestFormatSinvest).Name = "UploadOaRequestFormatSinvest"
	admin.GET("/oarequestdata/lasthistory/:key", controllers.GetLastHistoryOaRequestData).Name = "GetLastHistoryOaRequestData"

	//Admin Post
	admin.GET("/posts", controllers.GetAdminCmsPostList).Name = "GetAdminCmsPostList"
	admin.GET("/post/:key", controllers.GetAdminCmsPostData).Name = "GetAdminCmsPostData"
	admin.POST("/post/create", controllers.CreateAdminCmsPost).Name = "CreateAdminCmsPost"
	admin.POST("/post/update/:key", controllers.UpdateAdminCmsPost).Name = "UpdateAdminCmsPost"
	admin.POST("/post/delete", controllers.DeleteAdminCmsPost).Name = "DeleteAdminCmsPost"
	admin.GET("/typelist", controllers.AdminGetListPostType).Name = "AdminGetListPostType"
	admin.GET("/subtypelist", controllers.AdminGetListPostSubtype).Name = "AdminGetListPostSubtype"
	admin.GET("/subtypelist/:post_type", controllers.AdminGetListPostSubtypeByType).Name = "AdminGetListPostSubtypeByType"

	//Admin CMS Financial Calc
	admin.GET("/financialcalc", controllers.GetAdminCmsFinancialCalcList).Name = "GetAdminCmsFinancialCalcList"

	//Admin Transaction
	admin.GET("/transactionlist", controllers.GetTransactionApprovalList).Name = "GetTransactionApprovalList"
	admin.GET("/transactionlist/cutoff", controllers.GetTransactionCutOffList).Name = "GetTransactionCutOffList"
	admin.GET("/transactionlist/batch", controllers.GetTransactionBatchList).Name = "GetTransactionBatchList"
	admin.GET("/transactionlist/confirmation", controllers.GetTransactionConfirmationList).Name = "GetTransactionConfirmationList"
	admin.GET("/transactionlist/correctionadmin", controllers.GetTransactionCorrectionAdminList).Name = "GetTransactionCorrectionAdminList"
	admin.GET("/transactionlist/posting", controllers.GetTransactionPostingList).Name = "GetTransactionPostingList"
	admin.GET("/transaction/:key", controllers.GetTransactionDetail).Name = "GetTransactionDetail"
	admin.POST("/transactionapproval/cs", controllers.TransactionApprovalCs).Name = "TransactionApprovalCs"
	admin.POST("/transactionapproval/compliance", controllers.TransactionApprovalCompliance).Name = "TransactionApprovalCompliance"
	admin.POST("/transaction/updatenavdate", controllers.UpdateNavDate).Name = "UpdateNavDate"
	admin.POST("/transactionapproval/cutoff", controllers.TransactionApprovalCutOff).Name = "TransactionApprovalCutOff"
	admin.GET("/transaction/downloadformatsinvest", controllers.DownloadTransactionFormatSinvest).Name = "DownloadTransactionFormatSinvest"
	admin.POST("/transaction/downloadformatexcel", controllers.GetFormatExcelDownloadList).Name = "GetFormatExcelDownloadList"

	admin.POST("/transaction/uploadexcelconfirmation", controllers.UploadExcelConfirmation).Name = "UploadExcelConfirmation"

	admin.POST("/transactionapproval/posting", controllers.ProsesPosting).Name = "ProsesPosting"
	admin.POST("/transactionapproval/posting-all", controllers.ProsesPostingAll).Name = "ProsesPostingAll"

	admin.GET("/transaction/productbanklist/:key", controllers.GetProductBankAccount).Name = "GetProductBankAccount"
	admin.GET("/transaction/customerbanklist/:key", controllers.GetCustomerBankAccount).Name = "GetCustomerBankAccount"
	admin.GET("/transactionlist/unposting", controllers.GetTransactionUnpostingList).Name = "GetTransactionUnpostingList"
	admin.POST("/transactionapproval/unposting", controllers.ProsesUnposting).Name = "ProsesUnposting"
	admin.GET("/transaction/inquirylist", controllers.DataTransaksiInquiry).Name = "DataTransaksiInquiry"
	admin.GET("/transaction/inquiry/:key", controllers.DetailTransaksiInquiry).Name = "DetailTransaksiInquiry"
	admin.POST("/transactioncorrection", controllers.ProsesCorrection).Name = "ProsesCorrection"
	admin.GET("/transactionstamps", controllers.GetTransactionStampsAdmin).Name = "GetTransactionStamps"
	admin.GET("/transactiontypelist", controllers.GetTransactionType).Name = "GetTransactionType"
	admin.GET("/transactionstatuslist", controllers.GetTransactionStatus).Name = "GetTransactionStatus"
	admin.GET("/transactionothersfee", controllers.GetTransactionOthersFee)

	//Admin Product
	admin.GET("/productlist", controllers.GetListProductAdmin).Name = "GetListProductAdmin"
	admin.GET("/product/:key", controllers.GetProductDetailAdmin).Name = "GetProductDetailAdmin"
	admin.POST("/product/delete", controllers.DeleteProductAdmin).Name = "DeleteProductAdmin"
	admin.POST("/product/create", controllers.CreateAdminMsProduct).Name = "CreateAdminMsProduct"
	admin.POST("/product/update", controllers.UpdateAdminMsProduct).Name = "UpdateAdminMsProduct"

	//Admin Product
	admin.GET("/currencylist", controllers.GetListMsCurrency).Name = "GetListMsCurrency"

	//Admin Custodian Bank
	admin.GET("/custodianbanklist", controllers.GetListMsCustodianBank).Name = "GetListMsCustodianBank"
	admin.GET("/custodian-bank-list", controllers.AdminGetListMsCustodianBank).Name = "AdminGetListMsCustodianBank"
	admin.POST("/custodian-bank/delete", controllers.AdminDeleteMsCustodianBank).Name = "AdminDeleteMsCustodianBank"
	admin.POST("/custodian-bank/create", controllers.AdminCreateMsCustodianBank).Name = "AdminCreateMsCustodianBank"
	admin.POST("/custodian-bank/update", controllers.AdminUpdateMsCustodianBank).Name = "AdminUpdateMsCustodianBank"
	admin.GET("/custodian-bank/detail/:custodian_key", controllers.AdminDetailMsCustodianBank).Name = "AdminDetailMsCustodianBank"

	//Admin Fund Structure
	admin.GET("/fundstructurelist", controllers.GetListMsFundStructure).Name = "GetListMsFundStructure"

	//Admin Fund Type
	admin.GET("/fundtypelist", controllers.AdminGetListMsFundType).Name = "AdminGetListMsFundType"

	//Admin Product Category
	admin.GET("/productcategorylist", controllers.AdminGetListMsProductCategory).Name = "AdminGetListMsProductCategory"

	//Admin Product Type
	admin.GET("/producttypelist", controllers.AdminGetListMsProductType).Name = "AdminGetListMsProductType"

	//Admin Product Type
	admin.GET("/riskprofilelist", controllers.AdminGetListMsRiskProfile).Name = "AdminGetListMsRiskProfile"

	//Admin Product
	admin.GET("/productfeelist", controllers.GetListProductFeeAdmin).Name = "GetListProductFeeAdmin"
	admin.GET("/productfee/:key", controllers.GetProductFeeDetailAdmin).Name = "GetProductFeeDetailAdmin"
	admin.GET("/productlist/dropdown", controllers.GetListProductAdminDropdown).Name = "GetListProductAdminDropdown"
	admin.POST("/productfee/delete", controllers.DeleteProductFeeAdmin).Name = "DeleteProductFeeAdmin"
	// admin.POST("/productfee/create", controllers.CreateAdminMsProductFee).Name = "CreateAdminMsProductFee"
	admin.POST("/productfee/create", controllers.CreateProductFeeSettings)
	// admin.POST("/productfee/update", controllers.UpdateAdminMsProductFee).Name = "UpdateAdminMsProductFee"
	admin.POST("/productfee/update", controllers.UpdateProductFeeSettings)
	admin.POST("/productfee/feeitem/create", controllers.CreateAdminMsProductFeeItem).Name = "CreateAdminMsProductFeeItem"
	admin.POST("/productfee/feeitem/update", controllers.UpdateAdminMsProductFeeItem).Name = "UpdateAdminMsProductFeeItem"
	admin.POST("/productfee/feeitem/delete", controllers.DeleteAdminMsProductFeeItem).Name = "DeleteAdminMsProductFeeItem"
	admin.GET("/productfee/feeitem/:key", controllers.DetailAdminMsProductFeeItem).Name = "DetailAdminMsProductFeeItem"

	//Admin Product Bank Account
	admin.GET("/productbankaccountlist", controllers.GetListProductBankAccountAdmin).Name = "GetListProductBankAccountAdmin"
	admin.GET("/productbankaccount/:key", controllers.GetProductBankAccountDetailAdmin).Name = "GetProductBankAccountDetailAdmin"
	admin.POST("/productbankaccount/delete", controllers.DeleteProductBankAccountAdmin).Name = "DeleteProductBankAccountAdmin"
	admin.POST("/productbankaccount/create", controllers.CreateAdminMsProductBankAccount).Name = "CreateAdminMsProductBankAccount"
	admin.POST("/productbankaccount/update", controllers.UpdateAdminMsProductBankAccount).Name = "UpdateAdminMsProductBankAccount"

	//Admin User Management
	admin.GET("/logout", controllers.LogoutAdmin).Name = "LogoutAdmin"
	admin.GET("/usermanagementlist", controllers.GetListScUserLoginAdmin).Name = "GetListScUserLoginAdmin"
	admin.GET("/usermanagement/:key", controllers.GetDetailScUserLoginAdmin).Name = "GetDetailScUserLoginAdmin"
	admin.POST("/usermanagement/disableenable", controllers.DisableEnableUser).Name = "DisableEnableUser"
	admin.POST("/usermanagement/lockunlock", controllers.LockUnlockUser).Name = "LockUnlockUser"
	admin.GET("/rolelist", controllers.AdminGetListScRole).Name = "AdminGetListScRole"
	admin.GET("/usercategorylist", controllers.AdminGetListScUserCategory).Name = "AdminGetListScUserCategory"
	admin.GET("/userdeptlist", controllers.AdminGetListScUserDept).Name = "AdminGetListScUserDept"
	admin.POST("/usermanagement/create", controllers.CreateAdminScUserLogin).Name = "CreateAdminScUserLogin"
	admin.POST("/usermanagement/update", controllers.UpdateAdminScUserLogin).Name = "UpdateAdminScUserLogin"
	admin.POST("/usermanagement/changepassword", controllers.ChangePasswordUser).Name = "ChangePasswordUser"
	admin.POST("/usermanagement/changerole", controllers.ChangeRoleUser).Name = "ChangeRoleUser"
	admin.POST("/usermanagement/delete", controllers.DeleteUser).Name = "DeleteUser"
	admin.POST("/usermanagement/userdata-update", controllers.UpdateAuthScUserLogin).Name = "UpdateAuthScUserLogin"

	//Admin Role Management
	admin.GET("/rolemanagementlist", controllers.GetListRoleManagementAdmin).Name = "GetListRoleManagementAdmin"
	admin.GET("/rolemanagement/userlist", controllers.GetListUserByRole).Name = "GetListUserByRole"
	admin.GET("/rolemanagement/:key", controllers.GetDetailRoleManagement).Name = "GetDetailRoleManagement"
	admin.GET("/rolemanagement/menulist", controllers.GetDetailMenuRoleManagement).Name = "GetDetailMenuRoleManagement"
	admin.GET("/rolemanagement/rolecategorylist", controllers.GetListRoleCategory).Name = "GetListRoleCategory"
	admin.POST("/rolemanagement/delete", controllers.DeleteRoleManagement).Name = "DeleteRoleManagement"
	admin.POST("/rolemanagement/create", controllers.CreateAdminRoleManagement).Name = "CreateAdminRoleManagement"
	admin.POST("/rolemanagement/update", controllers.UpdateAdminRoleManagement).Name = "UpdateAdminRoleManagement"

	//Admin NAV
	admin.GET("/navlist", controllers.GetListTrNavAdmin).Name = "GetListTrNavAdmin"
	admin.GET("/nav/:key", controllers.GetNavDetailAdmin).Name = "GetNavDetailAdmin"
	admin.POST("/nav/delete", controllers.DeleteNavAdmin).Name = "DeleteNavAdmin"
	admin.POST("/nav/create", controllers.CreateAdminTrNav).Name = "CreateAdminTrNav"
	admin.POST("/nav/update", controllers.UpdateAdminTrNav).Name = "UpdateAdminTrNav"
	admin.GET("/nav/get-price", controllers.GetNavPrice).Name = ""

	//Admin NAV
	admin.GET("/menu", controllers.GetListMenuLogin).Name = "GetListMenuLogin"

	//Admin Customer
	admin.GET("/customer/individu/list", controllers.GetListCustomerIndividuInquiry).Name = "GetListCustomerIndividuInquiry"
	admin.GET("/customer/institution/list", controllers.GetListCustomerInstitutionInquiry).Name = "GetListCustomerInstitutionInquiry"
	admin.GET("/customer/individu/:key", controllers.GetDetailCustomerIndividu).Name = "GetDetailCustomerIndividu"
	admin.GET("/customer/institution/:key", controllers.GetDetailCustomerInstitution).Name = "GetDetailCustomerInstitution"
	admin.GET("/customer/detail/:key", controllers.GetDetailCustomerInquiry).Name = "GetDetailCustomerInquiry"
	admin.GET("/personaldata/individu/:key", controllers.DetailPersonalDataCustomerIndividu).Name = "DetailPersonalDataCustomerIndividu"
	admin.POST("/customer/create", controllers.AdminCreateCustomerIndividu).Name = "AdminCreateCustomerIndividu"
	admin.GET("/customer/redemption/dropdown", controllers.GetAdminListCustomerRedemption).Name = "GetAdminListCustomerRedemption"
	admin.GET("/customer/individu/data/:key", controllers.GetAdminOaRequestPersonalDataRiskProfile).Name = "GetAdminOaRequestPersonalDataRiskProfile"
	admin.POST("/customer/pengkinian/create", controllers.AdminSavePengkinianCustomerIndividu).Name = "AdminSavePengkinianCustomerIndividu"
	admin.POST("/customer/check-unique-email-nohp", controllers.CheckUniqueEmailNoHp).Name = "CheckUniqueEmailNoHp"
	admin.POST("/customer/check-unique-no-id", controllers.CheckUniqueNoId).Name = "CheckUniqueNoId"
	admin.GET("/customer/individu/send-account-statement", controllers.IndividuSendAccountStatement).Name = "SendAccount"

	//Admin Transaction Report
	admin.GET("/report/transactionhistorylist", controllers.GetTransactionHistory).Name = "GetTransactionHistory"
	admin.GET("/customer/dropdown", controllers.GetListCustomerDropDown).Name = "GetListCustomerDropDown"
	admin.GET("/customer/branch-agent/dropdown", controllers.GetListCustomerBranchAgentDropdown).Name = "GetListCustomerBranchAgentDropdown"
	admin.GET("/report/transactionhistory/:param", controllers.GetDetailCustomerProduct).Name = "GetDetailCustomerProduct"

	//Admin Transaction Report Daily
	admin.GET("/report/banktransaction", controllers.GetBankProductTransaction).Name = "GetBankProductTransaction"
	admin.GET("/report/daily-subscription", controllers.GetTransactionReportSubscribeDaily).Name = "GetTransactionReportSubscribeDaily"
	admin.GET("/report/daily-redemption", controllers.GetTransactionReportRedemptionDaily).Name = "GetTransactionReportRedemptionDaily"

	//Transaction Action Admin Subscription
	//subscribe
	admin.GET("/transaction/subscription", controllers.GetTransactionSubscription).Name = "GetTransactionSubscription"
	admin.GET("/product/subscription/:fund_type_key", controllers.AdminGetProductSubscription).Name = "AdminGetProductSubscription"
	admin.GET("/product/detail-subscription/:product_key", controllers.GetProductDetailTransactionSubscription).Name = "GetProductDetailTransactionSubscription"
	admin.GET("/bankproduct/subscription/:product_key", controllers.GetBankProductSubscription).Name = "GetBankProductSubscription"
	admin.POST("/createtransaction/subscription", controllers.CreateTransactionSubscription).Name = "CreateTransactionSubscription"
	admin.GET("/transaction/topupdata/:customer_key/:product_key", controllers.GetTopupData).Name = "GetTopupData"

	//Transaction Action Admin Redemption
	//redemption
	admin.GET("/transaction/redemption", controllers.GetTransactionRedemption).Name = "GetTransactionRedemption"
	// admin.GET("/product/redemption/:customer_key", controllers.AdminGetProductRedemption).Name = "AdminGetProductRedemption"
	admin.GET("/product/redemption", controllers.AdminGetProductListCanRedeem)
	admin.GET("/transaction/customerbankredemption/:customer_key", controllers.GetCustomerBankAccountRedemption).Name = "GetCustomerBankAccountRedemption"
	admin.GET("/metode-perhitungan", controllers.GetMetodePerhitungan).Name = "GetMetodePerhitungan"
	admin.POST("/createtransaction/redemption", controllers.CreateTransactionRedemption).Name = "CreateTransactionRedemption"

	//Transaction Action Admin Switching
	//switching
	admin.GET("/transaction/switching", controllers.GetTransactionSwitching).Name = "GetTransactionSwitching"
	admin.POST("/createtransaction/switching", controllers.CreateTransactionSwitching).Name = "CreateTransactionSwitching"
	admin.GET("/product/switchin/:customer_key/:product_switch_out_key", controllers.AdminGetProductSwitchIn).Name = "AdminGetProductSwitchIn"

	admin.POST("/transaction/delete", controllers.DeleteTransactionAdmin).Name = "DeleteTransactionAdmin"

	//Admin OA Pengkinian Risk Profile
	admin.GET("/pengkinian/profile-risiko", controllers.GetListPengkinianRiskProfile).Name = "GetListPengkinianRiskProfile"
	admin.GET("/pengkinian/detail/profile-risiko/:key", controllers.GetDetailPengkinianProfileRisiko).Name = "GetDetailPengkinianProfileRisiko"
	admin.GET("/pengkinian/detail/lasthistory/profile-risiko/:key", controllers.GetDetailPengkinianProfileRisikoLastHistory).Name = "GetDetailPengkinianProfileRisikoLastHistory"

	//Admin OA Pengkinian Personal Data
	admin.GET("/pengkinian/personal-data", controllers.GetListPengkinianPersonalData).Name = "GetListPengkinianPersonalData"
	admin.GET("/pengkinian/detail/personal-data/:key", controllers.GetDetailPengkinianPersonalData).Name = "GetDetailPengkinianPersonalData"
	admin.GET("/pengkinian/detail/lasthistory/personal-data/:key", controllers.GetDetailPengkinianPersonalDataLastHistory).Name = "GetDetailPengkinianPersonalDataLastHistory"

	//Admin Promo
	admin.GET("/promo", controllers.GetListPromo).Name = "GetListPromo"
	admin.POST("/promo/create", controllers.CreateAdminTrPromo).Name = "CreateAdminTrPromo"
	admin.POST("/promo/delete", controllers.DeletePromo).Name = "DeletePromo"
	admin.POST("/promo/update", controllers.UpdateAdminTrPromo).Name = "UpdateAdminTrPromo"
	admin.GET("/promo/detail/:key", controllers.DetailPromo).Name = "DetailPromo"
	admin.POST("/promo/check", controllers.CheckPromo).Name = "CheckPromo"

	//Admin Daily Transaction Report - J005
	admin.GET("/report/daily-transaction-report", controllers.GetDailyTransactionReport).Name = "GetDailyTransactionReport"

	//Admin MS BRANCH - J005
	admin.GET("/branchlist/dropdown", controllers.GetListBranchDropdown).Name = "GetListBranchDropdown"
	admin.GET("/branch-list", controllers.AdminGetListMsBranch).Name = "AdminGetListMsBranch"
	admin.POST("/branch/delete", controllers.AdminDeleteMsBranch).Name = "AdminDeleteMsBranch"
	admin.POST("/branch/create", controllers.AdminCreateMsBranch).Name = "AdminCreateMsBranch"
	admin.POST("/branch/update", controllers.AdminUpdateMsBranch).Name = "AdminUpdateMsBranch"
	admin.GET("/branch/detail/:branch_key", controllers.AdminDetailMsBranch).Name = "AdminDetailMsBranch"

	//Admin MS AGENT - J005
	admin.GET("/agentlist/dropdown", controllers.GetListAgentDropdown).Name = "GetListAgentDropdown"
	admin.GET("/agentlist/branch/:branch_key", controllers.GetListAgentLastBranch).Name = "GetListAgentLastBranch"
	admin.GET("/agent-list", controllers.AdminGetListMsAgent).Name = "AdminGetListMsAgent"
	admin.POST("/agent/delete", controllers.AdminDeleteMsAgent).Name = "AdminDeleteMsAgent"
	admin.POST("/agent/create", controllers.AdminCreateMsAgent).Name = "AdminCreateMsAgent"
	admin.POST("/agent/update", controllers.AdminUpdateMsAgent).Name = "AdminUpdateMsAgent"
	admin.GET("/agent/detail/:agent_key", controllers.AdminDetailMsAgent).Name = "AdminDetailMsAgent"
	admin.GET("/agent-customer-list", controllers.AdminGetListMsAgentCustomer).Name = "AdminGetListMsAgentCustomer"
	admin.POST("/agent-customer-movement", controllers.AdminAgentCustomerMovement).Name = "AdminAgentCustomerMovement"

	//Admin Subscription Batch Confirmation - J006
	admin.GET("/report/subscription-batch-confirmation", controllers.GetSubscriptionBatchConfirmation).Name = "GetSubscriptionBatchConfirmation"

	//Admin Redemption Batch Confirmation - J007
	admin.GET("/report/redemption-batch-confirmation", controllers.GetRedemptionBatchConfirmation).Name = "GetRedemptionBatchConfirmation"

	//Admin Promo
	admin.GET("/user-notif", controllers.GetListUserNotif).Name = "GetListUserNotif"
	admin.POST("/user-notif/create", controllers.CreateAdminScUserNotif).Name = "CreateAdminScUserNotif"
	admin.POST("/user-notif/delete", controllers.DeleteUserNotif).Name = "DeleteUserNotif"
	admin.POST("/user-notif/update", controllers.UpdateAdminScUserNotif).Name = "UpdateAdminScUserNotif"
	admin.GET("/user-notif/detail/:key", controllers.DetailUserNotif).Name = "DetailUserNotif"

	// MANAGEMENT USER YANG SEDANG LOGIN SAAT INI
	admin.GET("/user", controllers.GetDetailScUserLogin).Name = "GetDetailScUserLogin"
	admin.POST("/user/changepassword", controllers.AdminChangePasswordUserLogin).Name = "AdminChangePasswordUserLogin"
	admin.POST("/user/changedata", controllers.AdminChangeDataUserLogin).Name = "AdminChangeDataUserLogin"

	//Admin Data Suspend Status Customer (CIF)
	admin.GET("/customer/suspendstatuslist", controllers.GetListCustomerIndividuStatusSuspend).Name = "GetListCustomerIndividuStatusSuspend"
	admin.GET("/customer/detail/status-suspend/:customer_key", controllers.AdminGetDetailCustomer).Name = "AdminGetDetailCustomer"
	admin.POST("/customer/suspend-unsuspend", controllers.AdminSuspendUnsuspendCustomer).Name = "AdminSuspendUnsuspendCustomer"
	admin.POST("/customer/upload-suspend-unsuspend", controllers.UploadExcelCIFSuspendUnsuspend).Name = "UploadExcelCIFSuspendUnsuspend"

	//Admin Data Suspend Account
	admin.GET("/accountlist", controllers.GetListTrAccount).Name = "GetListTrAccount"
	admin.GET("/account/detail/:acc_key", controllers.AdminGetDetailAccount).Name = "AdminGetDetailAccount"
	admin.POST("/account/update", controllers.AdminUpdateTrAccount).Name = "AdminUpdateTrAccount"
	admin.GET("/account/customerlist/:product_key", controllers.AdminGetCustomerAccount).Name = "AdminGetCustomerAccount"

	//Admin Customer File
	admin.POST("/customer-file-update", controllers.CustomerUpdateFile).Name = "CustomerUpdateFile"
	admin.GET("/customer-file-detail/:customer_key", controllers.AdminGetDetailCustomerDocument).Name = "AdminGetDetailCustomerDocument"

	//Admin Menu
	admin.GET("/menu-list", controllers.AdminGetListMenu).Name = "AdminGetListMenu"
	admin.POST("/menu/delete", controllers.AdminDeleteMenu).Name = "AdminDeleteMenu"
	admin.POST("/menu/create", controllers.AdminCreateMenu).Name = "AdminCreateMenu"
	admin.POST("/menu/update", controllers.AdminUpdateMenu).Name = "AdminUpdateMenu"
	admin.GET("/menu/detail/:menu_key", controllers.AdminDetailMenu).Name = "AdminDetailMenu"

	//Admin Menu Type
	admin.GET("/menu-type-dropdown", controllers.AdminGetListScMenuTypeDropdown).Name = "AdminGetListScMenuTypeDropdown"

	//Admin App Module
	admin.GET("/app-module-dropdown", controllers.AdminGetListScAppModuleDropdown).Name = "AdminGetListScAppModuleDropdown"

	//Admin Menu
	admin.GET("/user-dept-list", controllers.GetListScUserDeptAdmin).Name = "GetListScUserDeptAdmin"
	admin.POST("/user-dept/delete", controllers.AdminDeleteScUserDept).Name = "AdminDeleteScUserDept"
	admin.POST("/user-dept/create", controllers.AdminCreateUserDept).Name = "AdminCreateUserDept"
	admin.POST("/user-dept/update", controllers.AdminUpdateUserDept).Name = "AdminUpdateUserDept"
	admin.GET("/user-dept/detail/:user_dept_key", controllers.AdminDetailUserDept).Name = "AdminDetailUserDept"

	//Admin Account Statement
	admin.GET("/account-statement-customer-product", controllers.AdminDetailAccountStatementCustomerProduct).Name = "AdminDetailAccountStatementCustomerProduct"
	admin.GET("/account-statement-customer-agent", controllers.AdminDetailAccountStatementCustomerAgent).Name = "AdminDetailAccountStatementCustomerAgent"

	//Admin Sc App Config
	admin.GET("/app-config-list", controllers.AdminGetListScAppConfig).Name = "AdminGetListScAppConfig"
	admin.POST("/app-config/delete", controllers.AdminDeleteScAppConfig).Name = "AdminDeleteScAppConfig"
	admin.POST("/app-config/create", controllers.AdminCreateScAppConfig).Name = "AdminCreateScAppConfig"
	admin.POST("/app-config/update", controllers.AdminUpdateScAppConfig).Name = "AdminUpdateScAppConfig"
	admin.GET("/app-config/detail/:app_config_key", controllers.AdminDetailScAppConfig).Name = "AdminDetailScAppConfig"

	//Admin Sc App Config Type
	admin.GET("/app-config-type-dropdown", controllers.AdminGetListDropdownScAppConfigType).Name = "AdminGetListDropdownScAppConfigType"

	//Admin Gen Lookup
	admin.GET("/lookup-list", controllers.AdminGetListLookup).Name = "AdminGetListLookup"
	admin.POST("/lookup/delete", controllers.AdminDeleteLookup).Name = "AdminDeleteLookup"
	admin.POST("/lookup/create", controllers.AdminCreateLookup).Name = "AdminCreateLookup"
	admin.POST("/lookup/update", controllers.AdminUpdateLookup).Name = "AdminUpdateLookup"
	admin.GET("/lookup/detail/:lookup_key", controllers.AdminDetailLookup).Name = "AdminDetailLookup"

	//Admin Gen Lookup Group
	admin.GET("/lookup-group-dropdown", controllers.AdminGetListDropdownLookupGroup).Name = "AdminGetListDropdownLookupGroup"

	//Admin Ms Participant
	admin.GET("/participant-dropdown", controllers.AdminGetListDropdownMsParticipant).Name = "AdminGetListDropdownMsParticipant"

	//Admin Bank
	admin.GET("/bank-list", controllers.AdminGetListMsBank).Name = "AdminGetListMsBank"
	admin.POST("/bank/delete", controllers.AdminDeleteMsBank).Name = "AdminDeleteMsBank"
	admin.POST("/bank/create", controllers.AdminCreateMsBank).Name = "AdminCreateMsBank"
	admin.POST("/bank/update", controllers.AdminUpdateMsBank).Name = "AdminUpdateMsBank"
	admin.GET("/bank/detail/:bank_key", controllers.AdminDetailBank).Name = "AdminDetailBank"

	//Admin Bank Charges
	admin.GET("/bank-charges-list", controllers.AdminGetListMsBankCharges).Name = "AdminGetListMsBankCharges"
	admin.POST("/bank-charges/delete", controllers.AdminDeleteMsBankCharges).Name = "AdminDeleteMsBankCharges"
	admin.POST("/bank-charges/create", controllers.AdminCreateMsBankCharges).Name = "AdminCreateMsBankCharges"
	admin.POST("/bank-charges/update", controllers.AdminUpdateMsBankCharges).Name = "AdminUpdateMsBankCharges"
	admin.GET("/bank-charges/detail/:bcharges_key", controllers.AdminDetailBankCharges).Name = "AdminDetailBankCharges"

	//Admin Currency
	admin.GET("/currency-list", controllers.AdminGetListMsCurrency).Name = "AdminGetListMsCurrency"
	admin.POST("/currency/delete", controllers.AdminDeleteMsCurrency).Name = "AdminDeleteMsCurrency"
	admin.POST("/currency/create", controllers.AdminCreateMsCurrency).Name = "AdminCreateMsCurrency"
	admin.POST("/currency/update", controllers.AdminUpdateMsCurrency).Name = "AdminUpdateMsCurrency"
	admin.GET("/currency/detail/:currency_key", controllers.AdminDetailMsCurrency).Name = "AdminDetailMsCurrency"

	//Admin Currency Rate
	admin.GET("/currency-rate-list", controllers.GetListTrCurrencyRate).Name = "GetListTrCurrencyRate"
	admin.POST("/currency-rate/delete", controllers.AdminDeleteTrCurrencyRate).Name = "AdminDeleteTrCurrencyRate"
	admin.POST("/currency-rate/create", controllers.AdminCreateTrCurrencyRate).Name = "AdminCreateTrCurrencyRate"
	admin.POST("/currency-rate/update", controllers.AdminUpdateTrCurrencyRate).Name = "AdminUpdateTrCurrencyRate"
	admin.GET("/currency-rate/detail/:curr_rate_key", controllers.AdminDetailTrCurrencyRate).Name = "AdminDetailTrCurrencyRate"

	//Admin Currency Rate
	admin.GET("/country-list", controllers.AdminGetListMsCountry).Name = "AdminGetListMsCountry"
	admin.POST("/country/delete", controllers.AdminDeleteMsCountry).Name = "AdminDeleteMsCountry"
	admin.POST("/country/create", controllers.AdminCreateMsCountry).Name = "AdminCreateMsCountry"
	admin.POST("/country/update", controllers.AdminUpdateMsCountry).Name = "AdminUpdateMsCountry"
	admin.GET("/country/detail/:country_key", controllers.AdminDetailMsCountry).Name = "AdminDetailMsCountry"

	//Admin Holiday
	admin.GET("/holiday-list", controllers.AdminGetListMsHoliday).Name = "AdminGetListMsHoliday"
	admin.POST("/holiday/delete", controllers.AdminDeleteMsHoliday).Name = "AdminDeleteMsHoliday"
	admin.POST("/holiday/create", controllers.AdminCreateMsHoliday).Name = "AdminCreateMsHoliday"
	admin.POST("/holiday/update", controllers.AdminUpdateMsHoliday).Name = "AdminUpdateMsHoliday"
	admin.GET("/holiday/detail/:holiday_key", controllers.AdminDetailMsHoliday).Name = "AdminDetailMsHoliday"
	admin.GET("/holiday/status", controllers.AdminCheckHoliday).Name = "AdminCheckHoliday"

	//Admin City
	admin.GET("/city-list", controllers.AdminGetListMsCity).Name = "AdminGetListMsCity"
	admin.GET("/city-level", controllers.GetCityLevel).Name = "GetCityLevel"
	admin.POST("/city/delete", controllers.AdminDeleteMsCity).Name = "AdminDeleteMsCity"
	admin.POST("/city/create", controllers.AdminCreateMsCity).Name = "AdminCreateMsCity"
	admin.POST("/city/update", controllers.AdminUpdateMsCity).Name = "AdminUpdateMsCity"
	admin.GET("/city/detail/:city_key", controllers.AdminDetailMsCity).Name = "AdminDetailMsCity"
	admin.GET("/city-parent", controllers.GetCityParent).Name = "GetCityParent"

	admin.GET("/mail-list", controllers.AdminGetListMmMailMaster).Name = "AdminGetListMmMailMaster"
	admin.POST("/mail/delete", controllers.AdminDeleteMmMailMaster).Name = "AdminDeleteMmMailMaster"
	admin.POST("/mail/create", controllers.AdminCreateMmMailMaster).Name = "AdminCreateMmMailMaster"
	admin.POST("/mail/update", controllers.AdminUpdateMmMailMaster).Name = "AdminUpdateMmMailMaster"
	admin.GET("/mail/detail/:mail_master_key", controllers.AdminDetailMmMailMaster).Name = "AdminDetailMmMailMaster"
	admin.POST("/tes-sent-email", controllers.TestSentEmail).Name = "TestSentEmail"

	//Admin OA Institusi
	admin.GET("/oa-institution-data-cs", controllers.GetListOaInstitusiDataCS).Name = "GetListOaInstitusiDataCS"
	admin.GET("/oa-institution-approve-cs", controllers.GetListOaInstitusiApproveCS).Name = "GetListOaInstitusiApproveCS"
	admin.GET("/oa-institution-data-branch", controllers.GetListOaInstitusiDataBranch).Name = "GetListOaInstitusiDataBranch"
	admin.GET("/oa-institution-approve-kyc", controllers.GetListOaInstitusiApproveKYC).Name = "GetListOaInstitusiApproveKYC"
	admin.GET("/oa-institution-approve-fund-admin", controllers.GetListOaInstitusiApproveFundAdmin).Name = "GetListOaInstitusiApproveFundAdmin"
	admin.GET("/oa-institution-detail/:request_key", controllers.GetDetailOaInstitusiAdmin).Name = "GetDetailOaInstitusiAdmin"

	admin.POST("/oa-institution-save", controllers.SaveOaInstitutionData).Name = "SaveOaInstitutionData"
	admin.POST("/oa-institution-docs-save", controllers.SaveDocsInstitution).Name = "SaveDocsInstitution"
	admin.POST("/oa-institution-docs-delete", controllers.DeleteDocsInstitution).Name = "DeleteDocsInstitution"
	admin.POST("/oa-institution-risk-profile-save", controllers.SaveRiskProfileInstitution).Name = "SaveRiskProfileInstitution"
	admin.POST("/oa-institution-user-save", controllers.SaveInstitutionUser).Name = "SaveInstitutionUser"
	admin.POST("/oa-institution-save-to-approver", controllers.SaveInstitutionToApprover).Name = "SaveInstitutionToApprover"

	admin.POST("/oa-institution-cs-approvedata", controllers.CSApproveOAInstitution).Name = "CSApproveOAInstitution"
	admin.POST("/oa-institution-kyc-approvedata", controllers.KYCAPproveOAInstitution).Name = "KYCAPproveOAInstitution"
	admin.GET("/oa-institution-download-sinvest", controllers.DownloadOaRequestInstitutionFormatSinvest).Name = "DownloadOaRequestInstitutionFormatSinvest"
	admin.POST("/oa-institution-upload-format-sinvest", controllers.UploadOaInstitutionRequestFormatSinvest).Name = "UploadOaInstitutionRequestFormatSinvest"

	//Admin Transaction Settlement Report - TSR
	admin.GET("/transaction-settlement-report", controllers.GetTransactionPaymentReport).Name = "GetTransactionPaymentReport"

	//Admin Auto Invest
	admin.GET("/list-autoinvest", controllers.ListAdminAutoInvest).Name = "ListAdminAutoInvest"
	admin.POST("/create-autoinvest", controllers.AdminCreateTrAutoInvest).Name = "AdminCreateTrAutoInvest"
	admin.POST("/update-autoinvest", controllers.AdminUpdateTrAutoInvest).Name = "AdminUpdateTrAutoInvest"
	admin.POST("/delete-autoinvest", controllers.AdminDeleteTrAutoInvest).Name = "AdminDeleteTrAutoInvest"
	admin.GET("/detail-autoinvest/:autoinvest_key", controllers.AdminDetailTrAutoInvest).Name = "AdminDetailTrAutoInvest"
	admin.POST("/generate-autoinvest", controllers.AdminGenerateTransactionFromTrAutoInvest).Name = "AdminGenerateTransactionFromTrAutoInvest"

	//admin correction
	admin.GET("/transaction/correction", controllers.GetListAdminTransactionCorrection).Name = "GetListAdminTransactionCorrection"
	admin.POST("/updatetransaction/subscription", controllers.UpdateTransactionSubscription).Name = "UpdateTransactionSubscription"
	admin.POST("/updatetransaction/redemption", controllers.UpdateTransactionRedemption).Name = "UpdateTransactionRedemption"

	// quiz
	admin.GET("/quiz", controllers.GetCmsQuiz).Name = "GetCmsQuiz"
	admin.GET("/getquizresult", controllers.GetQuizAnswer).Name = "GetQuizAnswer"
	admin.POST("/quizresult", controllers.PostQuizAnswer).Name = "PostQuizAnswer"

	admin.GET("/bankproduct/:product_key", controllers.GetBankProductSubscription).Name = "GetBankProductSubscription"
	admin.GET("/bankcustomer/:customer_key", controllers.GetCustomerBankAccountRedemption).Name = "GetCustomerBankAccountRedemption"

	//==================================================== NEW ENDPOINT ======================================================================

	e.GET("/get-lookup-by-group/:group_key", controllers.GetLookupByGroupKey)

	// OA DAN PENGKINIAN
	admin.GET("/new-oarequest-list", controllers.GetNewOAList)
	admin.GET("/pengkinian-personal-data-list", controllers.GetPengkinianPersonalDataList)
	admin.GET("/pengkinian-personal-data-details/:key", controllers.GetPengkinianPersonalDataDetails)
	admin.GET("/pengkinian-risk-profile-list", controllers.GetPengkinianRiskProfileList)
	admin.GET("/pengkinian-risk-profile-details/:key", controllers.GetPengkinianRiskProfileDetails)
	admin.GET("/oarequest-list/kyc-approved", controllers.GetOaRequestListKycApproved)
	admin.GET("/donwload-filetext-oa-request", controllers.DownloadOaRequestTextFile)
	admin.GET("/oarequest-list/customer-build", controllers.GetOaRequestListCustomerBuild)
	admin.POST("/revert-oa-status", controllers.RevertOAStatus)

	// SETTING SC USER
	admin.GET("/customer-individu-list", controllers.NewGetCustomerInquiryList)
	admin.GET("/user-bo-management-list", controllers.GetUserBackOfficeList)
	admin.GET("/user-customer-list", controllers.GetUserCustomerList)
	admin.POST("/user-customer-edit", controllers.UserCustomerEdit)

	// MANAGEMENT RISK PROFILE QUIZ
	admin.GET("/get-riskprofile-question-list", controllers.GetRiskProfileQuestionList)
	admin.GET("/get-question-detail/:question_key", controllers.GetQuestionDetail)
	admin.POST("/create-riskprofile-question", controllers.CreateQuizQuestion)
	admin.POST("/update-riskprofile-question", controllers.UpdateQuizQuestion)
	admin.POST("/delete-riskprofile-question", controllers.DeleteQuizQuestion)
	admin.GET("/get-option-list/:question_key", controllers.GetOptionListPerQuestion)
	admin.GET("/get-option-detail/:quiz_option_key", controllers.GetOptionDetail)
	admin.POST("/create-quiz-option", controllers.CreateQuizOption)
	admin.POST("/update-quiz-option", controllers.UpdateQuizOption)
	admin.POST("/delete-quiz-option", controllers.DeleteQuizOption)

	admin.POST("/save-product", controllers.SaveMasterProduct)

	//RISK-PROFILE
	admin.POST("/create-risk-profile", controllers.CreateRiskProfile)
	admin.GET("/get-risk-profile", controllers.GetriskProfileController)
	admin.POST("/update-risk-profile", controllers.UpdateRiskProfile)
	admin.POST("/delete-risk-profile", controllers.DeleteRiskProfile)
	admin.GET("/get-detail-risk-profile/:risk_profile_key", controllers.GetDetailRiskProfileController)

	//End Point
	admin.GET("/get-endpoint-sc", controllers.GetEndpointscController)
	admin.GET("/get-endpoint-detail/:endpoint_key", controllers.GetEndpointDetailController)
	admin.POST("/create-endpoint-sc", controllers.CreateEndpointController)
	admin.POST("/update-endpoint-sc", controllers.UpdateEndpointController)
	admin.POST("/delete-endpoint-sc", controllers.DeleteEndpointController)

	// PRODUCT PAYMENT CHANNEL
	admin.GET("/product-channel/:product_key", controllers.ProductPaymentChannelList)
	admin.POST("/product-channel/create", controllers.CreateProductPaymentChannels)
	admin.POST("/product-channel/update", controllers.UpdateProductPaymentChannels)
	admin.POST("/product-channel/delete", controllers.DeleteProductPaymentChannels)

	// UPDATE PRODUCT WITH APPROVAL
	admin.POST("/product/create-request", controllers.ProductCreateRequest)
	admin.POST("/product/update-request", controllers.ProductUpdateRequest)
	admin.GET("/product/approval-list", controllers.ProductApprovalList)
	admin.GET("/product/request-detail/", controllers.ProductApprovalDetails)

	// MASTER PAYMENT CHANNEL
	admin.GET("/get-payment-channel", controllers.GetmsPaymentChannelController)
	admin.GET("/get-payment-channel/:pchannel_key", controllers.GetMsPaymentDetailController)
	admin.POST("/delete-payment-channel", controllers.DeleteMsPaymentChannelController)
	admin.POST("/create-payment-channel", controllers.CreateMsPaymentChannelController)
	admin.POST("/update-payment-channel", controllers.UpdateMsPaymentChannelController)

	// Benchmark
	admin.GET("/get-benchmark-list", controllers.GetFfsBenchmarkController)
	admin.GET("/get-benchmark-detail/:benchmark_key", controllers.GetBenchmarkDetailController)
	admin.POST("/delete-benchmark", controllers.DeleteBenchmarkController)
	admin.POST("/create-benchmark", controllers.CreateFfsBenchmarkController)
	admin.POST("/update-benchmark", controllers.UpdateFfsBenchmarkController)

	//Product Benchmark
	admin.POST("/create-benchmark-product", controllers.CreateFfsBenchmarkProductController)
	admin.GET("/get-benchmark-product-list", controllers.GetFfsBenchmarkProductController)
	admin.GET("/get-benchmark-product/:bench_prod_key", controllers.GetBenchmarkProdDetailController)
	admin.POST("/delete-benchmark-product", controllers.DeleteBenchmarkProdController)
	admin.POST("/update-benchmark-product", controllers.UpdateBenchmarkProdController)

	//Securities List
	admin.POST("/create-securities-list", controllers.CreateMsSecuritiesController)
	admin.GET("/get-securities-list", controllers.GetMsSecuritiesController)
	admin.POST("/delete-securities-list", controllers.DeleteMsSecuritiesController)
	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL, "|", c.Request().UserAgent(), "|", c.RealIP())
		if c.Request().Method == "POST" {
			aa, _ := c.FormParams()
			log.Println(aa.Encode())
		}

		return next(c)
	}
}
