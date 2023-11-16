package controllers

import (
	"mf-bo-api/lib"
	"mf-bo-api/models"
	"net/http"
)

func IsMenuAccessAllowed(menu_key int) error {
	accessType := IsAccessAllowed(menu_key)
	if accessType.AccessEnable {
		return nil
	} else {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
}

func IsAccessAllowed(menu_key int) models.MenuPageAccess {

	var accessType models.MenuPageAccess
	accessType.AccessEnable = false
	accessType.AllowedCreate = false
	accessType.AllowedDelete = false
	accessType.AllowedRead = false
	accessType.AllowedUpdate = false
	accessType.AllowedApproval = false

	role_key := lib.Profile.RoleKey
	var access models.RoleMenuDetail
	_, err := models.GetRoleMenuAccess(&access, role_key, uint64(menu_key))
	if err == nil && role_key > 0 && menu_key > 0 {
		accessType.AccessEnable = access.AccessEnable
		accessType.AllowedCreate = access.AllowedCreate
		accessType.AllowedDelete = access.AllowedDelete
		accessType.AllowedRead = access.AllowedRead
		accessType.AllowedUpdate = access.AllowedUpdate
		accessType.AllowedApproval = access.AllowedApproval
	} else if menu_key == 888888888 { //jika ingin di set true semua tanpa lewat db
		accessType.AccessEnable = true
		accessType.AllowedCreate = true
		accessType.AllowedDelete = true
		accessType.AllowedRead = true
		accessType.AllowedUpdate = true
		accessType.AllowedApproval = true
	}
	return accessType
}
