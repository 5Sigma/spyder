package request

import (
	"net/http"
	"strings"
	"time"
)

// Response - Helper object returned from making a request. It holds all the
// relevant request/response data.
type (
	Response struct {
		Response    *http.Response
		RequestTime time.Duration
		Content     []byte
		Request     *http.Request
		Payload     []byte
	}
)

// IsResponseJSON - Returns true if the response is JSON. This is done by
// testing the Content-Type header.
func (res *Response) IsResponseJSON() bool {
	contentType := strings.ToLower(res.Response.Header.Get("Content-Type"))
	return strings.Contains(contentType, "json")
}
