/*
Package request provides a request object for accessing websites.
*/
package request

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"sync"
)

// Request is a wrapper around a url
// which makes it easy to access the data
// it points to.
type Request struct {
	URL string
	ctx context.Context

	request  *http.Request
	respMux  sync.Mutex
	response *http.Response
	text     string
	document *goquery.Document
}

// New creates a new request and initialises it
// with the provided url.
func New(url string) *Request {
	return NewWithContext(nil, url)
}

// NewWithContext creates a new request with the given context.
// Note that nil is a valid context.
func NewWithContext(ctx context.Context, url string) *Request {
	return &Request{URL: url, ctx: ctx}
}

// Close performs cleanup for the Request.
// This is a no-op if the Request doesn't need cleanup
func (req *Request) Close() {
	if req.response != nil {
		_ = req.response.Body.Close()
	}
}

// Reset closes the request and removes all cached data.
func (req *Request) Reset() {
	req.respMux.Lock()
	defer req.respMux.Unlock()

	req.Close()

	req.request = nil
	req.response = nil
	req.text = ""
	req.document = nil
}

// Context returns the context of the request.
func (req *Request) Context() context.Context {
	if req.ctx == nil {
		return context.Background()
	}

	return req.ctx
}

// Request creates an http.Request (GET) for the url
// and returns it. The request is internally cached
// so calling this method multiple times will return
// the same http.Request.
func (req *Request) Request() *http.Request {
	if req.request == nil {
		request, _ := http.NewRequest("GET", req.URL, nil)
		req.request = request.WithContext(req.Context())
	}

	return req.request
}

// Response performs the Request.Request and returns the response/error.
// The Response is internally cached and will be cleaned up by Request.Close.
func (req *Request) Response() (*http.Response, error) {
	req.respMux.Lock()
	defer req.respMux.Unlock()

	if req.response == nil {
		request := req.Request()

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}

		req.response = resp
	}

	return req.response, nil
}

// Text retrieves the text response.
// It reads the text from the Request.Response body.
// The text is internally cached
func (req *Request) Text() (string, error) {
	if req.document == nil {
		resp, err := req.Response()
		if err != nil {
			return "", err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		text := string(data)
		req.text = text
	}

	return req.text, nil
}

// Document returns a goquery.Document for the Request.Response
// The Document is internally cached.
func (req *Request) Document() (*goquery.Document, error) {
	if req.document == nil {
		resp, err := req.Response()
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		req.document = doc
	}

	return req.document, nil
}
