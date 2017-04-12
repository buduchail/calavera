package rest

import (
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/buduchail/go-skeleton/interfaces"
)

type (
	HttpRouterAPI struct {
		r      *httprouter.Router
		prefix string
	}
)

func NewHttpRouter(prefix string) (api HttpRouterAPI) {
	api = HttpRouterAPI{}
	api.r = httprouter.New()
	api.prefix = normalizePrefix(prefix)
	return api
}

func (api HttpRouterAPI) getBody(r *http.Request) interfaces.Payload {
	b, _ := ioutil.ReadAll(r.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api HttpRouterAPI) getQueryParameters(r *http.Request) interfaces.QueryParameters {
	return interfaces.QueryParameters(r.URL.Query())
}

func (api HttpRouterAPI) getParentIds(ps httprouter.Params, idParams []string) (ids []interfaces.ResourceID) {
	ids = make([]interfaces.ResourceID, 0)
	for _, id := range idParams {
		// prepend: /grandparent/1/parent/2/child/3 -> [2,1]
		ids = append([]interfaces.ResourceID{interfaces.ResourceID(ps.ByName(id))}, ids...)
	}
	return ids
}

func (api HttpRouterAPI) sendResponse(w http.ResponseWriter, code int, body interfaces.Payload, err error) {
	if code != http.StatusOK || err != nil {
		if err == nil {
			err = getHttpError(code)
		}
		http.Error(w, err.Error(), code)
	} else {
		w.Write(body)
	}
}

func (api HttpRouterAPI) AddResource(name string, handler interfaces.ResourceHandler) {

	path, parentIdParams, idParam := expandPath(name, ":%s")

	postRoute := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code, body, err := handler.Post(
			api.getParentIds(ps, parentIdParams),
			api.getBody(r),
		)
		api.sendResponse(w, code, body, err)
	}

	getRoute := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code, body, err := handler.Get(
			interfaces.ResourceID(ps.ByName(idParam)),
			api.getParentIds(ps, parentIdParams),
		)
		api.sendResponse(w, code, body, err)
	}

	getManyRoute := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code, body, err := handler.GetMany(
			api.getParentIds(ps, parentIdParams),
			api.getQueryParameters(r),
		)
		api.sendResponse(w, code, body, err)
	}

	putRoute := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code, body, err := handler.Put(
			interfaces.ResourceID(ps.ByName(idParam)),
			api.getParentIds(ps, parentIdParams),
			api.getBody(r),
		)
		api.sendResponse(w, code, body, err)
	}

	deleteRoute := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		code, body, err := handler.Delete(
			interfaces.ResourceID(ps.ByName(idParam)),
			api.getParentIds(ps, parentIdParams),
		)
		api.sendResponse(w, code, body, err)
	}

	fullPath := api.prefix + path

	api.r.POST(fullPath, postRoute)
	api.r.POST(fullPath+"/", postRoute)

	api.r.GET(fullPath+"/:"+idParam, getRoute)
	api.r.GET(fullPath+"", getManyRoute)
	api.r.GET(fullPath+"/", getManyRoute)

	api.r.PUT(fullPath+"/:"+idParam, putRoute)

	api.r.DELETE(fullPath+"/:"+idParam, deleteRoute)
}

func (api HttpRouterAPI) AddMiddleware(m interfaces.Middleware) {
	// NOT IMPLEMENTED
}

func (api HttpRouterAPI) Run(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), api.r)
}
