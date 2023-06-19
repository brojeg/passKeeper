package newccsecret

import (
	"fmt"
	app "passKeeper/internal/cmd/app"
	secret "passKeeper/internal/models/secret"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigTui starts the Bubbletea Configuration TUI

func EditCCTui(cc secret.CreditCard, meta string, id uint) error {
	finalModel, err := tea.NewProgram(InitialEditModel(cc, meta)).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}

	app := app.GetApplication()

	err = app.EditCCSecret(id, ans.Meta, ans.CCN, ans.EXP, ans.CVV, ans.CHolder)
	if err != nil {
		return err
	}

	return nil

}

func NewCCTui() error {
	finalModel, err := tea.NewProgram(InitialModel()).Run()
	if err != nil {
		return err
	}

	ans := finalModel.(Model)

	if !ans.Done {
		return nil
	}
	app := app.GetApplication()

	err = app.CreateCCSecret(ans.Meta, ans.CCN, ans.EXP, ans.CVV, ans.CHolder)
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

	inputs  []textinput.Model
	Meta    string
	CCN     string
	EXP     string
	CVV     string
	CHolder string
	Done    bool
	width   int
	height  int
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
				m.CCN = m.inputs[1].Value()
				m.EXP = m.inputs[2].Value()
				m.CVV = m.inputs[3].Value()
				m.CHolder = m.inputs[4].Value()
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

	title := "\n[:New Credit Card Secret:]\n"
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

func InitialEditModel(cc secret.CreditCard, meta string) Model {
	m := Model{
		inputs:  make([]textinput.Model, 5),
		CCN:     cc.Number,
		EXP:     cc.Expiration,
		Meta:    meta,
		CVV:     cc.CVV,
		CHolder: cc.Cardholder,
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 255
		t.Prompt = ""

		switch i {
		case 0:
			t.Placeholder = "Credit card metadata"
			t.TextStyle = focusedStyle
			t.SetValue(meta)
			t.CharLimit = 30

			t.Focus()
		case 1:
			t.Placeholder = "Credit card number"
			t.TextStyle = focusedStyle
			t.SetValue(cc.Number)
			t.CharLimit = 20
			t.Validate = ccnValidator
		case 2:
			t.Placeholder = "Exp"
			t.TextStyle = focusedStyle
			t.SetValue(cc.Expiration)
			t.Validate = expValidator
			t.CharLimit = 5
		case 3:
			t.Placeholder = "Cvv"
			t.TextStyle = focusedStyle
			t.SetValue(cc.CVV)
			t.Validate = cvvValidator
			t.CharLimit = 3
		case 4:
			t.Placeholder = "Card Holder"
			t.TextStyle = focusedStyle
			t.SetValue(cc.Cardholder)
			t.CharLimit = 20
		}

		m.inputs[i] = t
	}

	return m
}

func InitialModel() Model {
	m := Model{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 255
		t.Prompt = ""

		switch i {
		case 0:
			t.Placeholder = "Credit card metadata"
			t.TextStyle = focusedStyle
			t.CharLimit = 30

			t.Focus()
		case 1:
			t.Placeholder = "Credit card number"
			t.TextStyle = focusedStyle
			t.Validate = ccnValidator
			t.CharLimit = 20
		case 2:
			t.Placeholder = "Exp"
			t.TextStyle = focusedStyle
			t.Validate = expValidator
			t.CharLimit = 5
		case 3:
			t.Placeholder = "Cvv"
			t.TextStyle = focusedStyle
			t.Validate = cvvValidator
			t.CharLimit = 3
		case 4:
			t.Placeholder = "Card Holder"
			t.TextStyle = focusedStyle
			t.CharLimit = 20
		}

		m.inputs[i] = t
	}

	return m
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest thould be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}
	if len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}
