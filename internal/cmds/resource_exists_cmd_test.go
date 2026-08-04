package cmds

import (
	"github.com/bradfordwagner/ks/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("populatedNames", func() {
	It("returns only Probed+Populated entries, vote-sorted", func() {
		r := &resources.Resources{
			Names: []resources.ResourceEntry{
				{Name: "pods", Votes: 2, Probed: true, Populated: true},
				{Name: "randomsecrets.redhatcop.redhat.io", Votes: 5, Probed: true, Populated: true},
				{Name: "configmaps", Votes: 9, Probed: true, Populated: false}, // probed empty
				{Name: "secrets", Votes: 9, Probed: false, Populated: false},    // never probed
			},
		}
		Expect(populatedNames(r)).To(Equal([]string{
			"randomsecrets.redhatcop.redhat.io", // 5 votes
			"pods",                              // 2 votes
		}))
	})

	It("returns empty when nothing was probed as populated", func() {
		r := &resources.Resources{
			Names: []resources.ResourceEntry{
				{Name: "pods", Votes: 3, Probed: false},
				{Name: "services", Votes: 1, Probed: true, Populated: false},
			},
		}
		Expect(populatedNames(r)).To(BeEmpty())
	})
})
