package newkvsecret

import (
	"fmt"
	"strings"

	app "passKeeper/internal/cmd/app"
	secret "passKeeper/internal/models/secret"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func EditKVTui(kv secret.KeyValue, meta string, id uint) error {
	finalModel, err := tea.NewProgram(InitialEditModel(kv, meta)).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}
	app := app.GetApplication()

	err = app.EditKVSecret(id, ans.Meta, ans.Key, ans.Value)
	if err != nil {
		return err
	}

	return nil

}

// ConfigTui starts the Bubbletea Configuration TUI
func NewKVTui() error {
	finalModel, err := tea.NewProgram(InitialModel()).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}
	app := app.GetApplication()

	err = app.CreateKVSecret(ans.Meta, ans.Key, ans.Value)
	if err != nil {
		return err
	}

	return nil

}

var (
	focusedColor = lipgloss.AdaptiveColor{Light: "236", Dark: "248"}
	blurredColor = lipgloss.AdaptiveColor{Light: "238", Dark: "246"}

	focusedStyle = lipgloss.NewStyle().Foreground(focusedColor)
	blurredStyle = lipgloss.NewStyle().Foreground(blurredColor)
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Bold(true).Render("[ Save ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Save"))
)

type Model struct {
	focusIndex int

	inputs []textinput.Model
	Meta   string
	Key    string
	Value  string
	Done   bool
	width  int
	height int
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, store values provided
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.Meta = m.inputs[0].Value()
				m.Key = m.inputs[1].Value()
				m.Value = m.inputs[2].Value()
				m.Done = true
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	boderColor := lipgloss.AdaptiveColor{Light: "22", Dark: "42"}
	style := lipgloss.NewStyle().
		BorderForeground(boderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Width(80).
		BorderBottom(true)

	title := "\n[:New Key Value Secret:]\n"
	titleStyle := lipgloss.NewStyle().Foreground(boderColor).Bold(true)
	s := titleStyle.Render(title)

	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(style.Render(m.inputs[i].View()))
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			s,
			b.String(),
		),
	)

}

func InitialEditModel(kv secret.KeyValue, meta string) Model {
	m := Model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 255
		t.Prompt = ""

		switch i {
		case 0:
			t.Placeholder = "Meta"
			t.TextStyle = focusedStyle
			t.SetValue(meta)
			t.Focus()
		case 1:
			t.Placeholder = "Key"
			t.TextStyle = focusedStyle
			t.SetValue(kv.Key)
		case 2:
			t.Placeholder = "Value"
			t.TextStyle = focusedStyle
			t.SetValue(kv.Value)
		}

		m.inputs[i] = t
	}

	return m
}

func InitialModel() Model {
	m := Model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 255
		t.Prompt = ""

		switch i {
		case 0:
			t.Placeholder = "Meta"
			t.TextStyle = focusedStyle
			t.Focus()
		case 1:
			t.Placeholder = "Key"
			t.TextStyle = focusedStyle
		case 2:
			t.Placeholder = "Value"
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}

	return m
}
