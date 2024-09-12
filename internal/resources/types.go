package resources

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	CacheFile = ".ks.resources.json"
	Perms     = 0644 // rw-r--r--
)

// Resources is a struct that contains a list of resource names
type Resources struct {
	// TMUX env variable cache
	Cache map[string]Cache `json:"cache"`
	Names []string         `json:"names"`
}

type Cache struct {
	// TMUX_PANE env variable to resource id
	IdToResource map[string]string `json:"id_to_resource"`
}

func NewCache() Cache {
	return Cache{
		IdToResource: make(map[string]string),
	}
}

func (r *Resources) ResetCache() {
	r.Cache = make(map[string]Cache)
}

func (r *Resources) ResetPane() {
	// load tmux and pane from env
	tmux := os.Getenv("TMUX")
	pane := os.Getenv("TMUX_PANE")
	if tmux == "" || pane == "" {
		return
	}

	if _, ok := r.Cache[tmux]; !ok {
		return
	}
	c := r.Cache[tmux]
	delete(c.IdToResource, pane)
	r.Cache[tmux] = c
}

func (r *Resources) Upsert(resource string) {
	if r.Cache == nil {
		r.Cache = make(map[string]Cache)
	}

	// load tmux and pane from env
	tmux := os.Getenv("TMUX")
	pane := os.Getenv("TMUX_PANE")
	if tmux == "" || pane == "" {
		return
	}

	if _, ok := r.Cache[tmux]; !ok {
		r.Cache[tmux] = NewCache()
	}
	c := r.Cache[tmux]
	c.IdToResource[pane] = resource
}

func (r *Resources) Get() (resource string) {
	// load tmux and pane from env
	tmux := os.Getenv("TMUX")
	pane := os.Getenv("TMUX_PANE")
	if tmux == "" || pane == "" {
		return
	}

	if _, ok := r.Cache[tmux]; !ok {
		return
	}
	c := r.Cache[tmux]
	resource, _ = c.IdToResource[pane]
	return
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

	// if Cache is nil, create a new map
	if r.Cache == nil {
		r.Cache = make(map[string]Cache)
	}

	return
}
