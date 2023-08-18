package providers

import (
	"net/http"
	"queue-manager/internal/structures"
)

type RouterProviderInterface interface {
	Get(url string, handler http.Handler)
	Post(url string, handler http.Handler)
	GetRoutes() []structures.Route
}

type RouterProvider struct {
	routes []structures.Route
}

func (rp *RouterProvider) Get(url string, handler http.Handler) {
	rp.routes = append(rp.routes, structures.Route{Url: url, Handler: handler})
}

func (rp *RouterProvider) Post(url string, handler http.Handler) {
	rp.routes = append(rp.routes, structures.Route{Url: url, Handler: handler})
}

func (rp *RouterProvider) GetRoutes() []structures.Route {
	return rp.routes
}

func NewRouterProvider() RouterProviderInterface {
	return &RouterProvider{}
}
