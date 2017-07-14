package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var _ = Describe("UserLogin", func() {
	var page *agouti.Page
	var session *gexec.Session

	var username, password, atc, team, outputDir string

	BeforeEach(func() {
		var found bool
		username, found = os.LookupEnv("GITHUB_USERNAME")
		Expect(found).To(BeTrue())
		password, found = os.LookupEnv("GITHUB_PASSWORD")
		Expect(found).To(BeTrue())
		atc, found = os.LookupEnv("ATC_URL")
		Expect(found).To(BeTrue())
		team, found = os.LookupEnv("TEAM_NAME")
		Expect(found).To(BeTrue())
		outputDir, found = os.LookupEnv("OUTPUT_DIR")
		Expect(found).To(BeTrue())

		var err error
		Expect(agoutiDriver).ToNot(BeNil())
		page, err = agoutiDriver.NewPage(agouti.Debug)
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command("fly", "-t", "temp", "login", "-c", atc, "-k", "-n", team)
		_, err = command.StdinPipe()
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		session.Kill()
		Expect(page.Destroy()).To(Succeed())
	})

	It("can GitHub auth", func() {
		Eventually(session).Should(Say("navigate to the following URL in your browser"))
		flyOutput := string(session.Out.Contents())

		re := regexp.MustCompile("[0-9]*")
		match := re.FindStringSubmatch(flyOutput)

		Expect(page.Navigate(atc + "/auth/github?team_name=" + team + "&fly_local_port=" + match[0])).To(Succeed())
		Expect(page.Title()).To(Equal("Sign in to GitHub · GitHub"))

		Eventually(page.FindByID("login_field")).Should(BeFound())
		Expect(page.FindByID("login_field").Fill(username)).To(Succeed())
		Expect(page.FindByID("password").Fill(password)).To(Succeed())
		Expect(page.FindByName("commit").Submit()).To(Succeed())

		Expect(page.HTML()).ToNot(ContainSubstring("This application has made an unusually high number of requests to access your account. Please reauthorize the application to continue."))
		Expect(page.HTML()).ToNot(ContainSubstring("This site can’t be reached"))
		Expect(page.HTML()).To(ContainSubstring("Bearer"))

		re = regexp.MustCompile("Bearer [ \t]*([^\n\r]*)")
		html, err := page.HTML()
		Expect(err).ToNot(HaveOccurred())
		token := re.FindStringSubmatch(html)
		Expect(len(token)).To(Equal(2))

		ioutil.WriteFile(outputDir+"bearer-token", []byte(token[1]), 0644)
	})
})
