package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/tdewolff/minify"
	mjson "github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/xml"
	"strconv"
	"time"
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
		if i := strings.IndexAny(key, "&;"); i >= 0 {
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
