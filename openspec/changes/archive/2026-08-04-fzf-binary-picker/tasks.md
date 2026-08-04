## 1. Picker rewrite

- [x] 1.1 Rewrite `internal/choose` on the `fzf` binary: `One` (`--no-multi`), `Multi`, `MultiStream` (channel-fed via `StdinPipe`), `choose.ErrAbort`, `classify` (exit 1/130 → abort), prompt color, multi `--header`.
- [x] 1.2 Render in a centered `fzf --tmux` popup when `$TMUX` is set (full-screen fallback otherwise).

## 2. Callers

- [x] 2.1 Switch all 11 callers (`cmd/ks/main.go` + 10 `internal/cmds/*`) from `fzf.ErrAbort` to `choose.ErrAbort` and drop the go-fzf import.

## 3. Dependency & docs

- [x] 3.1 `go mod tidy` to remove `github.com/koki-develop/go-fzf`.
- [x] 3.2 Update `CLAUDE.md` (`internal/choose` role, `resource_live` streaming).

## 4. Verify

- [x] 4.1 `go build` / `go test ./...` / `go vet` / `go install`; non-interactive fzf smoke test (arg sets accepted; no-match→1, bad-flag→2).
- [ ] 4.2 Manually verify interactively (needs a tty): pickers open as tmux popups and honor `FZF_DEFAULT_OPTS` binds; `ks resource_live` streams populated types into the popup **without render artifacts**.
