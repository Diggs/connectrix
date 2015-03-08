package event

type Event struct {
	Namespace  string
	Source     string
	Type       string
	ParserName string
	Content    string
	Object     interface {}
}
