package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyPostServeActionDetails interface {
	GetAllPostServeActions() PostServeActionDetailsView
	SetLocalPostServeAction(string, string, string, int) error
	SetRemotePostServeAction(string, string, int) error
	DeletePostServeAction(string) error
}

type HoverflyPostServeActionDetailsHandler struct {
	Hoverfly HoverflyPostServeActionDetails
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {

	mux.Get("/api/v2/hoverfly/post-serve-action", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Get),
	))
	mux.Put("/api/v2/hoverfly/post-serve-action", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Put),
	))
	mux.Delete("/api/v2/hoverfly/post-serve-action/:actionName", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(postServeActionDetailsHandler.Delete),
	))
}
func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	postServeActionsDetailsView := postServeActionDetailsHandler.Hoverfly.GetAllPostServeActions()
	bytes, _ := json.Marshal(postServeActionsDetailsView)

	handlers.WriteResponse(w, bytes)
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var actionRequest ActionView
	err := handlers.ReadFromRequest(req, &actionRequest)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	if actionRequest.Remote != "" {
		err = postServeActionDetailsHandler.Hoverfly.SetRemotePostServeAction(actionRequest.ActionName, actionRequest.Remote, actionRequest.DelayInMs)
	} else {
		err = postServeActionDetailsHandler.Hoverfly.SetLocalPostServeAction(actionRequest.ActionName, actionRequest.Binary, actionRequest.ScriptContent, actionRequest.DelayInMs)
	}
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}
	postServeActionDetailsHandler.Get(w, req, next)
}

func (postServeActionDetailsHandler *HoverflyPostServeActionDetailsHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	actionName := bone.GetValue(req, "actionName")
	err := postServeActionDetailsHandler.Hoverfly.DeletePostServeAction(actionName)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	postServeActionDetailsHandler.Get(w, req, next)
}
