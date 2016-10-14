package v2

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
