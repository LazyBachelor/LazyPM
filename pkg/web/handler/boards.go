package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func BoardsHandler(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards" {
			handleNotFound(w, r)
			return
		}

		issues, err := svc.Beads.AllIssues(r.Context())
		if err != nil {
			http.Error(w, errRetrieveIssues, http.StatusInternalServerError)
			return
		}

		todo, inProgress, done := models.GroupIssuesByStatus(issues)
		models.SortIssuesByPriority(todo, false)
		models.SortIssuesByPriority(inProgress, false)
		models.SortIssuesByPriority(done, false)

		props := routes.BoardsProps{
			Todo:       todo,
			InProgress: inProgress,
			Done:       done,
		}
		routes.Boards(props).Render(r.Context(), w)
	}
}
