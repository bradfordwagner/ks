package resources

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func writeTempFile(dir, content string) {
	Expect(os.WriteFile(filepath.Join(dir, CacheFile), []byte(content), Perms)).To(Succeed())
}

var _ = Describe("Resources", func() {
	Describe("LoadResources", func() {
		Context("v1 migration", func() {
			It("migrates flat string names to ResourceEntry with zero votes", func() {
				dir := GinkgoT().TempDir()
				writeTempFile(dir, `{"names":["pods","services","deployments"]}`)

				r, err := LoadResources(dir)
				Expect(err).NotTo(HaveOccurred())
				Expect(r.Version).To(Equal(schemaV2))
				Expect(r.Names).To(HaveLen(3))
				for _, e := range r.Names {
					Expect(e.Votes).To(Equal(0), "entry %q should have 0 votes", e.Name)
				}
				Expect(r.Names).To(ContainElements(
					HaveField("Name", "pods"),
					HaveField("Name", "services"),
					HaveField("Name", "deployments"),
				))
			})

			It("is idempotent — loading a v2 file does not reset votes", func() {
				dir := GinkgoT().TempDir()
				v2 := Resources{
					Version: schemaV2,
					Names: []ResourceEntry{
						{Name: "pods", Votes: 5},
						{Name: "services", Votes: 2},
					},
				}
				_, err := v2.Write(dir)
				Expect(err).NotTo(HaveOccurred())

				r, err := LoadResources(dir)
				Expect(err).NotTo(HaveOccurred())
				Expect(r.Version).To(Equal(schemaV2))
				Expect(r.Names).To(HaveLen(2))
				Expect(r.Names[0].Votes).To(Equal(5))
				Expect(r.Names[1].Votes).To(Equal(2))
			})
		})

		Context("v2 JSON", func() {
			It("reads version and votes correctly", func() {
				dir := GinkgoT().TempDir()
				writeTempFile(dir, `{"version":2,"names":[{"name":"pods","votes":4},{"name":"services","votes":0}]}`)

				r, err := LoadResources(dir)
				Expect(err).NotTo(HaveOccurred())
				Expect(r.Version).To(Equal(schemaV2))
				Expect(r.Names).To(HaveLen(2))
				Expect(r.Names[0]).To(Equal(ResourceEntry{Name: "pods", Votes: 4}))
				Expect(r.Names[1]).To(Equal(ResourceEntry{Name: "services", Votes: 0}))
			})

			It("preserves votes across write+load round-trip", func() {
				dir := GinkgoT().TempDir()
				v2 := Resources{
					Version: schemaV2,
					Names: []ResourceEntry{
						{Name: "pods", Votes: 7},
						{Name: "services", Votes: 1},
					},
				}
				_, err := v2.Write(dir)
				Expect(err).NotTo(HaveOccurred())

				r, err := LoadResources(dir)
				Expect(err).NotTo(HaveOccurred())
				Expect(r.Names).To(ContainElement(Equal(ResourceEntry{Name: "pods", Votes: 7})))
				Expect(r.Names).To(ContainElement(Equal(ResourceEntry{Name: "services", Votes: 1})))
			})
		})
	})

	Describe("VoteFor", func() {
		It("increments the vote for the named resource", func() {
			r := Resources{
				Names: []ResourceEntry{
					{Name: "pods", Votes: 3},
					{Name: "services", Votes: 1},
				},
			}
			r.VoteFor("pods")
			Expect(r.Names[0].Votes).To(Equal(4))
			Expect(r.Names[1].Votes).To(Equal(1))
		})

		It("is a no-op for unknown resource names", func() {
			r := Resources{
				Names: []ResourceEntry{
					{Name: "pods", Votes: 2},
				},
			}
			r.VoteFor("nonexistent")
			Expect(r.Names[0].Votes).To(Equal(2))
		})
	})

	Describe("SortedNames", func() {
		It("sorts by votes descending then alphabetically ascending", func() {
			r := Resources{
				Names: []ResourceEntry{
					{Name: "services", Votes: 2},
					{Name: "pods", Votes: 5},
					{Name: "configmaps", Votes: 3},
					{Name: "deployments", Votes: 3},
				},
			}
			Expect(r.SortedNames()).To(Equal([]string{"pods", "configmaps", "deployments", "services"}))
		})

		It("sorts alphabetically when all votes are zero", func() {
			r := Resources{
				Names: []ResourceEntry{
					{Name: "services", Votes: 0},
					{Name: "pods", Votes: 0},
					{Name: "configmaps", Votes: 0},
				},
			}
			Expect(r.SortedNames()).To(Equal([]string{"configmaps", "pods", "services"}))
		})
	})
})
