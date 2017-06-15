package hoverfly

import (
	"strconv"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
	"time"
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

func Test_StoreLogsHook_GetLogs_LimitsLogs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Message: "log-" + strconv.Itoa(i),
		})
	}

	logs := unit.GetLogs(2, nil)
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

	logs := unit.GetLogs(1000, nil)
	Expect(logs).To(HaveLen(1))
	Expect(logs[0].Message).To(Equal("log-0"))
}

func Test_StoreLogsHook_GetLogs_FilteredByFromDateTime(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 2; i++ {
		unit.Fire(&logrus.Entry{
			Time: time.Date(2017, 6, 14, 10, 0, i, 0, time.Local),
			Message: "log-0",
		})
	}

	queryDate := time.Date(2017, 6, 14, 10, 0, 1, 0, time.Local)

	logs := unit.GetLogs(100, &queryDate)
	Expect(logs).To(HaveLen(1))
	expectPrecision, _ := time.ParseDuration("1s")
	Expect(logs[0].Time).To(BeTemporally("==",time.Date(2017, 6, 14, 10, 0, 2, 0, time.Local), expectPrecision))
}

func Test_StoreLogsHook_GetLogs_FilteredByFromDateTimeAndLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := NewStoreLogsHook()

	for i := 0; i <= 3; i++ {
		unit.Fire(&logrus.Entry{
			Time: time.Date(2017, 6, 14, 10, 0, i, 0, time.Local),
			Message: "log-" + strconv.Itoa(i),
		})
	}

	queryDate := time.Date(2017, 6, 14, 10, 0, 0, 0, time.Local)

	logs := unit.GetLogs(2, &queryDate)
	Expect(logs).To(HaveLen(2))
	expectPrecision, _ := time.ParseDuration("1s")

	Expect(logs[0].Message).To(Equal("log-2"))
	Expect(logs[0].Time).To(BeTemporally("==",time.Date(2017, 6, 14, 10, 0, 2, 0, time.Local), expectPrecision))

	Expect(logs[1].Message).To(Equal("log-3"))
	Expect(logs[1].Time).To(BeTemporally("==",time.Date(2017, 6, 14, 10, 0, 3, 0, time.Local), expectPrecision))

}