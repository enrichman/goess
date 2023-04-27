package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type model struct {
	list          list.Model
	listDelegate  list.ItemDelegate
	currentChoice *selection.Choice[*answer]
	quitting      bool
}

type question struct {
	title     string
	answers   []*answer
	selection *selection.Model[*answer]
}

type answer struct {
	value    string
	Selected bool
	Current  bool // is the current selection?
}

func (q question) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	selected, _ := m.SelectedItem().(*question)
	for _, listItem := range m.Items() {
		q, ok := listItem.(*question)
		if !ok {
			return nil
		}

		if selected.title == q.title {
			q.selection.Update(msg)
		}
	}
	return nil
}
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	q, ok := listItem.(*question)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, q.selection.View())

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.PaddingLeft(4).Render(str))
	} else {
		fmt.Fprint(w, itemStyle.Render(str))
	}

}

func (m model) Init() tea.Cmd {
	for _, listItem := range m.list.Items() {
		q, ok := listItem.(*question)
		if !ok {
			return nil
		}
		q.selection.Init()
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case msg.String() == tea.KeySpace.String():
			m.currentChoice.Value.Selected = !m.currentChoice.Value.Selected
			selectedQuestion, _ := m.list.SelectedItem().(*question)
			selectedQuestion.selection.Update(msg)

		case msg.String() == "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "enter":
			for _, listItem := range m.list.Items() {
				q, _ := listItem.(*question)
				for _, a := range q.answers {
					if a.Selected {
						fmt.Printf("'%s' - answer: '%s'\n", q.title, a.value)
					}
				}
			}

			return m, tea.Quit

		case key.Matches(msg, m.list.KeyMap.CursorUp):
			selectedQuestion, _ := m.list.SelectedItem().(*question)
			ch, _ := selectedQuestion.selection.ValueAsChoice()
			idx := ch.Index()
			if idx == 0 {
				m.list.CursorUp()
			}
			m.currentChoice.Value.Current = false

			selectedQuestion.selection.Update(msg)

			selectedQuestion, _ = m.list.SelectedItem().(*question)
			ch, _ = selectedQuestion.selection.ValueAsChoice()
			m.currentChoice = ch
			m.currentChoice.Value.Current = true

		case key.Matches(msg, m.list.KeyMap.CursorDown):
			selectedQuestion, _ := m.list.SelectedItem().(*question)
			ch, _ := selectedQuestion.selection.ValueAsChoice()
			idx := ch.Index()
			if idx >= len(selectedQuestion.answers)-1 {
				m.list.CursorDown()
			}
			m.currentChoice.Value.Current = false

			selectedQuestion.selection.Update(msg)

			selectedQuestion, _ = m.list.SelectedItem().(*question)
			ch, _ = selectedQuestion.selection.ValueAsChoice()
			m.currentChoice = ch
			m.currentChoice.Value.Current = true
		}

	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

var logfile io.Writer

func main() {
	var err error
	logfile, err = os.OpenFile("cli.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	q1 := newQuestion("Question 1?", []string{"Q1 A-1", "Q1 A-2"})
	q1.answers[0].Current = true

	items := []list.Item{
		q1,
		newQuestion("Question 2?", []string{"Q2 A-1", "Q2 A-2"}),
		// newQuestion("Question 3?", []string{"Q3 A-1", "Q3 A-2"}),
		// newQuestion("Question 4?", []string{"Q4 A-1", "Q4 A-2"}),
	}

	const defaultWidth = 20

	delegate := itemDelegate{}
	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	//	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{
		list:          l,
		currentChoice: &selection.Choice[*answer]{Value: q1.answers[0]},
		listDelegate:  delegate,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func newQuestion(title string, answers []string) *question {
	//blue := termenv.String().Foreground(termenv.ANSI256Color(32)) //nolint:gomnd

	answerss := []*answer{}
	for _, a := range answers {
		answerss = append(answerss, &answer{value: a})
	}

	sel := selection.New(title, answerss)
	sel.Filter = nil
	sel.Template = customTemplate
	sel.ResultTemplate = resultTemplate
	sel.SelectedChoiceStyle = func(c *selection.Choice[*answer]) string {
		current := " "
		if c.Value.Current {
			current = ">"
		}

		selected := "[ ]"
		if c.Value.Selected {
			selected = "[x]"
		}

		return fmt.Sprintf("%s %s %s \n", current, selected, c.Value.value)

		// return blue.Bold().Styled(c.Value.value)
	}
	sel.UnselectedChoiceStyle = func(c *selection.Choice[*answer]) string {
		if c.Value.Selected {
			return fmt.Sprintf("  [x] %s \n", c.Value.value)
		}
		return fmt.Sprintf("  [ ] %s \n", c.Value.value)
	}
	sel.ExtendedTemplateFuncs = map[string]interface{}{
		"name": func(c *selection.Choice[answer]) string { return c.Value.value },
	}

	return &question{
		title:     title,
		answers:   answerss,
		selection: selection.NewModel(sel),
	}
}

const (
	customTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if eq $.SelectedIndex $i }}
  	{{- print (Selected $choice)  }}
  {{- else }}
  	{{- print (Unselected $choice)  }}
  {{- end }}
{{- end}}`
	resultTemplate = `
	{{- print .Prompt " "  (name .FinalChoice) "\n" -}}
	`
)
