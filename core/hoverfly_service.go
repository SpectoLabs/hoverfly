package hoverfly

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/views"
	"regexp"
)

func (this Hoverfly) GetDestination() string {
	return this.Cfg.Destination
}

// UpdateDestination - updates proxy with new destination regexp
func (hf *Hoverfly) SetDestination(destination string) (err error) {
	_, err = regexp.Compile(destination)
	if err != nil {
		return fmt.Errorf("destination is not a valid regular expression string")
	}

	hf.mu.Lock()
	hf.StopProxy()
	hf.Cfg.Destination = destination
	err = hf.StartProxy()
	hf.mu.Unlock()
	return
}

func (this Hoverfly) GetMode() string {
	return this.Cfg.Mode
}

func (this *Hoverfly) SetMode(mode string) error {
	availableModes := map[string]bool{
		"simulate":   true,
		"capture":    true,
		"modify":     true,
		"synthesize": true,
	}

	if mode == "" || !availableModes[mode] {
		log.Error("Can't change mode to \"%d\"", mode)
		return fmt.Errorf("Not a valid mode")
	}

	if this.Cfg.Webserver {
		log.Error("Can't change state when configured as a webserver ")
		return fmt.Errorf("Can't change state when configured as a webserver")
	}
	this.Cfg.SetMode(mode)
	return nil
}

func (hf Hoverfly) GetMiddleware() string {
	return hf.Cfg.Middleware
}

func (hf Hoverfly) SetMiddleware(middleware string) error {
	if middleware == "" {
		hf.Cfg.Middleware = middleware
		return nil
	}
	originalPair := models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "www.test.com",
			Scheme:      "",
			Query:       "",
			Body:        "",
			Headers:     map[string][]string{"test_header": []string{"true"}},
		},
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "ok",
			Headers: map[string][]string{"test_header": []string{"true"}},
		},
	}
	c := NewConstructor(nil, originalPair)
	err := c.ApplyMiddleware(middleware)
	if err != nil {
		return err
	}

	hf.Cfg.Middleware = middleware
	return nil
}

func (hf Hoverfly) GetRequestCacheCount() (int, error) {
	return hf.RequestCache.RecordsCount()
}

func (this Hoverfly) GetMetadataCache() cache.Cache {
	return this.MetadataCache
}

func (hf Hoverfly) DeleteRequestCache() error {
	return hf.RequestCache.DeleteData()
}

func (this Hoverfly) GetTemplates() v1.RequestTemplateResponsePairPayload {
	return this.RequestMatcher.TemplateStore.GetPayload()
}

func (this *Hoverfly) ImportTemplates(pairPayload v1.RequestTemplateResponsePairPayload) error {
	return this.RequestMatcher.TemplateStore.ImportPayloads(pairPayload)
}

func (this *Hoverfly) DeleteTemplateCache() {
	this.RequestMatcher.TemplateStore.Wipe()
}

func (hf *Hoverfly) GetResponseDelays() []byte {
	return hf.ResponseDelays.Json()
}

func (hf *Hoverfly) SetResponseDelays(payloadView v1.ResponseDelayPayload) error {
	err := models.ValidateResponseDelayJson(payloadView)
	if err != nil {
		return err
	}

	var responseDelays models.ResponseDelayList

	for _, responseDelayView := range *payloadView.Data {
		responseDelays = append(responseDelays, models.ResponseDelay{
			UrlPattern: responseDelayView.UrlPattern,
			HttpMethod: responseDelayView.HttpMethod,
			Delay:      responseDelayView.Delay,
		})
	}

	payload := models.ResponseDelayPayload{
		Data: &responseDelays,
	}

	hf.ResponseDelays = payload.Data
	return nil
}

func (hf *Hoverfly) DeleteResponseDelays() {
	hf.ResponseDelays = &models.ResponseDelayList{}
}

func (hf Hoverfly) GetStats() metrics.Stats {
	return hf.Counter.Flush()
}

func (hf Hoverfly) GetRecords() ([]views.RequestResponsePairView, error) {
	records, err := hf.RequestCache.GetAllEntries()
	if err != nil {
		return nil, err
	}

	var pairViews []views.RequestResponsePairView

	for _, v := range records {
		if pair, err := models.NewRequestResponsePairFromBytes(v); err == nil {
			pairView := pair.ConvertToRequestResponsePairView()
			pairViews = append(pairViews, *pairView)
		} else {
			log.Error(err)
			return nil, err
		}
	}

	for _, v := range hf.RequestMatcher.TemplateStore {
		pairView := v.ConvertToRequestResponsePairView()
		pairViews = append(pairViews, pairView)
	}

	return pairViews, nil

}
