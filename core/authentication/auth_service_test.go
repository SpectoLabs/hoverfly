package authentication_test

import (
	"testing"
	"time"

	"github.com/SpectoLabs/hoverfly/core/authentication"
	. "github.com/onsi/gomega"
)

var timeLayout = "2006-01-02T15:04:05.000Z"

func Test_HasReachFailedAttemptsLimit_ReturnsFalseIfAttemptsIsBelowOrEqualToLimit(t *testing.T) {
	RegisterTestingT(t)

	authentication.Attempts.Count = 0
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeFalse())

	authentication.Attempts.Count = 1
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeFalse())

	authentication.Attempts.Count = 2
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeFalse())

	authentication.Attempts.Count = 3
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeFalse())
}

func Test_HasReachFailedAttemptsLimit_ReturnsTrueIfAttemptsIsAboveLimit_AndLastFailedIsWithinTheTimeoutPeriod(t *testing.T) {
	RegisterTestingT(t)

	authentication.Attempts.Count = 4
	authentication.Attempts.LastFailed = time.Now()
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeTrue())
}

func Test_HasReachFailedAttemptsLimit_IncreasesCountIfAttemptsIsAboveLimit_AndLastFailedIsWithinTheTimeoutPeriod(t *testing.T) {
	RegisterTestingT(t)

	authentication.Attempts.Count = 4
	authentication.Attempts.LastFailed = time.Now()

	authentication.HasReachedFailedAttemptsLimit(3, "10s")

	Expect(authentication.Attempts.Count).To(Equal(5))
}

func Test_HasReachFailedAttemptsLimit_ReturnsFalseIfAttemptsIsAboveLimit_ButLastFailedTimeIsAfterTimeoutPeriod(t *testing.T) {
	RegisterTestingT(t)

	authentication.Attempts.Count = 4

	authentication.Attempts.LastFailed, _ = time.Parse(timeLayout, "2017-04-01T10:45:26.371Z")
	Expect(authentication.HasReachedFailedAttemptsLimit(3, "10s")).To(BeFalse())
}

func Test_HasReachFailedAttemptsLimit_ResetsCountIfAttemptsIsAboveLimit_ButLastFailedTimeIsAfterTimeoutPeriod(t *testing.T) {
	RegisterTestingT(t)

	authentication.Attempts.Count = 4

	authentication.Attempts.LastFailed, _ = time.Parse(timeLayout, "2017-04-01T10:45:26.371Z")
	authentication.HasReachedFailedAttemptsLimit(3, "10s")

	Expect(authentication.Attempts.Count).To(Equal(0))
}
