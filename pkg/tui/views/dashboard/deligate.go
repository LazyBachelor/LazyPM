package dashboard

import (
	"fmt"
	"io"

	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type IssueListDelegate struct{}

func (d IssueListDelegate) Height() int                               { return 2 }
func (d IssueListDelegate) Spacing() int                              { return 1 }
func (d IssueListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d IssueListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	issue, ok := listItem.(ListIssue)

	if !ok {
		return
	}

	id := issue.ID
	title := issue.Title()
	description := issue.Description()

	str := fmt.Sprintf("ID: %s\t%s\nDescription:\t%s", id, title, description)

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render(s...)
		}
	}

	fmt.Fprint(w, fn(str))
}
