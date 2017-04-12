package rest

import (
	"bytes"
	"errors"
	"strconv"
	"io/ioutil"
	"net/http"

	"github.com/buduchail/catrina"
)

type (
	NetHTTP struct {
		root       *pathHandler
		prefix     string
		prefixLen  int
		middleware []catrina.Middleware
	}
)

func NewNetHTTP(prefix string) (api *NetHTTP) {
	api = &NetHTTP{}
	api.prefix = normalizePrefix(prefix)
	api.prefixLen = len(api.prefix)
	api.root = NewPathHandler(api.prefix)
	api.middleware = make([]catrina.Middleware, 0)
	return api
}

func (api *NetHTTP) getBody(r *http.Request) catrina.Payload {
	b, _ := ioutil.ReadAll(r.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api *NetHTTP) getQueryParameters(r *http.Request) catrina.QueryParameters {
	return catrina.QueryParameters(r.URL.Query())
}

func (api *NetHTTP) sendResponse(w http.ResponseWriter, code int, body catrina.Payload, err error) error {

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

func (api *NetHTTP) handleResource(method string, id catrina.ResourceID, parentIds []catrina.ResourceID, r *http.Request, handler catrina.ResourceHandler) (code int, body catrina.Payload, err error) {

	switch method {
	case "OPTIONS":
		return  handler.Options()
	case "POST":
		if id != "" {
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("POST requests must not provide an ID")
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
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("PUT method must provide an ID")
		}
		return handler.Put(id, parentIds, api.getBody(r))
	case "DELETE":
		if id == "" {
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("DELETE method must provide an ID")
		}
		return handler.Delete(id, parentIds)
	}

	return http.StatusMethodNotAllowed, catrina.EmptyBody, errors.New("Method not allowed")
}

func (api *NetHTTP) handle(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path) > api.prefixLen+1 {

		handler, id, parentIds := api.root.findHandler(r.URL.Path[api.prefixLen:])
		if handler == nil {
			api.sendResponse(w, http.StatusNotFound, catrina.EmptyBody, nil)
			return
		}

		// apply middleware
		for _, m := range api.middleware {
			err := m.Handle(w, r)
			if err != nil {
				api.sendResponse(w, http.StatusInternalServerError, catrina.EmptyBody, *err)
				return
			}
		}

		code, body, err := api.handleResource(r.Method, id, parentIds, r, handler)
		api.sendResponse(w, code, body, err)

	} else {
		api.sendResponse(w, http.StatusNotFound, catrina.EmptyBody, nil)
	}
}

func (api *NetHTTP) AddResource(name string, handler catrina.ResourceHandler) {
	api.root.addHandler(name, handler)
}

func (api *NetHTTP) AddMiddleware(m catrina.Middleware) {
	api.middleware = append(api.middleware, m)
}

func (api *NetHTTP) Run(port int) {

	mux := http.NewServeMux()

	mux.HandleFunc(api.prefix, api.handle)
	mux.HandleFunc(api.prefix+"/", api.handle)

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
