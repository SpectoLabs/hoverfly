package hoverfly

import (
	"testing"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_rebuildHashes_whenDataIsHashedForAProxy_andStillAProxy_keysAreNotChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := false

	db := cache.NewInMemoryCache()

	testPayload := models.Payload{
		Request: models.RequestDetails{
			Path: "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	testPayloadBytes, _ := testPayload.Encode()

	db.Set([]byte(testPayload.Id()), testPayloadBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(testPayload.Id()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(testPayloadBytes))
}

func Test_rebuildHashes_whenDataIsHashedForAWebserver_andStillAWebserver_keysAreNotChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := true

	db := cache.NewInMemoryCache()

	testPayload := models.Payload{
		Request: models.RequestDetails{
			Path: "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	testPayloadBytes, _ := testPayload.Encode()

	db.Set([]byte(testPayload.IdWithoutHost()), testPayloadBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(testPayload.IdWithoutHost()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(testPayloadBytes))
}

func Test_rebuildHashes_whenDataIsHashedForAProxy_andIsNowAWebserver_keysAreChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := true

	db := cache.NewInMemoryCache()

	testPayload := models.Payload{
		Request: models.RequestDetails{
			Path: "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	testPayloadBytes, _ := testPayload.Encode()

	db.Set([]byte(testPayload.Id()), testPayloadBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(testPayload.IdWithoutHost()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(testPayloadBytes))
}