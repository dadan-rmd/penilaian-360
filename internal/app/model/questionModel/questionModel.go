package questionModel

type Question struct {
	Id           int64  `json:"id"`
	EvaluationId int64  `json:"evaluation_id"`
	Question     string `json:"question"`
	Type         string `json:"type"`
}
