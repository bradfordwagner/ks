package resources

import (
	"encoding/json"
	"fmt"
	"os"
)

const CacheFile = ".resources.json"
const Perms = 0644 // rw-r--r--

// Resources is a struct that contains a list of resource names
type Resources struct {
	Names []string `json:"names"`
}

func (r Resources) Write(ksdir string) (err error) {
	// convert r to json and write to CacheFile
	bytes, err := json.Marshal(r)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s/%s", ksdir, CacheFile)
	err = os.WriteFile(fileName, bytes, Perms)
	return
}
