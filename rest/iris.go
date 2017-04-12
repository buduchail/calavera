package rest

import (
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"github.com/buduchail/calavera"
)

type (
	IrisAPI struct {
		i      *iris.Framework
		prefix string
	}
)

func NewIris(prefix string) (api IrisAPI) {
	api = IrisAPI{}
	api.i = iris.New()
	api.i.Adapt(httprouter.New())
	api.prefix = normalizePrefix(prefix)
	return api
}

func (api IrisAPI) getBody(c *iris.Context) calavera.Payload {
	b, _ := ioutil.ReadAll(c.Request.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api IrisAPI) getQueryParameters(c *iris.Context) calavera.QueryParameters {
	params := calavera.QueryParameters{}
	if nil != c.URLParamsAsMulti() {
		for k, v := range c.URLParamsAsMulti() {
			params[k] = v
		}
	}
	return params
}

func (api IrisAPI) getParentIds(c *iris.Context, idParams []string) (ids []calavera.ResourceID) {
	ids = make([]calavera.ResourceID, 0)
	for _, id := range idParams {
		// prepend: /grandparent/1/parent/2/child/3 -> [2,1]
		ids = append([]calavera.ResourceID{calavera.ResourceID(c.Param(id))}, ids...)
	}
	return ids
}

func (api IrisAPI) sendResponse(c *iris.Context, code int, body calavera.Payload, err error) {
	if code != http.StatusOK || err != nil {
		c.EmitError(code)
	} else {
		c.Writef(string(body))
	}
}

func (api IrisAPI) AddResource(name string, handler calavera.ResourceHandler) {

	path, parentIdParams, idParam := expandPath(name, ":%s")

	postRoute := func(c *iris.Context) {
		code, body, err := handler.Post(
			[]calavera.ResourceID{},
			api.getBody(c),
		)
		api.sendResponse(c, code, body, err)
	}

	getRoute := func(c *iris.Context) {
		code, body, err := handler.Get(
			calavera.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
		)
		api.sendResponse(c, code, body, err)
	}

	getManyRoute := func(c *iris.Context) {
		code, body, err := handler.GetMany(
			api.getParentIds(c, parentIdParams),
			api.getQueryParameters(c),
		)
		api.sendResponse(c, code, body, err)
	}

	putRoute := func(c *iris.Context) {
		code, body, err := handler.Put(
			calavera.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
			api.getBody(c),
		)
		api.sendResponse(c, code, body, err)
	}

	deleteRoute := func(c *iris.Context) {
		code, body, err := handler.Delete(
			calavera.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
		)
		api.sendResponse(c, code, body, err)
	}

	fullPath := api.prefix + path

	api.i.Post(fullPath, postRoute)
	api.i.Post(fullPath+"/", postRoute)

	api.i.Get(fullPath+"/:"+idParam, getRoute)
	api.i.Get(fullPath+"", getManyRoute)
	api.i.Get(fullPath+"/", getManyRoute)

	api.i.Put(fullPath+"/:"+idParam, putRoute)

	api.i.Delete(fullPath+"/:"+idParam, deleteRoute)
}

func (api IrisAPI) AddMiddleware(m calavera.Middleware) {
	// NOT IMPLEMENTED
}

func (api IrisAPI) Run(port int) {
	api.i.Listen(":" + strconv.Itoa(port))
}
