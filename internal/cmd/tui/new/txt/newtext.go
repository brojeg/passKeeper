package newtextsecret

import (
	"fmt"
	"strings"

	cmd "passKeeper/internal/cmd/app"
	secret "passKeeper/internal/models/secret"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

	input    textinput.Model
	text     textarea.Model
	Metadata string
	Data     string
	Done     bool
	width    int
	height   int
}

func EditTextTui(txt secret.Text, meta string, id uint) error {
	finalModel, err := tea.NewProgram(editTextSecretModel(txt, meta)).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}

	app := cmd.GetApplication()

	err = app.EditTextSecret(id, ans.Metadata, ans.Data)
	if err != nil {
		return err
	}
	return nil

}

func NewTextTui() error {
	finalModel, err := tea.NewProgram(newTextSecretModel()).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}

	app := cmd.GetApplication()

	err = app.CreateTextSecret(ans.Metadata, ans.Data)
	if err != nil {
		return err
	}
	return nil

}

func editTextSecretModel(txt secret.Text, meta string) Model {

	input := textinput.New()
	input.CursorStyle = cursorStyle
	input.CharLimit = 255
	input.Prompt = ""
	input.Placeholder = "Secret metadata"
	input.TextStyle = focusedStyle
	input.SetValue(meta)
	input.Focus()

	text := textarea.New()
	text.Prompt = ""
	text.Placeholder = "Put your secret text here"
	text.CharLimit = 255
	text.SetValue(txt.Value)

	return Model{
		focusIndex: 0,
		input:      input,
		text:       text,
	}
}

func newTextSecretModel() Model {

	input := textinput.New()
	input.CursorStyle = cursorStyle
	input.CharLimit = 255
	input.Prompt = ""
	input.Placeholder = "Secret metadata"
	input.TextStyle = focusedStyle
	input.Focus()

	ta := textarea.New()
	ta.Prompt = ""
	ta.Placeholder = "Put your secret text here"

	return Model{
		focusIndex: 0,
		input:      input,
		text:       ta,
	}
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

		case "enter":
			if m.focusIndex == 2 {
				m.Metadata = m.input.Value()
				m.Data = m.text.Value()
				m.Done = true
				return m, tea.Quit
			}
		// Set focus to next input
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, store values provided

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			cmds := []tea.Cmd{}

			if m.focusIndex == 0 {
				cmds = append(cmds, m.input.Focus())
				m.input.PromptStyle = focusedStyle
				m.input.TextStyle = focusedStyle

				// Remove focused state for textarea
				m.text.Blur()
			} else if m.focusIndex == 1 {
				// Set focused state for textarea
				cmds = append(cmds, m.text.Focus())

				// Remove focused state for input
				m.input.Blur()
				m.input.PromptStyle = noStyle
				m.input.TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	m.input, cmd = m.input.Update(msg)

	var cmdText tea.Cmd
	m.text, cmdText = m.text.Update(msg)

	return tea.Batch(cmd, cmdText)
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

	title := "\n[:New Text Secret:]\n"
	titleStyle := lipgloss.NewStyle().Foreground(boderColor).Bold(true)
	s := titleStyle.Render(title)

	var b strings.Builder
	b.WriteString(style.Render(m.input.View()))
	b.WriteRune('\n')
	b.WriteString(style.Render(m.text.View()))

	button := &blurredButton
	if m.focusIndex == 2 {
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
