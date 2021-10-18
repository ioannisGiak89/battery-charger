package factory

import (
	"ioannisGiak89/arenko/assets"
	"ioannisGiak89/arenko/nationalgrid"
	"log"
	"net/http"
	"os"
	"time"
)

// Factory abstracts the creation of instances.
type Factory interface {
	// CreateClient creates a National Grid Client
	CreateClient() nationalgrid.NGClient
	// CreateAssetCharger creates a ChargeRegulator
	CreateAssetCharger() assets.ChargeRegulator
	// CreateIntensityFetcher creates a CarbonIntensityFetcher
	CreateIntensityFetcher() nationalgrid.CarbonIntensityFetcher
	// CreateBattery creates a battery
	CreateBattery(ID string) *assets.Battery
}

// AppFactory creates new instances
type AppFactory struct{}

// NewAppFactory creates a Factory
func NewAppFactory() *AppFactory {
	return &AppFactory{}
}

func (f *AppFactory) CreateIntensityFetcher() nationalgrid.CarbonIntensityFetcher {
	return nationalgrid.NewIntensityFetcher(f.CreateClient())
}

func (f *AppFactory) CreateClient() nationalgrid.NGClient {
	baseURL := os.Getenv("NATIONAL_GRID_BASE_URL")
	if baseURL == "" {
		log.Fatal("base url is not set")
	}

	return nationalgrid.NewClient(f.CreateHttpClient(), baseURL)
}

func (f *AppFactory) CreateAssetCharger() assets.ChargeRegulator {
	return assets.NewAssetCharger(f.CreateIntensityFetcher())
}

func (f *AppFactory) CreateBattery(ID string) *assets.Battery {
	return assets.NewBattery(ID)
}

// CreateHttpClient creates an HTTP client
func (f *AppFactory) CreateHttpClient() nationalgrid.HTTPClient {
	return &http.Client{
		// The default client does not have a timeout set u, so
		// calls might hang forever if something goes wrong.
		// Setting a timeout is considered a good practice
		Timeout: time.Second * 15,
	}
}
