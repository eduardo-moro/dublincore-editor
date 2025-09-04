package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eduardo-moro/metadata-editor/dublincore"
)

var (
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurryStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	currentValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Italic(true)
	fieldLabelStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	placeholderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

type model struct {
	inputs    []textinput.Model
	focused   int
	dc        *dublincore.DublinCore
	done      bool
	cancelled bool
}

func initialModel(dc *dublincore.DublinCore) model {
	m := model{
		inputs: make([]textinput.Model, 4),
		dc:     dc,
	}

	// Title input
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "e.g., Senior Backend Developer"
	m.inputs[0].PlaceholderStyle = placeholderStyle
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle
	if len(dc.Title) > 0 && dc.Title[0] != "" {
		m.inputs[0].SetValue(dc.Title[0])
	}

	// Creator input
	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "e.g., João Silva, Maria Santos"
	m.inputs[1].PlaceholderStyle = placeholderStyle
	m.inputs[1].PromptStyle = blurryStyle
	if len(dc.Creator) > 0 {
		m.inputs[1].SetValue(strings.Join(dc.Creator, ", "))
	}

	// Keywords input
	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "e.g., Go, Backend, Microservices, PHP"
	m.inputs[2].PlaceholderStyle = placeholderStyle
	m.inputs[2].PromptStyle = blurryStyle
	if len(dc.Keywords) > 0 {
		m.inputs[2].SetValue(strings.Join(dc.Keywords, ", "))
	}

	// Description input
	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "e.g., Experienced backend developer with 6+ years in technology"
	m.inputs[3].PlaceholderStyle = placeholderStyle
	m.inputs[3].PromptStyle = blurryStyle
	m.inputs[3].CharLimit = 200
	if len(dc.Description) > 0 && dc.Description[0] != "" {
		m.inputs[3].SetValue(dc.Description[0])
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs) {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = blurryStyle
				m.inputs[i].TextStyle = blurryStyle
			}

			return m, tea.Batch(cmds...)

		case "enter":
			if m.focused == len(m.inputs) {
				m.done = true
				// Update the Dublin Core object with user input before quitting
				m.updateDublinCoreFromInputs()
				return m, tea.Quit
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *model) updateDublinCoreFromInputs() {
	// Title
	titleInput := strings.TrimSpace(m.inputs[0].Value())
	if titleInput != "" && titleInput != m.inputs[0].Placeholder {
		m.dc.SetTitle(titleInput)
	}

	// Creator
	creatorsInput := strings.TrimSpace(m.inputs[1].Value())
	if creatorsInput != "" && creatorsInput != m.inputs[1].Placeholder {
		newCreators := []string{}
		for _, creator := range strings.Split(creatorsInput, ",") {
			if trimmed := strings.TrimSpace(creator); trimmed != "" {
				newCreators = append(newCreators, trimmed)
			}
		}
		m.dc.Creator = newCreators
	}

	// Keywords
	keywordsInput := strings.TrimSpace(m.inputs[2].Value())
	if keywordsInput != "" && keywordsInput != m.inputs[2].Placeholder {
		newKeywords := []string{}
		for _, keyword := range strings.Split(keywordsInput, ",") {
			if trimmed := strings.TrimSpace(keyword); trimmed != "" {
				newKeywords = append(newKeywords, trimmed)
			}
		}
		m.dc.Keywords = newKeywords
	}

	// Description
	descriptionInput := strings.TrimSpace(m.inputs[3].Value())
	if descriptionInput != "" && descriptionInput != m.inputs[3].Placeholder {
		m.dc.SetDescription(descriptionInput)
	}

	// Always set category to "curriculo"
	m.dc.SetCategory()
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📄 Dublin Core Metadata Editor\n\n"))

	// Title field
	b.WriteString(fieldLabelStyle.Render("DC: Title") + "\n")
	b.WriteString(m.inputs[0].View())
	b.WriteString("\n\n")

	// Creator field
	b.WriteString(fieldLabelStyle.Render("DC: Creator (comma-separated)") + "\n")
	b.WriteString(m.inputs[1].View())
	b.WriteString("\n\n")

	// Keywords field
	b.WriteString(fieldLabelStyle.Render("CP: Keywords (comma-separated)") + "\n")
	b.WriteString(m.inputs[2].View())
	b.WriteString("\n\n")

	// Description field
	b.WriteString(fieldLabelStyle.Render("CP: Description") + "\n")
	b.WriteString(m.inputs[3].View())
	b.WriteString("\n\n")

	// Category field (read-only)
	b.WriteString(fieldLabelStyle.Render("CP: Category") + "\n")
	b.WriteString("curriculo (fixed value)\n\n")

	// Navigation help
	b.WriteString(helpStyle.Render("↑/↓: Navigate • Tab/Shift+Tab: Next/Previous • Enter: Submit • Esc: Cancel"))
	b.WriteString("\n\n")

	// Submit button
	button := blurryStyle
	if m.focused == len(m.inputs) {
		button = focusedStyle
	}
	b.WriteString(button.Render("[ Submit Changes ]"))

	return b.String()
}

// RunEditor starts the BubbleTea TUI and returns updated metadata
func RunEditor(dc *dublincore.DublinCore) (*dublincore.DublinCore, bool, error) {
	p := tea.NewProgram(initialModel(dc))

	finalModel, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	if m, ok := finalModel.(model); ok {
		if m.done {
			return m.dc, false, nil
		}
		if m.cancelled {
			return dc, true, nil
		}
	}

	return dc, false, nil
}
