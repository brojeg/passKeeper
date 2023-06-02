package list

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type modelList struct {
	table table.Model
}

func (m modelList) Init() tea.Cmd { return nil }

func (m modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m modelList) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewModel(columns []table.Column, rows []table.Row) modelList {

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return modelList{t}
}

func truncate(str string, num int) string {
	if len(str) <= num {
		return str
	}
	if num > 3 {
		num -= 3
	}
	return str[0:num] + "..."
}
