package ctxkeys

type contextKey string
type logger string

var (
	// CtxRequestID context key for X-Request-ID
	CtxRequestID contextKey = "X-Request-ID"

	// CtxLogger context key for logger
	CtxLogger logger = "Ctx-Logger"
)

func (c contextKey) String() string {
	return string(c)
}

func (c logger) String() string {
	return string(c)
}
