package internal

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
)

type fakeTransport struct {
	W io.Writer
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	io.WriteString(t.W, "\nrequest:\n")
	io.WriteString(t.W, "----------------------------------------\n")
	t.W.Write(b)
	io.WriteString(t.W, "\n")

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
	}, nil
}
