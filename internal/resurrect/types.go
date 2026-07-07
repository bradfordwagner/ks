package resurrect

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ResurrectFile   = ".ks.resurrect.json"
	currentVersion  = 1
	filePerms       = 0644
)

type ResurrectPane struct {
	Session    string `json:"session"`
	WindowIdx  int    `json:"window_idx"`
	PaneIdx    int    `json:"pane_idx"`
	Kubeconfig string `json:"kubeconfig"`
	Resource   string `json:"resource,omitempty"`
	Verb       string `json:"verb,omitempty"`
}

type ResurrectState struct {
	Version int             `json:"version"`
	Panes   []ResurrectPane `json:"panes"`
}

func (s ResurrectState) Write(ksdir string) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("%s/%s", ksdir, ResurrectFile), bytes, filePerms)
}

// Upsert loads the sidecar, replaces or appends the entry matching p's positional key, and writes it back.
func Upsert(ksdir string, p ResurrectPane) error {
	s, err := Load(ksdir)
	if err != nil {
		return err
	}
	s.Version = currentVersion
	type posKey struct{ session string; w, pane int }
	key := posKey{p.Session, p.WindowIdx, p.PaneIdx}
	for i, existing := range s.Panes {
		if (posKey{existing.Session, existing.WindowIdx, existing.PaneIdx}) == key {
			s.Panes[i] = p
			return s.Write(ksdir)
		}
	}
	s.Panes = append(s.Panes, p)
	return s.Write(ksdir)
}

func Load(ksdir string) (ResurrectState, error) {
	var s ResurrectState
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", ksdir, ResurrectFile))
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	s.Version = currentVersion
	err = json.Unmarshal(data, &s)
	return s, err
}
