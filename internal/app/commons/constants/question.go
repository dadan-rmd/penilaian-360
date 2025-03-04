package constants

type QuestionType string
type TypeOfCompetency string

const (
	QuestionTypeRate           = QuestionType("rate")
	QuestionTypeEssay          = QuestionType("essay")
	TypeOfCompetencyFunctional = TypeOfCompetency("functional")
	TypeOfCompetencyPersonal   = TypeOfCompetency("personal")
)
