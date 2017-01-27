package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/semka/hysteria/models"
)

var _ = Describe("Models", func() {

	project := &Project{Path: "./examples", TestTimeout: 30}

	Context("Project", func() {
		It("should be able to fetch test suites", func() {
			suites, err := project.GetTestSuites()

			Expect(err).To(BeNil())
			Expect(suites).To(Not(BeEmpty()))
			Expect(len(suites)).To(Equal(1))
			Expect([]string{"tests"}).To(ContainElement(suites[0].Name()))
		})
	})

	Context("Test suite", func() {
		It("should be able to fetch files", func() {
			suite := TestSuite{Path: "./examples/tests", Project: project}
			tests, err := suite.GetTests()
			Expect(err).To(BeNil())
			Expect(tests).To(Not(BeEmpty()))
			for _, test := range tests {
				Expect([]string{"passing.sh", "failing.sh"}).To(
					ContainElement(test.Name()),
				)
			}
		})

		It("should be able to run all tests and collect results", func() {
			suite := TestSuite{Path: "./examples/tests", Project: project}
			results, err := suite.Run()
			Expect(err).To(BeNil())
			Expect(len(results)).To(Equal(2))

			successCount := 0
			failureCount := 0

			for _, res := range results {
				if res.IsPassed {
					successCount++
				} else {
					failureCount++
				}
			}

			Expect(successCount).To(Equal(1))
			Expect(failureCount).To(Equal(1))
		})
	})

	Context("Test file", func() {
		It("should be able to run actual shell scripts", func() {
			passingTest := TestFile{Path: "./examples/tests/passing.sh"}
			passingRes, _ := passingTest.Run()
			Expect(passingRes.IsPassed).To(BeTrue())
			Expect(passingRes.StdoutStr).To(Equal("Pass\n"))
			Expect(passingRes.StderrStr).To(BeEmpty())

			failingTest := TestFile{Path: "./examples/tests/failing.sh"}
			failingRes, _ := failingTest.Run()
			Expect(failingRes.IsPassed).To(BeFalse())
			Expect(failingRes.StdoutStr).To(Equal("Testing\n"))
			Expect(failingRes.StderrStr).To(Equal("Test Failed\n"))
		})
	})
})
