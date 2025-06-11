package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

type state int

const (
	stateList state = iota
	stateEdit
)

type model struct {
	list   list.Model
	input  textinput.Model
	state  state
	choice string
}

func newModel(options []string) model {
	items := make([]list.Item, len(options))
	width := 20
	for i, opt := range options {
		items[i] = item{title: opt}
		if l := len(opt) + 2; l > width {
			width = l
		}
	}
	l := list.New(items, list.NewDefaultDelegate(), width, len(options))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = ""
	l.SetShowTitle(false)
	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 0
	ti.Blur()
	ti.Width = width
	return model{list: l, input: ti}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateList:
		m.list, cmd = m.list.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			if it, ok := m.list.SelectedItem().(item); ok {
				m.input.SetValue(it.title)
				m.input.Focus()
				m.state = stateEdit
			}
		}
	case stateEdit:
		m.input, cmd = m.input.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			m.choice = m.input.Value()
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.state == stateEdit {
		return m.input.View()
	}
	return m.list.View()
}

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
