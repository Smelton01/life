package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStep(t *testing.T) {

	testCases := []struct {
		desc  string
		model model
		cells []pos
		want  map[pos]alive
	}{
		{
			desc: "empty grid should not change",
			model: model{
				height: 20,
				width:  20,
				grid: grid{
					alive: make(map[pos]alive),
				},
			},
			cells: []pos{},
			want:  map[pos]alive{},
		},
		{
			desc: "single cell should die",
			model: model{
				height: 20,
				width:  20,
				grid: grid{
					alive: make(map[pos]alive),
				},
			},
			cells: []pos{{
				x: 1,
				y: 1,
			}},
			want: map[pos]alive{},
		},
		{
			desc: "oscillator should rotate",
			model: model{
				height: 20,
				width:  20,
				grid: grid{
					alive: make(map[pos]alive),
				},
			},
			cells: []pos{{2, 2}, {3, 2}, {4, 2}},
			want:  map[pos]alive{{3, 1}: {}, {3, 3}: {}, {3, 2}: {}},
		},
		{
			desc: "three points converge",
			model: model{
				height: 20,
				width:  20,
				grid: grid{
					alive: make(map[pos]alive),
				},
			},
			cells: []pos{{2, 2}, {3, 3}, {4, 1}},
			want:  map[pos]alive{{3, 2}: {}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.model.speed = time.Second
			for _, cell := range tC.cells {
				tC.model.grid.alive[cell] = alive{}
			}
			tC.model.nextState()

			assert.Equal(t, tC.want, tC.model.grid.alive)
		})
	}
}
