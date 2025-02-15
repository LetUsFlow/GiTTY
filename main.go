package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    choices  []connection
    selected int
    maxlen int
    quit bool
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
        connlen := len(fmt.Sprintf("%s@%s", choice.Username, choice.Hostname))
        if connlen > maxlen {
            maxlen = connlen
        }
    }

    return model{
        choices: choices,
        selected: 0,
        maxlen: maxlen,
        quit: false,
    }
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
        conn := fmt.Sprintf("%s@%s", choice.Username, choice.Hostname)
        number := fmt.Sprintf("%d.", i)
        if i > 9 {
            number = "  "
        }
        if i == m.selected {
            output += selectedStyle.Render(fmt.Sprintf("> %s %s", number, conn))
        } else {
            output += style.Render(fmt.Sprintf("  %s %s", number, conn))
        }
        output += fmt.Sprintf("%*s", m.maxlen - len(conn), "")
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
        if (m.quit) {
            os.Exit(0)
        }
        selected := m.choices[m.selected]
        fmt.Printf("Connecting...\n")

        conn := fmt.Sprintf("%s@%s", selected.Username, selected.Hostname)
        selected.Args = append(selected.Args, conn)
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
