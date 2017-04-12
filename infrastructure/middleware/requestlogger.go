package middleware

import (
	"net/http"
	"github.com/buduchail/go-skeleton/interfaces"
)

type (
	RequestLogger struct {
		logger    interfaces.Logger
		logHeader string
	}
)

func NewRequestLogger(logger interfaces.Logger, correlationIdHeader string) *RequestLogger {
	return &RequestLogger{logger, correlationIdHeader}
}

func (m RequestLogger) Handle(w http.ResponseWriter, r *http.Request) (err *error) {

	m.logger.Info(
		r.Method+" "+r.URL.String(),
		&interfaces.LoggerContext{m.logHeader: r.Header[m.logHeader]},
	)

	return
}
