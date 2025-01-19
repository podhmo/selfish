package internal

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// fprintJSON is pretty printed json output shorthand.
func fprintJSON(w io.Writer, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}

// PrintJSON is similar that a relation about fmt.Printf and fmt.Fprintf.
func PrintJSON(data interface{}) {
	fprintJSON(os.Stdout, data)
}
