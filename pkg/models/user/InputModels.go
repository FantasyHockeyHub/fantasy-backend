package user

type SignUpInput struct {
	Nickname string `json:"nickname" binding:"required,min=4,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Code     int    `json:"code" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email,max=64"`
}

type RefreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type NicknameInput struct {
	Nickname string `json:"nickname" binding:"required,min=4,max=64"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required,min=8,max=64"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=64"`
}

type ResetPasswordInput struct {
	Hash        string `json:"hash" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=64"`
}
