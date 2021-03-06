package generic_test

import (
	"fmt"
	"integration-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("MasterTlsCertificate", func() {

	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	DescribeTable("hostnames", func(hostname string) {
		url := fmt.Sprintf("https://%s", hostname)
		session := runner.RunKubectlCommandInNamespace("default", "run", "test-master-cert-via-curl", "--image=tutum/curl", "--restart=Never", "-ti", "--rm", "--", "curl", url, "--cacert", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		Eventually(session, "5m").Should(gexec.Exit(0))
	},
		Entry("kubernetes", "kubernetes"),
		Entry("kubernetes.default", "kubernetes.default"),
		Entry("kubernetes.default.svc", "kubernetes.default.svc"),
		Entry("kubernetes.default.svc.cluster.local", "kubernetes.default.svc.cluster.local"),
	)
})
