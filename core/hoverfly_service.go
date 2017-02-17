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
	"github.com/SpectoLabs/hoverfly/core/modes"
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
		modes.Simulate:   true,
		modes.Capture:    true,
		modes.Modify:     true,
		modes.Synthesize: true,
	}

	if mode == "" || !availableModes[mode] {
		log.Error("Can't change mode to \"%d\"", mode)
		return fmt.Errorf("Not a valid mode")
	}

	if this.Cfg.Webserver && mode == modes.Capture {
		log.Error("Can't change mode to when configured as a webserver")
		return fmt.Errorf("Can't change mode to capture when configured as a webserver")
	}

	this.Cfg.SetMode(mode)
	return nil
}

func (hf Hoverfly) GetMiddleware() (string, string, string) {
	script, _ := hf.Cfg.Middleware.GetScript()
	return hf.Cfg.Middleware.Binary, script, hf.Cfg.Middleware.Remote
}

func (hf *Hoverfly) SetMiddleware(binary, script, remote string) error {
	newMiddleware := Middleware{}

	if binary == "" && script == "" && remote == "" {
		hf.Cfg.Middleware = newMiddleware
		return nil
	}

	if binary == "" && script != "" {
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

	err = newMiddleware.SetRemote(remote)
	if err != nil {
		return err
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
	_, err = newMiddleware.Execute(testData)
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

	for _, v := range hf.Simulation.Templates {
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

	for _, v := range hf.Simulation.Templates {
		pairViews = append(pairViews, v.ConvertToRequestResponsePairView())
	}

	responseDelays := hf.ResponseDelays.ConvertToResponseDelayPayloadView()

	return v2.SimulationView{
		MetaView: v2.MetaView{
			HoverflyVersion: hf.version,
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
	var templates []models.RequestTemplateResponsePair
	this.Simulation.Templates = templates
	this.DeleteResponseDelays()
	this.DeleteRequestCache()
}

func (this Hoverfly) GetVersion() string {
	return this.version
}

func (this Hoverfly) GetUpstreamProxy() string {
	return this.Cfg.UpstreamProxy
}
