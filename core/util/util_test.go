package util

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GetRequestBody_GettingTheRequestBodyGetsTheCorrectData(t *testing.T) {
	RegisterTestingT(t)

	request := &http.Request{}
	request.Body = io.NopCloser(bytes.NewBuffer([]byte("test")))

	requestBody, err := GetRequestBody(request)
	Expect(err).To(BeNil())

	Expect(requestBody).To(Equal("test"))
}

func Test_GetRequestBody_GettingTheRequestBodySetsTheSameBodyAgain(t *testing.T) {
	RegisterTestingT(t)

	request := &http.Request{}
	request.Body = io.NopCloser(bytes.NewBuffer([]byte("test-preserve")))

	_, err := GetRequestBody(request)
	Expect(err).To(BeNil())

	newRequestBody, err := io.ReadAll(request.Body)
	Expect(err).To(BeNil())

	Expect(string(newRequestBody)).To(Equal("test-preserve"))
}

func Test_GetRequestBody_DecompressGzipContent(t *testing.T) {
	RegisterTestingT(t)

	request, _ := http.NewRequest("POST", "", nil)
	request.Header.Set("Content-Encoding", "gzip")
	originalBody := "hello_world"

	compressedBody, err := CompressGzip([]byte(originalBody))
	Expect(err).To(BeNil())
	Expect(string(compressedBody)).To(Not(Equal(originalBody)))
	request.Body = io.NopCloser(bytes.NewBuffer(compressedBody))

	_, err = GetRequestBody(request)
	Expect(err).To(BeNil())

	newRequestBody, err := io.ReadAll(request.Body)
	Expect(err).To(BeNil())

	Expect(string(newRequestBody)).To(Equal(originalBody))
}

func Test_GetResponseBody_GettingTheResponseBodyGetsTheCorrectData(t *testing.T) {
	RegisterTestingT(t)

	response := &http.Response{}
	response.Body = io.NopCloser(bytes.NewBuffer([]byte("test")))

	responseBody, err := GetResponseBody(response)
	Expect(err).To(BeNil())

	Expect(responseBody).To(Equal("test"))

}

func Test_GetResponseBody_GettingTheResponseBodySetsTheSameBodyAgain(t *testing.T) {
	RegisterTestingT(t)

	response := &http.Response{}
	response.Body = io.NopCloser(bytes.NewBuffer([]byte("test-preserve")))

	_, err := GetResponseBody(response)
	Expect(err).To(BeNil())

	newResponseBody, err := io.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(newResponseBody)).To(Equal("test-preserve"))
}

func Test_SortQueryString_ReordersQueryStringAlphabetically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=e&d=d")).To(Equal("d=d&e=e"))
}

func Test_SortQueryString_ReordersQueryValuesAlphabetically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=e&e=d")).To(Equal("e=d&e=e"))
}

func Test_SortQueryString_ReordersQueryValuesNumerically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=2&e=1")).To(Equal("e=1&e=2"))
}

func Test_SortQueryString_ReordersQueryValuesAlphanumerically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=2&e=d&e=1&e=e")).To(Equal("e=1&e=2&e=d&e=e"))
}

func Test_SortQueryString_KeepsAsteriskInTact(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("&e=*")).To(Equal("e=*"))
}

func Test_SortQueryString_PreservesEqualsAndEmptyValueQuery(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=")).To(Equal("e="))
}

func Test_SortQueryString_PreservesNoEqualsAndEmptyValueQuery(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e")).To(Equal("e"))
}

func Test_SortQueryString_PreservesBothEqualsAndNoEqualsWithEmptyValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("a&b&c=&d&e=&f=")).To(Equal("a&b&c=&d&e=&f="))
}

func Test_GetContentTypeFromHeaders_ReturnsEmptyStringIfHeadersAreNil(t *testing.T) {
	RegisterTestingT(t)

	Expect(GetContentTypeFromHeaders(nil)).To(Equal(""))
}

func Test_GetContentTypeFromHeaders_ReturnsEmptyStringIfHeadersAreEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(GetContentTypeFromHeaders(map[string][]string{})).To(Equal(""))
}

func Test_GetContentTypeFromHeaders_ReturnsJsonIfJson(t *testing.T) {
	RegisterTestingT(t)

	Expect(GetContentTypeFromHeaders(map[string][]string{
		"Content-Type": {"application/json"},
	})).To(Equal("json"))
}

func Test_GetContentTypeFromHeaders_ReturnsXmlIfXml(t *testing.T) {
	RegisterTestingT(t)

	Expect(GetContentTypeFromHeaders(map[string][]string{
		"Content-Type": {"application/xml"},
	})).To(Equal("xml"))
}

