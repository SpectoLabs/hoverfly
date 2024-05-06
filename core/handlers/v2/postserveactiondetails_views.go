package v2

type PostServeActionDetailsView struct {
	Actions []ActionView `json:"actions,omitempty"`
}

type ActionView struct {
	ActionName    string `json:"actionName,omitempty"`
	Binary        string `json:"binary,omitempty"`
	ScriptContent string `json:"script,omitempty"`
	Remote        string `json:"remote,omitempty"`
	DelayInMs     int    `json:"delayInMs,omitempty"`
}
