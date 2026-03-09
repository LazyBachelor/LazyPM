package components

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
)

type IssueDetailProps struct {
	Issue *models.Issue
}

func IssueDetail(props IssueDetailProps) templ.Component {
	return templruntime.GeneratedTemplate(func(in templruntime.GeneratedComponentInput) (err error) {
		w, ctx := in.Writer, in.Context
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		buf, isBuffer := templruntime.GetBuffer(w)
		if !isBuffer {
			defer func() {
				bufErr := templruntime.ReleaseBuffer(buf)
				if err == nil {
					err = bufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		children := templ.GetChildren(ctx)
		if children == nil {
			children = templ.NopComponent
		}
		_ = children
		ctx = templ.ClearChildren(ctx)

		if props.Issue == nil {
			return templruntime.WriteString(buf, 1, "<div class=\"alert alert-error\">Issue not found</div>")
		}

		// Title
		if err = templruntime.WriteString(buf, 2, "<div class=\"space-y-3\"><h2 class=\"text-2xl font-bold m-2\">"); err != nil {
			return err
		}
		var title string
		title, err = templ.JoinStringErrs(props.Issue.Title)
		if err != nil {
			return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 15, Col: 3}
		}
		if _, err = buf.WriteString(templ.EscapeString(title)); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 3, "</h2>"); err != nil {
			return err
		}

		// Badges row.
		if err = templruntime.WriteString(buf, 4, "<div class=\"flex gap-2 flex-wrap\">"); err != nil {
			return err
		}
		if err = StatusBadge(props.Issue.Status).Render(ctx, buf); err != nil {
			return err
		}
		if err = TypeBadge(props.Issue.IssueType).Render(ctx, buf); err != nil {
			return err
		}
		if err = PriorityBadge(props.Issue.Priority).Render(ctx, buf); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 5, "</div>"); err != nil {
			return err
		}

		// Description, if present.
		if props.Issue.Description != "" {
			if err = templruntime.WriteString(buf, 6, "<div class=\"mt-4\"><h3 class=\"text-sm font-semibold opacity-70 mb-1\">Description</h3><p class=\"whitespace-pre-wrap text-sm\">"); err != nil {
				return err
			}
			var desc string
			desc, err = templ.JoinStringErrs(props.Issue.Description)
			if err != nil {
				return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 25, Col: 4}
			}
			if _, err = buf.WriteString(templ.EscapeString(desc)); err != nil {
				return err
			}
			if err = templruntime.WriteString(buf, 7, "</p></div>"); err != nil {
				return err
			}
		}

		if err = templruntime.WriteString(buf, 8, "<div class=\"divider my-2\"></div><div class=\"grid grid-cols-2 gap-2 text-sm\">"); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 9, "<div class=\"bg-base-200 rounded-lg p-2\"><span class=\"text-xs opacity-70 block\">Created</span><p>"); err != nil {
			return err
		}
		var created string
		created, err = templ.JoinStringErrs(props.Issue.CreatedAt.Format("Jan 2, 2006"))
		if err != nil {
			return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 32, Col: 4}
		}
		if _, err = buf.WriteString(templ.EscapeString(created)); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 10, "</p></div>"); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 11, "<div class=\"bg-base-200 rounded-lg p-2\"><span class=\"text-xs opacity-70 block\">Updated</span><p>"); err != nil {
			return err
		}
		var updated string
		updated, err = templ.JoinStringErrs(props.Issue.UpdatedAt.Format("Jan 2, 2006"))
		if err != nil {
			return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 36, Col: 4}
		}
		if _, err = buf.WriteString(templ.EscapeString(updated)); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 12, "</p></div>"); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 13, "<div class=\"flex flex-row gap-1 bg-base-200 rounded-lg p-2\"><span class=\"text-xs opacity-70 block\">Created by</span><p class=\"text-xs opacity-70\">"); err != nil {
			return err
		}
		var createdBy string
		createdBy, err = templ.JoinStringErrs(props.Issue.CreatedBy)
		if err != nil {
			return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 40, Col: 4}
		}
		if _, err = buf.WriteString(templ.EscapeString(createdBy)); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 14, "</p></div>"); err != nil {
			return err
		}

		if err = templruntime.WriteString(buf, 15, "<a class=\"btn btn-primary btn-sm flex flex-row gap-1 items-center justify-center\" hx-get=\""); err != nil {
			return err
		}
		var assigneeURL string
		assigneeURL, err = templ.JoinStringErrs("/issues/" + props.Issue.ID + "/assignee")
		if err != nil {
			return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 44, Col: 4}
		}
		if _, err = buf.WriteString(templ.EscapeString(assigneeURL)); err != nil {
			return err
		}
		if err = templruntime.WriteString(buf, 16, "\" hx-target=\"#modal-container\" hx-swap=\"innerHTML\">"); err != nil {
			return err
		}
		if props.Issue.Assignee == "" {
			if err = templruntime.WriteString(buf, 17, "<span>Unassigned (click to assign) </span>"); err != nil {
				return err
			}
		} else {
			if err = templruntime.WriteString(buf, 18, "<span>Assignee: "); err != nil {
				return err
			}
			var assignee string
			assignee, err = templ.JoinStringErrs(props.Issue.Assignee)
			if err != nil {
				return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 51, Col: 5}
			}
			if _, err = buf.WriteString(templ.EscapeString(assignee)); err != nil {
				return err
			}
			if err = templruntime.WriteString(buf, 19, "</span>"); err != nil {
				return err
			}
		}
		if err = templruntime.WriteString(buf, 20, "</a></div></div>"); err != nil {
			return err
		}

		return nil
	})
}

