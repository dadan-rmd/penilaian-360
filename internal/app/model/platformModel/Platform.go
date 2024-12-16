package platformModel

import (
	"central-auth/internal/app/model/helperModel"
)

/* Table Definition */
type Platform struct {
	Id                                  int64  `json:"id"`
	UserId                              int64  `json:"user_id"`
	Name                                string `json:"name"`
	helperModel.DateAuditModelTimeStamp `json:"-"`
}

func (Platform) TableName() string {
	return "platforms"
}

const (
	MsgSuccessUpdatePlatform = "platform berhasil diperbaharui"
	MsgSuccessDeletePlatform = "platform berhasil dihapus"
	MsgFailedAddPlatform     = "platform gagal di tambahkan"
)