func Test_JSONMarshal_MarshalsIntoJson(t *testing.T) {
	RegisterTestingT(t)

	jsonBytes, err := JSONMarshal(map[string]string{
		"test": "testing",
	})

	Expect(err).To(BeNil())
	Expect(string(jsonBytes)).To(Equal(`{"test":"testing"}` + "\n"))
}

func Test_MinifyJson_MinifiesJsonString(t *testing.T) {
	RegisterTestingT(t)

	Expect(MinifyJson(`{
		"test": {
			"something": [
				1, 2, 3
			]
		}
	}`)).To(Equal(`{"test":{"something":[1,2,3]}}`))
}

func Test_MinifyJson_ErrorsOnInvalidJsonString(t *testing.T) {
	RegisterTestingT(t)

	_, err := MinifyJson(`{
		"test": {
			"something":
				1, 2, 3
			
		}
	}`)

	Expect(err).ToNot(BeNil())
}

func Test_MinifyXml_MinifiesXmlString(t *testing.T) {
	RegisterTestingT(t)

	Expect(MinifyXml(`<xml>
		<document  key="value">test</document>
	</xml>`)).To(Equal(`<xml><document key="value">test</document></xml>`))
}

func Test_MinifyXml_SimplifiesXmlString(t *testing.T) {
	RegisterTestingT(t)

	Expect(MinifyXml(`<xml>
		<document></document>
	</xml>`)).To(Equal(`<xml><document/></xml>`))
}

func Test_CopyMap(t *testing.T) {
	RegisterTestingT(t)

	originalMap := make(map[string]string)
	originalMap["first"] = "1"
	originalMap["second"] = "2"

	newMap := CopyMap(originalMap)

	delete(originalMap, "first")
	originalMap["second"] = ""
	originalMap["third"] = "3"

	Expect(newMap).To(HaveLen(2))
	Expect(newMap["first"]).To(Equal("1"))
	Expect(newMap["second"]).To(Equal("2"))
}

func Test_JsonPathMethod_WithBigFloatingNumber(t *testing.T) {

	RegisterTestingT(t)
	res := jsonPath("$.registrant", `{"registrant":2343534534534.345}`)
	Expect(res).To(Equal("2343534534534.345"))
}

func Test_JsonPathMethod_WithStringContainingLargeNumber(t *testing.T) {

	RegisterTestingT(t)
	res := jsonPath("$.registrant", `{"registrant":"0009007199254740999"}`)
	Expect(res).To(Equal("0009007199254740999"))
}

func Test_JsonPathMethod_WithWordContainingLetterE(t *testing.T) {

	RegisterTestingT(t)
	res := jsonPath("$.registrant", `{"registrant":"ETest"}`)
	Expect(res).To(Equal("ETest"))
}

func Test_Identical_ReturnsTrue_WithExactlySameArray(t *testing.T) {
	RegisterTestingT(t)
	first := [2]string{"q1", "q2"}
	second := [2]string{"q1", "q2"}

	Expect(Identical(first[:], second[:])).To(BeTrue())

}

func Test_Identical_ReturnsFalseWithDifferentArrayOfDifferentLength(t *testing.T) {
	RegisterTestingT(t)
	first := [2]string{"q1", "q2"}
	second := [3]string{"q1", "q2", "q3"}

	Expect(Identical(first[:], second[:])).To(BeFalse())

}

func Test_Identical_ReturnsFalseWithDifferentArrayOfSameLength(t *testing.T) {
	RegisterTestingT(t)
	first := [3]string{"q1", "q2", "q3"}
	second := [3]string{"q1", "q2", "q4"}

	Expect(Identical(first[:], second[:])).To(BeFalse())

}

func Test_Contains_ReturnsFalseWithEmptyArrayMatcher(t *testing.T) {

	RegisterTestingT(t)
	first := [0]string{}
	second := [2]string{"q1", "q2"}

	Expect(Contains(first[:], second[:])).To(BeFalse())

}

func Test_Contains_ReturnsTrueWithContainingBothValues(t *testing.T) {
	RegisterTestingT(t)
	first := [2]string{"q1", "q2"}
	second := [2]string{"q1", "q2"}

	Expect(Contains(first[:], second[:])).To(BeTrue())

}

func Test_Contains_ReturnsTrueWithContainingOneOfValues(t *testing.T) {
	RegisterTestingT(t)
	first := [3]string{"q1", "q2", "q3"}
	second := [4]string{"q5", "q6", "q7", "q1"}

	Expect(Contains(first[:], second[:])).To(BeTrue())

}

func Test_Contains_ReturnsFalseWithContainingNoneOfValuesSpecified(t *testing.T) {
	RegisterTestingT(t)
	first := [3]string{"q1", "q2", "q3"}
	second := [5]string{"q5", "q6", "q7", "q8", "q9"}

	Expect(Contains(first[:], second[:])).To(BeFalse())

}

