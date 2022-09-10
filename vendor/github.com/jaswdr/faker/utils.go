package faker

import (
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

const (
	maxRetries = 7
)

// HTTPClient does HTTP requests to remote servers
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// HTTPClientImpl is the default implementation of HTTPClient
type HTTPClientImpl struct{}

// Get do a GET request and returns a *http.Response
func (HTTPClientImpl) Get(url string) (resp *http.Response, err error) {
	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(url)
		if err == nil {
			return resp, err
		}
	}

	return resp, err
}

// TempFileCreator creates temporary files
type TempFileCreator interface {
	TempFile(prefix string) (f *os.File, err error)
}

// TempFileCreatorImpl is the default implementation of TempFileCreator
type TempFileCreatorImpl struct{}

// TempFile creates a temporary file
func (TempFileCreatorImpl) TempFile(prefix string) (f *os.File, err error) {
	return ioutil.TempFile(os.TempDir(), prefix)
}

// OSResolver returns the GOOS value for operating an operating system
type OSResolver interface {
	OS() string
}

// OSResolverImpl is the default implementation of OSResolver
type OSResolverImpl struct{}

// OS returns the runtime.GOOS value for the host operating system
func (OSResolverImpl) OS() string {
	return runtime.GOOS
}
