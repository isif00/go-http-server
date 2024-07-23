package types

type HttpRequest struct {
	Method    string
	Path      string
	Headers   map[string]string
	Body      string
	UserAgent string
}
