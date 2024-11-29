package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	encoding_xml "encoding/xml"
	"fmt"
	"github.com/ChrisTrenkamp/xsel/exec"
	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/parser"
	"github.com/ChrisTrenkamp/xsel/store"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/client-go/util/jsonpath"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"strconv"
	"time"

	xj "github.com/SpectoLabs/goxml2json"
	"github.com/tdewolff/minify/v2"
	mjson "github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	// mime types which will not be base 64 encoded when exporting as JSON
	SupportedMimeTypes = [...]string{"text", "plain", "css", "html", "json", "xml", "js", "javascript"}
)

// GetRequestBody will read the http.Request body io.ReadCloser
// and will also set the buffer to the original value as the
// buffer will be empty after reading it.
// It also decompress if any Content-Encoding is applied
func GetRequestBody(request *http.Request) (string, error) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	// Will add more compression support in the future
	if request.Header.Get("Content-Encoding") == "gzip" {
		decompressedBody, err := DecompressGzip(bodyBytes)
		if err == nil {
			bodyBytes = decompressedBody
		}
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}

// GetResponseBody will read the http.Response body io.ReadCloser
// and will also set the buffer to the original value as the
// buffer will be empty after reading it.
func GetResponseBody(response *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}

func GetResponseHeaders(response *http.Response) map[string][]string {

	// Make a copy of the response headers, preventing any changes to response being saved into the simulation
	headers := make(map[string][]string)
	for key, value := range response.Header {
		headers[key] = value
	}

	if response.Trailer == nil {
		return headers
	}

	var trailerKeys []string
	for key, value := range response.Trailer {
		headers[key] = value
		trailerKeys = append(trailerKeys, key)
	}

	headers["Trailer"] = trailerKeys
	return headers
}

func GetUnixTimeQueryParam(request *http.Request, paramName string) *time.Time {
	var timeQuery *time.Time
	epochValue, _ := strconv.Atoi(request.URL.Query().Get(paramName))
	if epochValue != 0 {
		timeValue := time.Unix(int64(epochValue), 0)
		timeQuery = &timeValue
	}
	return timeQuery
}

func StringToPointer(value string) *string {
	return &value
}

func PointerToString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

// SortQueryString will sort a http query string alphanumerically
// by key and then by value.
func SortQueryString(query string) string {
	keyValues := make(url.Values)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i+1], key[i+1:]
		}

		keyValues[key] = append(keyValues[key], value)
	}

	var queryBuffer bytes.Buffer
	keys := make([]string, 0, len(keyValues))
	for key := range keyValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		values := keyValues[key]
		sort.Strings(values)
		for _, value := range values {
			if queryBuffer.Len() > 0 {
				queryBuffer.WriteByte('&')
			}
			queryBuffer.WriteString(key)
			queryBuffer.WriteString(value)
		}
	}
	return queryBuffer.String()
}

func GetContentTypeFromHeaders(headers map[string][]string) string {
	if headers == nil {
		return ""
	}

	for _, v := range headers["Content-Type"] {
		if regexp.MustCompile("[/+]json$").MatchString(v) {
			return "json"
		}
		if regexp.MustCompile("[/+]xml$").MatchString(v) {
			return "xml"
		}
		if regexp.MustCompile(`form\-\w+$`).MatchString(v) {
			return "form"
		}
	}
	return ""
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

var minifier *minify.M

func GetMinifier() *minify.M {
	if minifier == nil {
		minifier = minify.New()
		minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), mjson.Minify)
		minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	}

	return minifier
}

func MinifyJson(toMinify string) (string, error) {
	minifier := GetMinifier()

	return minifier.String("application/json", toMinify)
}

func MinifyXml(toMinify string) (string, error) {
	minifier := GetMinifier()

	return minifier.String("application/xml", toMinify)
}

func CopyMap(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

// URL is regexp to match http urls
const urlPattern = `^((ftp|https?):\/\/)(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`

var rxURL = regexp.MustCompile(urlPattern)

func IsURL(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}

	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	if strings.HasPrefix(u.Host, ".") {
		return false
	}

	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return rxURL.MatchString(str)
}

func DecompressGzip(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		return body, err
	}
	defer reader.Close()
	body, err = ioutil.ReadAll(reader)
	if err != nil {
		return body, err
	}
	return body, err
}

func CompressGzip(body []byte) ([]byte, error) {
	var byteBuffer bytes.Buffer
	var err error
	gzWriter := gzip.NewWriter(&byteBuffer)
	if _, err := gzWriter.Write(body); err != nil {
		return body, err
	}
	if err := gzWriter.Flush(); err != nil {
		return body, err
	}
	if err := gzWriter.Close(); err != nil {
		return body, err
	}

	return byteBuffer.Bytes(), err
}

func Identical(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for i, v := range first {
		if v != second[i] {
			return false
		}
	}
	return true
}

func Contains(first, second []string) bool {

	set := make(map[string]bool)
	for _, value := range first {
		set[value] = true
	}

	for _, value := range second {

		if _, found := set[value]; found {
			return true
		}
	}

	return false
}

func ContainsOnly(first, second []string) bool {

	set := make(map[string]bool)

	for _, value := range first {
		set[value] = true
	}

	for _, value := range second {

		if _, found := set[value]; !found {
			return false
		}
	}
	return true
}

