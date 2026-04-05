package resources

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	CacheFile   = ".ks.resources.json"
	Perms       = 0644 // rw-r--r--
	schemaV1    = 1
	schemaV2    = 2
	CurrentVersion = schemaV2
)

// ResourceEntry holds a Kubernetes resource name and its selection vote count.
type ResourceEntry struct {
	Name  string `json:"name"`
	Votes int    `json:"votes"`
}

// Resources is a struct that contains a list of resource names
type Resources struct {
	Version int             `json:"version"`
	// TMUX env variable cache
	Cache   map[string]Cache `json:"cache"`
	Names   []ResourceEntry  `json:"names"`
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

// VoteFor increments the vote count for the named resource entry. Unknown names are a no-op.
func (r *Resources) VoteFor(name string) {
	for i := range r.Names {
		if r.Names[i].Name == name {
			r.Names[i].Votes++
			return
		}
	}
}

// SortedNames returns resource names sorted by votes descending, then alphabetically ascending.
func (r *Resources) SortedNames() []string {
	entries := make([]ResourceEntry, len(r.Names))
	copy(entries, r.Names)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Votes != entries[j].Votes {
			return entries[i].Votes > entries[j].Votes
		}
		return entries[i].Name < entries[j].Name
	})
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
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

// v1Raw is used to decode the old flat []string names format.
type v1Raw struct {
	Cache map[string]Cache `json:"cache"`
	Names []string         `json:"names"`
}

// LoadResources reads the json file at ksdir/.ks.resources.json and returns a Resources struct.
// Files without a "version" field (v1) are automatically migrated to v2.
func LoadResources(ksdir string) (r Resources, err error) {
	fileName := fmt.Sprintf("%s/%s", ksdir, CacheFile)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	// Peek at the version to decide how to decode.
	var versionProbe struct {
		Version int             `json:"version"`
		Names   json.RawMessage `json:"names"`
	}
	if err = json.Unmarshal(data, &versionProbe); err != nil {
		return
	}

	if versionProbe.Version < schemaV2 {
		// v1: names is []string — migrate to v2 ResourceEntry slice.
		var raw v1Raw
		if jsonErr := json.Unmarshal(data, &raw); jsonErr == nil {
			r.Cache = raw.Cache
			r.Names = make([]ResourceEntry, len(raw.Names))
			for i, n := range raw.Names {
				r.Names[i] = ResourceEntry{Name: n, Votes: 0}
			}
		}
		r.Version = schemaV2
	} else {
		if err = json.Unmarshal(data, &r); err != nil {
			return
		}
	}

	// if Cache is nil, create a new map
	if r.Cache == nil {
		r.Cache = make(map[string]Cache)
	}

	return
}
