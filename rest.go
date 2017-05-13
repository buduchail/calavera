package catrina

import "net/http"

var (
	EmptyBody = Payload([]byte(""))
)

type (
	RestAPI interface {
		AddResource(name string, handler ResourceHandler)
		AddMiddleware(m Middleware)
		Run(port int)
	}

	ResourceHandler interface {
		Options() (
			code int, body Payload, err error,
		)
		Post(parentIds []string, payload Payload) (
			code int, body Payload, err error,
		)
		Get(id string, parentIds []string) (
			code int, body Payload, err error,
		)
		GetMany(parentIds []string, query QueryParameters) (
			code int, body Payload, err error,
		)
		Put(id string, parentIds []string, payload Payload) (
			code int, body Payload, err error,
		)
		Delete(id string, parentIds []string) (
			code int, body Payload, err error,
		)
	}

	// Some syntactic sugar
	Payload []byte
	QueryParameters map[string][]string

	Middleware interface {
		Handle(w http.ResponseWriter, r *http.Request) *error
	}
)
