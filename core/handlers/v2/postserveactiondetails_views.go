package v2

type PostServeActionDetailsView struct {
	Actions []ActionView `json:"actions,omitempty"`
}

type ActionView struct {
	ActionName    string `json:"actionName"`
	Binary        string `json:"binary"`
	ScriptContent string `json:"script"`
	DelayInMs     int    `json:"delayInMs,omitempty"`
}
