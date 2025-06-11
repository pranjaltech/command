package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModelSize(t *testing.T) {
	m := newModel([]string{"a", "b"})
	if m.list.Width() <= 0 {
		t.Errorf("expected width > 0, got %d", m.list.Width())
	}
	if m.list.Height() != 2 {
		t.Errorf("expected height 2, got %d", m.list.Height())
	}
}

func TestModelEditFlow(t *testing.T) {
	m := newModel([]string{"first"})
	// select the item
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(model)
	if m.state != stateEdit {
		t.Fatalf("expected edit state, got %v", m.state)
	}
	// type space and -x
	for _, r := range " -x" {
		mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = mm.(model)
	}
	// confirm edited command
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(model)
	if m.choice != "first -x" {
		t.Errorf("expected 'first -x', got %q", m.choice)
	}
}
