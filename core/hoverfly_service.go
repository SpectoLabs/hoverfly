package hoverfly

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
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
	return this.SetModeWithArguments(v2.ModeView{
		Mode: mode,
	})
}

func (this *Hoverfly) SetModeWithArguments(modeView v2.ModeView) error {

	availableModes := map[string]bool{
		modes.Simulate:   true,
		modes.Capture:    true,
		modes.Modify:     true,
		modes.Synthesize: true,
	}

	if modeView.Mode == "" || !availableModes[modeView.Mode] {
		log.Error("Can't change mode to \"%d\"", modeView.Mode)
		return fmt.Errorf("Not a valid mode")
	}

	if this.Cfg.Webserver && modeView.Mode == modes.Capture {
		log.Error("Can't change mode to when configured as a webserver")
		return fmt.Errorf("Cannot change the mode of Hoverfly to capture when running as a webserver")
	}

	for _, header := range modeView.Arguments.Headers {
		if header == "*" {
			if len(modeView.Arguments.Headers) > 1 {
				return errors.New("Must provide a list containing only an asterix, or a list containing only headers names")
			}
		}
	}

	this.Cfg.SetMode(modeView.Mode)
	if this.Cfg.GetMode() == "capture" {
		this.CacheMatcher.FlushCache()
	}

	modeArguments := modes.ModeArguments{
		Headers: modeView.Arguments.Headers,
	}

	this.modeMap[this.Cfg.GetMode()].SetArguments(modeArguments)
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
	return len(hf.Simulation.Templates), nil
}

func (this Hoverfly) GetMetadataCache() cache.Cache {
	return this.MetadataCache
}

func (this Hoverfly) GetCache() ([]v2.RequestResponsePairViewV2, error) {
	return this.CacheMatcher.GetAllResponses()
}

func (hf Hoverfly) FlushCache() error {
	return hf.CacheMatcher.FlushCache()
}

func (hf *Hoverfly) GetResponseDelays() v1.ResponseDelayPayloadView {
	return hf.Simulation.ResponseDelays.ConvertToResponseDelayPayloadView()
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

	hf.Simulation.ResponseDelays = &responseDelays
	return nil
}

func (hf *Hoverfly) DeleteResponseDelays() {
	hf.Simulation.ResponseDelays = &models.ResponseDelayList{}
}

func (hf Hoverfly) GetStats() metrics.Stats {
	return hf.Counter.Flush()
}

func (hf Hoverfly) GetSimulation() (v2.SimulationViewV2, error) {
	pairViews := make([]v2.RequestResponsePairViewV2, 0)

	for _, v := range hf.Simulation.Templates {
		pairViews = append(pairViews, v.BuildView())
	}

	responseDelays := hf.Simulation.ResponseDelays.ConvertToResponseDelayPayloadView()

	return v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: pairViews,
			GlobalActions: v2.GlobalActionsView{
				Delays: responseDelays.Data,
			},
		},
		*v2.NewMetaView(hf.version),
	}, nil
}

func (this *Hoverfly) PutSimulation(simulationView v2.SimulationViewV2) error {
	err := this.ImportRequestResponsePairViews(simulationView.DataViewV2.RequestResponsePairs)
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
	this.FlushCache()
}

func (this Hoverfly) GetVersion() string {
	return this.version
}

func (this Hoverfly) GetUpstreamProxy() string {
	return this.Cfg.UpstreamProxy
}
