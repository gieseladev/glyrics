package models

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
)

// Request is a wrapper around a url
// which makes it easy to access the data
// it points to.
type Request struct {
	Url      string
	request  *http.Request
	response *http.Response
	text     string
	document *goquery.Document
}

// NewRequest creates a new request and initialises it
// with the provided url.
func NewRequest(url string) *Request {
	return &Request{Url: url}
}

// Close performs cleanup for the Request.
// This is a no-op if the Request doesn't need cleanup
func (req *Request) Close() {
	if req.response != nil {
		_ = req.response.Body.Close()
	}
}

// Request creates an http.Request (GET) for the url
// and returns it. The request is internally cached
// so calling this method multiple times will return
// the same http.Request.
func (req *Request) Request() *http.Request {
	if req.request == nil {
		request, _ := http.NewRequest("GET", req.Url, nil)
		req.request = request
	}
	return req.request
}

// Response performs the Request.Request and returns the response/error.
// The Response is internally cached and will be cleaned up by Request.Close.
func (req *Request) Response() (*http.Response, error) {
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
