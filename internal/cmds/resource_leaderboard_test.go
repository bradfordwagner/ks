package cmds_test

import (
	"os"
	"path/filepath"
	"io"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/bradfordwagner/ks/internal/resources"
)

func TestCmds(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmds Suite")
}

// captureLeaderboard redirects stdout and runs ResourceLeaderboard, returning the output.
func captureLeaderboard(dir string, all bool) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	a := args.Standard{Directory: dir}
	_ = cmds.ResourceLeaderboard(a, all)

	w.Close()
	os.Stdout = old

	data, _ := io.ReadAll(r)
	return string(data)
}

func writeResources(dir string, entries []resources.ResourceEntry) {
	r := resources.Resources{
		Version: 2,
		Names:   entries,
	}
	_, err := r.Write(dir)
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("ResourceLeaderboard", func() {
	var dir string

	BeforeEach(func() {
		dir = GinkgoT().TempDir()
	})

	Context("3.1 missing file", func() {
		It("prints a friendly message and returns nil", func() {
			// ensure no file exists
			_ = os.Remove(filepath.Join(dir, ".ks.resources.json"))

			a := args.Standard{Directory: dir}
			err := cmds.ResourceLeaderboard(a, false)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("3.2 default filter (votes >= 1 only)", func() {
		It("excludes zero-vote entries", func() {
			writeResources(dir, []resources.ResourceEntry{
				{Name: "pods", Votes: 3},
				{Name: "nodes", Votes: 0},
				{Name: "services", Votes: 1},
			})

			out := captureLeaderboard(dir, false)
			Expect(out).To(ContainSubstring("pods"))
			Expect(out).To(ContainSubstring("services"))
			Expect(out).NotTo(ContainSubstring("nodes"))
		})

		It("prints no-usage message when all votes are zero", func() {
			writeResources(dir, []resources.ResourceEntry{
				{Name: "pods", Votes: 0},
			})

			out := captureLeaderboard(dir, false)
			Expect(out).To(ContainSubstring("no usage data"))
		})
	})

	Context("3.3 --all flag", func() {
		It("includes zero-vote entries when all=true", func() {
			writeResources(dir, []resources.ResourceEntry{
				{Name: "pods", Votes: 2},
				{Name: "nodes", Votes: 0},
			})

			out := captureLeaderboard(dir, true)
			Expect(out).To(ContainSubstring("pods"))
			Expect(out).To(ContainSubstring("nodes"))
		})
	})

	Context("3.4 sort order", func() {
		It("sorts votes descending then name ascending", func() {
			writeResources(dir, []resources.ResourceEntry{
				{Name: "services", Votes: 2},
				{Name: "pods", Votes: 5},
				{Name: "configmaps", Votes: 3},
				{Name: "deployments", Votes: 3},
			})

			out := captureLeaderboard(dir, false)
			lines := strings.Split(strings.TrimSpace(out), "\n")
			// lines[0] = header, lines[1..] = data rows
			Expect(lines).To(HaveLen(5)) // header + 4 entries
			Expect(lines[1]).To(ContainSubstring("pods"))
			Expect(lines[2]).To(ContainSubstring("configmaps"))
			Expect(lines[3]).To(ContainSubstring("deployments"))
			Expect(lines[4]).To(ContainSubstring("services"))
		})
	})

	Context("3.5 table format", func() {
		It("has a header row and data rows with rank, name, and votes", func() {
			writeResources(dir, []resources.ResourceEntry{
				{Name: "pods", Votes: 7},
				{Name: "services", Votes: 2},
			})

			out := captureLeaderboard(dir, false)
			lines := strings.Split(strings.TrimSpace(out), "\n")
			Expect(lines[0]).To(ContainSubstring("RESOURCE"))
			Expect(lines[0]).To(ContainSubstring("VOTES"))
			Expect(lines[1]).To(ContainSubstring("1"))
			Expect(lines[1]).To(ContainSubstring("pods"))
			Expect(lines[1]).To(ContainSubstring("7"))
			Expect(lines[2]).To(ContainSubstring("2"))
			Expect(lines[2]).To(ContainSubstring("services"))
			Expect(lines[2]).To(ContainSubstring("2"))
		})
	})
})
