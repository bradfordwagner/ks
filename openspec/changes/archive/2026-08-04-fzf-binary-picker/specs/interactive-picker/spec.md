## Purpose

Defines how ks presents interactive selection to the user — backed by the `fzf` binary so pickers honor the user's own fzf configuration — with single-select, multi-select, and live-streaming modes and consistent abort handling.

## ADDED Requirements

### Requirement: Pickers use the fzf binary and honor user config

ks pickers SHALL run the `fzf` binary, feeding candidates on stdin and reading the selection from stdout. Because `fzf` reads `FZF_DEFAULT_OPTS`, the user's own options and key binds SHALL apply to ks pickers.

#### Scenario: User key binds apply
- **WHEN** the user has `FZF_DEFAULT_OPTS` binds (e.g. `ctrl-a:select-all`) and opens any ks picker
- **THEN** those binds are active in the picker

### Requirement: Popup rendering inside tmux

When ks runs inside tmux, pickers SHALL render in a centered tmux popup (`fzf --tmux`) so tmux owns the popup's screen region. Outside tmux, pickers SHALL render normally (full-screen fzf).

#### Scenario: Picker opens as a popup in tmux
- **WHEN** any ks picker is invoked from inside a tmux session
- **THEN** it renders in a centered tmux popup

#### Scenario: Picker renders normally outside tmux
- **WHEN** a picker is invoked outside tmux
- **THEN** it renders as a normal full-screen fzf without requiring tmux

### Requirement: Streaming multi-select

A streaming multi-select picker SHALL accept candidates emitted over time and display them as they arrive, without waiting for the full set, and return the chosen values.

#### Scenario: Options appear as they are produced
- **WHEN** a caller streams candidates into the picker as a background task produces them
- **THEN** the picker shows candidates as they arrive and returns those the user selects

### Requirement: Single-select is forced for One

A single-select picker SHALL return exactly one value and SHALL force single-select even when the user's `FZF_DEFAULT_OPTS` enables multi-select (`-m`).

#### Scenario: One returns a single value despite a multi default
- **WHEN** `FZF_DEFAULT_OPTS` contains `-m` and a single-select picker is shown
- **THEN** the user can select only one item and that one value is returned

### Requirement: Consistent abort handling

Dismissing a picker (Esc / Ctrl-C) or matching nothing SHALL return a single sentinel error (`choose.ErrAbort`) that callers treat as a no-op; any other failure SHALL be surfaced as a real error.

#### Scenario: Escape is a no-op
- **WHEN** the user presses Esc in any ks picker
- **THEN** the caller receives `ErrAbort` and performs no action

#### Scenario: A genuine failure is surfaced
- **WHEN** `fzf` fails for a reason other than abort or no-match
- **THEN** the error is returned to the caller rather than being treated as an abort
