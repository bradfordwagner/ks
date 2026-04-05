package cmds

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/resources"
)

func ResourceLeaderboard(a args.Standard, all bool) error {
	r, err := resources.LoadResources(a.Directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("no resource data found — run 'ks resource_load' first")
			return nil
		}
		return err
	}

	entries := make([]resources.ResourceEntry, len(r.Names))
	copy(entries, r.Names)

	// filter
	if !all {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Votes >= 1 {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if len(entries) == 0 {
		fmt.Println("no usage data yet — select resources with 'ks resource' to build the leaderboard")
		return nil
	}

	// sort: votes desc, name asc
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Votes != entries[j].Votes {
			return entries[i].Votes > entries[j].Votes
		}
		return entries[i].Name < entries[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tRESOURCE\tVOTES")
	for i, e := range entries {
		fmt.Fprintf(w, "%d\t%s\t%d\n", i+1, e.Name, e.Votes)
	}
	w.Flush()

	return nil
}
