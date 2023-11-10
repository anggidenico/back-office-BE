package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func GetUserInfo(c echo.Context) error {

	var user models.User

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = user

	return c.JSON(http.StatusOK, response)
}

func GetMMMenus(c echo.Context) error {

	var ret_data models.SideBarMenuModel
	_, err := models.GetRootMenu(&ret_data)
	if err != nil {
		log.Println("Fail to load root menu")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Fail to load root menu")
	}

	var mnu_parent []models.ScMenuModel
	_, err = models.GetMenuTree(&mnu_parent, ret_data.RootMenuKey)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Fail to load root menu Tree")
	}

	var base_url string = ret_data.RootHomeURL
	if ret_data.RootMenuPage != nil {
		base_url = ret_data.RootHomeURL + *ret_data.RootMenuPage
	}
	ret_data.RootHomeURL = base_url
	base_url = ret_data.RootHomeURL

	if len(mnu_parent) > 0 {
		ret_data.MenuList = nil
		for _, _row := range mnu_parent {
			var _mnu models.ScMenuModel
			_mnu.MenuKey = _row.MenuKey
			_mnu.MenuParent = _row.MenuParent
			_mnu.MenuCode = _row.MenuCode
			_mnu.MenuName = _row.MenuName
			//_mnu.MenuPage = _row.MenuPage

			var str_folder string = ""
			if _row.MenuPage != nil && strings.Trim(*_row.MenuPage, " ") != "#" {
				str_folder = strings.Trim(*_row.MenuPage, " ")
				page_full_url := base_url + str_folder
				_mnu.MenuUrl = page_full_url
			} else {
				_mnu.MenuUrl = "#"
			}
			_mnu.MenuDesc = _row.MenuDesc
			_mnu.RecOrder = _row.RecOrder
			_mnu.MenuIcon = _row.MenuIcon

			var mnu_child []models.ScMenuModel
			_, err := models.GetMenuTree(&mnu_child, _row.MenuKey)
			if err == nil {
				if len(mnu_child) > 0 {
					var arr_child []models.ScMenuModel
					for _, _child := range mnu_child {
						var _child_row models.ScMenuModel
						_child_row.MenuKey = _child.MenuKey
						_child_row.MenuParent = _child.MenuParent
						_child_row.MenuCode = _child.MenuCode
						_child_row.MenuName = _child.MenuName
						_child_row.MenuPage = _child.MenuPage
						if _row.MenuPage != nil {
							str_folder = strings.Trim(*_child.MenuPage, " ")
						}
						page_full_url := base_url + str_folder
						_child_row.MenuUrl = page_full_url
						_child_row.MenuDesc = _child.MenuDesc
						_child_row.RecOrder = _child.RecOrder
						_child_row.MenuIcon = _child.MenuIcon
						arr_child = append(arr_child, _child_row)
					}
					_mnu.ChildList = arr_child
				} else {
					log.Println("Fail to load child menu")
					_mnu.ChildList = nil
				}
			} else {
				log.Println("Fail to load child menu")
				_mnu.ChildList = nil
			}

			ret_data.MenuList = append(ret_data.MenuList, _mnu)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ret_data

	return c.JSON(http.StatusOK, response)

}
