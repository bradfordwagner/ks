package cmds

import (
	"os"

	"github.com/bradfordwagner/ks/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("singleResource", func() {
	var r *resources.Resources

	BeforeEach(func() {
		r = &resources.Resources{
			Version: 2,
			Names: []resources.ResourceEntry{
				{Name: "pods", Votes: 3},
				{Name: "services", Votes: 1},
			},
		}
		os.Unsetenv("KS_RESOURCE")
		os.Unsetenv("TMUX")
		os.Unsetenv("TMUX_PANE")
	})

	Context("KS_RESOURCE env var is set", func() {
		BeforeEach(func() {
			os.Setenv("KS_RESOURCE", "deployments")
			os.Setenv("TMUX", "/tmp/tmux-test,1234,0")
			os.Setenv("TMUX_PANE", "%99")
		})

		AfterEach(func() {
			os.Unsetenv("KS_RESOURCE")
			os.Unsetenv("TMUX")
			os.Unsetenv("TMUX_PANE")
		})

		It("returns the env var value and caches it to the pane", func() {
			result, ok := singleResource(r)
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("deployments"))
			Expect(r.Get()).To(Equal("deployments"))
		})
	})

	Context("KS_RESOURCE is not set and cache has a hit", func() {
		BeforeEach(func() {
			os.Setenv("TMUX", "/tmp/tmux-test,1234,0")
			os.Setenv("TMUX_PANE", "%42")
			r.Upsert("services")
		})

		AfterEach(func() {
			os.Unsetenv("TMUX")
			os.Unsetenv("TMUX_PANE")
		})

		It("returns the cached resource", func() {
			result, ok := singleResource(r)
			Expect(ok).To(BeTrue())
			Expect(result).To(Equal("services"))
		})
	})

	Context("no env var and no cache hit", func() {
		It("reports ok=false so the caller prompts", func() {
			_, ok := singleResource(r)
			Expect(ok).To(BeFalse())
		})
	})
})

var _ = Describe("resourceVerb / resourcePaneCommand", func() {
	It("maps the all flag to the ks verb", func() {
		Expect(resourceVerb(false)).To(Equal("resource"))
		Expect(resourceVerb(true)).To(Equal("resource_all"))
	})

	It("composes the KS_RESOURCE launch line", func() {
		Expect(resourcePaneCommand("randomsecrets.redhatcop.redhat.io", resourceVerb(true))).
			To(Equal("KS_RESOURCE=randomsecrets.redhatcop.redhat.io ks resource_all"))
	})
})
