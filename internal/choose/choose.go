package choose

import (
	"fmt"
	"github.com/koki-develop/go-fzf"
	"log"
	"strings"
)

func One(opts []string) (res string, err error) {
	choices, err := choose(opts, fzf.WithLimit(1))
	if err != nil {
		return
	} else if len(choices) > 1 {
		err = fmt.Errorf("expected 1 selection, got %s", strings.Join(choices, ", "))
		return
	}
	res = choices[0]

	return
}

func Multi(opts []string) (res []string, err error) {
	return choose(opts,
		fzf.WithNoLimit(true),
		fzf.WithKeyMap(fzf.KeyMap{
			Toggle: []string{"tab", "ctrl+s"},
		}),
	)
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
