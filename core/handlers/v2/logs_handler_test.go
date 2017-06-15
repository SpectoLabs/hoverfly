package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
	"time"
)

type HoverflyLogsStub struct {
	limit int
	from *time.Time
}

func (this *HoverflyLogsStub) GetLogs(limit int, from *time.Time) []*logrus.Entry {
	this.limit = limit
	this.from  = from
	return []*logrus.Entry{&logrus.Entry{
		Level:   logrus.InfoLevel,
		Message: "a line of logs",
		Data: map[string]interface{}{
			"custom-field": "field value",
		},
	}}
}

func Test_LogsHandler_Get_ReturnsLogsView(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	logsView, err := unmarshalLogsView(response.Body)
	Expect(err).To(BeNil())

	Expect(logsView.Logs[0]["msg"]).To(Equal("a line of logs"))
}

func Test_LogsHandler_Get_SetsTheDefaultLimitIfNoneIsSpecified(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.limit).To(Equal(500))
	Expect(stubHoverfly.from).To(BeNil())
}

func Test_LogsHandler_Get_SetsTheLimitIfLimitQueryProvided(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?limit=20", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.limit).To(Equal(20))
}

func Test_LogsHandler_Get_SetsTheFromTimeIfFromQueryProvided(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?from=1497521986", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.from.Unix()).To(Equal(int64(1497521986)))
}

func Test_LogsHandler_Get_DoesNotSetTimeIfFromQueryIsBadTime(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?from=bad-time", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.from).To(BeNil())
}


func Test_LogsHandler_Get_ReturnsLogsInPlaintext_UsingAcceptHeader(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	request.Header.Set("Accept", "text/plain")

	response := makeRequestOnHandler(unit.Get, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	logs, _ := ioutil.ReadAll(response.Body)

	Expect(string(logs)).To(ContainSubstring("INFO"))
	Expect(string(logs)).To(ContainSubstring("a line of logs"))

	Expect(string(logs)).To(ContainSubstring("custom-field"))
	Expect(string(logs)).To(ContainSubstring("field value"))
}

func Test_LogsHandler_Get_ReturnsLogsInPlaintext_UsingContentTypeHeader(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	request.Header.Set("Content-Type", "text/plain")

	response := makeRequestOnHandler(unit.Get, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	logs, _ := ioutil.ReadAll(response.Body)

	Expect(string(logs)).To(ContainSubstring("INFO"))
	Expect(string(logs)).To(ContainSubstring("a line of logs"))

	Expect(string(logs)).To(ContainSubstring("custom-field"))
	Expect(string(logs)).To(ContainSubstring("field value"))
}

func Test_LogsHandler_Get_SetsTheDefaultLimitIfNoneIsSpecified_InPlaintext(t *testing.T) {

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	request.Header.Set("Content-Type", "text/plain")

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.limit).To(Equal(500))
}

func Test_LogsHandler_Get_SetsTheLimitIfLimitQueryProvided_InPlaintext(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyLogsStub{}
	unit := LogsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?limit=20", nil)
	Expect(err).To(BeNil())

	request.Header.Set("Content-Type", "text/plain")

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.limit).To(Equal(20))
}

func Test_LogsHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyLogsStub
	unit := LogsHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/logs", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET"))
}

func unmarshalLogsView(buffer *bytes.Buffer) (LogsView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return LogsView{}, err
	}

	var logsView LogsView

	err = json.Unmarshal(body, &logsView)
	if err != nil {
		return LogsView{}, err
	}

	return logsView, nil
}
