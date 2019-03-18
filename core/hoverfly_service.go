package hoverfly

import (
	"errors"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/delay"
	"regexp"

	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/middleware"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/state"
	"github.com/SpectoLabs/hoverfly/core/util"
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

func (this Hoverfly) GetMode() v2.ModeView {
	return this.modeMap[this.Cfg.Mode].View()
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
		modes.Spy:        true,
		modes.Diff:       true,
	}

	if modeView.Mode == "" || !availableModes[modeView.Mode] {
		log.WithFields(log.Fields{
			"mode": modeView.Mode,
		}).Error("Unknown mode")
		return fmt.Errorf("Not a valid mode")
	}

	if this.Cfg.Webserver && modeView.Mode == modes.Capture {
		log.Error("Cannot change the mode of Hoverfly to capture when running as a webserver")
		return fmt.Errorf("Cannot change the mode of Hoverfly to capture when running as a webserver")
	}

	for _, header := range modeView.Arguments.Headers {
		if header == "*" {
			if len(modeView.Arguments.Headers) > 1 {
				return errors.New("Must provide a list containing only an asterix, or a list containing only headers names")
			}
		}
	}

	matchingStrategy := modeView.Arguments.MatchingStrategy
	if modeView.Mode == modes.Simulate {
		if matchingStrategy == nil {
			matchingStrategy = util.StringToPointer("strongest")
		}

		if strings.ToLower(*matchingStrategy) != "strongest" && strings.ToLower(*matchingStrategy) != "first" {
			return errors.New("Only matching strategy of 'first' or 'strongest' is permitted")
		}
	}

	this.Cfg.SetMode(modeView.Mode)
	if this.Cfg.GetMode() == "capture" {
		this.CacheMatcher.FlushCache()
	} else if this.Cfg.GetMode() == "simulate" {
		this.CacheMatcher.PreloadCache(*this.Simulation)
	} else if this.Cfg.GetMode() == "spy" {
		this.CacheMatcher.PreloadCache(*this.Simulation)
	}

	modeArguments := modes.ModeArguments{
		Headers:          modeView.Arguments.Headers,
		MatchingStrategy: matchingStrategy,
		Stateful:         modeView.Arguments.Stateful,
	}

	this.modeMap[this.Cfg.GetMode()].SetArguments(modeArguments)

	log.WithFields(log.Fields{
		"mode": this.Cfg.GetMode(),
	}).Info("Mode has been changed")

	return nil
}

func (hf Hoverfly) GetMiddleware() (string, string, string) {
	script, _ := hf.Cfg.Middleware.GetScript()
	return hf.Cfg.Middleware.Binary, script, hf.Cfg.Middleware.Remote
}

func (hf *Hoverfly) SetMiddleware(binary, script, remote string) error {
	newMiddleware := &middleware.Middleware{}
	if binary == "" && script == "" && remote == "" {
		hf.Cfg.Middleware = *newMiddleware
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
		return err
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
			Query:       map[string][]string{},
			Body:        "",
			Headers:     map[string][]string{"test_header": {"true"}},
		},
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "ok",
			Headers: map[string][]string{"test_header": {"true"}},
		},
	}

	_, err = newMiddleware.Execute(testData)
	if err != nil {
		return err
	}
	hf.Cfg.Middleware = *newMiddleware
	return nil
}

func (hf Hoverfly) GetRequestCacheCount() (int, error) {
	return len(hf.Simulation.GetMatchingPairs()), nil
}

func (this Hoverfly) GetCache() (v2.CacheView, error) {
	return this.CacheMatcher.GetAllResponses()
}

