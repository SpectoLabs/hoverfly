package v2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyDiff interface {
	GetDiff() map[SimpleRequestDefinitionView][]DiffReport
	GetFilteredDiff(diffFilterView DiffFilterView) map[SimpleRequestDefinitionView][]DiffReport
	ClearDiff()
}

type DiffHandler struct {
	Hoverfly HoverflyDiff
}

func (this *DiffHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/diff", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Delete("/api/v2/diff", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Post("/api/v2/diff", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.GetFilteredData),
	))
	mux.Options("/api/v2/diff", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *DiffHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	diffsToReturn := convertToResponseDiffView(this.Hoverfly.GetDiff())
	marshal, err := json.Marshal(DiffView{
		Diff: diffsToReturn,
	})
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.WriteResponse(w, marshal)
}

func (this *DiffHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.ClearDiff()

	handlers.WriteResponse(w, []byte(""))
}

func (this *DiffHandler) GetFilteredData(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var diffFilterView DiffFilterView
	err = json.Unmarshal(requestBody, &diffFilterView)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	diffsToReturn := convertToResponseDiffView(this.Hoverfly.GetFilteredDiff(diffFilterView))
	marshal, err := json.Marshal(DiffView{
		Diff: diffsToReturn,
	})
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	handlers.WriteResponse(w, marshal)
}

func convertToResponseDiffView(responsesDiff map[SimpleRequestDefinitionView][]DiffReport) []ResponseDiffForRequestView {

	var diffsToReturn []ResponseDiffForRequestView
	for request, value := range responsesDiff {
		diffsToReturn = append(diffsToReturn, ResponseDiffForRequestView{
			Request:    request,
			DiffReport: value,
		})
	}
	return diffsToReturn
}

func (this *DiffHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, DELETE")
	handlers.WriteResponse(w, []byte(""))
}
