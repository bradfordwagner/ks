package cmds

import (
	"os"

	"github.com/bradfordwagner/ks/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resolveResourceType", func() {
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

		It("returns the env var value without fzf", func() {
			result, err := resolveResourceType(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("deployments"))
		})

		It("populates the pane cache via Upsert", func() {
			_, err := resolveResourceType(r)
			Expect(err).NotTo(HaveOccurred())
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

		It("returns the cached resource without fzf", func() {
			result, err := resolveResourceType(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("services"))
		})
	})
})
