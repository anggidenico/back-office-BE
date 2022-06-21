package controllers

import (
	"bufio"
	"database/sql"
	"mf-bo-api/config"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func DownloadOaRequestInstitutionFormatSinvest(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
	// var limit uint64

	var oaRequestDB []models.OaRequest
	status, err = models.GetOaInstitutionDoTransaction(&oaRequestDB)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var oaRequestLookupIds []string
	var oaRequestIds []string
	var customerIds []string
	for _, oareq := range oaRequestDB {
		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}

		if oareq.OaRiskLevel != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
			}
		}

		if oareq.CustomerKey != nil {
			if _, ok := lib.Find(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10)); !ok {
				customerIds = append(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10))
			}
		}
	}

	//mapping personal data
	var institutionData []models.OaInstitutionData
	if len(oaRequestIds) > 0 {
		status, err = models.GetOaInstitutionDataIn(&institutionData, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	pdData := make(map[uint64]models.OaInstitutionData)
	var postalAddressIds []string
	var nasionalityIds []string

	for _, insData := range institutionData {
		pdData[insData.OaRequestKey] = insData

		if insData.Nationality != nil {
			if _, ok := lib.Find(nasionalityIds, strconv.FormatUint(*insData.Nationality, 10)); !ok {
				nasionalityIds = append(nasionalityIds, strconv.FormatUint(*insData.Nationality, 10))
			}
		}

		if insData.CorrespondenceKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*insData.CorrespondenceKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*insData.CorrespondenceKey, 10))
			}
		}

		if insData.IntitutionType != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*insData.IntitutionType, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*insData.IntitutionType, 10))
			}
		}

		if insData.IntitutionCharacteristic != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*insData.IntitutionCharacteristic, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*insData.IntitutionCharacteristic, 10))
			}
		}

		if insData.InstiAnnuallyIncome != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*insData.InstiAnnuallyIncome, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*insData.InstiAnnuallyIncome, 10))
			}
		}

		if insData.InstiInvestmentPurpose != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*insData.InstiInvestmentPurpose, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*insData.InstiInvestmentPurpose, 10))
			}
		}

		if insData.InstiSourceOfIncome != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*insData.InstiSourceOfIncome, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*insData.InstiSourceOfIncome, 10))
			}
		}
	}

	//gen lookup
	var lookupOaReq []models.GenLookup
	if len(oaRequestLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, oaRequestLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	//postal data
	var oaPostalAddressList []models.OaPostalAddress
	if len(postalAddressIds) > 0 {
		status, err = models.GetOaPostalAddressIn(&oaPostalAddressList, postalAddressIds, "postal_address_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	postalData := make(map[uint64]models.OaPostalAddress)
	var oaCityIds []string
	for _, posAdd := range oaPostalAddressList {
		postalData[posAdd.PostalAddressKey] = posAdd

		if posAdd.KabupatenKey != nil {
			if _, ok := lib.Find(oaCityIds, strconv.FormatUint(*posAdd.KabupatenKey, 10)); !ok {
				oaCityIds = append(oaCityIds, strconv.FormatUint(*posAdd.KabupatenKey, 10))
			}
		}
		if posAdd.KecamatanKey != nil {
			if _, ok := lib.Find(oaCityIds, strconv.FormatUint(*posAdd.KecamatanKey, 10)); !ok {
				oaCityIds = append(oaCityIds, strconv.FormatUint(*posAdd.KecamatanKey, 10))
			}
		}
	}

	var cityList []models.MsCity
	if len(oaCityIds) > 0 {
		status, err = models.GetMsCityIn(&cityList, oaCityIds, "city_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	cityData := make(map[uint64]models.MsCity)
	for _, city := range cityList {
		cityData[city.CityKey] = city

		if _, ok := lib.Find(nasionalityIds, strconv.FormatUint(city.CountryKey, 10)); !ok {
			nasionalityIds = append(nasionalityIds, strconv.FormatUint(city.CountryKey, 10))
		}
	}

	var countryList []models.MsCountry
	status, err = models.GetMsCountryIn(&countryList, nasionalityIds, "country_key")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	countryData := make(map[uint64]models.MsCountry)
	for _, country := range countryList {
		countryData[country.CountryKey] = country
	}

	//customer
	var customer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&customer, customerIds, "customer_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, cus := range customer {
		customerData[cus.CustomerKey] = cus
	}

	var responseData []models.OaRequestCsvFormatFiksTxt

	var scApp models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp, "LAST_CLIENT_CODE")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data LAST_CLIENT_CODE")
	}

	log.Println("last = " + *scApp.AppConfigValue)

	last, _ := strconv.ParseUint(*scApp.AppConfigValue, 10, 64)
	if last == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var lastValue string

	txtHeader := "Type|SA Code|SID|Company Name|Country of Domicile|SIUP No.|SIUP Expiration Date|SKD No.|SKD Expiration Date|NPWP No.|NPWP Registration Date|Country of Establishment |Place of Establishment |Date of Establishment|Articles of Association No.|Company Type|Company Characteristic|Income Level (IDR)|Investor’s Risk  Profile|Investment Objective|Source of Fund|Asset Owner|Company Address|Company City Code|Company City Name|Company Postal Code|Country of Company|Office Phone|Facsimile|Email|Statement Type|Authorized Person 1 - First Name|Authorized Person 1 - Middle Name|Authorized Person 1 - Last Name|Authorized Person 1 - Position|Authorized Person 1 - Mobile  Phone|Authorized Person 1 - Email|Authorized Person 1 - NPWP No.|Authorized Person 1 - KTP No.|Authorized Person 1 - KTP Expiration Date|Authorized Person 1 - Passport No.|Authorized Person 1 - Passport Expiration Date|Authorized Person 2 - First Name|Authorized Person 2 - Middle Name|Authorized Person 2 - Last Name|Authorized Person 2 - Position|Authorized Person 2 - Mobile  Phone|Authorized Person 2 - Email|Authorized Person 2 - NPWP No.|Authorized Person 2 - KTP No.|Authorized Person 2 - KTP Expiration Date|Authorized Person 2 - Passport No.|Authorized Person 2 - Passport Expiration Date|Asset Information for the Past 3 Years (IDR) - Last Year|Asset Information for the Past 3 Years (IDR) - 2 Years Ago|Asset Information for the Past 3 Years (IDR) - 3 Years Ago|Profit Information for the Past 3 Years (IDR) - Last Year|Profit Information for the Past 3 Years (IDR) - 2 Years Ago|Profit Information for the Past 3 Years (IDR) - 3 Years Ago|FATCA (Status)|TIN / Foreign TIN|TIN / Foreign TIN Issuance Country|GIIN|Substantial U.S. Owner Name|Substantial U.S. Owner Address|Substantial U.S. Owner TIN|REDM Payment Bank BIC Code 1|REDM Payment Bank BI Member Code 1|REDM Payment Bank Name 1|REDM Payment Bank Country 1|REDM Payment Bank Branch 1|REDM Payment A/C CCY 1|REDM Payment A/C No. 1|REDM Payment A/C Name 1|REDM Payment Bank BIC Code 2|REDM Payment Bank BI Member Code 2|REDM Payment Bank Name 2|REDM Payment Bank Country 2|REDM Payment Bank Branch 2|REDM Payment A/C CCY 2|REDM Payment A/C No. 2|REDM Payment A/C Name 2|REDM Payment Bank BIC Code 3|REDM Payment Bank BI Member Code 3|REDM Payment Bank Name 3|REDM Payment Bank Country 3|REDM Payment Bank Branch 3|REDM Payment A/C CCY 3|REDM Payment A/C No. 3|REDM Payment A/C Name 3|Client Code"
	var dataRow models.OaRequestCsvFormatFiksTxt
	dataRow.DataRow = txtHeader
	responseData = append(responseData, dataRow)

	for _, oareq := range oaRequestDB {
		if n, ok := pdData[oareq.OaRequestKey]; ok {
			var data models.OaRequestInstitutionCsvFormat

			strType := strconv.FormatUint(*oareq.OaRequestType, 10)
			if strType == "127" {
				data.Type = "1"
			} else {
				data.Type = "2"
			}
			data.SACode = "EP002"
			data.SID = ""

			if oareq.CustomerKey != nil {
				if g, ok := customerData[*oareq.CustomerKey]; ok {
					if g.SidNo != nil {
						io := *g.SidNo
						data.SID = io
					}
				}
			}
			data.CompanyName = *n.FullName
			data.CountryOfDomicile = "ID"
			if n.Nationality != nil {
				if g, ok := countryData[*n.Nationality]; ok {
					data.CountryOfDomicile = g.CouCode
				}
			}

			layout := "2006-01-02 15:04:05"
			newLayout := "20060102"

			data.SiupNo = ""
			data.SiupExpirationDate = ""
			if data.CountryOfDomicile != "ID" {
				if n.NibNo != nil {
					data.SiupNo = *n.NibNo
				}
				if n.NibDate != nil {
					date, _ := time.Parse(layout, *n.NibDate)
					data.SiupExpirationDate = date.Format(newLayout)
				}
			}

			data.SkdNo = ""
			data.SkdExpirationDate = ""
			if n.SkdLicenseNo != nil && *n.SkdLicenseNo != "" {
				data.SkdNo = *n.SkdLicenseNo
				if n.SkdLicenseDate != nil {
					date, _ := time.Parse(layout, *n.SkdLicenseDate)
					data.SkdExpirationDate = date.Format(newLayout)
				}
			} else {
				if n.LastChangeAaNo != nil && *n.LastChangeAaNo != "" {
					data.SkdNo = *n.LastChangeAaNo
					if n.LastChangeAaDate != nil {
						date, _ := time.Parse(layout, *n.LastChangeAaDate)
						data.SkdExpirationDate = date.Format(newLayout)
					}
				}
			}

			data.NpwpNo = *n.TinNumber
			data.NpwpRegistrationDate = ""
			data.CountryOfEstablishment = ""

			data.PlaceOfEstablishment = ""
			if n.EstablishedCity != nil {
				data.PlaceOfEstablishment = *n.EstablishedCity
			}

			data.DateOfEstablishment = ""
			if n.EstablishedDate != nil {
				date, _ := time.Parse(layout, *n.EstablishedDate)
				data.DateOfEstablishment = date.Format(newLayout)
			}

			data.ArticlesOfAssociationNo = ""
			if n.LastChangeAaNo != nil {
				data.ArticlesOfAssociationNo = *n.LastChangeAaNo
			}

			data.CompanyType = ""
			if n.IntitutionType != nil {
				if g, ok := gData[*n.IntitutionType]; ok {
					sof := *g.LkpCode
					data.CompanyType = sof
				}
			}

			data.CompanyCharacteristic = ""
			if n.IntitutionCharacteristic != nil {
				if g, ok := gData[*n.IntitutionCharacteristic]; ok {
					sof := *g.LkpCode
					data.CompanyCharacteristic = sof
				}
			}

			data.IncomeLevelIDR = ""
			if n.InstiAnnuallyIncome != nil {
				if g, ok := gData[*n.InstiAnnuallyIncome]; ok {
					sof := *g.LkpCode
					data.IncomeLevelIDR = sof
				}
			}

			data.InvestorsRiskProfile = ""
			if oareq.OaRiskLevel != nil {
				if g, ok := gData[*oareq.OaRiskLevel]; ok {
					sof := *g.LkpCode
					data.InvestorsRiskProfile = sof
				}
			}

			data.InvestmentObjective = ""
			if n.InstiInvestmentPurpose != nil {
				if g, ok := gData[*n.InstiInvestmentPurpose]; ok {
					sof := *g.LkpCode
					data.InvestmentObjective = sof
				}
			}

			data.SourceOfFund = ""
			if n.InstiSourceOfIncome != nil {
				if g, ok := gData[*n.InstiSourceOfIncome]; ok {
					sof := *g.LkpCode
					data.SourceOfFund = sof
				}
			}

			data.AssetOwner = "1"

			data.CompanyAddress = ""
			data.CompanyCityCode = ""
			data.CompanyCityName = ""
			data.CompanyPostalCode = ""

			data.CountryOfCompany = ""
			if n.Nationality != nil {
				if g, ok := countryData[*n.Nationality]; ok {
					data.CountryOfCompany = g.CouCode
				}
			}

			data.OfficePhone = ""
			if n.PhoneNo != nil {
				data.OfficePhone = *n.PhoneNo
			}

			data.Facsimile = ""
			if n.FaxNo != nil {
				data.Facsimile = *n.FaxNo
			}

			data.Email = ""
			if n.EmailAddress != nil {
				data.Email = *n.EmailAddress
			}

			data.StatementType = "2"

			//======================
			data.AuthorizedPerson1FirstName = ""
			data.AuthorizedPerson1MiddleName = ""
			data.AuthorizedPerson1LastName = ""
			data.AuthorizedPerson1Position = ""
			data.AuthorizedPerson1MobilePhone = ""
			data.AuthorizedPerson1Email = ""
			data.AuthorizedPerson1NpwpNo = ""
			data.AuthorizedPerson1KTPNo = ""
			data.AuthorizedPerson1KTPExpirationDate = ""
			data.AuthorizedPerson1PassportNo = ""
			data.AuthorizedPerson1PassportExpirationDate = ""

			data.AuthorizedPerson2FirstName = ""
			data.AuthorizedPerson2MiddleName = ""
			data.AuthorizedPerson2LastName = ""
			data.AuthorizedPerson2Position = ""
			data.AuthorizedPerson2MobilePhone = ""
			data.AuthorizedPerson2Email = ""
			data.AuthorizedPerson2NpwpNo = ""
			data.AuthorizedPerson2KTPNo = ""
			data.AuthorizedPerson2KTPExpirationDate = ""
			data.AuthorizedPerson2PassportNo = ""
			data.AuthorizedPerson2PassportExpirationDate = ""
			var authPerson []models.OaInstitutionAuthPersonDetail
			_, err = models.GetOaInstitutionAuthPersonRequest(&authPerson, strconv.FormatUint(oareq.OaRequestKey, 10))
			if err == nil && len(authPerson) > 0 {
				ky := 1
				for _, authP := range authPerson {
					if ky > 2 {
						break
					}
					sliceName := strings.Fields(*authP.FullName)
					fName := ""
					mName := ""
					lName := ""
					if len(sliceName) > 0 {
						if len(sliceName) == 1 {
							fName = sliceName[0]
							lName = sliceName[0]
						}
						if len(sliceName) == 2 {
							fName = sliceName[0]
							lName = sliceName[1]
						}
						if len(sliceName) > 2 {
							ln := len(sliceName)
							fName = sliceName[0]
							mName = sliceName[1]
							lastName := strings.Join(sliceName[2:ln], " ")
							lName = lastName
						}
					}
					noHp := ""
					if authP.PhoneNo != nil {
						noHp = *authP.PhoneNo
					}
					emailAuth := ""
					if authP.EmailAddress != nil {
						emailAuth = *authP.EmailAddress
					}
					posi := ""
					if authP.Position != nil {
						posi = *authP.Position
					}
					if ky == 1 {
						data.AuthorizedPerson1FirstName = fName
						data.AuthorizedPerson1MiddleName = mName
						data.AuthorizedPerson1LastName = lName
						data.AuthorizedPerson1Position = posi
						data.AuthorizedPerson1MobilePhone = noHp
						data.AuthorizedPerson1Email = emailAuth
						data.AuthorizedPerson1NpwpNo = ""

						if authP.IdcardType != nil {
							if *authP.IdcardType == uint64(11) { //KTP
								if authP.IdcardNo != nil {
									data.AuthorizedPerson1KTPNo = *authP.IdcardNo
								}
								if authP.IdcardNeverExpired == uint8(0) {
									if authP.IdcardExpiredDate != nil {
										date, _ := time.Parse(layout, *authP.IdcardExpiredDate)
										data.AuthorizedPerson1KTPExpirationDate = date.Format(newLayout)
									}
								}
							} else { //NON KTP
								if authP.IdcardNo != nil {
									data.AuthorizedPerson1PassportNo = *authP.IdcardNo
								}
								if authP.IdcardNeverExpired == uint8(0) {
									if authP.IdcardExpiredDate != nil {
										date, _ := time.Parse(layout, *authP.IdcardExpiredDate)
										data.AuthorizedPerson1PassportExpirationDate = date.Format(newLayout)
									}
								}
							}
						}
					}
					if ky == 2 {
						data.AuthorizedPerson2FirstName = fName
						data.AuthorizedPerson2MiddleName = mName
						data.AuthorizedPerson2LastName = lName
						data.AuthorizedPerson2Position = posi
						data.AuthorizedPerson2MobilePhone = noHp
						data.AuthorizedPerson2Email = emailAuth
						data.AuthorizedPerson2NpwpNo = ""
						if authP.IdcardType != nil {
							if *authP.IdcardType == uint64(11) { //KTP
								if authP.IdcardNo != nil {
									data.AuthorizedPerson2KTPNo = *authP.IdcardNo
								}
								if authP.IdcardNeverExpired == uint8(0) {
									if authP.IdcardExpiredDate != nil {
										date, _ := time.Parse(layout, *authP.IdcardExpiredDate)
										data.AuthorizedPerson2KTPExpirationDate = date.Format(newLayout)
									}
								}
							} else { //NON KTP
								if authP.IdcardNo != nil {
									data.AuthorizedPerson2PassportNo = *authP.IdcardNo
								}
								if authP.IdcardNeverExpired == uint8(0) {
									if authP.IdcardExpiredDate != nil {
										date, _ := time.Parse(layout, *authP.IdcardExpiredDate)
										data.AuthorizedPerson2PassportExpirationDate = date.Format(newLayout)
									}
								}
							}
						}
					}
					ky++
				}
			}
			mil100 := decimal.NewFromInt(100000000000)
			mil500 := decimal.NewFromInt(500000000000)
			mil1000 := decimal.NewFromInt(1000000000000)
			mil5000 := decimal.NewFromInt(5000000000000)

			data.AssetInformationforThePast3YearsIDRLastYear = ""
			if n.AssetY1 != nil {
				if n.AssetY1.Cmp(mil100) == -1 || n.AssetY1.Cmp(mil100) == 0 {
					data.AssetInformationforThePast3YearsIDRLastYear = "1"
				} else if (n.AssetY1.Cmp(mil100) == 1 && n.AssetY1.Cmp(mil500) == -1) || n.AssetY1.Cmp(mil500) == 0 {
					data.AssetInformationforThePast3YearsIDRLastYear = "2"
				} else if (n.AssetY1.Cmp(mil500) == 1 && n.AssetY1.Cmp(mil1000) == -1) || n.AssetY1.Cmp(mil1000) == 0 {
					data.AssetInformationforThePast3YearsIDRLastYear = "3"
				} else if (n.AssetY1.Cmp(mil1000) == 1 && n.AssetY1.Cmp(mil5000) == -1) || n.AssetY1.Cmp(mil5000) == 0 {
					data.AssetInformationforThePast3YearsIDRLastYear = "4"
				} else if n.AssetY1.Cmp(mil5000) == 1 {
					data.AssetInformationforThePast3YearsIDRLastYear = "5"
				}
			}

			data.AssetInformationforThePast3YearsIDR2YearsAgo = ""
			if n.AssetY2 != nil {
				if n.AssetY2.Cmp(mil100) == -1 || n.AssetY2.Cmp(mil100) == 0 {
					data.AssetInformationforThePast3YearsIDR2YearsAgo = "1"
				} else if (n.AssetY2.Cmp(mil100) == 1 && n.AssetY2.Cmp(mil500) == -1) || n.AssetY2.Cmp(mil500) == 0 {
					data.AssetInformationforThePast3YearsIDR2YearsAgo = "2"
				} else if (n.AssetY2.Cmp(mil500) == 1 && n.AssetY2.Cmp(mil1000) == -1) || n.AssetY2.Cmp(mil1000) == 0 {
					data.AssetInformationforThePast3YearsIDR2YearsAgo = "3"
				} else if (n.AssetY2.Cmp(mil1000) == 1 && n.AssetY2.Cmp(mil5000) == -1) || n.AssetY2.Cmp(mil5000) == 0 {
					data.AssetInformationforThePast3YearsIDR2YearsAgo = "4"
				} else if n.AssetY2.Cmp(mil5000) == 1 {
					data.AssetInformationforThePast3YearsIDR2YearsAgo = "5"
				}
			}

			data.AssetInformationforThePast3YearsIDR3YearsAgo = ""
			if n.AssetY3 != nil {
				if n.AssetY3.Cmp(mil100) == -1 || n.AssetY3.Cmp(mil100) == 0 {
					data.AssetInformationforThePast3YearsIDR3YearsAgo = "1"
				} else if (n.AssetY3.Cmp(mil100) == 1 && n.AssetY3.Cmp(mil500) == -1) || n.AssetY3.Cmp(mil500) == 0 {
					data.AssetInformationforThePast3YearsIDR3YearsAgo = "2"
				} else if (n.AssetY3.Cmp(mil500) == 1 && n.AssetY3.Cmp(mil1000) == -1) || n.AssetY3.Cmp(mil1000) == 0 {
					data.AssetInformationforThePast3YearsIDR3YearsAgo = "3"
				} else if (n.AssetY3.Cmp(mil1000) == 1 && n.AssetY3.Cmp(mil5000) == -1) || n.AssetY3.Cmp(mil5000) == 0 {
					data.AssetInformationforThePast3YearsIDR3YearsAgo = "4"
				} else if n.AssetY3.Cmp(mil5000) == 1 {
					data.AssetInformationforThePast3YearsIDR3YearsAgo = "5"
				}
			}

			data.ProfitInformationforThePast3YearsIDRLastYear = ""
			if n.OpsProfitY1 != nil {
				if n.OpsProfitY1.Cmp(mil100) == -1 || n.OpsProfitY1.Cmp(mil100) == 0 {
					data.ProfitInformationforThePast3YearsIDRLastYear = "1"
				} else if (n.OpsProfitY1.Cmp(mil100) == 1 && n.OpsProfitY1.Cmp(mil500) == -1) || n.OpsProfitY1.Cmp(mil500) == 0 {
					data.ProfitInformationforThePast3YearsIDRLastYear = "2"
				} else if (n.OpsProfitY1.Cmp(mil500) == 1 && n.OpsProfitY1.Cmp(mil1000) == -1) || n.OpsProfitY1.Cmp(mil1000) == 0 {
					data.ProfitInformationforThePast3YearsIDRLastYear = "3"
				} else if (n.OpsProfitY1.Cmp(mil1000) == 1 && n.OpsProfitY1.Cmp(mil5000) == -1) || n.OpsProfitY1.Cmp(mil5000) == 0 {
					data.ProfitInformationforThePast3YearsIDRLastYear = "4"
				} else if n.OpsProfitY1.Cmp(mil5000) == 1 {
					data.ProfitInformationforThePast3YearsIDRLastYear = "5"
				}
			}

			data.ProfitInformationforThePast3YearsIDR2YearsAgo = ""
			if n.OpsProfitY2 != nil {
				if n.OpsProfitY2.Cmp(mil100) == -1 || n.OpsProfitY2.Cmp(mil100) == 0 {
					data.ProfitInformationforThePast3YearsIDR2YearsAgo = "1"
				} else if (n.OpsProfitY2.Cmp(mil100) == 1 && n.OpsProfitY2.Cmp(mil500) == -1) || n.OpsProfitY2.Cmp(mil500) == 0 {
					data.ProfitInformationforThePast3YearsIDR2YearsAgo = "2"
				} else if (n.OpsProfitY2.Cmp(mil500) == 1 && n.OpsProfitY2.Cmp(mil1000) == -1) || n.OpsProfitY2.Cmp(mil1000) == 0 {
					data.ProfitInformationforThePast3YearsIDR2YearsAgo = "3"
				} else if (n.OpsProfitY2.Cmp(mil1000) == 1 && n.OpsProfitY2.Cmp(mil5000) == -1) || n.OpsProfitY2.Cmp(mil5000) == 0 {
					data.ProfitInformationforThePast3YearsIDR2YearsAgo = "4"
				} else if n.OpsProfitY2.Cmp(mil5000) == 1 {
					data.ProfitInformationforThePast3YearsIDR2YearsAgo = "5"
				}
			}

			data.ProfitInformationforThePast3YearsIDR3YearsAgo = ""
			if n.OpsProfitY3 != nil {
				if n.OpsProfitY3.Cmp(mil100) == -1 || n.OpsProfitY3.Cmp(mil100) == 0 {
					data.ProfitInformationforThePast3YearsIDR3YearsAgo = "1"
				} else if (n.OpsProfitY3.Cmp(mil100) == 1 && n.OpsProfitY3.Cmp(mil500) == -1) || n.OpsProfitY3.Cmp(mil500) == 0 {
					data.ProfitInformationforThePast3YearsIDR3YearsAgo = "2"
				} else if (n.OpsProfitY3.Cmp(mil500) == 1 && n.OpsProfitY3.Cmp(mil1000) == -1) || n.OpsProfitY3.Cmp(mil1000) == 0 {
					data.ProfitInformationforThePast3YearsIDR3YearsAgo = "3"
				} else if (n.OpsProfitY3.Cmp(mil1000) == 1 && n.OpsProfitY3.Cmp(mil5000) == -1) || n.OpsProfitY3.Cmp(mil5000) == 0 {
					data.ProfitInformationforThePast3YearsIDR3YearsAgo = "4"
				} else if n.OpsProfitY3.Cmp(mil5000) == 1 {
					data.ProfitInformationforThePast3YearsIDR3YearsAgo = "5"
				}
			}

			data.FATCAStatus = ""
			data.TINForeignTIN = ""
			data.TINForeignTINIssuance_Country = ""
			data.GIIN = ""
			data.SubstantialUSOwnerName = ""
			data.SubstantialUSOwnerAddress = ""
			data.SubstantialUSOwnerTIN = ""

			data.REDMPaymentBankBICCode1 = ""
			data.REDMPaymentBankBIMemberCode1 = ""
			data.REDMPaymentBankName1 = ""
			data.REDMPaymentBankCountry1 = ""
			data.REDMPaymentBankBranch1 = ""
			data.REDMPaymentACCcy1 = ""
			data.REDMPaymentACNo1 = ""
			data.REDMPaymentACName1 = ""

			data.REDMPaymentBankBICCode2 = ""
			data.REDMPaymentBankBIMemberCode2 = ""
			data.REDMPaymentBankName2 = ""
			data.REDMPaymentBankCountry2 = ""
			data.REDMPaymentBankBranch2 = ""
			data.REDMPaymentACCcy2 = ""
			data.REDMPaymentACNo2 = ""
			data.REDMPaymentACName2 = ""

			data.REDMPaymentBankBICCode3 = ""
			data.REDMPaymentBankBIMemberCode3 = ""
			data.REDMPaymentBankName3 = ""
			data.REDMPaymentBankCountry3 = ""
			data.REDMPaymentBankBranch3 = ""
			data.REDMPaymentACCcy3 = ""
			data.REDMPaymentACNo3 = ""
			data.REDMPaymentACName3 = ""

			if strType == "127" {
				var bank []models.CustomerBankAccountSinvest
				_, err = models.GetCustomerBankAccountSinvest(&bank, strconv.FormatUint(*oareq.CustomerKey, 10))

				if err == nil && len(bank) > 0 {
					bk := 1
					for _, b := range bank {
						bankCode := ""
						if b.BankCode != nil {
							bankCode = *b.BankCode
						}
						swift := ""
						if b.SwiftCode != nil {
							swift = *b.SwiftCode
						}
						bankname := ""
						if b.BankName != nil {
							bankname = *b.BankName
						}
						branchName := ""
						if b.BranchName != nil {
							branchName = *b.BranchName
						}
						code := ""
						if b.Code != nil {
							code = *b.Code
						}
						accountNo := ""
						if b.AccountNo != nil {
							accountNo = *b.AccountNo
						}
						accountName := ""
						if b.AccountHolderName != nil {
							accountName = *b.AccountHolderName
						}

						if bk == 1 {
							data.REDMPaymentBankBICCode1 = bankCode
							data.REDMPaymentBankBIMemberCode1 = swift
							data.REDMPaymentBankName1 = bankname
							data.REDMPaymentBankCountry1 = "ID"
							data.REDMPaymentBankBranch1 = branchName
							data.REDMPaymentACCcy1 = code
							data.REDMPaymentACNo1 = accountNo
							data.REDMPaymentACName1 = accountName
						}

						if bk == 2 {
							data.REDMPaymentBankBICCode2 = bankCode
							data.REDMPaymentBankBIMemberCode2 = swift
							data.REDMPaymentBankName2 = bankname
							data.REDMPaymentBankCountry2 = "ID"
							data.REDMPaymentBankBranch2 = branchName
							data.REDMPaymentACCcy2 = code
							data.REDMPaymentACNo2 = accountNo
							data.REDMPaymentACName2 = accountName
						}

						if bk == 3 {
							data.REDMPaymentBankBICCode3 = bankCode
							data.REDMPaymentBankBIMemberCode3 = swift
							data.REDMPaymentBankName3 = bankname
							data.REDMPaymentBankCountry3 = "ID"
							data.REDMPaymentBankBranch3 = branchName
							data.REDMPaymentACCcy3 = code
							data.REDMPaymentACNo3 = accountNo
							data.REDMPaymentACName3 = accountName
						}
						bk++
					}
				}
			}

			data.ClientCode = ""

			//start update client_code if new customer
			if strType == "127" { //type NEW
				last = last + 1
				paramsCustomer := make(map[string]string)
				var convLast string
				convLast = strconv.FormatUint(uint64(last), 10)
				clientCode := lib.PadLeft(convLast, "0", 6)
				paramsCustomer["client_code"] = clientCode
				dateLayout := "2006-01-02 15:04:05"
				paramsCustomer["rec_modified_date"] = time.Now().Format(dateLayout)
				strKeyLogin := strconv.FormatUint(lib.Profile.UserID, 10)
				paramsCustomer["rec_modified_by"] = strKeyLogin
				strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
				paramsCustomer["customer_key"] = strCustomerKey
				_, err = models.UpdateMsCustomer(paramsCustomer)
				if err != nil {
					log.Error("Error update oa request")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}
				lastValue = paramsCustomer["client_code"]

				data.ClientCode = paramsCustomer["client_code"]
			}
			//end

			txtData := data.Type + "|" +
				data.SACode + "|" +
				data.SID + "|" +
				data.CompanyName + "|" +
				data.CountryOfDomicile + "|" +
				data.SiupNo + "|" +
				data.SiupExpirationDate + "|" +
				data.SkdNo + "|" +
				data.SkdExpirationDate + "|" +
				data.NpwpNo + "|" +
				data.NpwpRegistrationDate + "|" +
				data.CountryOfEstablishment + "|" +
				data.PlaceOfEstablishment + "|" +
				data.DateOfEstablishment + "|" +
				data.ArticlesOfAssociationNo + "|" +
				data.CompanyType + "|" +
				data.CompanyCharacteristic + "|" +
				data.IncomeLevelIDR + "|" +
				data.InvestorsRiskProfile + "|" +
				data.InvestmentObjective + "|" +
				data.SourceOfFund + "|" +
				data.AssetOwner + "|" +
				data.CompanyAddress + "|" +
				data.CompanyCityCode + "|" +
				data.CompanyCityName + "|" +
				data.CompanyPostalCode + "|" +
				data.CountryOfCompany + "|" +
				data.OfficePhone + "|" +
				data.Facsimile + "|" +
				data.Email + "|" +
				data.StatementType + "|" +
				data.AuthorizedPerson1FirstName + "|" +
				data.AuthorizedPerson1MiddleName + "|" +
				data.AuthorizedPerson1LastName + "|" +
				data.AuthorizedPerson1Position + "|" +
				data.AuthorizedPerson1MobilePhone + "|" +
				data.AuthorizedPerson1Email + "|" +
				data.AuthorizedPerson1NpwpNo + "|" +
				data.AuthorizedPerson1KTPNo + "|" +
				data.AuthorizedPerson1KTPExpirationDate + "|" +
				data.AuthorizedPerson1PassportNo + "|" +
				data.AuthorizedPerson1PassportExpirationDate + "|" +
				data.AuthorizedPerson2FirstName + "|" +
				data.AuthorizedPerson2MiddleName + "|" +
				data.AuthorizedPerson2LastName + "|" +
				data.AuthorizedPerson2Position + "|" +
				data.AuthorizedPerson2MobilePhone + "|" +
				data.AuthorizedPerson2Email + "|" +
				data.AuthorizedPerson2NpwpNo + "|" +
				data.AuthorizedPerson2KTPNo + "|" +
				data.AuthorizedPerson2KTPExpirationDate + "|" +
				data.AuthorizedPerson2PassportNo + "|" +
				data.AuthorizedPerson2PassportExpirationDate + "|" +
				data.AssetInformationforThePast3YearsIDRLastYear + "|" +
				data.AssetInformationforThePast3YearsIDR2YearsAgo + "|" +
				data.AssetInformationforThePast3YearsIDR3YearsAgo + "|" +
				data.ProfitInformationforThePast3YearsIDRLastYear + "|" +
				data.ProfitInformationforThePast3YearsIDR2YearsAgo + "|" +
				data.ProfitInformationforThePast3YearsIDR3YearsAgo + "|" +
				data.FATCAStatus + "|" +
				data.TINForeignTIN + "|" +
				data.TINForeignTINIssuance_Country + "|" +
				data.GIIN + "|" +
				data.SubstantialUSOwnerName + "|" +
				data.SubstantialUSOwnerAddress + "|" +
				data.SubstantialUSOwnerTIN + "|" +
				data.REDMPaymentBankBICCode1 + "|" +
				data.REDMPaymentBankBIMemberCode1 + "|" +
				data.REDMPaymentBankName1 + "|" +
				data.REDMPaymentBankCountry1 + "|" +
				data.REDMPaymentBankBranch1 + "|" +
				data.REDMPaymentACCcy1 + "|" +
				data.REDMPaymentACNo1 + "|" +
				data.REDMPaymentACName1 + "|" +
				data.REDMPaymentBankBICCode2 + "|" +
				data.REDMPaymentBankBIMemberCode2 + "|" +
				data.REDMPaymentBankName2 + "|" +
				data.REDMPaymentBankCountry2 + "|" +
				data.REDMPaymentBankBranch2 + "|" +
				data.REDMPaymentACCcy2 + "|" +
				data.REDMPaymentACNo2 + "|" +
				data.REDMPaymentACName2 + "|" +
				data.REDMPaymentBankBICCode3 + "|" +
				data.REDMPaymentBankBIMemberCode3 + "|" +
				data.REDMPaymentBankName3 + "|" +
				data.REDMPaymentBankCountry3 + "|" +
				data.REDMPaymentBankBranch3 + "|" +
				data.REDMPaymentACCcy3 + "|" +
				data.REDMPaymentACNo3 + "|" +
				data.REDMPaymentACName3 + "|" +
				data.ClientCode

			var txt models.OaRequestCsvFormatFiksTxt
			txt.DataRow = txtData

			responseData = append(responseData, txt)
		}
	}

	//value awal = 009995 ----------------- update app_config
	if lastValue != "" {
		paramsConfig := make(map[string]string)
		paramsConfig["app_config_value"] = lastValue
		dateLayout := "2006-01-02 15:04:05"
		paramsConfig["rec_modified_date"] = time.Now().Format(dateLayout)
		strKeyLogin := strconv.FormatUint(lib.Profile.UserID, 10)
		paramsConfig["rec_modified_by"] = strKeyLogin
		_, err = models.UpdateMsCustomerByConfigCode(paramsConfig, "LAST_CLIENT_CODE")
		if err != nil {
			log.Error("Error update App Config")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}
	}
	//end

	if len(responseData) > 0 {
		paramsUpdate := make(map[string]string)

		strOaStatus := strconv.FormatUint(261, 10) //customer build, proses upload data to Sinvest
		paramsUpdate["oa_status"] = strOaStatus
		dateLayout := "2006-01-02 15:04:05"
		paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)
		strKey := strconv.FormatUint(lib.Profile.UserID, 10)
		paramsUpdate["rec_modified_by"] = strKey

		_, err = models.UpdateOaRequestByKeyIn(paramsUpdate, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error("Error update oa request")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UploadOaInstitutionRequestFormatSinvest(c echo.Context) error {

	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error

	err = os.MkdirAll(config.BasePath+"/oa_upload/sinvest/", 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("file")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			log.Println(extension)
			roles := []string{".txt", ".TXT"}
			_, found := lib.Find(roles, extension)
			if !found {
				return lib.CustomError(http.StatusUnauthorized, "Format file must .txt", "Format file must .txt")
			}
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename+extension)
			// Upload txt and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/oa_upload/sinvest/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}

			fileTxt, err := os.Open(config.BasePath + "/oa_upload/sinvest/" + filename + extension)

			if err != nil {
				log.Println("failed to open txt")
				log.Println(err)
				// log.Fatalf("failed to open")
			}

			scanner := bufio.NewScanner(fileTxt)

			scanner.Split(bufio.ScanLines)
			var text []string

			for scanner.Scan() {
				text = append(text, scanner.Text())
			}

			fileTxt.Close()

			dateLayout := "2006-01-02 15:04:05"

			var customerIds []string
			for idx, ea := range text {
				if idx > 0 {

					s := strings.Split(ea, "|")

					sidNo := strings.TrimSpace(s[2])
					ifuaNo := strings.TrimSpace(s[3])
					ifuaName := strings.TrimSpace(s[4])
					clientCode := strings.TrimSpace(s[5])

					//get ms_customer by clientCode
					var customer models.MsCustomer
					_, err := models.GetMsCustomerByClientCode(&customer, clientCode)
					if err != nil {
						log.Error("get customer error : client_code = " + clientCode + ". Error : " + err.Error())
						continue
					}

					strCusKey := strconv.FormatUint(customer.CustomerKey, 10)
					if _, ok := lib.Find(customerIds, strCusKey); !ok {
						customerIds = append(customerIds, strCusKey)
					}

					//update ms_customer
					params := make(map[string]string)
					params["sid_no"] = sidNo
					params["customer_key"] = strCusKey
					params["rec_modified_date"] = time.Now().Format(dateLayout)
					params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					_, err = models.UpdateMsCustomer(params)
					if err != nil {
						log.Error("Error update sid_no ms_customer")
						continue
					}

					//update tr_account_all
					paramsTrAccount := make(map[string]string)
					paramsTrAccount["ifua_no"] = ifuaNo
					paramsTrAccount["ifua_name"] = ifuaName
					paramsTrAccount["rec_modified_date"] = time.Now().Format(dateLayout)
					paramsTrAccount["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					_, err = models.UpdateTrAccountUploadSinvest(paramsTrAccount, "customer_key", strCusKey)
					if err != nil {
						log.Error("Error update ifua_no, ifua_name tr_account")
						continue
					}
				}
			}

			//update oa_status di oa_request by customer_key
			if len(customerIds) > 0 {
				paramsOa := make(map[string]string)
				paramsOa["oa_status"] = "262"
				paramsOa["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsOa["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				_, err := models.UpdateOaRequestByFieldIn(paramsOa, customerIds, "customer_key")
				if err != nil {
					log.Error("Error update oa_status in oa_request : " + err.Error())
				}
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}
