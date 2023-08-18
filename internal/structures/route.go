package structures

import "net/http"

type Route struct {
	Url     string
	Handler http.Handler
}
