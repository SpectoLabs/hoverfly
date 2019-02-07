package matching

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/errors"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	goCache "github.com/patrickmn/go-cache"
)

type CacheMatcher struct {
	RequestCache cache.Cache
	Webserver    bool
	NewRequestCache *goCache.Cache
}

// getResponse returns stored response from cache
func (this *CacheMatcher) GetCachedResponse(req *models.RequestDetails) (*models.CachedResponse, *errors.HoverflyError) {
	if this.NewRequestCache == nil {
		return nil, errors.NoCacheSetError()
	}

	log.Debug("Checking cache for request")

	var key string

	if this.Webserver {
		key = req.HashWithoutHost()
	} else {
		key = req.Hash()
	}

	cachedResponse, found := this.NewRequestCache.Get(key)

	if !found {
		log.WithFields(log.Fields{
			"key":         key,
			//"error":       err.Error(),
			"query":       req.Query,
			"path":        req.Path,
			"destination": req.Destination,
			"method":      req.Method,
		}).Debug("Failed to retrieve response from cache")

		return nil, errors.RecordedRequestNotInCacheError()
	}

	// getting cache response
	//cachedResponse, err := models.NewCachedResponseFromBytes(pairBytes)
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"error": err.Error(),
	//		"value": string(pairBytes),
	//		"key":   key,
	//	}).Debug("Failed to decode payload from cache")
	//	return nil, errors.DecodePayloadError()
	//}

	log.WithFields(log.Fields{
		"key":         key,
		"path":        req.Path,
		"rawQuery":    req.Query,
		"method":      req.Method,
		"destination": req.Destination,
	}).Info("Response found interface{} cache")

	response := cachedResponse.(models.CachedResponse)
	return &response, nil
}

func (this *CacheMatcher) GetAllResponses() (v2.CacheView, error) {
	cacheView := v2.CacheView{}

	if this.NewRequestCache == nil {
		return cacheView, errors.NoCacheSetError()
	}

	records := this.NewRequestCache.Items()
	//if err != nil {
	//	return cacheView, err
	//}

	for key, v := range records {
		cachedResponse := v.Object.(models.CachedResponse)
		//if cachedResponse, err := models.NewCachedResponseFromBytes(v); err == nil {

		var pair *v2.RequestMatcherResponsePairViewV5
		var closestMiss *v2.ClosestMissView

		if cachedResponse.MatchingPair != nil {
			pairView := cachedResponse.MatchingPair.BuildView()
			pair = &pairView
		}

		if cachedResponse.ClosestMiss != nil {
			closestMiss = cachedResponse.ClosestMiss.BuildView()
		}

		cachedResponseView := v2.CachedResponseView{
			Key:          key,
			MatchingPair: pair,
			ClosestMiss:  closestMiss,
		}

		cacheView.Cache = append(cacheView.Cache, cachedResponseView)

		//} else {
		//	log.Error(err)
		//	return cacheView, err
		//}
	}

	return cacheView, nil
}

// TODO: This would be easier to reason about if we had two methods, "CacheHit" and "CacheHit" in order to reduce bloating
func (this *CacheMatcher) SaveRequestMatcherResponsePair(request models.RequestDetails, pair *models.RequestMatcherResponsePair, matchError *models.MatchError) error {
	if this.NewRequestCache == nil {
		return errors.NoCacheSetError()
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
	}

	if matchError != nil {
		cachedResponse.ClosestMiss = matchError.ClosestMiss
	}

	//pairBytes, err := cachedResponse.Encode()

	//if err != nil {
	//	return err
	//}

	this.NewRequestCache.Set(key, cachedResponse, goCache.DefaultExpiration)
	return nil
}

func (this *CacheMatcher) FlushCache() error {
	if this.NewRequestCache == nil {
		return errors.NoCacheSetError()
	}

	this.NewRequestCache.Flush()
	return nil
}

func (this *CacheMatcher) PreloadCache(simulation models.Simulation) error {
	if this.NewRequestCache == nil {
		return errors.NoCacheSetError()
	}
	for _, pair := range simulation.GetMatchingPairs() {
		if requestDetails := pair.RequestMatcher.ToEagerlyCachable(); requestDetails != nil {
			this.SaveRequestMatcherResponsePair(*requestDetails, &pair, nil)
		}
	}

	return nil
}