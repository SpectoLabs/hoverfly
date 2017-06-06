package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
)

var EnableCors bool

type ErrorView struct {
	Error string `json:"error"`
}

type AdminHandler interface {
	RegisterRoutes(*bone.Mux, *AuthHandler)
}

func ReadFromRequest(request *http.Request, v interface{}) error {
	defer request.Body.Close()

	body, _ := ioutil.ReadAll(request.Body)

	err := json.Unmarshal(body, &v)
	if err != nil {
		return errors.New("Malformed JSON")
	}

	return nil
}

func WriteResponse(response http.ResponseWriter, bytes []byte) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if EnableCors {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		response.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS, DELETE")
		response.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	response.Write(bytes)
}

func WriteErrorResponse(response http.ResponseWriter, message string, code int) {
	errorView := &ErrorView{Error: message}
	errorBytes, err := json.Marshal(errorView)
	if err != nil {
		response.WriteHeader(500)
		return
	}
	response.WriteHeader(code)
	WriteResponse(response, errorBytes)
}

type WebSocketHandler func() ([]byte, error)

func NewWebsocket(handler WebSocketHandler, w http.ResponseWriter, r *http.Request) {

	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to upgrade websocket")
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		log.WithFields(log.Fields{
			"message": string(p),
		}).Debug("Got message...")

		for _ = range time.Tick(1 * time.Second) {

			updateBytes, err := handler()

			if err = conn.WriteMessage(messageType, updateBytes); err != nil {
				log.WithFields(log.Fields{
					"message": p,
					"error":   err.Error(),
				}).Debug("Got error when writing message...")
				continue
			}
		}
	}
}
