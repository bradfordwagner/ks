package cmds

import (
	"github.com/bradfordwagner/ks/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prioritizeCandidates", func() {
	entries := []resources.ResourceEntry{
		{Name: "pods", Votes: 5, Populated: true, Probed: true},                              // priority
		{Name: "randomsecrets.redhatcop.redhat.io", Votes: 2, Populated: true, Probed: true}, // priority
		{Name: "services", Votes: 9, Populated: false, Probed: true},                         // voted but not prev-populated -> rest
		{Name: "configmaps", Votes: 0, Populated: true, Probed: true},                        // prev-populated but no votes -> rest
	}

	It("dispatches voted + previously-populated first, each group by votes then name", func() {
		candidates := []string{"configmaps", "services", "pods", "randomsecrets.redhatcop.redhat.io"}
		Expect(prioritizeCandidates(entries, candidates)).To(Equal([]string{
			// priority group (votes>0 AND previously populated), votes desc
			"pods",
			"randomsecrets.redhatcop.redhat.io",
			// rest, votes desc then name
			"services",   // 9 votes
			"configmaps", // 0 votes
		}))
	})

	It("treats candidates absent from the cache as zero-vote rest, sorted by name", func() {
		candidates := []string{"newcrd.vendor.io", "pods", "another.vendor.io"}
		Expect(prioritizeCandidates(entries, candidates)).To(Equal([]string{
			"pods",              // priority
			"another.vendor.io", // rest, name asc
			"newcrd.vendor.io",
		}))
	})
})
