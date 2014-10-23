package parsers

type Parser interface {
	ParseContent(*[]byte) (interface{}, error)
}
