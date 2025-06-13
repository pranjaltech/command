package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type Picker interface {
	Pick(options []string) (int, error)
}

type pickItem struct{ title string }

func (i pickItem) Title() string       { return i.title }
func (i pickItem) Description() string { return "" }
func (i pickItem) FilterValue() string { return i.title }

// pickModel wraps bubbles list for simple item selection.
// Digits 1-9 can be used to select directly.
type pickModel struct {
	list   list.Model
	choice int
	done   bool
}

func newPickModel(options []string) pickModel {
	items := make([]list.Item, len(options))
	width := 0
	for i, opt := range options {
		title := fmt.Sprintf("%d. %s", i+1, opt)
		items[i] = pickItem{title: title}
		if len(title) > width {
			width = len(title)
		}
	}
	delegate := simpleDelegate{}
	l := list.New(items, delegate, width+2, len(items))
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	return pickModel{list: l}
}

func (m pickModel) Init() tea.Cmd { return nil }

func (m pickModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			m.choice = m.list.Index()
			m.done = true
			return m, tea.Quit
		}
		if msg.Type == tea.KeyRunes {
			if len(msg.Runes) == 1 && msg.Runes[0] >= '1' && msg.Runes[0] <= '9' {
				idx := int(msg.Runes[0] - '1')
				if idx < len(m.list.Items()) {
					m.list.Select(idx)
					m.choice = idx
					m.done = true
					return m, tea.Quit
				}
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickModel) View() string { return m.list.View() }

type bubblePicker struct{}

type simplePicker struct{}

// NewPicker creates a Picker implemented with Bubbletea.
func NewPicker() Picker {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return simplePicker{}
	}
	return bubblePicker{}
}

func (bubblePicker) Pick(options []string) (int, error) {
	m := newPickModel(options)
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return 0, err
	}
	final := res.(pickModel)
	return final.choice, nil
}

func (simplePicker) Pick(options []string) (int, error) {
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, opt)
	}
	fmt.Fprint(os.Stderr, "> ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, nil
	}
	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(options) {
		return 0, fmt.Errorf("invalid choice")
	}
	return idx - 1, nil
}
