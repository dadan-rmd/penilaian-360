package userModel

import (
	"central-auth/internal/app/model/helperModel"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

/* Table Definition */
type User struct {
	Id                                  int64  `json:"id"`
	Username                            string `json:"username"`
	Email                               string `json:"email"`
	Password                            string `json:"-"`
	helperModel.DateAuditModelTimeStamp `json:"-"`
}

func (User) TableName() string {
	return "users"
}

type CreateUserReq struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Platform []string `json:"platform"`
}

type ResLogin struct {
	User  User
	Token string `json:"token"`
}

const (
	MsgSuccessUpdateUser = "user berhasil diperbaharui"
	MsgSuccessDeleteUser = "user berhasil dihapus"
	MsgFailedAddUser     = "user gagal di tambahkan"
)

func validPlatform() validation.Rule {
	validPlatforms := map[string]struct{}{
		"cms kbr id":      {},
		"cms kbr prime":   {},
		"prime analytics": {},
	}

	return validation.By(func(value interface{}) error {
		platforms, ok := value.([]string)
		if !ok {
			return fmt.Errorf("invalid type: expected []string")
		}

		for _, platform := range platforms {
			normalized := strings.ToLower(strings.TrimSpace(platform))
			if _, found := validPlatforms[normalized]; !found {
				return fmt.Errorf("Invalid platform value: %s", platform)
			}
		}
		return nil
	})
}

func (m *CreateUserReq) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Invalid email format"),
		),
		validation.Field(&m.Password,
			validation.Required.Error("Password is required"),
		),
		validation.Field(&m.Username,
			validation.Required.Error("Username is required"),
		),
		validation.Field(&m.Platform,
			validation.Required.Error("Platform is required"),
			validPlatform(),
		),
	)
}
