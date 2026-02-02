package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

type Route struct {
	Pattern string
	Handler http.Handler
}

func GetRoutes(svc *service.Services) []Route {
	var routes []Route

	routes = append(routes, PagesRoutes(svc)...)
	routes = append(routes, IssuesRoutes(svc)...)

	return routes
}
