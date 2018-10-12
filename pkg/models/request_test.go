package models

import (
	"testing"
)

func TestRequest_Response(t *testing.T) {
	request := NewRequest("https://www.google.com/")
	defer request.Close()

	resp, err := request.Response()
	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error("Got empty response")
	}
}

func TestRequest_Text(t *testing.T) {
	request := NewRequest("https://httpbin.org/base64/VGVzdA==")
	defer request.Close()

	text, err := request.Text()
	if err != nil {
		t.Error(err)
	}

	if text != "Test" {
		t.Errorf("Incorrect text response %s", text)
	}
}

func TestRequest_Document(t *testing.T) {
	request := NewRequest("https://www.google.com/")
	defer request.Close()

	document, err := request.Document()
	if err != nil {
		t.Error(err)
	}

	src, exists := document.Find("#hplogo[alt=\"Google\"]").Attr("src")

	if !(exists && src != "") {
		t.Error("Couldn't find google logo source in html")
	}
}
