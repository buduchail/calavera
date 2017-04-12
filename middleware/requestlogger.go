package middleware

import (
	"net/http"
	"github.com/buduchail/catrina"
)

type (
	RequestLogger struct {
		logger    catrina.Logger
		logHeader string
	}
)

func NewRequestLogger(logger catrina.Logger, correlationIdHeader string) *RequestLogger {
	return &RequestLogger{logger, correlationIdHeader}
}

func (m RequestLogger) Handle(w http.ResponseWriter, r *http.Request) (err *error) {

	m.logger.Info(
		r.Method+" "+r.URL.String(),
		&catrina.LoggerContext{m.logHeader: r.Header[m.logHeader]},
	)

	return
}
