package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Selector chooses an option from a list.
type Selector interface {
	Select(options []string) (string, error)
}

type item struct{ title string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice string
}

func newModel(options []string) model {
	items := make([]list.Item, len(options))
	for i, opt := range options {
		items[i] = item{title: opt}
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
		if it, ok := m.list.SelectedItem().(item); ok {
			m.choice = it.title
		}
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() string { return m.list.View() }

type bubbleSelector struct{}

// NewSelector returns a Bubbletea-based selector.
func NewSelector() Selector { return bubbleSelector{} }

func (bubbleSelector) Select(options []string) (string, error) {
	m := newModel(options)
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return "", err
	}
	final := res.(model)
	return final.choice, nil
}
