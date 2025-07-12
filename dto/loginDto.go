package dto

type LoginDto struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email"`
	Password string `form:"password" binding:"required"`
}
