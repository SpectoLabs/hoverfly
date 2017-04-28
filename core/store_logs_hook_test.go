package hoverfly

import (
	"strconv"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
)

func Test_NewStoreLogsHook_CreatesNewStructWithInitializedEntryArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	Expect(unit.Entries).ToNot(BeNil())
	Expect(unit.Entries).To(HaveLen(0))
}

func Test_StoreLogsHook_Fire_SavesEntryToEntriesArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Fire(&logrus.Entry{
		Message: "test entry",
	})

	Expect(unit.Entries).To(HaveLen(1))
	Expect(unit.Entries[0].Message).To(Equal("test entry"))
}

func Test_StoreLogsHook_Fire_LogsArePrepended(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Fire(&logrus.Entry{
		Message: "oldest",
	})

	unit.Fire(&logrus.Entry{
		Message: "newest",
	})

	Expect(unit.Entries).To(HaveLen(2))
	Expect(unit.Entries[0].Message).To(Equal("newest"))
	Expect(unit.Entries[1].Message).To(Equal("oldest"))
}

func Test_StoreLogsHook_GetLogsView_BuildsLogsViewFromEntriesArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Entries = append(unit.Entries, &logrus.Entry{
		Message: "test entry",
	})

	logs := unit.GetLogsView()

	Expect(logs.Logs).To(HaveLen(1))
	Expect(logs.Logs[0]["msg"]).To(Equal("test entry"))
}

func Test_StoreLogsHook_GetLogsView_LogsContainFields(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Entries = append(unit.Entries, &logrus.Entry{
		Message: "test entry",
		Data: logrus.Fields{
			"test": "field",
		},
	})

	logs := unit.GetLogsView()

	Expect(logs.Logs).To(HaveLen(1))
	Expect(logs.Logs[0]["test"]).To(Equal("field"))
}

func Test_StoreLogsHook_GetLogsView_LimitsTo500Logs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 501; i++ {
		unit.Fire(&logrus.Entry{
			Message: strconv.Itoa(i),
		})
	}

	logs := unit.GetLogsView()

	Expect(logs.Logs).To(HaveLen(500))
	Expect(logs.Logs[0]["msg"]).To(Equal("501"))
	Expect(logs.Logs[499]["msg"]).To(Equal("2"))
}

func Test_StoreLogsHook_GetFilteredLogsView_CanSetLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 3; i++ {
		unit.Fire(&logrus.Entry{
			Message: strconv.Itoa(i),
		})
	}

	logs := unit.GetFilteredLogsView(2)

	Expect(logs.Logs).To(HaveLen(2))
	Expect(logs.Logs[0]["msg"]).To(Equal("3"))
	Expect(logs.Logs[1]["msg"]).To(Equal("2"))
}

func Test_StoreLogsHook_GetFilteredLogsView_DoesNotErrorWhenLimitIsLargerThanArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Message: strconv.Itoa(i),
		})
	}

	logs := unit.GetFilteredLogsView(1)

	Expect(logs.Logs).To(HaveLen(1))
	Expect(logs.Logs[0]["msg"]).To(Equal("2"))
}

func Test_StoreLogsHook_GetLogs_BuildsStringFromEntriesArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Entries = append(unit.Entries, &logrus.Entry{
		Message: "test entry",
	})

	logs := unit.GetLogs()

	Expect(logs).To(ContainSubstring(`time="0001-01-01T00:00:00Z" level=panic msg="test entry"`))
}

func Test_StoreLogsHook_GetLogs_LogsContainFields(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Entries = append(unit.Entries, &logrus.Entry{
		Message: "test entry",
		Data: logrus.Fields{
			"test": "field",
		},
	})

	logs := unit.GetLogs()

	Expect(logs).To(ContainSubstring(`time="0001-01-01T00:00:00Z" level=panic msg="test entry" test=field `))
}

func Test_StoreLogsHook_GetLogs_LimitsTo500Logs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 501; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs := unit.GetLogs()

	Expect(logs).To(ContainSubstring("log-501"))
	Expect(logs).To(ContainSubstring("log-1"))
	Expect(logs).ToNot(ContainSubstring("logs-0"))
}

func Test_StoreLogsHook_GetFilteredLogs_CanSetLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 3; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs := unit.GetFilteredLogs(2)

	Expect(logs).To(ContainSubstring("log-3"))
	Expect(logs).To(ContainSubstring("log-2"))
	Expect(logs).ToNot(ContainSubstring("logs-1"))
}

func Test_StoreLogsHook_GetFilteredLogs_DoesNotErrorWhenLimitIsLargerThanArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs := unit.GetFilteredLogs(1)

	Expect(logs).To(ContainSubstring("log-2"))
}
