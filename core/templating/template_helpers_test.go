package templating

import (
	"encoding/base64"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// mockRaymondOptions is a minimal mock for raymond.Options for testing
type mockRaymondOptions struct {
	internalVars map[string]interface{}
}

func (m *mockRaymondOptions) ValueFromAllCtx(key string) interface{} {
	if key == "InternalVars" {
		return m.internalVars
	}
	return nil
}

func testNow() time.Time {
	parsedTime, _ := time.Parse("2006-01-02T15:04:05Z", "2018-01-01T00:00:00Z")
	return parsedTime
}

func Test_now_withEmptyOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_now_withEmptyOffsetAndUnixFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "unix")).To(Equal("1514764800"))
}

func Test_now_withEmptyOffsetAndUnixMillisFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "epoch")).To(Equal("1514764800000"))
}

func Test_now_withEmptyOffsetAndCustomFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:00 UTC 2018"))
}

func Test_now_withPositiveOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("1d", "")).To(Equal("2018-01-02T00:00:00Z"))
}

func Test_now_withNegativeOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("-1d", "")).To(Equal("2017-12-31T00:00:00Z"))
}

func Test_now_withInvalidOffset(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("cat", "")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_now_withInvalidFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "dog")).To(Equal("dog"))
}

func Test_replace(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.replace("oink, oink, oink", "oink", "moo")).To(Equal("moo, moo, moo"))
}

func Test_split(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.split("one,two,three", ",")).To(ConsistOf("one", "two", "three"))
}

func Test_concat(t *testing.T) {

	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.concat("one", " two")).To(Equal("one two"))
}

func Test_concatWithManyStrings(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.concat("one", " two", " three", " four")).To(Equal("one two three four"))
}

func Test_length(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.length("onelongstring")).To(Equal("13"))
}

func Test_substring(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.substring("onelongstring", "3", "7")).To(Equal("long"))
}

func Test_substring_withInvalidStart(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.substring("onelongstring", "-3", "6")).To(Equal(""))
}

func Test_substring_withInvalidEnd(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.substring("onelongstring", "3", "the end")).To(Equal(""))
}

func Test_rightmostCharacters(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.rightmostCharacters("onelongstring", "3")).To(Equal("ing"))
}

func Test_rightmostCharacters_withInvalidCount(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.rightmostCharacters("onelongstring", "30")).To(Equal(""))
}

func Test_isNumeric_withInteger(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isNumeric("123")).To(Equal(true))
}

func Test_isNumeric_withFloat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isNumeric("45.67")).To(Equal(true))
}

func Test_isNumeric_withScientific(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isNumeric("1e10")).To(Equal(true))
}

func Test_isNumeric_withNegative(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isNumeric("-5")).To(Equal(true))
}

func Test_isNumeric_withString(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isNumeric("hello")).To(Equal(false))
}

func Test_isAlphanumeric_withAlphanumeric(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isAlphanumeric("ABC123")).To(Equal(true))
}

func Test_isAlphanumeric_withNumeric(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isAlphanumeric("123")).To(Equal(true))
}

func Test_isAlphanumeric_withAlpha(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isAlphanumeric("ABC")).To(Equal(true))
}

func Test_isAlphanumeric_withInvalidAlphanumeric(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isAlphanumeric("ABC!@123")).To(Equal(false))
}

func Test_isBool_withtrue(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBool("true")).To(Equal(true))
}

func Test_isBool_withfalse(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBool("false")).To(Equal(true))
}

func Test_isBool_with1(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBool("1")).To(Equal(true))
}

func Test_isBool_with0(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBool("0")).To(Equal(true))
}

func Test_isBool_withInvalidValue(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBool("maybe")).To(Equal(false))
}

func Test_isGreaterThan_withPositiveResult(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isGreaterThan("11", "10")).To(Equal(true))
}

func Test_isGreaterThan_withNegativeResult(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isGreaterThan("10", "11")).To(Equal(false))
}

func Test_isGreaterThan_withInvalidNumber(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isGreaterThan("abc", "11")).To(Equal(false))
}

func Test_isLessThan_withPositiveResult(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isLessThan("10", "11")).To(Equal(true))
}

func Test_isLessThan_withNegativeResult(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isLessThan("11", "10")).To(Equal(false))
}

func Test_isLessThan_withInvalidNumber(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isLessThan("abc", "11")).To(Equal(false))
}

func Test_isBetween_withPositiveOutcome(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBetween("5", "3", "7")).To(Equal(true))
}

func Test_isBetween_withNegativeOutcome(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBetween("5", "6", "7")).To(Equal(false))
}

func Test_isBetween_withInvalidArgument(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.isBetween("e", "6", "7")).To(Equal(false))
}

func Test_matchesRegex_withPositiveOutcome(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.matchesRegex("{\"someField\": \"a\", \"transactionId\": 1000, \"anotherField\": \"b\", \"store\": \"c\", \"clientUniqueId\": \"12345\", \"items\": [\"item1\", \"item2\", \"item3\"], \"extraField\": \"d\"}", "(?s).*(\"transactionId\": 1000).*store.*clientUniqueId.*items.*")).To(Equal(true))
}
func Test_matchesRegex_withNegativeOutcome(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.matchesRegex("{\"someField\": \"a\", \"transactionNumber\": 1000, \"anotherField\": \"b\", \"store\": \"c\", \"clientUniqueId\": \"12345\", \"items\": [\"item1\", \"item2\", \"item3\"], \"extraField\": \"d\"}", "(?s).*(\"transactionId\": 1000).*store.*clientUniqueId.*items.*")).To(Equal(false))
}

func Test_matchesRegex_withInvalidArgument(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.matchesRegex("I am looking for this string", "&^%$£!@:<>+_!¬")).To(Equal(false))
}

func Test_faker(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.faker("JobTitle")[0].String()).To(Not(BeEmpty()))
}

func Test_jsonFromJWT_basicClaim(t *testing.T) {
	RegisterTestingT(t)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user1","roles":["a","b"],"num":1234567890}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte("sig"))
	token := header + "." + payload + "." + sig

	unit := templateHelpers{}
	Expect(unit.jsonFromJWT("$.payload.sub", token)).To(Equal("user1"))
}

func Test_jsonFromJWT_arrayClaim(t *testing.T) {
	RegisterTestingT(t)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"roles":["a","b","c"]}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte("sig"))
	token := header + "." + payload + "." + sig

	unit := templateHelpers{}
	result := unit.jsonFromJWT("$.payload.roles", token)
	arr, ok := result.([]interface{})
	Expect(ok).To(BeTrue())
	Expect(arr).To(ConsistOf("a", "b", "c"))
}

func Test_jsonFromJWT_bearerPrefix(t *testing.T) {
	RegisterTestingT(t)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"userX"}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte("sig"))
	token := "Bearer " + header + "." + payload + "." + sig

	unit := templateHelpers{}
	Expect(unit.jsonFromJWT("$.payload.sub", token)).To(Equal("userX"))
}

func Test_jsonFromJWT_invalidToken(t *testing.T) {
	RegisterTestingT(t)
	unit := templateHelpers{}
	Expect(unit.jsonFromJWT("$.payload.sub", "not-a-jwt")).To(Equal(""))
}

func Test_jsonFromJWT_missingClaim(t *testing.T) {
	RegisterTestingT(t)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user1"}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte("sig"))
	token := header + "." + payload + "." + sig

	unit := templateHelpers{}
	Expect(unit.jsonFromJWT("$.payload.missing", token)).To(Equal(""))
}
