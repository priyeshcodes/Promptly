package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/priyeshcodes/smart-task-cli/internal/task"
)

type model struct {
	tasks          []task.Task
	cursor         int
	quitting       bool
	showingDetails bool
	themeDark      bool
	sortBy         string
	filter         string
	ts             *task.TaskStore
	err            error
	showHelp       bool
}

func initialModel() model {
	ts, err := task.NewTaskStore("tasks.db")
	if err != nil {
		return model{err: err}
	}
	tasks, err := ts.GetAllTasks()
	if err != nil {
		return model{err: err}
	}
	return model{
		tasks:  tasks,
		ts:     ts,
		sortBy: "default",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if !m.showingDetails && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if !m.showingDetails && m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "enter":
			m.showingDetails = !m.showingDetails

		case "b":
			m.showingDetails = false

		case "c":
			if len(m.tasks) > 0 {
				t := &m.tasks[m.cursor]
				t.Completed = !t.Completed
				if err := m.ts.UpdateTask(t); err == nil {
					// m.tasks[m.cursor] = t
				}
			}

		case "d":
			if len(m.tasks) > 0 {
				id := m.tasks[m.cursor].ID
				_ = m.ts.DeleteTask(id)
				m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "r":
			tasks, err := m.ts.GetAllTasks()
			if err == nil {
				m.tasks = tasks
				if m.cursor >= len(tasks) && m.cursor > 0 {
					m.cursor--
				}
			}

		case "s":
			if m.sortBy == "priority" {
				m.sortBy = "deadline"
				sort.SliceStable(m.tasks, func(i, j int) bool {
					if m.tasks[i].Deadline == nil {
						return false
					}
					if m.tasks[j].Deadline == nil {
						return true
					}
					return m.tasks[i].Deadline.Before(*m.tasks[j].Deadline)
				})
			} else {
				m.sortBy = "priority"
				sort.SliceStable(m.tasks, func(i, j int) bool {
					return m.tasks[i].Priority < m.tasks[j].Priority
				})
			}

		case "/":
			m.filter = "high" // Just example. Later: dynamic filter input.
			var filtered []task.Task
			for _, t := range m.tasks {
				if strings.ToLower(t.Priority.String()) == m.filter {
					filtered = append(filtered, t)
				}
			}
			m.tasks = filtered
			m.cursor = 0

		case "t":
			m.themeDark = !m.themeDark

		case "?":
			m.showHelp = !m.showHelp
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("❌ Error: %v\n", m.err)
	}
	if m.quitting {
		return "👋 Goodbye from Smart Task CLI!\n"
	}
	if m.showHelp {
		return helpView()
	}
	if m.showingDetails && len(m.tasks) > 0 {
		t := m.tasks[m.cursor]
		deadline := "None"
		if t.Deadline != nil {
			deadline = t.Deadline.Format(time.RFC1123)
		}
		return fmt.Sprintf(
			"📄 Task Details\n\nTitle: %s\nDescription: %s\nPriority: %s\nCompleted: %t\nDeadline: %s\n\nPress 'b' to go back\n",
			t.Title, t.Description, t.Priority.String(), t.Completed, deadline,
		)
	}
	var b strings.Builder
	if m.themeDark {
		b.WriteString("🌑 Smart Task CLI - Dark Mode\n\n")
	} else {
		b.WriteString("🌕 Smart Task CLI - Light Mode\n\n")
	}
	if len(m.tasks) == 0 {
		b.WriteString("🎉 No tasks yet. Add some from the CLI!\n")
	} else {
		for i, t := range m.tasks {
			cursor := "  "
			if m.cursor == i {
				cursor = "👉"
			}
			status := "❌"
			if t.Completed {
				status = "✅"
			}
			line := fmt.Sprintf("%s [%s] %s (Priority: %s)\n", cursor, status, t.Title, t.Priority.String())
			b.WriteString(line)
		}
	}
	b.WriteString("\n↑/↓ to move | enter = details | c = ✓ toggle | d = delete | r = refresh\n")
	b.WriteString("s = sort | / = filter | t = theme | ? = help | q = quit\n")
	return b.String()
}

func helpView() string {
	return `
🆘 Help - Smart Task CLI Commands

↑ / ↓     : Navigate tasks
Enter     : View task details
b         : Go back from details
c         : Toggle complete/incomplete
d         : Delete task
r         : Refresh task list
s         : Sort by priority/deadline
/         : Filter high priority tasks
t         : Toggle light/dark theme
?         : Toggle help menu
q / Ctrl+C: Quit
`
}

func Start() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
