package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {

	testCases := []struct {
		desc  string
		model model
		cells []pos
		want  map[pos]alive
	}{
		{
			desc: "simple case",
			model: model{
				height:  20,
				width:   20,
				timeout: time.Now().Add(time.Second * 20),
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
			desc: "oscillator",
			model: model{
				height:  20,
				width:   20,
				timeout: time.Now().Add(time.Second * 20),
				grid: grid{
					alive: make(map[pos]alive),
				},
			},
			cells: []pos{{2, 2}, {4, 2}},
			want:  map[pos]alive{pos{3, 1}: alive{}, pos{3, 3}: alive{}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for _, cell := range tC.cells {
				tC.model.grid.alive[cell] = alive{}
			}
			// log.Println(tC.model.grid.alive)
			tC.model.applyRules()
			// log.Println(tC.model.grid.alive)
			assert.Equal(t, tC.model.grid.alive, tC.want)
		})
	}
}
