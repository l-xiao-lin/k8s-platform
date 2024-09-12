package model

type ParamSignup struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type User struct {
	UserID   int64  `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Token    string `db:"token" json:"token"`
}

type ParamLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
