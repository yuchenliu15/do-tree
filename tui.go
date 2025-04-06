package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int
const (
	Selecting State = iota
	Staging
	Result
)

type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
	state State 
	textInput textinput.Model
	output string
}

func initalModel(choices []string) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		choices: choices,
		selected: make(map[int]struct{}),
		state: Selecting,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink 
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Printf("msg: %T\n", msg)
	switch state := m.state; state {
	case Selecting:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit	
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case " ":
				m.state = Staging
			}
		}
	case Staging:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				log.Printf("Input: %s", m.textInput.Value())
				selected := []string{}
				for i, choice := range m.choices {
					if _, ok := m.selected[i]; ok {
						selected = append(selected, choice)
					}
				}
				commandInput := append(strings.Split(m.textInput.Value(), " "), selected...)
				m.state = Result
				command := exec.Command(commandInput[0], commandInput[1:]...)
				output, err := command.CombinedOutput()
				if err != nil {
					log.Printf("Error: %v\n", err)
					return m, tea.Quit
				}
				m.output = string(output)
				m.textInput.Reset()
				
				return m, nil
			case "esc":
				m.state = Selecting
				m.textInput.Reset()
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
		}
	case Result:
		return m, tea.Quit
	}
	return m, cmd 
}

func (m model) View() string {
	var s string
	if m.state == Selecting {
		s = "Select file/dir to apply command to:\n"
		s += "(Hit space to enter comand)\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	} else if m.state == Staging{
		selected := []string{}
		for i, choice := range m.choices {
			if _, ok := m.selected[i]; ok {
				selected = append(selected, choice)
			}
		}
		s = "Enter command to apply to selected files:\n"
		s += fmt.Sprintln(selected)
		s += "(Hit space to enter comand)\n"
		s += m.textInput.View()
	} else {
		s = "Result:\n"
		s += m.output
	}
	return s
}

func main() {	
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Printf("Error could not log: %s", err)
	}
	defer f.Close()
	
	entries, err := tree("./")
	if err != nil {
		log.Println(err)
	}	
	
	p := tea.NewProgram(initalModel(entries))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occured: %v\n", err)
		os.Exit(1)
	}
}
