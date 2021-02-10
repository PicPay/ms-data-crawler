package v2

type HttpConfiguration struct {
	URL         string
	Method      string
	Headers     map[string]string
	ContentType string
	QueryParams map[string]string
	Data        interface{}
}

type HttpResponse struct {
	Status int
	Body   []byte
}
