package history

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

// saveFile :
func saveFile(filename string, write func(w io.Writer) error) (rerr error) {
	fp, err := ioutil.TempFile(".", "ch")
	if err != nil {
		return errors.Wrap(err, "save, create tempfile")
	}
	defer func() {
		fp.Close()
		// log.Println("rename", fp.Name(), "->", filename)
		if err := os.Rename(fp.Name(), filename); err != nil {
			rerr = err
		}
	}()
	if err := write(fp); err != nil {
		return err
	}
	log.Printf("save. %q\n", filename)
	rp, err := os.Open(filename)
	if err != nil {
		finfo, serr := os.Stat(filename)
		if serr == nil && finfo != nil {
			return errors.Wrap(err, "save, open file")
		}
		if serr != nil {
			log.Printf("create. %q\n", filename)
			fp, cerr := os.Create(filename)
			if cerr != nil {
				return err
			}
			rp = fp
		}
	}
	defer rp.Close()
	io.Copy(fp, rp)
	return
}
