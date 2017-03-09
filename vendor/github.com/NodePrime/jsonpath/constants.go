package jsonpath

const (
	BadStructure         = "Bad Structure"
	NoMoreResults        = "No more results"
	UnexpectedToken      = "Unexpected token in evaluation"
	AbruptTokenStreamEnd = "Token reader is not sending anymore tokens"
)

var (
	bytesTrue  = []byte{'t', 'r', 'u', 'e'}
	bytesFalse = []byte{'f', 'a', 'l', 's', 'e'}
	bytesNull  = []byte{'n', 'u', 'l', 'l'}
)
