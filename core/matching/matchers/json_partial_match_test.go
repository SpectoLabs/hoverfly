package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonPartialMatch_MatchesTrueWithEqualsJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithNotOrderedJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"minified":true,"json":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithAbsentNode(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"minified":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithAbsentObject(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"minified":true}}`, `{"test":{"json":true,"minified":true,"someObject":{"fieldA":"valueA"}}}`)).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesFalseWithAbsentNode(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"minified":true}}`)).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesFalseWithAbsentObject(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true,"someObject":{"fieldA":"valueA"}}}`, `{"test":{"minified":true}}`)).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesTrueEmptyJson(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{}`, `{}`)).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesFalseInvalidJson(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":}}`)).To(BeFalse())
}
