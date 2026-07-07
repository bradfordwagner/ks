package resurrect

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ReadEnv reads the environment variables of process pid from /proc/<pid>/environ.
// Returns an empty map on non-Linux platforms or if the file is unreadable.
func ReadEnv(pid int) (map[string]string, error) {
	env := make(map[string]string)
	if runtime.GOOS != "linux" {
		return env, nil
	}
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return env, nil // process may have exited; treat as empty
	}
	for _, entry := range strings.Split(string(data), "\x00") {
		k, v, ok := strings.Cut(entry, "=")
		if ok && k != "" {
			env[k] = v
		}
	}
	return env, nil
}

// ReadCmdline reads the command-line arguments of process pid from /proc/<pid>/cmdline.
// Returns nil on non-Linux platforms or if the file is unreadable.
func ReadCmdline(pid int) ([]string, error) {
	if runtime.GOOS != "linux" {
		return nil, nil
	}
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, nil
	}
	var args []string
	for _, arg := range strings.Split(string(data), "\x00") {
		if arg != "" {
			args = append(args, arg)
		}
	}
	return args, nil
}

type procInfo struct {
	name string
	ppid int
}

// FindDescendantByName walks /proc to find a process whose name matches target
// and whose ancestor chain includes ancestorPid. Returns 0 if not found.
func FindDescendantByName(ancestorPid int, target string) (int, error) {
	if runtime.GOOS != "linux" {
		return 0, nil
	}

	procs := make(map[int]procInfo)
	entries, err := filepath.Glob("/proc/[0-9]*/status")
	if err != nil {
		return 0, err
	}
	for _, path := range entries {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var name string
		var ppid int
		for _, line := range strings.Split(string(data), "\n") {
			switch {
			case strings.HasPrefix(line, "Name:\t"):
				name = strings.TrimPrefix(line, "Name:\t")
			case strings.HasPrefix(line, "PPid:\t"):
				ppid, _ = strconv.Atoi(strings.TrimPrefix(line, "PPid:\t"))
			}
		}
		if name == "" {
			continue
		}
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			continue
		}
		pid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}
		procs[pid] = procInfo{name: name, ppid: ppid}
	}

	for pid, info := range procs {
		if info.name != target {
			continue
		}
		if hasAncestor(pid, ancestorPid, procs) {
			return pid, nil
		}
	}
	return 0, nil
}

func hasAncestor(pid, ancestor int, procs map[int]procInfo) bool {
	visited := make(map[int]bool)
	current := pid
	for {
		info, ok := procs[current]
		if !ok || visited[current] {
			return false
		}
		if info.ppid == ancestor {
			return true
		}
		visited[current] = true
		current = info.ppid
	}
}
