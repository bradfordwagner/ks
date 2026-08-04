## Context

`internal/choose` wrapped `go-fzf` (a Go reimplementation). It ignored `FZF_DEFAULT_OPTS`, exposed only 5 fixed key actions, and its hot-reload races its own renderer. The user drives everything else through the real `fzf` binary (with rich `FZF_DEFAULT_OPTS` binds), which is always on their `PATH`. See proposal.md.

## Goals / Non-Goals

**Goals:**
- ks pickers behave exactly like the user's other fzf usages (config + binds inherited for free).
- Clean rendering (including incremental/streaming updates) in the target terminal.
- No churn for callers beyond the abort sentinel.

**Non-Goals:**
- A custom in-picker filter hotkey or bespoke keymap (the user's `FZF_DEFAULT_OPTS` owns that now).
- Reimplementing fzf features; we shell out.

## Decisions

### D1: Shell out; stdin = candidates, stdout = selection; popup in tmux

`run` pipes candidates to `fzf` on stdin and captures the selection from stdout — the canonical `list | fzf` contract. `MultiStream` uses `StdinPipe` + a writer goroutine so candidates can be produced over time. When `$TMUX` is set, `fzf --tmux center,80%,80%` renders the finder in a tmux popup.
- **Why the popup:** full-screen fzf drives the terminal's alternate screen directly; in the target stack (WSL2 + Windows Terminal + tmux) that corrupts on incremental redraws (streaming). A tmux popup makes *tmux* own and composite the popup's region, so redraws render cleanly. Outside tmux, fzf falls back to full-screen.

### D2: Exit-code → error mapping

`classify` maps exit 1 (no match) and 130 (interrupt) to `choose.ErrAbort`; other non-zero exits surface as real errors. `One` forces `--no-multi` (overrides a default `-m`); prompt keeps the Catppuccin blue via `--color prompt:#89b4fa`; multi shows a `tab: select · enter: confirm` header.
- **Why:** preserves the prior `ErrAbort` no-op semantics callers already rely on, and keeps the prompt styling.

### D3: Streaming writer drains after fzf exits

`MultiStream`'s writer goroutine keeps ranging the channel after fzf exits (e.g. abort), discarding once the pipe breaks, so producers never block; the caller cancels its context to stop work.
- **Why:** avoids a goroutine/producer deadlock when the user aborts mid-scan.
- **Note:** streaming previously corrupted under full-screen fzf; it is only reliable now because of the popup (D1). If `fzf --tmux` buffers stdin rather than reading it incrementally, the worst case is a clean non-live render (the popup appears once the scan finishes) — never the earlier corruption.

## Risks / Trade-offs

- **Hard dependency on the `fzf` binary.** → Already required by the user's environment; a missing binary surfaces a clear exec error.
- **Cannot unit-test the interactive TUI (needs a tty).** → Flag/exit-code plumbing is validated non-interactively (`fzf --filter`, exit codes); interactive behavior verified manually.
- **`FZF_DEFAULT_OPTS` now influences ks pickers.** → Intended; `--no-multi` and explicit flags guard the cases that must not vary.
