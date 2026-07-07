## 1. KS_RESOURCE env var override in Resource()

- [x] 1.1 In `internal/cmds/resource_cmd.go`, read `os.Getenv("KS_RESOURCE")` before `r.Get()`; if non-empty use it as `resourceType` and call `r.Upsert(resourceType)` then proceed normally
- [x] 1.2 Write a test in `internal/cmds/` (or `internal/resources/`) verifying that `KS_RESOURCE` bypasses cache lookup and still calls Upsert/VoteFor/Write

## 2. Sidecar types in internal/resurrect

- [x] 2.1 Create `internal/resurrect/types.go` with `ResurrectPane{Session, WindowIdx, PaneIdx, Kubeconfig, Resource, Verb string}` and `ResurrectState{Version int, Panes []ResurrectPane}`
- [x] 2.2 Add `Write(ksdir string)` and `Load(ksdir string)` to `ResurrectState` — JSON serde to `~/<ksdir>/.ks.resurrect.json`

## 3. Proc helpers in internal/resurrect

- [x] 3.1 Create `internal/resurrect/proc.go` with `ReadEnv(pid int) (map[string]string, error)` that reads `/proc/<pid>/environ` (null-delimited) into a map; return empty map + nil error on non-Linux (build tag or runtime check)
- [x] 3.2 Add `FindDescendantByName(pid int, name string) (int, error)` that walks `/proc/*/status` to find a process whose `PPid` chain leads to `pid` and whose `Name` matches; return 0 if not found
- [x] 3.3 Add `ReadCmdline(pid int) ([]string, error)` that reads `/proc/<pid>/cmdline` (null-delimited args)

## 4. Tmux pane listing helper in internal/resurrect

- [x] 4.1 Create `internal/resurrect/tmux.go` with `ListPanes() ([]TmuxPane, error)` that runs `tmux list-panes -a -F "#{session_name}\t#{window_index}\t#{pane_index}\t#{pane_id}\t#{pane_pid}\t#{pane_current_command}"` and parses output into `TmuxPane{Session, WindowIdx, PaneIdx int, PaneID, PanePID, CurrentCommand string}`
- [x] 4.2 Add `SendKeys(target, command string) error` that runs `tmux send-keys -t <target> <command> Enter`

## 5. ks save command

- [x] 5.1 Create `internal/cmds/resurrect_save_cmd.go` with `ResurrectSave(a args.Standard) error` that: lists panes, reads KUBECONFIG from `/proc`, looks up pane ID in resource cache, detects verb from k9s cmdline, builds `ResurrectState`, calls `state.Write(a.Directory)`
- [x] 5.2 Create `cmd/ks/resurrect_save.go` registering `ks save` Cobra command wired to `cmds.ResurrectSave`; add to `main.go`
- [x] 5.3 Verify `ks save` exits cleanly (zero) with no output when run outside tmux
- [x] 5.4 Verify `ks save` exits cleanly (zero) on non-Linux at runtime

## 6. ks restore command

- [x] 6.1 Create `internal/cmds/resurrect_restore_cmd.go` with `ResurrectRestore(a args.Standard) error` that: loads sidecar, lists current panes, matches by position, sends `KS_RESOURCE=<r> KUBECONFIG=<path> ks <verb>` (with resource) or `export KUBECONFIG=<path>` (without resource) via `SendKeys`
- [x] 6.2 Create `cmd/ks/resurrect_restore.go` registering `ks restore` Cobra command wired to `cmds.ResurrectRestore`; add to `main.go`
- [x] 6.3 Verify `ks restore` exits cleanly when sidecar file does not exist

## 7. Documentation

- [x] 7.1 Add `resurrect_save` and `resurrect_restore` to the Command Reference table in `CLAUDE.md`
- [x] 7.2 Add a note on tmux.conf wiring in `CLAUDE.md`: `set -g @resurrect-hook-pre-save 'ks save'` and `set -g @resurrect-hook-post-restore-all 'ks restore'`
- [x] 7.3 Add tmux-resurrect integration section to `README.md` covering setup, how it works, and caveats
