package httpclient

type Response struct {
	Status int
	Body   []byte
	Method string
	URL    string
	Curl   string
	Err    error
}

func NewResponse(status int, body []byte, url, method, curl string, err error) Response {
	return Response{
		URL:    url,
		Method: method,
		Curl:   curl,
		Status: status,
		Body:   body,
		Err:    err,
	}
}

func responseError(err error, url, method, curl string) Response {
	return Response{
		URL:    url,
		Method: method,
		Curl:   curl,
		Err:    err,
	}
}

// IsStatusOk returns true if the status code is between 200 and 299
func (r Response) IsStatusOk() bool {
	return r.Status >= 200 && r.Status < 300
}
