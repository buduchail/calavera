package rest

import (
	"errors"
	"strconv"
	"net/http"
	"github.com/valyala/fasthttp"
	"github.com/buduchail/calavera"
)

type (
	FastAPI struct {
		root      *pathHandler
		prefix    string
		prefixLen int
	}
)

func NewFast(prefix string) (api FastAPI) {
	api = FastAPI{}
	api.prefix = normalizePrefix(prefix)
	api.prefixLen = len(api.prefix)
	api.root = NewPathHandler(api.prefix)
	return api
}

func (api FastAPI) getBody(ctx *fasthttp.RequestCtx) calavera.Payload {
	return ctx.Request.Body()
}

func (api FastAPI) getQueryParameters(ctx *fasthttp.RequestCtx) calavera.QueryParameters {
	params := calavera.QueryParameters{}
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = []string{string(value)}
	})
	return params
}

func (api FastAPI) sendResponse(ctx *fasthttp.RequestCtx, code int, body calavera.Payload, err error) error {

	if code == http.StatusOK {
		_, err = ctx.Write(body)
	} else {
		if err == nil {
			err = getHttpError(code)
		}
		ctx.Error(err.Error(), code)
	}

	return err
}

func (api FastAPI) handleResource(method string, id calavera.ResourceID, parentIds []calavera.ResourceID, ctx *fasthttp.RequestCtx, handler calavera.ResourceHandler) (code int, body calavera.Payload, err error) {

	switch method {
	case "POST":
		if id != "" {
			return http.StatusBadRequest, calavera.EmptyBody, errors.New("POST requests must not provide an ID")
		}
		return handler.Post(parentIds, api.getBody(ctx))
	case "GET":
		if id != "" {
			return handler.Get(id, parentIds)
		} else {
			return handler.GetMany(parentIds, api.getQueryParameters(ctx))
		}
	case "PUT":
		if id == "" {
			return http.StatusBadRequest, calavera.EmptyBody, errors.New("PUT method must provide an ID")
		}
		return handler.Put(id, parentIds, api.getBody(ctx))
	case "DELETE":
		if id == "" {
			return http.StatusBadRequest, calavera.EmptyBody, errors.New("DELETE method must provide an ID")
		}
		return handler.Delete(id, parentIds)
	}

	return http.StatusMethodNotAllowed, calavera.EmptyBody, errors.New("Method not allowed")
}

func (api FastAPI) handle(ctx *fasthttp.RequestCtx) {

	path := string(ctx.Request.URI().Path())

	if len(path) > api.prefixLen+1 {

		handler, id, parentIds := api.root.findHandler(path[api.prefixLen:])
		if handler == nil {
			api.sendResponse(ctx, http.StatusNotFound, calavera.EmptyBody, nil)
			return
		}

		code, body, err := api.handleResource(string(ctx.Method()), id, parentIds, ctx, handler)
		api.sendResponse(ctx, code, body, err)

	} else {
		api.sendResponse(ctx, http.StatusNotFound, calavera.EmptyBody, nil)
	}
}

func (api FastAPI) AddResource(name string, handler calavera.ResourceHandler) {
	api.root.addHandler(name, handler)
}

func (api FastAPI) AddMiddleware(m calavera.Middleware) {
	// NOT IMPLEMENTED
}

func (api FastAPI) Run(port int) {

	fasthttp.ListenAndServe(":"+strconv.Itoa(port), api.handle)
}
