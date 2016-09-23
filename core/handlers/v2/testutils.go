package v2

import (
	"github.com/codegangsta/negroni"
	"net/http/httptest"
	"bytes"
	"io/ioutil"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"encoding/json"
	"net/http"
)

func makeRequestOnHandler(handlerFunc negroni.HandlerFunc, request *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	negroni := negroni.New(handlerFunc)
	negroni.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

func unmarshalErrorView(buffer *bytes.Buffer) (handlers.ErrorView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return handlers.ErrorView{}, err
	}

	var errorView handlers.ErrorView

	err = json.Unmarshal(body, &errorView)
	if err != nil {
		return handlers.ErrorView{}, err
	}

	return errorView, nil
}
