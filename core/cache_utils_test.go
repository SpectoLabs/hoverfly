package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_rebuildHashes_whenDataIsHashedForAProxy_andStillAProxy_keysAreNotChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := false

	db := cache.NewInMemoryCache()

	pair := models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	pairBytes, _ := pair.Encode()

	db.Set([]byte(pair.Id()), pairBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(pair.Id()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(pairBytes))
}

func Test_rebuildHashes_whenDataIsHashedForAWebserver_andStillAWebserver_keysAreNotChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := true

	db := cache.NewInMemoryCache()

	pair := models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	pairBytes, _ := pair.Encode()

	db.Set([]byte(pair.IdWithoutHost()), pairBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(pair.IdWithoutHost()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(pairBytes))
}

func Test_rebuildHashes_whenDataIsHashedForAProxy_andIsNowAWebserver_keysAreChanged(t *testing.T) {
	RegisterTestingT(t)
	webserver := true

	db := cache.NewInMemoryCache()

	pair := models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "/hello",
			Destination: "a-host.com",
		},
		Response: models.ResponseDetails{
			Body: "a body",
		},
	}

	pairBytes, _ := pair.Encode()

	db.Set([]byte(pair.Id()), pairBytes)

	rebuildHashes(db, webserver)

	result, err := db.Get([]byte(pair.IdWithoutHost()))

	Expect(err).To(BeNil())
	Expect(result).To(Equal(pairBytes))
}
