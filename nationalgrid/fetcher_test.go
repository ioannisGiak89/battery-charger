package nationalgrid_test

import (
	"encoding/json"
	"ioannisGiak89/arenko/nationalgrid"
	"net/url"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Implements client interface. This struct is used to mock the FsaRestClient
type client struct {
	baseUrl *url.URL
	MockGet func(path string) (*nationalgrid.Response, error)
}

func (m *client) Get(path string) (*nationalgrid.Response, error) {
	return m.MockGet(path)
}

func TestIntensityFetcher_FetchCurrentCarbonIntensity(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)

	t.Run("should return an IntensityData object", func(t *testing.T) {
		expectedResponse := getFakeIntensityResponse()
		jsonResponse, err := json.Marshal(expectedResponse)
		require.NoError(t, err)

		mockedResponse := &nationalgrid.Response{
			StatusCode:   200,
			ResponseBody: jsonResponse,
		}

		fetcher := nationalgrid.NewIntensityFetcher(&client{
			MockGet: func(path string) (*nationalgrid.Response, error) {
				return mockedResponse, nil
			},
		})

		response, err := fetcher.FetchCurrentCarbonIntensity()
		assert.Nil(t, err)
		assert.Equal(t, &expectedResponse.Data[0], response)
	})

	t.Run("should return an error if there is no intensity data", func(t *testing.T) {
		expectedResponse := &nationalgrid.IntensityResponse{
			Data: []nationalgrid.IntensityData{},
		}

		jsonResponse, err := json.Marshal(expectedResponse)
		require.NoError(t, err)

		mockedResponse := &nationalgrid.Response{
			StatusCode:   200,
			ResponseBody: jsonResponse,
		}

		fetcher := nationalgrid.NewIntensityFetcher(&client{
			MockGet: func(path string) (*nationalgrid.Response, error) {
				return mockedResponse, nil
			},
		})

		response, err := fetcher.FetchCurrentCarbonIntensity()
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})

	t.Run("should return an error if there intensity data is nil", func(t *testing.T) {
		expectedResponse := &nationalgrid.IntensityResponse{}

		jsonResponse, err := json.Marshal(expectedResponse)
		require.NoError(t, err)

		mockedResponse := &nationalgrid.Response{
			StatusCode:   200,
			ResponseBody: jsonResponse,
		}

		fetcher := nationalgrid.NewIntensityFetcher(&client{
			MockGet: func(path string) (*nationalgrid.Response, error) {
				return mockedResponse, nil
			},
		})

		response, err := fetcher.FetchCurrentCarbonIntensity()
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})

	t.Run("should return an error if the client fails", func(t *testing.T) {
		fetcher := nationalgrid.NewIntensityFetcher(&client{
			MockGet: func(path string) (*nationalgrid.Response, error) {
				return nil, &nationalgrid.HttpError{
					StatusCode: 404,
					Message:    "Not found",
				}
			},
		})
		response, err := fetcher.FetchCurrentCarbonIntensity()
		assert.Nil(t, response)
		assert.Equal(t, err, &nationalgrid.HttpError{
			StatusCode: 404,
			Message:    "Not found",
		})
	})

	t.Run("should return an error if unmarshal fails", func(t *testing.T) {
		fetcher := nationalgrid.NewIntensityFetcher(&client{
			MockGet: func(path string) (*nationalgrid.Response, error) {
				return &nationalgrid.Response{
					StatusCode:   200,
					ResponseBody: []byte{12, 12},
				}, nil
			},
		})
		response, err := fetcher.FetchCurrentCarbonIntensity()
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})
}

func getFakeIntensityResponse() *nationalgrid.IntensityResponse {
	return &nationalgrid.IntensityResponse{
		Data: []nationalgrid.IntensityData{nationalgrid.IntensityData{
			From: "2018-01-20T12:00Z",
			To:   "2018-01-20T12:30Z",
			Intensity: nationalgrid.Intensity{
				Forecast: 266,
				Actual:   263,
				Index:    "low",
			},
		}},
	}
}
