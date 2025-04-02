package main 

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
    "fmt"
    "os"
	"log"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
}

func initalModel() model {
	return model{
		choices: []string{"carrorts", "potatoes", "onions"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("msg: %T\n", msg)
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
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "here's the market\n"
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
	return s
}

func main() {
	p := tea.NewProgram(initalModel())
	
	if len(os.Getenv("DEBUG")) > 0 {
		fmt.Printf("Running in debug mode\n")
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Printf("Error could not log: %s", err)
		}
		defer f.Close()
	}
	
	entries, err := tree(".")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Starting tree")
		for _, entry := range entries {
			log.Printf("%s\n", entry)
		}
	}

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occured: %v\n", err)
		os.Exit(1)
	}
}
