package evaluationModel

type EvaluationAnswer struct {
	Id                   int64  `json:"id"`
	EvaluationId         int64  `json:"evaluation_id"`
	EvaluationEmployeeId int64  `json:"evaluation_employee_id"`
	QuestionId           int64  `json:"question_id"`
	Answer               string `json:"answer"`
	FinalPoint           int    `json:"final_point"`
}
