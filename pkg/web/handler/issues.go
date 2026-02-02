package handler

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"net/http"
)

type IssuesHandler struct {
	Services *service.Services
}

func NewIssuesHandler(services *service.Services) *IssuesHandler {
	return &IssuesHandler{
		Services: services,
	}
}

func (h *IssuesHandler) HandleGetAllIssues() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
