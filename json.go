package selfish

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

func ppJSON(target interface{}) {
	b, err := json.Marshal(target)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	out.WriteTo(os.Stdout)
}
