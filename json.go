package selfish

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

// FprintJSON is pretty printed json output shorthand.
func FprintJSON(w io.Writer, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	out.WriteTo(w)
}

// PrintJSON is similar that a relation about fmt.Printf and fmt.Fprintf.
func PrintJSON(data interface{}) {
    FprintJSON(os.Stdout, data)
}