func StatusBadge(status models.Status) templ.Component {
	return templruntime.GeneratedTemplate(func(in templruntime.GeneratedComponentInput) (err error) {
		w, ctx := in.Writer, in.Context
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		buf, isBuffer := templruntime.GetBuffer(w)
		if !isBuffer {
			defer func() {
				bufErr := templruntime.ReleaseBuffer(buf)
				if err == nil {
					err = bufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		children := templ.GetChildren(ctx)
		if children == nil {
			children = templ.NopComponent
		}
		_ = children
		ctx = templ.ClearChildren(ctx)

		switch string(status) {
		case "open":
			return templruntime.WriteString(buf, 1, `<span class="badge badge-info badge-sm">Open</span>`)
		case "in_progress":
			return templruntime.WriteString(buf, 2, `<span class="badge badge-warning badge-sm">In Progress</span>`)
		case "ready_to_sprint":
			return templruntime.WriteString(buf, 3, `<span class="badge badge-warning badge-sm">Ready to sprint</span>`)
		case "closed":
			return templruntime.WriteString(buf, 4, `<span class="badge badge-success badge-sm">Closed</span>`)
		case "blocked":
			return templruntime.WriteString(buf, 5, `<span class="badge badge-error badge-sm">Blocked</span>`)
		case "deferred":
			return templruntime.WriteString(buf, 6, `<span class="badge badge-ghost badge-sm">Deferred</span>`)
		default:
			var text string
			text, err = templ.JoinStringErrs(string(status))
			if err != nil {
				return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 70, Col: 4}
			}
			if err = templruntime.WriteString(buf, 7, `<span class="badge badge-ghost badge-sm">`); err != nil {
				return err
			}
			if _, err = buf.WriteString(templ.EscapeString(text)); err != nil {
				return err
			}
			return templruntime.WriteString(buf, 8, "</span>")
		}
	})
}

func TypeBadge(issueType models.IssueType) templ.Component {
	return templruntime.GeneratedTemplate(func(in templruntime.GeneratedComponentInput) (err error) {
		w, ctx := in.Writer, in.Context
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		buf, isBuffer := templruntime.GetBuffer(w)
		if !isBuffer {
			defer func() {
				bufErr := templruntime.ReleaseBuffer(buf)
				if err == nil {
					err = bufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		children := templ.GetChildren(ctx)
		if children == nil {
			children = templ.NopComponent
		}
		_ = children
		ctx = templ.ClearChildren(ctx)

		switch string(issueType) {
		case "bug":
			return templruntime.WriteString(buf, 1, `<span class="badge badge-outline badge-error badge-sm">Bug</span>`)
		case "feature":
			return templruntime.WriteString(buf, 2, `<span class="badge badge-outline badge-success badge-sm">Feature</span>`)
		case "task":
			return templruntime.WriteString(buf, 3, `<span class="badge badge-outline badge-info badge-sm">Task</span>`)
		case "chore":
			return templruntime.WriteString(buf, 4, `<span class="badge badge-outline badge-warning badge-sm">Chore</span>`)
		case "epic":
			return templruntime.WriteString(buf, 5, `<span class="badge badge-outline badge-secondary badge-sm">Epic</span>`)
		default:
			var text string
			text, err = templ.JoinStringErrs(string(issueType))
			if err != nil {
				return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 87, Col: 4}
			}
			if err = templruntime.WriteString(buf, 6, `<span class="badge badge-outline badge-sm">`); err != nil {
				return err
			}
			if _, err = buf.WriteString(templ.EscapeString(text)); err != nil {
				return err
			}
			return templruntime.WriteString(buf, 7, "</span>")
		}
	})
}

func PriorityBadge(priority int) templ.Component {
	return templruntime.GeneratedTemplate(func(in templruntime.GeneratedComponentInput) (err error) {
		w, ctx := in.Writer, in.Context
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		buf, isBuffer := templruntime.GetBuffer(w)
		if !isBuffer {
			defer func() {
				bufErr := templruntime.ReleaseBuffer(buf)
				if err == nil {
					err = bufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		children := templ.GetChildren(ctx)
		if children == nil {
			children = templ.NopComponent
		}
		_ = children
		ctx = templ.ClearChildren(ctx)

		switch priority {
		case 0:
			return templruntime.WriteString(buf, 1, `<span class="badge badge-error badge-sm">Critical</span>`)
		case 1:
			return templruntime.WriteString(buf, 2, `<span class="badge badge-warning badge-sm">High</span>`)
		case 2:
			return templruntime.WriteString(buf, 3, `<span class="badge badge-info badge-sm">Medium</span>`)
		case 3:
			return templruntime.WriteString(buf, 4, `<span class="badge badge-ghost badge-sm">Low</span>`)
		case 4:
			return templruntime.WriteString(buf, 5, `<span class="badge badge-outline badge-sm">Minimal</span>`)
		default:
			var text string
			text, err = templ.JoinStringErrs(fmt.Sprintf("P%d", priority))
			if err != nil {
				return templ.Error{Err: err, FileName: `pkg/web/components/issue_detail.templ`, Line: 104, Col: 4}
			}
			if err = templruntime.WriteString(buf, 6, `<span class="badge badge-ghost badge-sm">`); err != nil {
				return err
			}
			if _, err = buf.WriteString(templ.EscapeString(text)); err != nil {
				return err
			}
			return templruntime.WriteString(buf, 7, "</span>")
		}
	})
}

var _ = templruntime.GeneratedTemplate

