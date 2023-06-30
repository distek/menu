package menu

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SingleModel struct {
	Selected  string
	Interrupt bool // If ctrl-c or q was pressed
	Choices   []string
	cursor    int
	Message   string
	Title     string
}

type MultipleModel struct {
	Selected  map[int]struct{}
	Interrupt bool // If ctrl-c or q was pressed
	Choices   []string
	cursor    int
	Message   string
	Title     string
}

func NewSingle(choices []string, title, message string) *SingleModel {
	return &SingleModel{
		Choices:  choices,
		Title:    title,
		Message:  message,
		Selected: "",
		cursor:   0,
	}
}

func NewMultiple(choices []string, title, message string) *MultipleModel {
	return &MultipleModel{
		Choices:  choices,
		Title:    title,
		Message:  message,
		Selected: make(map[int]struct{}),
		cursor:   0,
	}
}

// Run runs the tea.Model. Returns the final model, and error if any
func Run(m tea.Model, altScreen bool) (tea.Model, error) {
	var p *tea.Program
	if altScreen {
		p = tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	} else {
		p = tea.NewProgram(m, tea.WithOutput(os.Stderr))
	}

	m, err := p.Run()

	return m, err

}

// Single

func (m SingleModel) Init() tea.Cmd {
	return nil
}

func (m SingleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Interrupt = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}

		case tea.KeyEnter.String(), " ":
			m.Selected = m.Choices[m.cursor]
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m SingleModel) View() string {
	s := fmt.Sprintf("%s\n%s\n\n", m.Title, m.Message)

	for i, choice := range m.Choices {
		cursor := " "

		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q or Ctrl+C to quit, Enter to select\n"

	return s
}

// Multiple
func (m MultipleModel) Init() tea.Cmd {
	return nil
}

func (m MultipleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Interrupt = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}

		case tea.KeyEnter.String():
			return m, tea.Quit
		case " ":
			if _, ok := m.Selected[m.cursor]; ok {
				delete(m.Selected, m.cursor)
			} else {
				m.Selected[m.cursor] = struct{}{}
			}

		case tea.KeyEsc.String():
			somethingSelected := false
			for i := range m.Selected {
				if _, ok := m.Selected[i]; ok {
					somethingSelected = true
				}
			}

			for i := range m.Choices {
				if !somethingSelected {
					m.Selected[i] = struct{}{}
				} else {
					delete(m.Selected, i)
				}
			}
		}
	}

	return m, nil
}

func (m MultipleModel) View() string {
	s := fmt.Sprintf("%s\n%s\n\n", m.Title, m.Message)

	for i, choice := range m.Choices {
		cursor := " "

		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += fmt.Sprintf("\nSelected: %s\n", m.Message)

	s += "\nPress q or Ctrl+C to quit, Space to select, Esc to select all, Enter to finalize selection\n"

	return s
}

type InputModel struct {
	Input  textinput.Model
	Prompt string
	Error  error
}

func NewInput(prompt, placeholder string, charLimit, width int) InputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = charLimit
	ti.Width = width

	return InputModel{
		Input:  ti,
		Prompt: prompt,
		Error:  nil,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.Error = msg
		return m, nil
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		m.Prompt,
		m.Input.View(),
		"(esc to cancel, enter to accept)",
	)
}
