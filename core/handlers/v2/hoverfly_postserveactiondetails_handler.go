package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyPostServeActionDetails interface {
	GetPostServeActionDetails() PostServeActionDetailsView
	RegisterPostServeActionHook(string, string, string, int) error
	DeletePostServeActionHook(string) error
}

type HoverflyPostServeActionDetailsHandler struct {
	Hoverfly HoverflyPostServeActionDetails
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {

	mux.Get("/api/v2/hoverfly/post-serve-actions", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Get),
	))
	mux.Post("/api/v2/hoverfly/post-serve-actions/hook", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Post),
	))
	mux.Delete("/api/v2/hoverfly/post-serve-actions/hook", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Delete),
	))
}
func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	postSimulationDetailsView := postServeActionDetailsHandler.Hoverfly.GetPostServeActionDetails()
	bytes, _ := json.Marshal(postSimulationDetailsView)

	handlers.WriteResponse(w, bytes)
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Post(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var hookReq HookView
	err := handlers.ReadFromRequest(req, &hookReq)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	err = postServeActionDetailsHandler.Hoverfly.RegisterPostServeActionHook(hookReq.HookName, hookReq.Binary, hookReq.ScriptContent, hookReq.DelayInMilliSeconds)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
	}

	postServeActionDetailsHandler.Get(w, req, next)
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	queryParams := req.URL.Query()
	hookName := queryParams.Get("name")

	err := postServeActionDetailsHandler.Hoverfly.DeletePostServeActionHook(hookName)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
	}

	postServeActionDetailsHandler.Get(w, req, next)
}
