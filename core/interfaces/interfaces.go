package interfaces

type RequestResponsePair interface {
	GetRequest() Request
	GetResponse() Response
}

type Request interface {
	GetPath() *string
	GetMethod() *string
	GetDestination() *string
	GetScheme() *string
	GetQuery() *string
	GetBody() *string
	GetHeaders() map[string][]string
}

type Response interface {
	GetStatus() int
	GetBody() string
	GetEncodedBody() bool
	GetTemplated() bool
	GetHeaders() map[string][]string
	GetTransitionsState() map[string]string
	GetRemovesState() []string
}
