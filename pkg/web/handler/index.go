package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	svc := Services(r)

	issues, err := svc.Beads.AllIssues(r.Context())
	if err != nil {
		http.Error(w, "failed to retrieve issues", http.StatusInternalServerError)
		return
	}

	props := routes.IndexProps{
		IssueTable: components.IssueTableProps{
			Issues: issues,
		},
	}

	routes.Index(props).Render(r.Context(), w)
}
