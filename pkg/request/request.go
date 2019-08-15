/*
Package request provides a request object for accessing websites.
*/
package request

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// Requester is an interface which provides
type Requester interface {
	// Context returns the context used by the request.
	Context() context.Context

	// Close closes the request.
	Close() error

	// URL returns the url of the request.
	URL() *url.URL
	// Body returns a reader for the request's body.
	// To close the body, use the Close() method.
	Body() (io.Reader, error)
	// Document returns a goquery document for the body.
	Document() (*goquery.Document, error)
}

const (
	userAgent = "gLyrics/3 (https://github.com/gieseladev/glyrics)"
)

// httpRequest is an implementation of Requester based on the default http
// library.
type httpRequest struct {
	ctx      context.Context
	mux      sync.Mutex
	url      *url.URL
	response *http.Response
	document *goquery.Document
}

// New creates a new request and initialises it
// with the provided url.
func New(u *url.URL) Requester {
	return NewWithContext(nil, u)
}

// NewWithContext creates a new request with the given context.
// Note that nil is a valid context.
func NewWithContext(ctx context.Context, u *url.URL) Requester {
	return &httpRequest{url: u, ctx: ctx}
}

// Close performs cleanup for the Request.
// This is a no-op if the Request doesn't need cleanup
func (req *httpRequest) Close() error {
	if req.response != nil {
		return req.response.Body.Close()
	}

	return nil
}

// Context returns the context of the request.
func (req *httpRequest) Context() context.Context {
	if req.ctx == nil {
		return context.Background()
	}

	return req.ctx
}

// URL returns the url of the request.
func (req *httpRequest) URL() *url.URL {
	return req.url
}

// getResponse returns the http.Response.
// If the response is already cached, it is returned directly.
// Otherwise a new http request is started.
// This function DOESN'T use the lock, it is assumed that the caller holds the
// lock!
func (req *httpRequest) getResponse() (*http.Response, error) {
	if req.response != nil {
		return req.response, nil
	}

	header := make(http.Header)
	header.Set("User-Agent", userAgent)

	request := (&http.Request{
		Method: "GET",
		URL:    req.url,
		Header: header,
	}).WithContext(req.Context())

	resp, err := http.DefaultClient.Do(request)
	if err == nil {
		req.response = resp
	}

	return resp, err
}

// Body returns the http request's body.
func (req *httpRequest) Body() (io.Reader, error) {
	req.mux.Lock()
	defer req.mux.Unlock()

	resp, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// Document returns the goquery.Document for the body.
func (req *httpRequest) Document() (*goquery.Document, error) {
	req.mux.Lock()
	defer req.mux.Unlock()

	if req.document != nil {
		return req.document, nil
	}

	resp, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err == nil {
		doc.Url = req.url
	}

	return doc, err
}
