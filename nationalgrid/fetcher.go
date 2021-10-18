package nationalgrid

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// CarbonIntensityFetcher defines the carbon intensity fetcher interface
type CarbonIntensityFetcher interface {
	// FetchCurrentCarbonIntensity fetches the current carbon intensity
	FetchCurrentCarbonIntensity() (*IntensityData, error)
}

// IntensityFetcher fetches data from National Grid
type IntensityFetcher struct {
	nationalGridClient NGClient
}

// NewIntensityFetcher creates a new IntensityFetcher
func NewIntensityFetcher(nationalGridClient NGClient) *IntensityFetcher {
	return &IntensityFetcher{nationalGridClient: nationalGridClient}
}

func (f IntensityFetcher) FetchCurrentCarbonIntensity() (*IntensityData, error) {
	ie := os.Getenv("INTENSITY_ENDPOINT")
	if ie == "" {
		// Log and exit the app as all the requests will fail
		// We could handle it differently if we had more routines that do other tasks
		log.Fatal("carbon intensity endpoint is not set")
	}

	res, err := f.nationalGridClient.Get(ie)
	if err != nil {
		return nil, err
	}

	var intensityResponse IntensityResponse
	err = json.Unmarshal(res.ResponseBody, &intensityResponse)

	if err != nil {
		return nil, err
	}

	if intensityResponse.Data == nil || len(intensityResponse.Data) == 0 {
		return nil, errors.New("no intensity data")
	}

	return &intensityResponse.Data[0], nil
}
