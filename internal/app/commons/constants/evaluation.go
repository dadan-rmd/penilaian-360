package constants

type EvaluationStatus string

const (
	EvaluationStatusDraft   = EvaluationStatus("draft")
	EvaluationStatusPublish = EvaluationStatus("publish")
)
