package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func IndexHandler(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handleNotFound(w, r)
			return
		}

		issues, err := svc.Beads.AllIssues(r.Context())
		if err != nil {
			http.Error(w, errRetrieveIssues, http.StatusInternalServerError)
			return
		}

		openCount, closedCount := models.CountByStatus(issues)
		recentCount := len(models.RecentIssuesByID(issues, 5))
		recentIssues := models.RecentIssuesByID(issues, 15)
		models.SortIssuesByPriority(recentIssues, true)

		calendarRange, calendarDays, prevURL, nextURL := buildCalendarWeek(r.URL.Query().Get("week"))

		props := routes.IndexProps{
			OpenCount:       openCount,
			ClosedCount:     closedCount,
			AssignedCount:   openCount,
			DueSoonCount:    0,
			RecentCount:     recentCount,
			TodayCount:      0,
			WeekCount:       0,
			TeamMembers:     []routes.TeamMember{},
			RecentIssues:    recentIssues,
			CalendarRange:   calendarRange,
			CalendarDays:    calendarDays,
			CalendarPrevURL: prevURL,
			CalendarNextURL: nextURL,
		}

		if r.Header.Get("HX-Request") == "true" {
			routes.DashboardWidgets(props).Render(r.Context(), w)
			return
		}
		routes.Index(props).Render(r.Context(), w)
	}
}

func NewDashboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboards/new" {
			handleNotFound(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func buildCalendarWeek(weekParam string) (string, []routes.CalendarDay, string, string) {
	now := time.Now()
	var start time.Time

	if weekParam != "" {
		if t, err := time.Parse("2006-01-02", weekParam); err == nil {
			weekday := int(t.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			start = t.AddDate(0, 0, -(weekday - 1))
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		} else {
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			start = now.AddDate(0, 0, -(weekday - 1))
		}
	} else {
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = now.AddDate(0, 0, -(weekday - 1))
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())

	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var calendarDays []routes.CalendarDay
	for i := 0; i < 7; i++ {
		d := start.AddDate(0, 0, i)
		label := fmt.Sprintf("%02d %s", d.Day(), dayNames[d.Weekday()])
		calendarDays = append(calendarDays, routes.CalendarDay{
			Date:      d.Format("2006-01-02"),
			Label:     label,
			IsToday:   d.YearDay() == now.YearDay() && d.Year() == now.Year(),
			HasEvents: false,
		})
	}

	rangeStr := fmt.Sprintf("%s %02d – %02d", start.Month().String(), start.Day(), start.AddDate(0, 0, 6).Day())

	prevStart := start.AddDate(0, 0, -7)
	nextStart := start.AddDate(0, 0, 7)
	prevURL := fmt.Sprintf("/?week=%s", prevStart.Format("2006-01-02"))
	nextURL := fmt.Sprintf("/?week=%s", nextStart.Format("2006-01-02"))

	return rangeStr, calendarDays, prevURL, nextURL
}
