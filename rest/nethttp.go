package rest

import (
	"bytes"
	"errors"
	"strconv"
	"io/ioutil"
	"net/http"

	"github.com/buduchail/go-skeleton/interfaces"
)

type (
	NetHTTP struct {
		root       *pathHandler
		prefix     string
		prefixLen  int
		middleware []interfaces.Middleware
	}
)

func NewNetHTTP(prefix string) (api *NetHTTP) {
	api = &NetHTTP{}
	api.prefix = normalizePrefix(prefix)
	api.prefixLen = len(api.prefix)
	api.root = NewPathHandler(api.prefix)
	api.middleware = make([]interfaces.Middleware, 0)
	return api
}

func (api *NetHTTP) getBody(r *http.Request) interfaces.Payload {
	b, _ := ioutil.ReadAll(r.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api *NetHTTP) getQueryParameters(r *http.Request) interfaces.QueryParameters {
	return interfaces.QueryParameters(r.URL.Query())
}

func (api *NetHTTP) sendResponse(w http.ResponseWriter, code int, body interfaces.Payload, err error) error {

	if code == http.StatusOK {
		_, err = w.Write(body)
	} else {
		if err == nil {
			err = getHttpError(code)
		}
		http.Error(w, err.Error(), code)
	}

	return err
}

func (api *NetHTTP) handleResource(method string, id interfaces.ResourceID, parentIds []interfaces.ResourceID, r *http.Request, handler interfaces.ResourceHandler) (code int, body interfaces.Payload, err error) {

	switch method {
	case "OPTIONS":
		return  handler.Options()
	case "POST":
		if id != "" {
			return http.StatusBadRequest, interfaces.EmptyBody, errors.New("POST requests must not provide an ID")
		}
		return handler.Post(parentIds, api.getBody(r))
	case "GET":
		if id != "" {
			return handler.Get(id, parentIds)
		} else {
			return handler.GetMany(parentIds, api.getQueryParameters(r))
		}
	case "PUT":
		if id == "" {
			return http.StatusBadRequest, interfaces.EmptyBody, errors.New("PUT method must provide an ID")
		}
		return handler.Put(id, parentIds, api.getBody(r))
	case "DELETE":
		if id == "" {
			return http.StatusBadRequest, interfaces.EmptyBody, errors.New("DELETE method must provide an ID")
		}
		return handler.Delete(id, parentIds)
	}

	return http.StatusMethodNotAllowed, interfaces.EmptyBody, errors.New("Method not allowed")
}

func (api *NetHTTP) handle(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path) > api.prefixLen+1 {

		handler, id, parentIds := api.root.findHandler(r.URL.Path[api.prefixLen:])
		if handler == nil {
			api.sendResponse(w, http.StatusNotFound, interfaces.EmptyBody, nil)
			return
		}

		// apply middleware
		for _, m := range api.middleware {
			err := m.Handle(w, r)
			if err != nil {
				api.sendResponse(w, http.StatusInternalServerError, interfaces.EmptyBody, *err)
				return
			}
		}

		code, body, err := api.handleResource(r.Method, id, parentIds, r, handler)
		api.sendResponse(w, code, body, err)

	} else {
		api.sendResponse(w, http.StatusNotFound, interfaces.EmptyBody, nil)
	}
}

func (api *NetHTTP) AddResource(name string, handler interfaces.ResourceHandler) {
	api.root.addHandler(name, handler)
}

func (api *NetHTTP) AddMiddleware(m interfaces.Middleware) {
	api.middleware = append(api.middleware, m)
}

func (api *NetHTTP) Run(port int) {

	mux := http.NewServeMux()

	mux.HandleFunc(api.prefix, api.handle)
	mux.HandleFunc(api.prefix+"/", api.handle)

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
