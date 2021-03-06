package rest

import (
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/emicklei/go-restful"
	"github.com/buduchail/catrina"
)

type (
	GoRestfulAPI struct {
		container *restful.Container
		prefix    string
	}
)

func NewGoRestful(prefix string) (api GoRestfulAPI) {
	api = GoRestfulAPI{}
	api.container = restful.NewContainer()
	api.prefix = normalizePrefix(prefix)
	return api
}

func (api GoRestfulAPI) getBody(rq *restful.Request) catrina.Payload {
	b, _ := ioutil.ReadAll(rq.Request.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api GoRestfulAPI) getQueryParameters(rq *restful.Request) catrina.QueryParameters {
	return catrina.QueryParameters(rq.Request.URL.Query())
}

func (api GoRestfulAPI) getParentIds(rq *restful.Request, idParams []string) (ids []string) {
	ids = make([]string, 0)
	for _, id := range idParams {
		// prepend: /grandparent/1/parent/2/child/3 -> [2,1]
		ids = append([]string{rq.Request.URL.Query().Get(id)}, ids...)
	}
	return ids
}

func (api GoRestfulAPI) sendResponse(rp *restful.Response, code int, body catrina.Payload, err error) {

	if code != http.StatusOK || err != nil {
		if err == nil {
			err = getHttpError(code)
		}
		rp.WriteErrorString(code, err.Error())
	} else {
		rp.Write(body)
	}
}

func (api GoRestfulAPI) AddResource(name string, handler catrina.ResourceHandler) {

	path, parentIdParams, idParam := expandPath(name, "{%s}")

	postRoute := func(rq *restful.Request, rp *restful.Response) {
		code, body, err := handler.Post(
			api.getParentIds(rq, parentIdParams),
			api.getBody(rq),
		)
		api.sendResponse(rp, code, body, err)
	}

	getRoute := func(rq *restful.Request, rp *restful.Response) {
		code, body, err := handler.Get(
			rq.PathParameter("id"),
			api.getParentIds(rq, parentIdParams),
		)
		api.sendResponse(rp, code, body, err)
	}

	getManyRoute := func(rq *restful.Request, rp *restful.Response) {
		code, body, err := handler.GetMany(
			api.getParentIds(rq, parentIdParams),
			api.getQueryParameters(rq),
		)
		api.sendResponse(rp, code, body, err)
	}

	putRoute := func(rq *restful.Request, rp *restful.Response) {
		code, body, err := handler.Put(
			rq.PathParameter("id"),
			api.getParentIds(rq, parentIdParams),
			api.getBody(rq),
		)
		api.sendResponse(rp, code, body, err)
	}

	deleteRoute := func(rq *restful.Request, rp *restful.Response) {
		code, body, err := handler.Delete(
			rq.PathParameter("id"),
			api.getParentIds(rq, parentIdParams),
		)
		api.sendResponse(rp, code, body, err)
	}

	ws := new(restful.WebService)

	ws.Path(api.prefix + path + "/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(postRoute))
	ws.Route(ws.POST("/").To(postRoute))

	ws.Route(ws.GET("/{" + idParam + "}").To(getRoute))
	ws.Route(ws.GET("").To(getManyRoute))
	ws.Route(ws.GET("/").To(getManyRoute))

	ws.Route(ws.PUT("/{" + idParam + "}").To(putRoute))

	ws.Route(ws.DELETE("/{" + idParam + "}").To(deleteRoute))

	api.container.Add(ws)
}

func (api GoRestfulAPI) AddMiddleware(m catrina.Middleware) {
	// NOT IMPLEMENTED
}

func (api GoRestfulAPI) Run(port int) {
	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: api.container}
	server.ListenAndServe()
}
