package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

// ProviderInfo holds metadata about an AI provider
type ProviderInfo struct {
	Name string // Display name (e.g., "OpenAI")
	Key  string // Internal key (e.g., "openai")
}

// ConfigResult holds the result of the configuration flow
type ConfigResult struct {
	Provider string // Provider key
	APIKey   string // API key entered by user
	Canceled bool   // True if user canceled
}

// ExistingConfig holds current configuration for display
type ExistingConfig struct {
	Provider string
	APIKey   string
	Model    string
}

// configStep represents the current step in the config flow
type configStep int

const (
	stepProvider configStep = iota
	stepAPIKey
	stepDone
)

// configItem implements list.Item for provider selection
type configItem struct {
	name string
	key  string
}

func (i configItem) Title() string       { return i.name }
func (i configItem) Description() string { return "" }
func (i configItem) FilterValue() string { return i.name }

// ConfigUI is a Bubbletea model for the configuration flow
type ConfigUI struct {
	step      configStep
	providers []ProviderInfo
	list      list.Model
	apiInput  textinput.Model
	selected  int
	result    ConfigResult
	existing  *ExistingConfig
	err       error
	width     int
	height    int
}

// NewConfigUI creates a new configuration UI
func NewConfigUI(providers []ProviderInfo, existing *ExistingConfig) ConfigUI {
	// Create list items
	items := make([]list.Item, len(providers))
	width := 0
	for i, p := range providers {
		title := fmt.Sprintf("%d. %s", i+1, p.Name)
		items[i] = configItem{name: title, key: p.Key}
		if len(title) > width {
			width = len(title)
		}
	}

	// Create list model
	delegate := simpleDelegate{}
	l := list.New(items, delegate, width+4, len(items)+2)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	// Pre-select current provider if existing
	if existing != nil {
		for i, p := range providers {
			if p.Key == existing.Provider {
				l.Select(i)
				break
			}
		}
	}

	// Create text input for API key
	ti := textinput.New()
	ti.Placeholder = ""
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '\u2022' // bullet
	ti.CharLimit = 256
	ti.Width = 40

	return ConfigUI{
		step:      stepProvider,
		providers: providers,
		list:      l,
		apiInput:  ti,
		existing:  existing,
	}
}

func (m ConfigUI) Init() tea.Cmd {
	return nil
}

func (m ConfigUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.step == stepAPIKey {
				// Go back to provider selection
				m.step = stepProvider
				m.apiInput.Blur()
				return m, nil
			}
			// Cancel and quit
			m.result.Canceled = true
			return m, tea.Quit

		case tea.KeyEnter:
			switch m.step {
			case stepProvider:
				// Provider selected, move to API key input
				m.selected = m.list.Index()
				m.result.Provider = m.providers[m.selected].Key
				m.step = stepAPIKey
				m.apiInput.Focus()
				return m, textinput.Blink

			case stepAPIKey:
				// API key entered
				key := strings.TrimSpace(m.apiInput.Value())
				if key == "" {
					// Don't allow empty key
					return m, nil
				}
				m.result.APIKey = key
				m.step = stepDone
				return m, tea.Quit
			}
		}

		// Handle number keys for quick selection in provider step
		if m.step == stepProvider && msg.Type == tea.KeyRunes {
			if len(msg.Runes) == 1 && msg.Runes[0] >= '1' && msg.Runes[0] <= '9' {
				idx := int(msg.Runes[0] - '1')
				if idx < len(m.providers) {
					m.list.Select(idx)
					m.selected = idx
					m.result.Provider = m.providers[idx].Key
					m.step = stepAPIKey
					m.apiInput.Focus()
					return m, textinput.Blink
				}
			}
		}
	}

	// Update child models
	var cmd tea.Cmd
	switch m.step {
	case stepProvider:
		m.list, cmd = m.list.Update(msg)
	case stepAPIKey:
		m.apiInput, cmd = m.apiInput.Update(msg)
	}
	return m, cmd
}

func (m ConfigUI) View() string {
	switch m.step {
	case stepProvider:
		return m.viewProviderSelect()
	case stepAPIKey:
		return m.viewAPIKeyInput()
	case stepDone:
		return m.viewComplete()
	}
	return ""
}

func (m ConfigUI) viewProviderSelect() string {
	var s strings.Builder

	s.WriteString(BoldStyle.Render("Select AI Provider:"))
	s.WriteString("\n\n")
	s.WriteString(m.list.View())
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("[↑↓] Navigate  [Enter] Select  [Esc] Cancel"))

	return s.String()
}

func (m ConfigUI) viewAPIKeyInput() string {
	var s strings.Builder

	providerName := m.providers[m.selected].Name
	s.WriteString(Checkmark() + " Provider: " + BoldStyle.Render(providerName))
	s.WriteString("\n\n")

	// Show existing key hint if updating
	if m.existing != nil && m.existing.APIKey != "" && m.existing.Provider == m.result.Provider {
		s.WriteString(MutedStyle.Render("Current: " + MaskAPIKey(m.existing.APIKey)))
		s.WriteString("\n\n")
	}

	s.WriteString("Enter your API key:\n")
	s.WriteString("> " + m.apiInput.View())
	s.WriteString("\n\n")
	s.WriteString(MutedStyle.Render("[Enter] Confirm  [Esc] Back"))

	return s.String()
}

func (m ConfigUI) viewComplete() string {
	var s strings.Builder

	s.WriteString(Checkmark() + " Provider: " + BoldStyle.Render(m.providers[m.selected].Name))
	s.WriteString("\n")
	s.WriteString(Checkmark() + " API Key:  " + MaskAPIKey(m.result.APIKey))

	return s.String()
}

// Result returns the configuration result
func (m ConfigUI) Result() ConfigResult {
	return m.result
}

// RunConfigUI runs the configuration UI and returns the result
func RunConfigUI(providers []ProviderInfo, existing *ExistingConfig) (*ConfigResult, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return runSimpleConfigUI(providers, existing)
	}

	m := NewConfigUI(providers, existing)
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return nil, err
	}

	final := res.(ConfigUI)
	result := final.Result()
	return &result, nil
}

// runSimpleConfigUI is a fallback for non-TTY environments
func runSimpleConfigUI(providers []ProviderInfo, existing *ExistingConfig) (*ConfigResult, error) {
	picker := NewPicker()

	// Provider selection
	names := make([]string, len(providers))
	for i, p := range providers {
		names[i] = p.Name
	}

	fmt.Println("Select AI Provider:")
	idx, err := picker.Pick(names)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(providers) {
		return &ConfigResult{Canceled: true}, nil
	}

	selected := providers[idx]

	// API key input
	fmt.Printf("\n%s Provider: %s\n", Checkmark(), selected.Name)
	fmt.Print("\nEnter your API key: ")

	// Use term.ReadPassword for secure input
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, err
	}

	key := strings.TrimSpace(string(b))
	if key == "" {
		return &ConfigResult{Canceled: true}, nil
	}

	return &ConfigResult{
		Provider: selected.Key,
		APIKey:   key,
	}, nil
}
