// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"fmt"
)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []Track
}

// String returns a string representation of a splice drump pattern
func (p Pattern) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Saved with HW Version: %s\nTempo: %g\n", p.Version, p.Tempo))
	for _, track := range p.Tracks {
		buf.WriteString(fmt.Sprintf("%s\n", track))
	}
	return buf.String()
}

// Track represents a single piece of the drum track
type Track struct {
	ID    uint8
	Name  string
	Steps []byte
}

const barSeparator = "|"

// String returns a string representation of a track.
// 'x' represents when the piece makes a sound.
func (t Track) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("(%d) %s\t", t.ID, t.Name))
	buf.WriteString(barSeparator)
	for i := 0; i < 16; i += 4 {
		notes := t.Steps[i : i+4]
		buf.WriteString(string(notes) + barSeparator)
	}
	return buf.String()
}
