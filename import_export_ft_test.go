package hoverfly_test

import (
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capturing, exporting, importing and simulating", func() {


	Describe("Import, Export GZip", func() {
		Context("The Gzipped response should be returned after exporting and importing", func() {

			BeforeEach(func() {
				// Spin up a fake server which returns hello world gzipped
				// Make a request to the endpoint
				// Export the data into a file
				// Wipe the records
				// Import the data from the file
				// Switch to simulate mode
				// Make the request
			})


			It("Returns a status code of 200", func() {

			})

			It("Returns hello world ungzipped", func() {

			})
		})
	})
});
