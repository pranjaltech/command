package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfigUI(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
		{Name: "Anthropic", Key: "anthropic"},
	}

	// Test without existing config
	m := NewConfigUI(providers, nil)
	if m.step != stepProvider {
		t.Errorf("NewConfigUI() step = %v, expected stepProvider", m.step)
	}
	if len(m.providers) != 2 {
		t.Errorf("NewConfigUI() providers len = %d, expected 2", len(m.providers))
	}

	// Test with existing config
	existing := &ExistingConfig{
		Provider: "anthropic",
		APIKey:   "test-key",
		Model:    "claude-3",
	}
	m2 := NewConfigUI(providers, existing)
	if m2.existing != existing {
		t.Error("NewConfigUI() existing config not set")
	}
}

func TestConfigUIResult(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
	}

	m := NewConfigUI(providers, nil)

	// Initial result should be empty
	result := m.Result()
	if result.Provider != "" {
		t.Errorf("Result().Provider = %q, expected empty", result.Provider)
	}
	if result.APIKey != "" {
		t.Errorf("Result().APIKey = %q, expected empty", result.APIKey)
	}
	if result.Canceled {
		t.Error("Result().Canceled = true, expected false")
	}
}

func TestConfigUIUpdate_Escape(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
	}

	m := NewConfigUI(providers, nil)

	// Pressing Escape should cancel
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(ConfigUI)

	if !m.result.Canceled {
		t.Error("Update(Esc) should set Canceled = true")
	}
	if cmd == nil {
		t.Error("Update(Esc) should return tea.Quit command")
	}
}

func TestConfigUIUpdate_ProviderSelection(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
		{Name: "Anthropic", Key: "anthropic"},
	}

	m := NewConfigUI(providers, nil)

	// Select with Enter
	m.list.Select(1) // Select Anthropic
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(ConfigUI)

	if m.step != stepAPIKey {
		t.Errorf("After Enter, step = %v, expected stepAPIKey", m.step)
	}
	if m.result.Provider != "anthropic" {
		t.Errorf("After Enter, Provider = %q, expected 'anthropic'", m.result.Provider)
	}
}

func TestConfigUIUpdate_NumberKeySelection(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
		{Name: "Anthropic", Key: "anthropic"},
	}

	m := NewConfigUI(providers, nil)

	// Press '2' to select second provider
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m = newModel.(ConfigUI)

	if m.step != stepAPIKey {
		t.Errorf("After '2' key, step = %v, expected stepAPIKey", m.step)
	}
	if m.result.Provider != "anthropic" {
		t.Errorf("After '2' key, Provider = %q, expected 'anthropic'", m.result.Provider)
	}
}

func TestConfigUIView(t *testing.T) {
	providers := []ProviderInfo{
		{Name: "OpenAI", Key: "openai"},
	}

	m := NewConfigUI(providers, nil)

	// Provider selection view
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string")
	}

	// Should contain "Select AI Provider"
	if len(view) < 10 {
		t.Errorf("View() too short: %q", view)
	}
}
