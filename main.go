package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	choices  []connection
	selected int
	quit     bool
}

type connection struct {
	Username string   `json:"username"`
	Hostname string   `json:"hostname"`
	Comment  string   `json:"comment"`
	Args     []string `json:"args"`
	Command  string   `json:"command"`
}

func initialModel() model {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var choices []connection
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&choices); err != nil {
		fmt.Println("Error decoding config:", err)
		os.Exit(1)
	}

	return model{
		choices:  choices,
		selected: 0,
		quit:     false,
	}
}

func GetLoginTarget(conn connection) string {
	return fmt.Sprintf("%s@%s", conn.Username, conn.Hostname)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quit = true
			return m, tea.Quit
		case "up":
			m.selected = (m.selected - 1 + len(m.choices)) % len(m.choices)
		case "down":
			m.selected = (m.selected + 1) % len(m.choices)
		case "enter":
			return m, tea.Quit
		}

		for i := 0; i < len(m.choices) && i < 10; i++ {
			if msg.String() == fmt.Sprintf("%d", i) {
				m.selected = i
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	baseStyle := lipgloss.NewStyle().Padding(0, 1)

	selectedStyle := baseStyle.Foreground(lipgloss.Color("10")).Bold(true)
	targetStyle := baseStyle.Foreground(lipgloss.Color("205"))
	commandStyle := baseStyle.Foreground(lipgloss.Color("240"))
	commentStyle := baseStyle

	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 0)
		}).Border(lipgloss.HiddenBorder())

	for i, choice := range m.choices {
		target := GetLoginTarget(choice)
		number := fmt.Sprintf("%d.", i)

		var indicatorCell, targetCell string

		if i == m.selected {
			indicatorCell = selectedStyle.Render("> " + number)
			targetCell = selectedStyle.Render(target)
		} else {
			indicatorCell = targetStyle.Render("  " + number)
			targetCell = targetStyle.Render(target)
		}

		commandCell := commandStyle.Render(choice.Command)
		commentCell := commentStyle.Render(choice.Comment)

		t.Row(
			indicatorCell,
			targetCell,
			commandCell,
			commentCell,
		)
	}

	var b strings.Builder
	b.WriteString(t.Render())
	b.WriteString("\n")
	b.WriteString("Press 'enter' to select, 'q' to quit.")

	return b.String()
}

func main() {
	final, err := tea.NewProgram(initialModel()).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	if m, ok := final.(model); ok {
		if m.quit {
			os.Exit(0)
		}
		selected := m.choices[m.selected]
		fmt.Printf("Connecting...\n")

		target := GetLoginTarget(selected)
		selected.Args = append(selected.Args, target, selected.Command)
		cmd := exec.Command("ssh", selected.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Error: Unexpected model type")
		os.Exit(1)
	}
}
