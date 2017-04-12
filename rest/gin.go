package rest

import (
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/buduchail/catrina"
)

type (
	GinAPI struct {
		g      *gin.Engine
		prefix string
	}
)

func NewGin(prefix string) (api GinAPI) {
	gin.SetMode(gin.ReleaseMode)
	api = GinAPI{}
	api.g = gin.New()
	api.prefix = normalizePrefix(prefix)
	return api
}

func (api GinAPI) getBody(c *gin.Context) catrina.Payload {
	b, _ := ioutil.ReadAll(c.Request.Body)
	return bytes.NewBuffer(b).Bytes()
}

func (api GinAPI) getQueryParameters(c *gin.Context) catrina.QueryParameters {
	return catrina.QueryParameters(c.Request.URL.Query())
}

func (api GinAPI) getParentIds(c *gin.Context, idParams []string) (ids []catrina.ResourceID) {
	ids = make([]catrina.ResourceID, 0)
	for _, id := range idParams {
		// prepend: /grandparent/1/parent/2/child/3 -> [2,1]
		ids = append([]catrina.ResourceID{catrina.ResourceID(c.Param(id))}, ids...)
	}
	return ids
}

func (api GinAPI) sendResponse(c *gin.Context, code int, body catrina.Payload, err error) {

	if code != http.StatusOK || err != nil {
		if err == nil {
			err = getHttpError(code)
		}
		c.String(code, err.Error())
	} else {
		c.String(code, string(body))
	}
}

func (api GinAPI) AddResource(name string, handler catrina.ResourceHandler) {

	path, parentIdParams, idParam := expandPath(name, ":%s")

	postRoute := func(c *gin.Context) {
		code, body, err := handler.Post(
			api.getParentIds(c, parentIdParams),
			api.getBody(c),
		)
		api.sendResponse(c, code, body, err)
	}

	getRoute := func(c *gin.Context) {
		code, body, err := handler.Get(
			catrina.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
		)
		api.sendResponse(c, code, body, err)
	}

	getManyRoute := func(c *gin.Context) {
		code, body, err := handler.GetMany(
			api.getParentIds(c, parentIdParams),
			api.getQueryParameters(c),
		)
		api.sendResponse(c, code, body, err)
	}

	putRoute := func(c *gin.Context) {
		code, body, err := handler.Put(
			catrina.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
			api.getBody(c),
		)
		api.sendResponse(c, code, body, err)
	}

	deleteRoute := func(c *gin.Context) {
		code, body, err := handler.Delete(
			catrina.ResourceID(c.Param(idParam)),
			api.getParentIds(c, parentIdParams),
		)
		api.sendResponse(c, code, body, err)
	}

	fullPath := api.prefix + path

	api.g.POST(fullPath, postRoute)
	api.g.POST(fullPath+"/", postRoute)

	api.g.GET(fullPath+"/:"+idParam, getRoute)
	api.g.GET(fullPath, getManyRoute)
	api.g.GET(fullPath+"/", getManyRoute)

	api.g.PUT(fullPath+"/:"+idParam, putRoute)

	api.g.DELETE(fullPath+"/:"+idParam, deleteRoute)
}

func (api GinAPI) AddMiddleware(m catrina.Middleware) {
	// NOT IMPLEMENTED
}

func (api GinAPI) Run(port int) {
	api.g.Run(":" + strconv.Itoa(port))
}