func Test_ContainsOnly_ReturnsFalseWithEmptyArrayMatcher(t *testing.T) {

	RegisterTestingT(t)
	first := [0]string{}
	second := [2]string{"q1", "q2"}

	Expect(Contains(first[:], second[:])).To(BeFalse())

}

func Test_ContainsOnly_ReturnsTrueWithArrayContainingOnlyValuesWithDups(t *testing.T) {
	RegisterTestingT(t)
	first := [4]string{"a", "b", "c", "d"}
	second := [5]string{"a", "b", "b", "c", "b"}

	Expect(ContainsOnly(first[:], second[:])).To(BeTrue())

}

func Test_ContainsOnly_ReturnsTrueWithArrayInDifferentOrder(t *testing.T) {
	RegisterTestingT(t)
	first := [4]string{"c", "b", "a", "d"}
	second := [4]string{"a", "b", "c", "d"}

	Expect(ContainsOnly(first[:], second[:])).To(BeTrue())

}

func Test_ContainsOnly_ReturnsTrueWithIdenticalArray(t *testing.T) {
	RegisterTestingT(t)
	first := [3]string{"a", "b", "c"}
	second := [3]string{"a", "b", "c"}

	Expect(ContainsOnly(first[:], second[:])).To(BeTrue())

}

func Test_ContainsOnly_ReturnsTrueWithSubsetOfValues(t *testing.T) {

	RegisterTestingT(t)
	first := [3]string{"a", "b", "c"}
	second := [3]string{"a", "b", "a"}

	Expect(ContainsOnly(first[:], second[:])).To(BeTrue())

}

func Test_ContainsOnly_ReturnFalseWithOneExtraValue(t *testing.T) {

	RegisterTestingT(t)
	first := [3]string{"a", "b", "c"}
	second := [4]string{"a", "b", "a", "d"}

	Expect(ContainsOnly(first[:], second[:])).To(BeFalse())

}

func TestResolveAndValidatePath(t *testing.T) {
	RegisterTestingT(t)

	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	tests := []struct {
		basePath     string
		relativePath string
		expected     string
		shouldErr    bool
	}{
		{"/home/user/project", "subdir/file.txt", "/home/user/project/subdir/file.txt", false},
		{"/home/user/project", "../subdir/file.txt", "", true},
		{"/home/user/project", "../../etc/passwd", "", true},
		{"/home/user/project", "./subdir/file.txt", "/home/user/project/subdir/file.txt", false},
		{"/home/user/project", "subdir/../file.txt", "/home/user/project/file.txt", false},
		{"/home/user/project", ".", "/home/user/project", false},
		{"", "subdir/file.txt", filepath.Join(cwd, "subdir/file.txt"), false},
		{"", "", cwd, false},
		{"", ".", cwd, false},
		{"home/user/project", "subdir/file.txt", filepath.Join(cwd, "home/user/project/subdir/file.txt"), false},
		{"./home/user/project", "subdir/file.txt", filepath.Join(cwd, "home/user/project/subdir/file.txt"), false},
	}

	for _, test := range tests {
		t.Run(test.relativePath, func(t *testing.T) {
			result, err := ResolveAndValidatePath(test.basePath, test.relativePath)
			if test.shouldErr {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				// Convert paths to absolute to avoid issues with relative paths in tests.
				expectedAbs, _ := filepath.Abs(test.expected)
				resultAbs, _ := filepath.Abs(result)
				Expect(resultAbs).To(Equal(expectedAbs))
			}
		})
	}
}


func TestTruncateStringWithEllipsis(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		maxSize  int
		expected string
	}{
		{
			name:     "No truncation required",
			input:    "Hello",
			maxSize:  10,
			expected: "Hello",
		},
		{
			name:     "Truncate with ellipsis",
			input:    "Hello, World!",
			maxSize:  10,
			expected: "Hello, ...",
		},
		{
			name:     "String exactly maxSize",
			input:    "Hello, World!",
			maxSize:  13,
			expected: "Hello, World!",
		},
		{
			name:     "Small maxSize adds ellipsis only",
			input:    "Hello",
			maxSize:  3,
			expected: "...",
		},
		{
			name:     "UTF-8 truncation valid",
			input:    "你好，世界", // "Hello, World" in Chinese
			maxSize:  8,
			expected: "你...",
		},
		{
			name:     "Empty string",
			input:    "",
			maxSize:  5,
			expected: "",
		},
		{
			name:     "UTF-8 truncation with invalid byte",
			input:    "Hello, 世界\xef\xbf\xbd", // Invalid UTF-8 at the end
			maxSize:  10,
			expected: "Hello, ...",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			result := truncateStringWithEllipsis(test.input, test.maxSize)
			g.Expect(result).To(Equal(test.expected), "Expected %q but got %q for input %q with maxSize %d", test.expected, result, test.input, test.maxSize)
		})
	}
}


