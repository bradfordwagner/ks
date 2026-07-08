package choose

import (
	"fmt"
	"github.com/koki-develop/go-fzf"
	"log"
	"strings"
)

// catppuccinMochaBlue is the "Blue" accent color from the Catppuccin Mocha palette,
// used to highlight the context indicator in picker prompts.
const catppuccinMochaBlue = "#89b4fa"

func One(opts []string, header string) (res string, err error) {
	choices, err := choose(opts, fzf.WithLimit(1), fzf.WithPrompt(prompt(header)), promptStyle())
	if err != nil {
		return
	} else if len(choices) > 1 {
		err = fmt.Errorf("expected 1 selection, got %s", strings.Join(choices, ", "))
		return
	}
	res = choices[0]

	return
}

func Multi(opts []string, header string) (res []string, err error) {
	return choose(opts,
		fzf.WithNoLimit(true),
		fzf.WithPrompt(prompt(header)),
		promptStyle(),
		fzf.WithKeyMap(fzf.KeyMap{
			Toggle: []string{"tab", "ctrl+s"},
		}),
	)
}

// prompt formats a picker header as a "[<header>] > " prompt prefix.
func prompt(header string) string {
	return fmt.Sprintf("[%s] > ", header)
}

// promptStyle colors the picker prompt in Catppuccin Mocha blue.
func promptStyle() fzf.Option {
	return fzf.WithStyles(fzf.WithStylePrompt(fzf.Style{ForegroundColor: catppuccinMochaBlue}))
}

func choose(opts []string, fzfOpts ...fzf.Option) (res []string, err error) {
	f, err := fzf.New(fzfOpts...)
	if err != nil {
		log.Fatal(err)
	}

	idxs, err := f.Find(opts, func(i int) string { return opts[i] })
	if err != nil {
		return
	}

	// convert index to value
	for _, idx := range idxs {
		res = append(res, opts[idx])
	}

	return
}
