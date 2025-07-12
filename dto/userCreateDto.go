package dto

type UserCreateDto struct {
	Username   string `form:"username" binding:"required,min=5,max=30"`
	Email      string `form:"email" binding:"required,email"`
	Password   string `form:"password" binding:"required,min=6,max=30"`
	RePassword string `form:"repassword"  binding:"required,min=6,max=30"`
}