func GetStringArray(data interface{}) ([]string, bool) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil, false
	}
	var dataArr []string
	for i := 0; i < val.Len(); i++ {
		currentValue := val.Index(i)
		if currentValue.Kind() == reflect.Interface {
			dataArr = append(dataArr, currentValue.Elem().String())
		} else {
			dataArr = append(dataArr, currentValue.String())
		}

	}
	return dataArr, true
}

func GetBoolOrDefault(data map[string]interface{}, key string, defaultValue bool) bool {
	if data == nil {
		return defaultValue
	}

	genericValue, found := data[key]
	if !found {
		return defaultValue
	}
	return genericValue.(bool)
}

func RandStringFromTimestamp(length int) string {
	timestamp := time.Now()
	randomBytes := make([]byte, length)

	// Combine timestamp seconds and microseconds
	seed := timestamp.UnixNano()

	// Use the timestamp as a seed for the random number generator
	rand.Seed(seed)

	// Generate random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		return strconv.FormatInt(seed, 10)
	}

	// Encode the random bytes as a string
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)

	return randomString
}

func FetchFromRequestBody(queryType, query, toMatch string) interface{} {

	if queryType == "jsonpath" {
		return jsonPath(query, toMatch)
	} else if queryType == "xpath" {
		return xPath(query, toMatch)
	} else if queryType == "jsonpathfromxml" {
		xmlReader := strings.NewReader(toMatch)
		jsonBytes, err := xj.Convert(xmlReader)
		if err != nil {
			return ""
		}
		return jsonPath(query, jsonBytes.String())
	}
	log.Errorf("Unknown query type \"%s\" for templating Request.Body", queryType)
	return ""
}

func jsonPath(query, toMatch string) interface{} {
	query = PrepareJsonPathQuery(query)

	result, err := JsonPathExecution(query, toMatch)
	if err != nil {
		return ""
	}

	//// Jsonpath library converts large int into a string with scientific notion, the following
	//// reverts that process to avoid mismatching when using the jsonpath result for csv data lookup
	//floatResult, err := strconv.ParseFloat(result, 64)
	//// if the string is a float and a whole number
	//if err == nil && floatResult == float64(int64(floatResult)) {
	//	intResult := int(floatResult)
	//	result = strconv.Itoa(intResult)
	//}

	// convert to array data if applicable
	var data interface{}
	err = json.Unmarshal([]byte(result), &data)

	arrayData, ok := data.([]interface{})

	if err != nil || !ok {
		return result
	}
	return arrayData
}

func xPath(query, toMatch string) string {
	result, err := XpathExecution(query, toMatch)
	if err != nil {
		return ""
	}
	return result.String()
}

func PrepareJsonPathQuery(query string) string {
	if query[0:1] != "{" && query[len(query)-1:] != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}

func JsonPathExecution(matchString, toMatch string) (string, error) {
	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchString, err.Error())
		return "", err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return "", err
	}

	buf := new(bytes.Buffer)

	err = jsonPath.Execute(buf, data)
	if err != nil {
		log.Errorf("err to execute json path match: %s", err.Error())
		return "", err
	}

	return buf.String(), nil
}

func XpathExecution(matchString, toMatch string) (exec.Result, error) {

	contextSettings := func(c *exec.ContextSettings) {
		xmlns := xmlns{}
		_ = encoding_xml.Unmarshal([]byte(toMatch), &xmlns)
		for key, value := range xmlns.Namespaces {
			c.NamespaceDecls[key] = value
		}
	}
	xpath := grammar.MustBuild(matchString)
	parsedXml := parser.ReadXml(bytes.NewBufferString(toMatch))
	cursor, _ := store.CreateInMemory(parsedXml)

	results, err := exec.Exec(cursor, &xpath, contextSettings)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return nil, err
	}

	return results, nil
}

type xmlns struct {
	Namespaces map[string]string
}

func (a *xmlns) UnmarshalXML(_ *encoding_xml.Decoder, start encoding_xml.StartElement) error {
	a.Namespaces = map[string]string{}
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			a.Namespaces[attr.Name.Local] = attr.Value
		}
	}
	return nil
}

func NeedsEncoding(headers map[string][]string, body string) bool {
	needsEncoding := false

	// Check headers for gzip
	contentEncodingValues := headers["Content-Encoding"]
	if len(contentEncodingValues) > 0 {
		needsEncoding = true
	} else {
		mimeType := http.DetectContentType([]byte(body))
		needsEncoding = true
		for _, v := range SupportedMimeTypes {
			if strings.Contains(mimeType, v) {
				needsEncoding = false
				break
			}
		}
	}
	return needsEncoding
}

// Resolves a relative path from an absolute basePath, and fails if the relative path starts with ".."
func ResolveAndValidatePath(absBasePath, relativePath string) (string, error) {

	cleanRelativePath := filepath.Clean(relativePath)

	// Check if the relative path starts with ".."
	if strings.HasPrefix(cleanRelativePath, "..") {
		return "", fmt.Errorf("relative path is invalid as it attempts to backtrack")
	}

	resolvedPath := filepath.Join(absBasePath, cleanRelativePath)

	// Verify that the resolved path is indeed a subpath of the base path
	finalPath, err := filepath.Rel(absBasePath, resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %v", err)
	}

	if strings.HasPrefix(finalPath, "..") {
		return "", fmt.Errorf("resolved path is outside the base path")
	}

	return resolvedPath, nil
}
