package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPickModelDigit(t *testing.T) {
	m := newPickModel([]string{"a", "b"})
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m = mm.(pickModel)
	if !m.done || m.choice != 1 {
		t.Fatalf("expected choice 1 done, got %d %v", m.choice, m.done)
	}
}

func TestPickModelEnter(t *testing.T) {
	m := newPickModel([]string{"a", "b"})
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = mm.(pickModel)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(pickModel)
	if m.choice != 1 || !m.done {
		t.Fatalf("expected choice 1 done, got %d %v", m.choice, m.done)
	}
}
