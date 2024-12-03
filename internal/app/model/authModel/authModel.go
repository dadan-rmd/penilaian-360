package authModel

import (
	"central-auth/internal/app/model/helperModel"
	"central-auth/internal/app/model/userModel"
)

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginSSOReq struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Image string `json:"image"`
}

type LoginRes struct {
	User  userModel.User `json:"user"`
	Token string         `json:"token"`
}

type ChangePasswordReq struct {
	User        userModel.User `json:"-"`
	OldPassword string         `json:"old_password" binding:"required"`
	NewPassword string         `json:"new_password" binding:"required"`
}

type ChangePasswordRes struct {
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required"`
}

type ForgotPasswordRes struct {
}

type VerifyForgotPasswordReq struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type VerifyForgotPasswordRes struct {
}

type GetUserFromJWTAndRoleReq struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type GetUserFromJWTAndRoleRes struct {
	User userModel.User
}

type TokenValidityReq struct {
	Token string `json:"token"`
}

type TokenValidityRes struct {
	IsValid bool `json:"is_valid"`
}

type ChangePasswordFromForgotPassReq struct {
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ChangePasswordFromForgotPassRes struct {
}

type TokenSession struct {
	ID        int64  `json:"id"`
	TokenUID  string `json:"token_uid"`
	IsRevoked bool   `json:"is_revoked"`
	helperModel.DateAuditModel
	helperModel.UserAuditModel
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}
