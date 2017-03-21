package matching

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type CacheMatcher struct {
	RequestCache cache.Cache
	Webserver    bool
}

// getResponse returns stored response from cache
func (this *CacheMatcher) GetCachedResponse(req *models.RequestDetails) (*models.CachedResponse, *MatchingError) {
	if this.RequestCache == nil {
		return nil, &MatchingError{
			Description: "No cache set",
		}
	}

	log.Debug("Checking cache for request")

	var key string

	if this.Webserver {
		key = req.HashWithoutHost()
	} else {
		key = req.Hash()
	}

	pairBytes, err := this.RequestCache.Get([]byte(key))

	if err != nil {
		log.WithFields(log.Fields{
			"key":         key,
			"error":       err.Error(),
			"query":       req.Query,
			"path":        req.Path,
			"destination": req.Destination,
			"method":      req.Method,
		}).Debug("Failed to retrieve response from cache")

		return nil, &MatchingError{
			StatusCode:  412,
			Description: "Could not find recorded request in cache",
		}
	}

	// getting cache response
	cachedResponse, err := models.NewCachedResponseFromBytes(pairBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"value": string(pairBytes),
			"key":   key,
		}).Debug("Failed to decode payload from cache")
		return nil, &MatchingError{
			StatusCode:  500,
			Description: "Failed to decode payload from cache",
		}
	}

	log.WithFields(log.Fields{
		"key":         key,
		"path":        req.Path,
		"rawQuery":    req.Query,
		"method":      req.Method,
		"destination": req.Destination,
		"status":      cachedResponse.MatchingPair.Response.Status,
	}).Info("Response found interface{} cache")

	return cachedResponse, nil
}

func (this CacheMatcher) GetAllResponses() ([]v2.RequestResponsePairViewV2, error) {
	if this.RequestCache == nil {
		return nil, &MatchingError{
			Description: "No cache set",
		}
	}

	records, err := this.RequestCache.GetAllEntries()
	if err != nil {
		return []v2.RequestResponsePairViewV2{}, err
	}

	pairViews := []v2.RequestResponsePairViewV2{}

	for _, v := range records {
		if cachedResponse, err := models.NewCachedResponseFromBytes(v); err == nil {
			pairView := cachedResponse.MatchingPair.BuildView()
			pairViews = append(pairViews, pairView)
		} else {
			log.Error(err)
			return []v2.RequestResponsePairViewV2{}, err
		}
	}

	return pairViews, nil
}

func (this *CacheMatcher) SaveRequestTemplateResponsePair(request models.RequestDetails, pair *models.RequestTemplateResponsePair) error {
	if this.RequestCache == nil {
		return errors.New("No cache set")
	}

	var key string

	if this.Webserver {
		key = request.HashWithoutHost()
	} else {
		key = request.Hash()
	}

	log.WithFields(log.Fields{
		"path":          request.Path,
		"rawQuery":      request.Query,
		"requestMethod": request.Method,
		"bodyLen":       len(request.Body),
		"destination":   request.Destination,
		"hashKey":       key,
	}).Debug("Saving response to cache")

	cachedResponse := models.CachedResponse{
		Request:      request,
		MatchingPair: pair,
		HeaderMatch:  len(pair.RequestTemplate.Headers) > 0,
	}

	pairBytes, err := cachedResponse.Encode()

	if err != nil {
		return err
	}

	return this.RequestCache.Set([]byte(key), pairBytes)
}

func (this CacheMatcher) FlushCache() error {
	if this.RequestCache == nil {
		return errors.New("No cache set")
	}

	return this.RequestCache.DeleteData()
}

func (this CacheMatcher) PreloadCache(simulation models.Simulation) error {
	if this.RequestCache == nil {
		return errors.New("No cache set")
	}
	for _, pair := range simulation.Templates {
		if requestDetails := pair.RequestTemplate.BuildRequestDetailsFromExactMatches(); requestDetails != nil {
			this.SaveRequestTemplateResponsePair(*requestDetails, &pair)
		}
	}

	return nil
}
