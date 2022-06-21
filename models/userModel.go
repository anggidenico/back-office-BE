package models

type User struct {
	ID          uint64  `db:"user_login_key"		json:"id"`
	UserName    *string `db:"ulogin_name"      	json:"username"`
	Password    *string `db:"ulogin_password"  	json:"password"`
	DisplayName *string `db:"ulogin_full_name"  json:"ulogin_full_name"`
	Role        *string `db:"role_key"         	json:"role"`
}
