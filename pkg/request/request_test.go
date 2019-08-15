package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/url"
	"testing"
)

func MustParseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return u
}

func TestRequest_Text(t *testing.T) {
	request := New(MustParseURL("https://httpbin.org/base64/VGVzdA=="))
	defer func() { _ = request.Close() }()

	body, err := request.Body()
	require.NoError(t, err)
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	assert.Equal(t, "Test", string(data))
}

func TestRequest_Document(t *testing.T) {
	request := New(MustParseURL("http://example.com/"))
	defer func() { _ = request.Close() }()

	document, err := request.Document()
	require.NoError(t, err)

	headline := document.Find("div > h1").Text()
	assert.Equal(t, "Example Domain", headline)
}
