package helperModel

type UserAuditModel struct {
	CreatedBy int `json:"created_by"`
	UpdatedBy int `json:"updated_by"`
	DeletedBy int `json:"deleted_by"`
}
