package middleware

import (
	"net/http"
	"github.com/buduchail/calavera"
)

type (
	RequestLogger struct {
		logger    calavera.Logger
		logHeader string
	}
)

func NewRequestLogger(logger calavera.Logger, correlationIdHeader string) *RequestLogger {
	return &RequestLogger{logger, correlationIdHeader}
}

func (m RequestLogger) Handle(w http.ResponseWriter, r *http.Request) (err *error) {

	m.logger.Info(
		r.Method+" "+r.URL.String(),
		&calavera.LoggerContext{m.logHeader: r.Header[m.logHeader]},
	)

	return
}
