package drum

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func debug(d io.Reader) {
	dumper := hex.Dumper(os.Stdout)
	bs, err := ioutil.ReadAll(d)
	if err != nil {
		log.Fatal(err)
	}
	dumper.Write(bs)
	dumper.Close()
}