func (hf Hoverfly) FlushCache() error {
	return hf.CacheMatcher.FlushCache()
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

func (hf *Hoverfly) SetResponseDelaysLogNormal(payloadView v1.ResponseDelayLogNormalPayloadView) error {
	err := models.ValidateResponseDelayLogNormalPayload(payloadView)
	if err != nil {
		return err
	}

	var responseDelaysLogNormal models.ResponseDelayLogNormalList

	for _, responseDelayView := range payloadView.Data {
		responseDelaysLogNormal = append(responseDelaysLogNormal, models.ResponseDelayLogNormal{
			UrlPattern: responseDelayView.UrlPattern,
			HttpMethod: responseDelayView.HttpMethod,
			Min:        responseDelayView.Min,
			Max:        responseDelayView.Max,
			Mean:       responseDelayView.Mean,
			Median:     responseDelayView.Median,
			DelayGenerator: delay.NewLogNormalGenerator(
				responseDelayView.Min,
				responseDelayView.Max,
				responseDelayView.Mean,
				responseDelayView.Median,
			),
		})
	}

	hf.Simulation.ResponseDelaysLogNormal = &responseDelaysLogNormal
	return nil
}

func (hf *Hoverfly) DeleteResponseDelays() {
	hf.Simulation.ResponseDelays = &models.ResponseDelayList{}
}

func (hf *Hoverfly) DeleteResponseDelaysLogNormal() {
	hf.Simulation.ResponseDelaysLogNormal = &models.ResponseDelayLogNormalList{}
}

func (hf Hoverfly) GetStats() metrics.Stats {
	return hf.Counter.Flush()
}

func (hf Hoverfly) GetSimulation() (v2.SimulationViewV5, error) {
	pairViews := make([]v2.RequestMatcherResponsePairViewV5, 0)

	for _, v := range hf.Simulation.GetMatchingPairs() {
		pairViews = append(pairViews, v.BuildView())
	}

	return v2.BuildSimulationView(pairViews,
		hf.Simulation.ResponseDelays.ConvertToResponseDelayPayloadView(),
		hf.Simulation.ResponseDelaysLogNormal.ConvertToResponseDelayLogNormalPayloadView(),
		hf.version), nil
}

func (hf Hoverfly) GetFilteredSimulation(urlPattern string) (v2.SimulationViewV5, error) {
	pairViews := make([]v2.RequestMatcherResponsePairViewV5, 0)
	regexPattern, err := regexp.Compile(urlPattern)

	if err != nil {
		return v2.SimulationViewV5{}, err
	}

	for _, v := range hf.Simulation.GetMatchingPairs() {

		var urlStringToMatch string
		if v.RequestMatcher.Destination != nil && len(v.RequestMatcher.Destination) != 0 && v.RequestMatcher.Destination[0].Matcher == matchers.Exact {
			urlStringToMatch += v.RequestMatcher.Destination[0].Value.(string)
		}
		if v.RequestMatcher.Path != nil && len(v.RequestMatcher.Path) != 0 && v.RequestMatcher.Path[0].Matcher == matchers.Exact {
			urlStringToMatch += v.RequestMatcher.Path[0].Value.(string)
		}

		if regexPattern.MatchString(urlStringToMatch) {
			pairViews = append(pairViews, v.BuildView())
		}
	}

	return v2.BuildSimulationView(pairViews,
		hf.Simulation.ResponseDelays.ConvertToResponseDelayPayloadView(),
		hf.Simulation.ResponseDelaysLogNormal.ConvertToResponseDelayLogNormalPayloadView(),
		hf.version), nil
}

func (this *Hoverfly) PutSimulation(simulationView v2.SimulationViewV5) v2.SimulationImportResult {
	result := this.importRequestResponsePairViews(simulationView.DataViewV5.RequestResponsePairs)

	result.AddError(this.SetResponseDelays(v1.ResponseDelayPayloadView{Data: simulationView.GlobalActions.Delays}))
	result.AddError(this.SetResponseDelaysLogNormal(v1.ResponseDelayLogNormalPayloadView{Data: simulationView.GlobalActions.DelaysLogNormal}))

	return result
}

func (this *Hoverfly) DeleteSimulation() {
	this.Simulation.DeleteMatchingPairs()
	this.DeleteResponseDelays()
	this.DeleteResponseDelaysLogNormal()
	this.FlushCache()
}

func (this Hoverfly) GetVersion() string {
	return this.version
}

func (this Hoverfly) GetUpstreamProxy() string {
	return this.Cfg.UpstreamProxy
}

func (this Hoverfly) IsWebServer() bool {

	return this.Cfg.Webserver
}

func (this Hoverfly) IsMiddlewareSet() bool {
	return this.Cfg.Middleware.IsSet()
}

func (this *Hoverfly) GetState() map[string]string {
	return this.state.State
}

func (this *Hoverfly) SetState(state map[string]string) {
	this.state.SetState(state)
}

func (this *Hoverfly) PatchState(toPatch map[string]string) {
	this.state.PatchState(toPatch)
}

func (this *Hoverfly) ClearState() {
	this.state = state.NewState()
}

func (this *Hoverfly) GetDiff() map[v2.SimpleRequestDefinitionView][]v2.DiffReport {
	return this.responsesDiff
}

func (this *Hoverfly) ClearDiff() {
	this.responsesDiff = make(map[v2.SimpleRequestDefinitionView][]v2.DiffReport)
}

func (this *Hoverfly) AddDiff(requestView v2.SimpleRequestDefinitionView, diffReport v2.DiffReport) {
	if len(diffReport.DiffEntries) > 0 {
		diffs := this.responsesDiff[requestView]
		this.responsesDiff[requestView] = append(diffs, diffReport)
	}
}

func (this *Hoverfly) GetPACFile() []byte {
	return this.Cfg.PACFile
}

func (this *Hoverfly) SetPACFile(pacFile []byte) {
	if len(pacFile) == 0 {
		pacFile = nil
	}
	this.Cfg.PACFile = pacFile
}

func (this *Hoverfly) DeletePACFile() {
	this.Cfg.PACFile = nil
}
