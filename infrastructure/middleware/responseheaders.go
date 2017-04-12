package middleware

import (
	"net/http"
)

type (
	ResponseHeaders struct {
		headers map[string]string
	}
)

func NewResponseHeaders(headers map[string]string) *ResponseHeaders {
	return &ResponseHeaders{headers: headers}
}

func (m ResponseHeaders) Handle(w http.ResponseWriter, r *http.Request) (err *error) {

	for k, v := range m.headers {
		w.Header().Set(k, v)
	}

	return
}
