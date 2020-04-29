package pputil

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// FprintJSON is pretty printed json output shorthand.
func FprintJSON(w io.Writer, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}

// PrintJSON is similar that a relation about fmt.Printf and fmt.Fprintf.
func PrintJSON(data interface{}) {
	FprintJSON(os.Stdout, data)
}
