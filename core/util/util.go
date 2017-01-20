package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// GetRequestBody will read the http.Request body io.ReadCloser
// and will also set the buffer to the original value as the
// buffer will be empty after reading it.
func GetRequestBody(request *http.Request) (string, error) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
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
