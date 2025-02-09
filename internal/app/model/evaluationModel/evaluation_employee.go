package evaluationModel

import (
	"penilaian-360/internal/app/commons/constants"
	"time"
)

const (
	ActionDraft  = "draft"
	ActionSubmit = "submit"
)

type EvaluationEmployee struct {
	Id           int64                  `json:"id"`
	EvaluationId int64                  `json:"evaluation_id"`
	EmployeeId   int64                  `json:"employee_id"`
	Type         constants.EmployeeType `json:"type"`
	Avg          float64                `json:"avg"`
	EmailSent    string                 `json:"email_sent"`
	CreatedAt    time.Time              `json:"created_at" gorm:"autoCreateTime"`
}
