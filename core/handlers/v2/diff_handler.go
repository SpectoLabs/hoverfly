package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type DiffHandler struct {
	Hoverfly Hoverfly
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
	mux.Options("/api/v2/diff", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *DiffHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var diffsToReturn []ResponseDiffForRequestView
	for request, value := range this.Hoverfly.GetDiff() {
		diffsToReturn = append(diffsToReturn, ResponseDiffForRequestView{
			Request:    request,
			DiffReport: value,
		})
	}

	marshal, err := json.Marshal(DiffView{
		Diff: diffsToReturn,
	})
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.WriteResponse(w, marshal)
	w.WriteHeader(http.StatusOK)
}

func (this *DiffHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.ClearDiff()
	w.WriteHeader(http.StatusOK)
}

func (this *DiffHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, DELETE")
	handlers.WriteResponse(w, []byte(""))
}
