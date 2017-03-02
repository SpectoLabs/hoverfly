package hoverfly

import (
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
