package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type simpleDelegate struct{}

func (simpleDelegate) Height() int                             { return 1 }
func (simpleDelegate) Spacing() int                            { return 0 }
func (simpleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (simpleDelegate) Render(w io.Writer, m list.Model, index int, it list.Item) {
	var title string
	switch v := it.(type) {
	case pickItem:
		title = v.title
	case item:
		title = v.title
	default:
		if t, ok := any(it).(interface{ Title() string }); ok {
			title = t.Title()
		}
	}
	cursor := "  "
	if index == m.Index() {
		cursor = "> "
	}
	fmt.Fprintf(w, "%s%s", cursor, title)
}
