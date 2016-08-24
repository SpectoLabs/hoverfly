package hoverfly

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func rebuildHashes(db cache.Cache, webserver bool) {
	log.Info("Checking if keys in cache need rehashing")

	entries, err := db.GetAllEntries()
	if err != nil {
		log.Fatal("Unable to read from BoltDB cache")
	}

	for key, bytes := range entries {
		pair, err := models.NewRequestResponsePairFromBytes(bytes)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"value": string(bytes),
				"key":   key,
			}).Error("Failed to decode payload")
		}
		var newKey string
		if webserver {
			newKey = pair.IdWithoutHost()
		} else {
			newKey = pair.Id()
		}

		if key != newKey {
			db.Delete([]byte(key))
			db.Set([]byte(newKey), bytes)
		}
	}
}
