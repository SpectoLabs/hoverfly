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

func SortQueryString(query string) string {
	m := make(url.Values)
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
			key, value = key[:i], key[i+1:]
		}

		m[key] = append(m[key], value)
	}

	var buf bytes.Buffer
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := m[k]
		prefix := k + "="
		sort.Strings(vs)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
