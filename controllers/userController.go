package controllers

import (
	"log"
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"

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

/*
build menu on side bar
*/
func GetMFBOMenus(c echo.Context) error {
	var user_role_Key uint16 = 15
	var ret_data models.SideBarMenuModel
	_, err := models.GetRootMenu(&ret_data, "MFBO", user_role_Key)
	if err != nil {
		log.Println("Fail to load root menu")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Fail to load root menu")
	}

	var mnu_parent []models.ScMenuModel
	_, err = models.GetMenuTree(&mnu_parent, user_role_Key, ret_data.RootMenuKey)
	if err != nil {
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Fail to load root menu Tree")
	}

	if len(mnu_parent) > 0 {
		ret_data.MenuList = nil
		for _, _row := range mnu_parent {
			var _mnu models.ScMenuModel
			_mnu.MenuKey = _row.MenuKey
			_mnu.MenuParent = _row.MenuParent
			_mnu.MenuCode = _row.MenuCode
			_mnu.MenuName = _row.MenuName
			_mnu.MenuPage = _row.MenuPage
			_mnu.MenuUrl = _row.MenuUrl

			_mnu.MenuDesc = _row.MenuDesc
			_mnu.RecOrder = _row.RecOrder
			_mnu.MenuIcon = _row.MenuIcon

			var mnu_child []models.ScMenuModel
			_, err := models.GetMenuTree(&mnu_child, user_role_Key, _row.MenuKey)
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
						_child_row.MenuUrl = _child.MenuUrl
						_child_row.MenuDesc = _child.MenuDesc
						_child_row.RecOrder = _child.RecOrder
						_child_row.MenuIcon = _child.MenuIcon
						arr_child = append(arr_child, _child_row)
					}
					_mnu.ChildList = arr_child
				} else {
					log.Println("node menu has no child")
					_mnu.ChildList = nil
				}
			} else {
				log.Println("node menu has no child")
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
