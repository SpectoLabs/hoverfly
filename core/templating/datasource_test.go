package templating

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_NewDataSourceMethod(t *testing.T) {
	RegisterTestingT(t)

	dataSource, err := NewCsvDataSource("test-csv", "id,name,marks\n1,Test1,55\n2,Test2,56")

	Expect(err).To(BeNil())
	Expect(dataSource).NotTo(BeNil())
	Expect(dataSource.Name).To(Equal("test-csv"))
	Expect(dataSource.Data).To(HaveLen(3))

	Expect(dataSource.Data[0][0]).To(Equal("id"))
	Expect(dataSource.Data[0][1]).To(Equal("name"))
	Expect(dataSource.Data[0][2]).To(Equal("marks"))

	Expect(dataSource.Data[1][0]).To(Equal("1"))
	Expect(dataSource.Data[1][1]).To(Equal("Test1"))
	Expect(dataSource.Data[1][2]).To(Equal("55"))

	Expect(dataSource.Data[2][0]).To(Equal("2"))
	Expect(dataSource.Data[2][1]).To(Equal("Test2"))
	Expect(dataSource.Data[2][2]).To(Equal("56"))
}
