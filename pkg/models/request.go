package models

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Url      string
	request  *http.Request
	response *http.Response
	text     string
	document *goquery.Document
}

func NewRequest(url string) *Request {
	return &Request{Url: url}
}

func (req *Request) Close() {
	if req.response != nil {
		req.response.Body.Close()
	}
}

func (req *Request) Request() *http.Request {
	if req.request == nil {
		request, _ := http.NewRequest("GET", req.Url, nil)
		req.request = request
	}
	return req.request
}

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
