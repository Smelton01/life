package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type model struct {
	height int
	width  int
	grid   grid
	speed  time.Duration
	start  bool
}

type pos struct {
	x int
	y int
}
type grid struct {
	alive map[pos]alive
}

type alive struct{}

var style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#25A065")).
	Padding(0, 1)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithMouseAllMotion(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		height: 20,
		width:  20,
		grid:   initCell(),
		speed:  time.Millisecond * 500,
	}
}

func initCell() grid {
	grid := grid{alive: make(map[pos]alive)}
	init := []pos{{5, 5}, {6, 5}, {7, 5}}
	for _, pos := range init {
		grid.alive[pos] = alive{}
	}
	return grid
}

// Get adjacent cells within board height, width
func (p pos) getAdjacent(height, width int) []pos {
	adj := []pos{}
	candidates := [][2]int{{0, 1}, {0, -1}, {1, 0}, {1, -1}, {1, 1}, {-1, 0}, {-1, 1}, {-1, -1}}
	for _, n := range candidates {
		newPos := pos{
			x: p.x + n[0],
			y: p.y + n[1],
		}
		if newPos.x < 0 || newPos.x >= width || newPos.y < 0 || newPos.y >= height {
			continue
		}
		adj = append(adj, newPos)
	}
	return adj
}

func (m model) Init() tea.Cmd {
	return tick(m.speed)
}

func tick(speed time.Duration) tea.Cmd {
	return tea.Tick(time.Duration(speed), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.start = !m.start
			return m, tick(m.speed)
		}

	case tickMsg:
		if !m.start {
			return m, nil
		}

		m.nextState()
		return m, tick(m.speed)
	case tea.MouseMsg:
		e := tea.MouseEvent(msg)
		m.toggleCell(e)
	}

	return m, nil
}

func (m *model) nextState() {
	counter := map[pos]int{}
	// iterate through all live cells
	for k := range m.grid.alive {
		if _, ok := counter[k]; !ok {
			counter[k] = 0
		}
		// get adjacent cells of each
		adjList := k.getAdjacent(m.height, m.width)
		for _, adj := range adjList {
			// increment adjacent counter for each neighbor
			counter[adj]++
		}
	}
	for cell, numAdj := range counter {
		// live cells which die
		if numAdj < 2 || numAdj > 3 {
			delete(m.grid.alive, cell)
		}
		// dead cell which live
		if numAdj == 3 {
			m.grid.alive[cell] = alive{}
		}
	}
}

func (m model) View() string {
	var s string
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			var val string
			if _, ok := m.grid.alive[pos{x: x, y: y}]; !ok {
				val = " "
			} else {
				val = "█"
			}
			s += val
		}
		s += "\n"
	}

	return style.Render(s)
}

// Toggle a cell state with a left click
func (m model) toggleCell(event tea.MouseEvent) {
	if event.Type != tea.MouseLeft {
		return
	}
	if event.X < 0 || event.X >= m.width || event.Y < 0 || event.Y >= m.height {
		return
	}
	cell := pos{
		x: event.X,
		y: event.Y,
	}
	if _, ok := m.grid.alive[cell]; ok {
		delete(m.grid.alive, cell)
		return
	}
	m.grid.alive[cell] = alive{}
}
