package authRequest

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email,max=50"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
}
