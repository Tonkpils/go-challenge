// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import "fmt"

const measureCount = 16

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []Track
}

func (p Pattern) String() string {
	out := fmt.Sprintf("Saved with HW Version: %s\nTempo: %g\n", p.Version, p.Tempo)
	for _, track := range p.Tracks {
		out += track.String() + "\n"
	}
	return out
}

type Track struct {
	ID    uint8
	Name  string
	Steps []byte
}

func (t Track) String() string {
	out := fmt.Sprintf("(%d) %s", t.ID, t.Name)
	measure := "|"
	for i := 0; i < 16; i += 4 {
		notes := t.Steps[i : i+4]
		measure += string(notes) + "|"
	}
	out += "\t" + measure
	return out
}
