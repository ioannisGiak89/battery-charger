package nationalgrid_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"ioannisGiak89/arenko/nationalgrid"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockedHttpClient is used to mock any functions from http.Client
type mockedHttpClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (cl *mockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	return cl.MockDo(req)
}

func TestClient_Get(t *testing.T) {

	t.Run("should return an error if the request fails", func(t *testing.T) {
		client := nationalgrid.NewClient(
			&mockedHttpClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("network request failed")
				},
			},
			"fakeURL",
		)

		responseBody, err := client.Get("path/to/endpoint")

		assert.Nil(t, responseBody)
		assert.Equal(t, errors.New("network request failed"), err)
	})

	t.Run("should return an error if status code is 404", func(t *testing.T) {
		client := nationalgrid.NewClient(
			&mockedHttpClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("not found"))),
						StatusCode: http.StatusNotFound,
					}, nil
				},
			},
			"fakeURL",
		)

		responseBody, err := client.Get("path/to/endpoint")

		assert.Equal(t, &nationalgrid.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "not found",
		}, err)
		assert.Nil(t, responseBody)
	})

	t.Run("should return an error if status code is 400", func(t *testing.T) {
		client := nationalgrid.NewClient(
			&mockedHttpClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("Bad request"))),
						StatusCode: http.StatusBadRequest,
					}, nil
				},
			},
			"fakeURL",
		)

		responseBody, err := client.Get("path/to/endpoint")

		assert.Equal(t, &nationalgrid.HttpError{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
		}, err)
		assert.Nil(t, responseBody)
	})

	t.Run("should return an error if status code is 500", func(t *testing.T) {
		client := nationalgrid.NewClient(
			&mockedHttpClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("Interval Server"))),
						StatusCode: http.StatusInternalServerError,
					}, nil
				},
			},
			"fakeURL",
		)
		responseBody, err := client.Get("path/to/endpoint")

		assert.Equal(t, &nationalgrid.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Interval Server",
		}, err)
		assert.Nil(t, responseBody)
	})

	t.Run("should return a response when the status code is 200", func(t *testing.T) {
		client := nationalgrid.NewClient(
			&mockedHttpClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("A valid response"))),
						StatusCode: http.StatusOK,
					}, nil
				},
			},
			"fakeURL",
		)
		mockedBodyToBytes, err := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte("A valid response"))))
		require.NoError(t, err)

		cr, err := client.Get("path/to/endpoint")

		assert.Equal(t, &nationalgrid.Response{
			StatusCode:   http.StatusOK,
			ResponseBody: mockedBodyToBytes,
		}, cr)
		assert.Nil(t, err)
	})
}
