package rest

import (
	"net/http"
	"github.com/buduchail/calavera"
)

type (
	// Base implementation of ResourceHandler interface, to be used by
	// RestAPI instances. This implementation can be used by concrete
	// resource handlers to provide default behaviour for HTTP verbs
	// that are not implemented by the handler. By embedding this type,
	// concrete handlers will only need to implement methods for verbs
	// they handle.
	ResourceHandler struct {
	}
)

func (s ResourceHandler) Options() (code int, body calavera.Payload, err error) {
	return http.StatusOK, calavera.EmptyBody, nil
}

func (s ResourceHandler) Post(parentIds []calavera.ResourceID, payload calavera.Payload) (code int, body calavera.Payload, err error) {
	return http.StatusMethodNotAllowed, calavera.EmptyBody, nil
}

func (s ResourceHandler) Get(id calavera.ResourceID, parentIds []calavera.ResourceID) (code int, body calavera.Payload, err error) {
	return http.StatusMethodNotAllowed, calavera.EmptyBody, nil
}

func (s ResourceHandler) GetMany(parentIds []calavera.ResourceID, params calavera.QueryParameters) (code int, body calavera.Payload, err error) {
	return http.StatusMethodNotAllowed, calavera.EmptyBody, nil
}

func (s ResourceHandler) Put(id calavera.ResourceID, parentIds []calavera.ResourceID, payload calavera.Payload) (code int, body calavera.Payload, err error) {
	return http.StatusMethodNotAllowed, calavera.EmptyBody, nil
}

func (s ResourceHandler) Delete(id calavera.ResourceID, parentIds []calavera.ResourceID) (code int, body calavera.Payload, err error) {
	return http.StatusMethodNotAllowed, calavera.EmptyBody, nil
}
