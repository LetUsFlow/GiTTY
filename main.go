package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []connection
	selected int
	maxlen   int
	quit     bool
}

type connection struct {
	Username string   `json:"username"`
	Hostname string   `json:"hostname"`
	Comment  string   `json:"comment"`
	Args     []string `json:"args"`
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

	maxlen := 0
	for _, choice := range choices {
		targetlen := len(GetLoginTarget(choice))
		if targetlen > maxlen {
			maxlen = targetlen
		}
	}

	return model{
		choices:  choices,
		selected: 0,
		maxlen:   maxlen,
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
	var output string
	style := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("205"))
	selectedStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("10")).Bold(true)
	grayedOutStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("240"))

	for i, choice := range m.choices {
		target := GetLoginTarget(choice)
		number := fmt.Sprintf("%d.", i)
		if i > 9 {
			number = "  "
		}
		if i == m.selected {
			output += selectedStyle.Render(fmt.Sprintf("> %s %s", number, target))
		} else {
			output += style.Render(fmt.Sprintf("  %s %s", number, target))
		}
		output += fmt.Sprintf("%*s", m.maxlen-len(target), "")
		output += grayedOutStyle.Render(choice.Comment) + "\n"
	}
	output += "Press 'enter' to select, 'q' to quit."
	return output
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
		selected.Args = append(selected.Args, target)
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
