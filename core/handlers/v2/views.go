package v2

import (
	"github.com/SpectoLabs/hoverfly/core/metrics"
)

type DestinationView struct {
	Destination string `json:"destination"`
}

type UsageView struct {
	Usage metrics.Stats `json:"usage"`
}

type MiddlewareView struct {
	Binary string `json:"binary"`
	Script string `json:"script"`
	Remote string `json:"remote"`
}

type CORSView struct {
	Enabled          bool   `json:"enabled"`
	AllowOrigin      string `json:"allowOrigin,omitempty"`
	AllowMethods     string `json:"allowMethods,omitempty"`
	AllowHeaders     string `json:"allowHeaders,omitempty"`
	PreflightMaxAge  int64  `json:"preflightMaxAge,omitempty"`
	AllowCredentials bool   `json:"allowCredentials,omitempty"`
	ExposeHeaders    string `json:"exposeHeaders,omitempty"`
}

type ModeView struct {
	Mode      string            `json:"mode"`
	Arguments ModeArgumentsView `json:"arguments,omitempty"`
}

type ModeArgumentsView struct {
	Headers            []string `json:"headersWhitelist,omitempty"`
	MatchingStrategy   *string  `json:"matchingStrategy,omitempty"`
	Stateful           bool     `json:"stateful,omitempty"`
	OverwriteDuplicate bool     `json:"overwriteDuplicate,omitempty"`
	CaptureOnMiss      bool     `json:"captureOnMiss,omitempty"`
}

type IsWebServerView struct {
	IsWebServer bool `json:"isWebServer"`
}

type VersionView struct {
	Version string `json:"version"`
}

type UpstreamProxyView struct {
	UpstreamProxy string `json:"upstreamProxy"`
}

type HoverflyView struct {
	CORSView `json:"cors"`
	DestinationView
	MiddlewareView `json:"middleware"`
	ModeView
	IsWebServerView
	UsageView
	VersionView
	UpstreamProxyView
}

type LogsView struct {
	Logs []map[string]interface{} `json:"logs"`
}

type CacheView struct {
	Cache []CachedResponseView `json:"cache"`
}

type CachedResponseView struct {
	Key          string                            `json:"key"`
	MatchingPair *RequestMatcherResponsePairViewV5 `json:"matchingPair,omitempty"`
	HeaderMatch  bool                              `json:"headerMatch"`
	ClosestMiss  *ClosestMissView                  `json:"closestMiss"`
}

type ClosestMissView struct {
	Response       ResponseDetailsViewV5 `json:"response"`
	RequestMatcher RequestMatcherViewV5  `json:"requestMatcher"`
	MissedFields   []string              `json:"missedFields"`
}

type JournalView struct {
	Journal []JournalEntryView `json:"journal"`
	Index   []JournalIndexView `json:"indexes"`
	Offset  int                `json:"offset"`
	Limit   int                `json:"limit"`
	Total   int                `json:"total"`
}

type JournalEntryView struct {
	Request              RequestDetailsView        `json:"request"`
	Response             ResponseDetailsView       `json:"response"`
	Mode                 string                    `json:"mode"`
	TimeStarted          string                    `json:"timeStarted"`
	Latency              float64                   `json:"latency"`
	Id                   string                    `json:"id"`
	PostServeActionEntry *PostServeActionEntryView `json:"postServeAction,omitEmpty"`
}

type PostServeActionEntryView struct {
	ActionName    string `json:"name"`
	InvokedTime   string `json:"invoked"`
	CompletedTime string `json:"completed"`
	CorrelationId string `json:"correlationId,omitempty"`
	HttpStatus    int    `json:"status,omitempty"`
}

type JournalEntryFilterView struct {
	Request *RequestMatcherViewV5 `json:"request"`
}

type StateView struct {
	State map[string]string `json:"state"`
}

type DiffView struct {
	Diff []ResponseDiffForRequestView `json:"diff"`
}

type ResponseDiffForRequestView struct {
	Request    SimpleRequestDefinitionView `json:"request"`
	DiffReport []DiffReport                `json:"diffReports"`
}

type SimpleRequestDefinitionView struct {
	Method string `json:"method"`
	Host   string `json:"host"`
	Path   string `json:"path"`
	Query  string `json:"query"`
}

type DiffReport struct {
	Timestamp   string            `json:"timestamp"`
	DiffEntries []DiffReportEntry `json:"diffEntries"`
}

type DiffReportEntry struct {
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

type DiffFilterView struct {
	ExcludedHeaders        []string `json:"excludedHeaders"`
	ExcludedResponseFields []string `json:"excludedResponseFields"`
}

type JournalIndexView struct {
	Name    string                  `json:"name"`
	Entries []JournalIndexEntryView `json:"entries,omitempty"`
}

type JournalIndexEntryView struct {
	Key            string `json:"key"`
	JournalEntryId string `json:"journalEntryId"`
}

type JournalIndexRequestView struct {
	Name string `json:"name"`
}
