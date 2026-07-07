package resurrect

import (
	"os/exec"
	"strconv"
	"strings"
)

type TmuxPane struct {
	Session        string
	WindowIdx      int
	PaneIdx        int
	PaneID         string
	PanePID        int
	CurrentCommand string
}

// ListPanes enumerates all panes across all tmux sessions.
// Returns nil if tmux is not running or not inside a session.
func ListPanes() ([]TmuxPane, error) {
	out, err := exec.Command("tmux", "list-panes", "-a", "-F",
		"#{session_name}\t#{window_index}\t#{pane_index}\t#{pane_id}\t#{pane_pid}\t#{pane_current_command}",
	).Output()
	if err != nil {
		return nil, nil // tmux not running; treat as empty
	}

	var panes []TmuxPane
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 6 {
			continue
		}
		windowIdx, _ := strconv.Atoi(parts[1])
		paneIdx, _ := strconv.Atoi(parts[2])
		panePID, _ := strconv.Atoi(parts[4])
		panes = append(panes, TmuxPane{
			Session:        parts[0],
			WindowIdx:      windowIdx,
			PaneIdx:        paneIdx,
			PaneID:         parts[3],
			PanePID:        panePID,
			CurrentCommand: parts[5],
		})
	}
	return panes, nil
}

// CurrentPane returns the positional info for the pane the caller is running in.
func CurrentPane() (TmuxPane, error) {
	out, err := exec.Command("tmux", "display-message", "-p",
		"#{session_name}\t#{window_index}\t#{pane_index}\t#{pane_id}\t#{pane_pid}\t#{pane_current_command}",
	).Output()
	if err != nil {
		return TmuxPane{}, err
	}
	parts := strings.Split(strings.TrimSpace(string(out)), "\t")
	if len(parts) != 6 {
		return TmuxPane{}, nil
	}
	windowIdx, _ := strconv.Atoi(parts[1])
	paneIdx, _ := strconv.Atoi(parts[2])
	panePID, _ := strconv.Atoi(parts[4])
	return TmuxPane{
		Session:        parts[0],
		WindowIdx:      windowIdx,
		PaneIdx:        paneIdx,
		PaneID:         parts[3],
		PanePID:        panePID,
		CurrentCommand: parts[5],
	}, nil
}

// SendKeys sends a command to a tmux pane identified by target (pane ID or index).
func SendKeys(target, command string) error {
	return exec.Command("tmux", "send-keys", "-t", target, command, "Enter").Run()
}
