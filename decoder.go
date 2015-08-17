package drum

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	p := &Pattern{}
	if err := NewDecoder(file).Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

// Decoder reads and decodes splice patterns from an input stream
type Decoder struct {
	io.Reader
}

// NewDecoder returns a new decoder that reads from r
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

// Decode reads the input stream and decodes it's values into a Pattern
func (d *Decoder) Decode(p *Pattern) error {
	n, err := d.spliceHeaderInfo()
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("invalid .splice header: %s", err)
		}
	}

	if err := d.decodeBody(p, n); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("unable to decode .splice format: %s", err)
		}
	}
	return nil
}

// spliceHeaderInfo will read the inital bites
// and returns the size of the file contents or error if it's
// unable check the SPLICE portion or read the headers
func (d *Decoder) spliceHeaderInfo() (uint64, error) {
	hdr := make([]byte, 6)
	if _, err := io.ReadFull(d, hdr); err != nil {
		return 0, err
	}
	if string(hdr) != "SPLICE" {
		return 0, errors.New("unable to decode non SPLICE files")
	}

	var size uint64
	if err := binary.Read(d, binary.BigEndian, &size); err != nil {
		return 0, err
	}

	return size, nil
}

// decodeBody reads the contents of the input stream up to size
// and into p. It returns EOF or ErrUnexpectedEOF if it reads over
// size.
func (d *Decoder) decodeBody(p *Pattern, size uint64) error {
	version := make([]byte, 32)
	if _, err := io.ReadFull(d, version); err != nil {
		return errors.New("unable to decode hw version: " + err.Error())
	}
	p.Version = strings.Trim(string(version), "\x00")

	if err := binary.Read(d, binary.LittleEndian, &p.Tempo); err != nil {
		return errors.New("unable to decode tempo: " + err.Error())
	}

	// version and tempo
	size -= (32 + 4)

	// decode the tracks
	for size >= 0 {
		t := Track{}
		if err := binary.Read(d, binary.BigEndian, &t.ID); err != nil {
			return err
		}

		var nameLen int32
		if err := binary.Read(d, binary.BigEndian, &nameLen); err != nil {
			return err
		}

		name := make([]byte, nameLen)
		if _, err := io.ReadFull(d, name); err != nil {
			return err
		}
		t.Name = string(name)

		steps := make([]byte, measureCount)
		if err := binary.Read(d, binary.BigEndian, steps); err != nil {
			return err
		}

		// convert steps into 'x' or '-' runes
		for idx, step := range steps {
			if step == 1 {
				steps[idx] = 'x'
			} else {
				steps[idx] = '-'
			}
		}
		t.Steps = steps

		p.Tracks = append(p.Tracks, t)
		// ID, nameLen, measureCount, and size of name
		size -= (1 + 4 + measureCount + uint64(len(name)))
	}

	return nil
}
