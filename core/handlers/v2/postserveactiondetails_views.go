package v2

type PostServeActionDetailsView struct {
	Hooks []HookView `json:"hooks,omitempty"`
}

type HookView struct {
	HookName            string `json:"hookName"`
	Binary              string `json:"binary"`
	ScriptContent       string `json:"script"`
	DelayInMilliSeconds int    `json:"delayInMilliSeconds,omitempty"`
}
