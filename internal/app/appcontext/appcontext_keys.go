package appcontext

type clientName string
type externalID string

const (
	// CtxKeyClientName saving clientName from basic AUTH
	CtxKeyClient clientName = "Ctx-Client"

	CtxExtID externalID = "Ctx-External-ID"
)

func (c clientName) String() string {
	return string(c)
}

func (e externalID) String() string {
	return string(e)
}
