package workload_test

import (
	"fmt"
	"net/http"
	"time"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var ipAddressError = "No IP address found for service"

func getServiceIP() string {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	validIP := regexp.MustCompile(numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock) 

	for {
		select {
			case <-timeout:
				return ipAddressError
			case <-tick:
				getServiceIp := runner.RunKubectlCommand("get", "service", "nginx", "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'")
	                        Eventually(getServiceIp, "5s").Should(gexec.Exit())
        	                serviceIP := string(getServiceIp.Out.Contents())
				// Remove quotes from the IP address
				serviceIP = serviceIP[1:len(serviceIP)-1]
				if validIP.MatchString(serviceIP) {
					return serviceIP
				}
		}
	}
}

var _ = Describe("Deploy workload", func() {

	It("exposes routes via LBs", func() {

		if (iaas == "gcp") {
	                deployNginx := runner.RunKubectlCommand("create", "-f", nginxLBSpec)
	                Eventually(deployNginx, "60s").Should(gexec.Exit(0))
	                rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
	                Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
	
	                serviceIP := getServiceIP()
	                Expect(serviceIP).To(Not(Equal(ipAddressError)))
	                appUrl := fmt.Sprintf("http://%s", serviceIP)
	
	                timeout := time.Duration(5 * time.Second)
	                httpClient := http.Client{
	                        Timeout: timeout,
	                }
	
	                Eventually(func() int {
	                       	result, err := httpClient.Get(appUrl)
	                       	if err != nil {
	                       	        return -1
	                       	}
	                       	return result.StatusCode
                	}, "120s", "5s").Should(Equal(200))
		} else {
			// TODO Once we have inabled cloud provider packages on 
			// vSphere and AWS, this else block can go away
			appUrl := fmt.Sprintf("http://%s:%s", workerAddress, nodePort)
	
			timeout := time.Duration(5 * time.Second)
			httpClient := http.Client{
				Timeout: timeout,
			}
	
			_, err := httpClient.Get(appUrl)
			Expect(err).To(HaveOccurred())
	
			deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
			Eventually(deployNginx, "60s").Should(gexec.Exit(0))
			rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
			Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
	
			Eventually(func() int {
				result, err := httpClient.Get(appUrl)
				if err != nil {
					return -1
				}
				return result.StatusCode
			}, "120s", "5s").Should(Equal(200))
		}

	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("30s")
	})

})
