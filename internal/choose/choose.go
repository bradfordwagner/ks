package choose

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ErrAbort is returned when the picker is dismissed (Esc / Ctrl-C / Ctrl-D) or
// nothing matches. Callers treat it as a no-op.
var ErrAbort = errors.New("choose: aborted")

// catppuccinMochaBlue is the "Blue" accent from Catppuccin Mocha, used to color
// the context indicator in the picker prompt.
const catppuccinMochaBlue = "#89b4fa"

// multiHeader is shown above the list for multi-select pickers as a key hint.
const multiHeader = "tab: select · enter: confirm"

// One shows a single-select picker and returns the chosen value. It forces
// single-select even when the user's FZF_DEFAULT_OPTS enables multi (`-m`).
func One(opts []string, header string) (string, error) {
	sel, err := run(strings.NewReader(joinLines(opts)), false, header)
	if err != nil {
		return "", err
	}
	if len(sel) == 0 {
		return "", ErrAbort
	}
	return sel[0], nil
}

// Multi shows a multi-select picker and returns the chosen values.
func Multi(opts []string, header string) ([]string, error) {
	return run(strings.NewReader(joinLines(opts)), true, header)
}

// MultiStream shows a multi-select picker whose candidates are streamed from names
// as they arrive; fzf reads its input incrementally. Inside tmux this renders in a
// popup (see fzfArgs). The caller closes names when done. The writer keeps draining
// names after fzf exits (e.g. the user aborts) so producers never block.
func MultiStream(names <-chan string, header string) ([]string, error) {
	cmd := exec.Command("fzf", fzfArgs(true, header)...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		w := bufio.NewWriter(stdin)
		broken := false
		for name := range names {
			if broken {
				continue // fzf is gone; keep draining so the producer never blocks
			}
			if _, werr := w.WriteString(name + "\n"); werr != nil {
				broken = true
				continue
			}
			if werr := w.Flush(); werr != nil {
				broken = true
			}
		}
		_ = stdin.Close()
	}()

	if err := cmd.Wait(); err != nil {
		return nil, classify(err)
	}
	return splitLines(out.String()), nil
}

// run pipes opts into fzf and returns the selected lines.
func run(stdin *strings.Reader, multi bool, header string) ([]string, error) {
	cmd := exec.Command("fzf", fzfArgs(multi, header)...)
	cmd.Stdin = stdin
	cmd.Stderr = os.Stderr
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, classify(err)
	}
	return splitLines(out.String()), nil
}

// fzfArgs builds the fzf flags. FZF_DEFAULT_OPTS is inherited automatically; these
// flags take precedence, so `--no-multi` overrides a default `-m` for One(). Inside
// tmux the picker renders in a centered popup (`--tmux`), which tmux composites into
// its own screen region — sidestepping the full-screen alt-screen redraw corruption
// some terminals show, especially for incremental (streaming) updates.
func fzfArgs(multi bool, header string) []string {
	args := []string{
		"--prompt", prompt(header),
		"--color", "prompt:" + catppuccinMochaBlue,
	}
	if os.Getenv("TMUX") != "" {
		args = append(args, "--tmux", "center,80%,80%")
	}
	if multi {
		args = append(args, "--multi", "--header", multiHeader)
	} else {
		args = append(args, "--no-multi")
	}
	return args
}

// classify maps fzf's exit codes to errors: 1 (no match) and 130 (interrupt) are
// treated as an abort; anything else is a real error.
func classify(err error) error {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		switch exit.ExitCode() {
		case 1, 130:
			return ErrAbort
		}
	}
	return err
}

// prompt formats a picker header as a "[<header>] > " prompt prefix.
func prompt(header string) string {
	return fmt.Sprintf("[%s] > ", header)
}

// joinLines renders options as newline-terminated fzf input ("" when empty).
func joinLines(opts []string) string {
	if len(opts) == 0 {
		return ""
	}
	return strings.Join(opts, "\n") + "\n"
}

// splitLines splits fzf's output into selected values, dropping the trailing newline.
func splitLines(s string) []string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
