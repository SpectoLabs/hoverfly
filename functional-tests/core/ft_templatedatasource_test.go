package hoverfly_test

import (
	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manage template data source in hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("get template data source", func() {

		Context("hoverfly with template data source", func() {

			BeforeEach(func() {
				hoverfly.Start("-templating-data-source", "test-csv testdata/test-student-data.csv")
			})

			It("Should return template data sources", func() {
				templateDataSourceView := hoverfly.GetAllDataSources()
				Expect(templateDataSourceView).NotTo(BeNil())
				Expect(templateDataSourceView.DataSources).To(HaveLen(1))
				Expect(templateDataSourceView.DataSources[0].Name).To(Equal("test-csv"))
				Expect(templateDataSourceView.DataSources[0].Data).To(Equal("Id,Name,Marks\n1,Test1,45\n2,Test2,55\n3,Test3,67\n4,Test4,89\n*,NA,ABSENT\n"))
			})
		})
	})
})
