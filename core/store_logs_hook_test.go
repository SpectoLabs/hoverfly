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

func Test_StoreLogsHook_NewGetLogs_LimitsLogs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs := unit.GetLogs(2)
	Expect(logs).To(HaveLen(2))
	Expect(logs[0].Message).To(Equal("log-2"))
	Expect(logs[1].Message).To(Equal("log-1"))
}

func Test_StoreLogsHook_NewGetLogs_AcceptsALimitThatIsTooLarge(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Fire(&logrus.Entry{
		Message: "log-0",
	})

	logs := unit.GetLogs(1000)
	Expect(logs).To(HaveLen(1))
	Expect(logs[0].Message).To(Equal("log-0"))
}
