package hoverfly

import (
	"strconv"
	"testing"

	"time"

	"github.com/sirupsen/logrus"
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

func Test_StoreLogsHook_Fire_LogsAreInAscendingOrder(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Fire(&logrus.Entry{
		Message: "oldest",
	})

	unit.Fire(&logrus.Entry{
		Message: "newest",
	})

	Expect(unit.Entries).To(HaveLen(2))
	Expect(unit.Entries[0].Message).To(Equal("oldest"))
	Expect(unit.Entries[1].Message).To(Equal("newest"))
}

func Test_StoreLogsHook_Fire_RespectLogsLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()
	unit.LogsLimit = 5

	for i := 1; i < 8; i++ {
		unit.Fire(&logrus.Entry{
			Message: strconv.Itoa(i),
		})
	}

	Expect(unit.Entries).To(HaveLen(5))
	Expect(unit.Entries[0].Message).To(Equal("3"))
	Expect(unit.Entries[1].Message).To(Equal("4"))
	Expect(unit.Entries[2].Message).To(Equal("5"))
	Expect(unit.Entries[3].Message).To(Equal("6"))
	Expect(unit.Entries[4].Message).To(Equal("7"))
}

func Test_StoreLogsHook_Fire_IfLogsLimitIs0StoreLogsHookIsDisabled(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()
	unit.LogsLimit = 0

	for i := 1; i < 8; i++ {
		unit.Fire(&logrus.Entry{
			Message: strconv.Itoa(i),
		})
	}

	Expect(unit.Entries).To(HaveLen(0))

	entries, err := unit.GetLogs(100, nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Logs disabled"))

	Expect(entries).To(HaveLen(0))
}

func Test_StoreLogsHook_GetLogs_LimitsLogs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs, err := unit.GetLogs(2, nil)
	Expect(err).To(BeNil())

	Expect(logs).To(HaveLen(2))
	Expect(logs[0].Message).To(Equal("log-1"))
	Expect(logs[1].Message).To(Equal("log-2"))
}

func Test_StoreLogsHook_GetLogs_AcceptsALimitThatIsTooLarge(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	unit.Fire(&logrus.Entry{
		Message: "log-0",
	})

	logs, err := unit.GetLogs(1000, nil)
	Expect(err).To(BeNil())

	Expect(logs).To(HaveLen(1))
	Expect(logs[0].Message).To(Equal("log-0"))
}

func Test_StoreLogsHook_GetLogs_FilteredByFromDateTime(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Time:    time.Date(2017, 6, 14, 10, 0, i, 0, time.Local),
			Message: "log-0",
		})
	}

	queryDate := time.Date(2017, 6, 14, 10, 0, 1, 0, time.Local)

	logs, err := unit.GetLogs(100, &queryDate)
	Expect(err).To(BeNil())

	Expect(logs).To(HaveLen(1))
	expectPrecision, _ := time.ParseDuration("1s")
	Expect(logs[0].Time).To(BeTemporally("==", time.Date(2017, 6, 14, 10, 0, 2, 0, time.Local), expectPrecision))
}

func Test_StoreLogsHook_GetLogs_FilteredByFromDateTimeAndLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 3; i++ {
		unit.Fire(&logrus.Entry{
			Time:    time.Date(2017, 6, 14, 10, 0, i, 0, time.Local),
			Message: "log-" + strconv.Itoa(i),
		})
	}

	queryDate := time.Date(2017, 6, 14, 10, 0, 0, 0, time.Local)

	logs, err := unit.GetLogs(2, &queryDate)
	Expect(err).To(BeNil())

	Expect(logs).To(HaveLen(2))
	expectPrecision, _ := time.ParseDuration("1s")

	Expect(logs[0].Message).To(Equal("log-2"))
	Expect(logs[0].Time).To(BeTemporally("==", time.Date(2017, 6, 14, 10, 0, 2, 0, time.Local), expectPrecision))

	Expect(logs[1].Message).To(Equal("log-3"))
	Expect(logs[1].Time).To(BeTemporally("==", time.Date(2017, 6, 14, 10, 0, 3, 0, time.Local), expectPrecision))

}
