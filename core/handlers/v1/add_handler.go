package v1

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
	"strconv"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
)

type AddHandler struct {
	Hoverfly HoverflyRecords
}

func (this *AddHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Post("/api/add", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Post),
	))
}

// ManualAddHandler - manually add new request/responses, using a form
func (this *AddHandler) Post(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := req.ParseForm()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Got error while parsing form")
	}

	// request details
	destination := req.PostFormValue("inputDestination")
	method := req.PostFormValue("inputMethod")
	path := req.PostFormValue("inputPath")
	query := req.PostFormValue("inputQuery")
	reqBody := req.PostFormValue("inputRequestBody")

	preq := RequestDetailsView{
		Destination: util.StringToPointer(destination),
		Method:      util.StringToPointer(method),
		Path:        util.StringToPointer(path),
		Query:       util.StringToPointer(query),
		Body:        util.StringToPointer(reqBody),
	}

	// response
	respStatusCode := req.PostFormValue("inputResponseStatusCode")
	respBody := req.PostFormValue("inputResponseBody")
	contentType := req.PostFormValue("inputContentType")

	headers := make(map[string][]string)

	// getting content type
	if contentType == "xml" {
		headers["Content-Type"] = []string{"application/xml"}
	} else if contentType == "json" {
		headers["Content-Type"] = []string{"application/json"}
	} else {
		headers["Content-Type"] = []string{"text/html"}
	}

	sc, _ := strconv.Atoi(respStatusCode)

	presp := ResponseDetailsView{
		Status:  sc,
		Headers: headers,
		Body:    respBody,
	}

	log.WithFields(log.Fields{
		"respBody":    respBody,
		"contentType": contentType,
	}).Info("manually adding request/response")

	p := RequestResponsePairView{Request: preq, Response: presp}

	pairViews := make([]interfaces.RequestResponsePair, 0)

	pairViews = append(pairViews, p)

	err = this.Hoverfly.ImportRequestResponsePairViews(pairViews)

	w.Header().Set("Content-Type", "application/json")
	var response MessageResponse

	if err != nil {
		response.Message = fmt.Sprintf("Got error: %s", err.Error())
		w.WriteHeader(400)

	} else {
		// redirecting to home
		response.Message = "Record added successfuly"
		w.WriteHeader(201)
	}
	b, err := response.Encode()
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not encode response body!")
		http.Error(w, "Failed to encode response", 500)
		return
	}
	w.Write(b)

}
