package hoverfly

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
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
	return hf.Cfg.Middleware.FullCommand
}

func (hf Hoverfly) GetMiddlewareV2() (string, string) {
	script, _ := hf.Cfg.Middleware.GetScript()
	return hf.Cfg.Middleware.Binary, script
}

func (hf Hoverfly) SetMiddleware(middleware string) error {
	if middleware == "" {
		hf.Cfg.Middleware.FullCommand = middleware
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

	middlewareObject := &Middleware{
		FullCommand: middleware,
	}

	_, err := middlewareObject.executeMiddlewareLocally(originalPair)
	if err != nil {
		return err
	}

	hf.Cfg.Middleware.FullCommand = middleware
	return nil
}

func (hf *Hoverfly) SetMiddlewareV2(binary, script string) error {
	newMiddleware := Middleware{}

	if binary == "" && script == "" {
		hf.Cfg.Middleware = newMiddleware
		return nil
	} else if binary == "" {
		return fmt.Errorf("Cannot run script with no binary")
	}

	err := newMiddleware.SetBinary(binary)
	if err != nil {
		return err
	}

	err = newMiddleware.SetScript(script)
	if err != nil {
		return nil
	}

	testData := models.RequestResponsePair{
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
	_, err = newMiddleware.executeMiddlewareLocally(testData)
	if err != nil {
		return err
	}

	hf.Cfg.Middleware = newMiddleware
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

func (hf *Hoverfly) GetResponseDelays() v1.ResponseDelayPayloadView {
	return hf.ResponseDelays.ConvertToResponseDelayPayloadView()
}

func (hf *Hoverfly) SetResponseDelays(payloadView v1.ResponseDelayPayloadView) error {
	err := models.ValidateResponseDelayPayload(payloadView)
	if err != nil {
		return err
	}

	var responseDelays models.ResponseDelayList

	for _, responseDelayView := range payloadView.Data {
		responseDelays = append(responseDelays, models.ResponseDelay{
			UrlPattern: responseDelayView.UrlPattern,
			HttpMethod: responseDelayView.HttpMethod,
			Delay:      responseDelayView.Delay,
		})
	}

	hf.ResponseDelays = &responseDelays
	return nil
}

func (hf *Hoverfly) DeleteResponseDelays() {
	hf.ResponseDelays = &models.ResponseDelayList{}
}

func (hf Hoverfly) GetStats() metrics.Stats {
	return hf.Counter.Flush()
}

func (hf Hoverfly) GetRecords() ([]v1.RequestResponsePairView, error) {
	records, err := hf.RequestCache.GetAllEntries()
	if err != nil {
		return nil, err
	}

	var pairViews []v1.RequestResponsePairView

	for _, v := range records {
		if pair, err := models.NewRequestResponsePairFromBytes(v); err == nil {
			pairView := pair.ConvertToV1RequestResponsePairView()
			pairViews = append(pairViews, *pairView)
		} else {
			log.Error(err)
			return nil, err
		}
	}

	for _, v := range hf.RequestMatcher.TemplateStore {
		pairView := v.ConvertToV1RequestResponsePairView()
		pairViews = append(pairViews, pairView)
	}

	return pairViews, nil

}

func (hf Hoverfly) GetSimulation() (v2.SimulationView, error) {
	records, err := hf.RequestCache.GetAllEntries()
	if err != nil {
		return v2.SimulationView{}, err
	}

	pairViews := make([]v2.RequestResponsePairView, 0)

	for _, v := range records {
		if pair, err := models.NewRequestResponsePairFromBytes(v); err == nil {
			pairView := pair.ConvertToRequestResponsePairView()
			pairViews = append(pairViews, pairView)
		} else {
			log.Error(err)
			return v2.SimulationView{}, err
		}
	}

	for _, v := range hf.RequestMatcher.TemplateStore {
		pairViews = append(pairViews, v.ConvertToRequestResponsePairView())
	}

	responseDelays := hf.ResponseDelays.ConvertToResponseDelayPayloadView()

	return v2.SimulationView{
		MetaView: v2.MetaView{
			HoverflyVersion: "v0.9.2",
			SchemaVersion:   "v1",
			TimeExported:    time.Now().Format(time.RFC3339),
		},
		DataView: v2.DataView{
			RequestResponsePairs: pairViews,
			GlobalActions: v2.GlobalActionsView{
				Delays: responseDelays.Data,
			},
		},
	}, nil
}

func (this *Hoverfly) PutSimulation(simulationView v2.SimulationView) error {
	requestResponsePairViews := make([]interfaces.RequestResponsePair, len(simulationView.RequestResponsePairs))
	for i, v := range simulationView.RequestResponsePairs {
		requestResponsePairViews[i] = v
	}

	err := this.ImportRequestResponsePairViews(requestResponsePairViews)
	if err != nil {
		return err
	}

	err = this.SetResponseDelays(v1.ResponseDelayPayloadView{Data: simulationView.GlobalActions.Delays})
	if err != nil {
		return err
	}

	return nil
}

func (this *Hoverfly) DeleteSimulation() {
	this.DeleteTemplateCache()
	this.DeleteResponseDelays()
	this.DeleteRequestCache()
}
