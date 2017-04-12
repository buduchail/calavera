package rest

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"net/http"

	"github.com/buduchail/go-skeleton/interfaces"
)

var (
	unknownErr = errors.New("Unknown error")

	routers = map[string]string{
		"n": "nethttp",
		"i": "iris",
		"h": "httprouter",
		"e": "echo",
		"f": "fasthttp",
		"g": "gin",
		"r": "go-restful",
	}
)

func NewApi(prefix, router string) interfaces.RestAPI {

	switch router {
	case "n", routers["n"]:
		return NewNetHTTP(prefix)
	case "i", routers["i"]:
		return NewIris(prefix)
	case "h", routers["h"]:
		return NewHttpRouter(prefix)
	case "e", routers["e"]:
		return NewEcho(prefix)
	case "f", routers["f"]:
		return NewFast(prefix)
	case "g", routers["g"]:
		return NewGin(prefix)
	case "r", routers["r"]:
		return NewGoRestful(prefix)
	}

	return nil
}

func getHttpError(code int) error {

	status := http.StatusText(code)
	if status != "" {
		return errors.New(status)
	}

	return unknownErr
}

func normalizePrefix(prefix string) string {
	normalized := strings.TrimLeft(strings.TrimRight(prefix, "/"), "/")
	switch normalized {
	case "", "/":
		return "/"
	default:
		return "/" + normalized + "/"
	}
}

func expandPath(path, idTemplate string) (fullPath string, parentIds []string, idParam string) {
	var i = 0
	parts := strings.Split(path, "/*/")
	parentIds = make([]string, 0, len(parts))
	fullPath = parts[0]
	l := len(parts)
	if l > 1 {
		for i = range parts[1:] {
			id := "id" + strconv.Itoa(i+1)
			parentIds = append(parentIds, id)
			fullPath += "/" + fmt.Sprintf(idTemplate, id) + "/" + parts[i+1]
		}
		i += 1
	}
	return fullPath, parentIds, "id" + strconv.Itoa(i+1)
}
