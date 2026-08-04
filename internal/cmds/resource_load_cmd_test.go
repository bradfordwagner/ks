package cmds

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resource_load naming", func() {
	Describe("isStandardGroup", func() {
		DescribeTable("classifying API groups",
			func(group string, expected bool) {
				Expect(isStandardGroup(group)).To(Equal(expected))
			},
			Entry("core (empty) group is standard", "", true),
			Entry("apps is standard", "apps", true),
			Entry("batch is standard", "batch", true),
			Entry("autoscaling is standard", "autoscaling", true),
			Entry("extensions is standard", "extensions", true),
			Entry("policy is standard", "policy", true),
			Entry("*.k8s.io groups are standard", "networking.k8s.io", true),
			Entry("rbac.authorization.k8s.io is standard", "rbac.authorization.k8s.io", true),
			Entry("vault-config-operator group is not standard", "redhatcop.redhat.io", false),
			Entry("cert-manager group is not standard", "cert-manager.io", false),
			Entry("argo group is not standard", "argoproj.io", false),
		)
	})

	Describe("qualifiedResourceName", func() {
		It("keeps core resources bare", func() {
			Expect(qualifiedResourceName("pods", "")).To(Equal("pods"))
		})

		It("keeps standard short-group resources bare", func() {
			Expect(qualifiedResourceName("deployments", "apps")).To(Equal("deployments"))
		})

		It("keeps *.k8s.io resources bare", func() {
			Expect(qualifiedResourceName("ingresses", "networking.k8s.io")).To(Equal("ingresses"))
		})

		It("qualifies CRD/vendor-group resources with the group", func() {
			Expect(qualifiedResourceName("randomsecrets", "redhatcop.redhat.io")).
				To(Equal("randomsecrets.redhatcop.redhat.io"))
		})
	})
})
