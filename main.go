package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

const (
	interval = time.Second
	duration = time.Second * 10
)

type model struct {
	height  int
	width   int
	timeout time.Time
	// playing bool
	grid grid
}

type pos struct {
	x int
	y int
}
type grid struct {
	alive map[pos]alive
}

type alive struct{}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		height:  20,
		width:   20,
		grid:    initCell(),
		timeout: time.Now().Add(duration),
		// playing: false,
	}
}

func initCell() grid {
	grid := grid{alive: make(map[pos]alive)}
	pos1 := pos{x: 5, y: 5}
	grid.alive[pos1] = alive{}
	return grid
}

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
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
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
			m.height++
		}

	case tickMsg:
		// m.grid.alive[pos{
		// 	x: rando(m.width),
		// 	y: rando(m.height),
		// }] = alive{}
		m.applyRules()
		return m, tick()
	}
	return m, nil
}

func rando(n int) int {
	return rand.Intn(n)
}

func (m *model) applyRules() {
	counter := map[pos]int{}
	// iterate through all alive cells
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
		if numAdj < 2 || numAdj > 3 {
			delete(m.grid.alive, cell)
		}
		if numAdj == 3 {
			m.grid.alive[cell] = alive{}
		}
	}
}

func (m model) View() string {
	// The header
	s := []rune("What should we buy at the market?\n\n")

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			var val string
			if _, ok := m.grid.alive[pos{x: x, y: y}]; !ok {
				val = " "
			} else {
				val = string([]byte{254})
			}
			s = append(s, []rune(val)...)
		}
		s = append(s, rune('\n'))
	}

	return string(s)
}
