package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
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
	l := list.New(items, simpleDelegate{}, width, len(options))
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
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m, nil
		}
		m.list, cmd = m.list.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			if it, ok := m.list.SelectedItem().(item); ok {
				m.input.SetValue(it.title)
				m.input.Focus()
				m.state = stateEdit
			}
		}
	case stateEdit:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.input.Width = msg.Width
		}
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
type simpleSelector struct{}

// NewSelector returns a selector that uses Bubbletea when a TTY is available and
// falls back to basic stdin prompts otherwise.
func NewSelector() Selector {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return simpleSelector{}
	}
	return bubbleSelector{}
}

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

func (simpleSelector) Select(options []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, opt)
	}
	fmt.Fprint(os.Stderr, "> ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	idx := 0
	if line != "" {
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > len(options) {
			return "", fmt.Errorf("invalid choice")
		}
		idx = n - 1
	}
	choice := options[idx]
	fmt.Fprintf(os.Stderr, "%s\n> ", choice)
	edit, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	edit = strings.TrimSpace(edit)
	if edit != "" {
		choice = edit
	}
	return choice, nil
}
