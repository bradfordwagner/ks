package resources

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	CacheFile = ".resources.json"
	Perms     = 0644 // rw-r--r--
)

// Resources is a struct that contains a list of resource names
type Resources struct {
	Names []string `json:"names"`
}

func (r Resources) Write(ksdir string) (fileName string, err error) {
	// convert r to json and write to CacheFile
	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return
	}

	fileName = fmt.Sprintf("%s/%s", ksdir, CacheFile)
	err = os.WriteFile(fileName, bytes, Perms)
	return
}

// LoadResources reads the json file at ksdir/.resources.json and returns a Resources struct
func LoadResources(ksdir string) (r Resources, err error) {
	fileName := fmt.Sprintf("%s/%s", ksdir, CacheFile)
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &r)
	return
}
