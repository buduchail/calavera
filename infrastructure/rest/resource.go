package rest

import (
	"net/http"
	"github.com/buduchail/go-skeleton/interfaces"
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

func (s ResourceHandler) Options() (code int, body interfaces.Payload, err error) {
	return http.StatusOK, interfaces.EmptyBody, nil
}

func (s ResourceHandler) Post(parentIds []interfaces.ResourceID, payload interfaces.Payload) (code int, body interfaces.Payload, err error) {
	return http.StatusMethodNotAllowed, interfaces.EmptyBody, nil
}

func (s ResourceHandler) Get(id interfaces.ResourceID, parentIds []interfaces.ResourceID) (code int, body interfaces.Payload, err error) {
	return http.StatusMethodNotAllowed, interfaces.EmptyBody, nil
}

func (s ResourceHandler) GetMany(parentIds []interfaces.ResourceID, params interfaces.QueryParameters) (code int, body interfaces.Payload, err error) {
	return http.StatusMethodNotAllowed, interfaces.EmptyBody, nil
}

func (s ResourceHandler) Put(id interfaces.ResourceID, parentIds []interfaces.ResourceID, payload interfaces.Payload) (code int, body interfaces.Payload, err error) {
	return http.StatusMethodNotAllowed, interfaces.EmptyBody, nil
}

func (s ResourceHandler) Delete(id interfaces.ResourceID, parentIds []interfaces.ResourceID) (code int, body interfaces.Payload, err error) {
	return http.StatusMethodNotAllowed, interfaces.EmptyBody, nil
}
