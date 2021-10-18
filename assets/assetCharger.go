package assets

import (
	"fmt"
	"ioannisGiak89/arenko/nationalgrid"
	"log"
	"os"
	"strconv"
	"time"
)

// ChargeRegulator defines the ChargeRegulator interface
type ChargeRegulator interface {
	// CheckCarbonIntensity checks the current carbon intensity and sends a charge/discharge signal
	CheckCarbonIntensity(signal chan<- float64)
	// Sends the charging signal to the assets
	UpdateAssets(signal <-chan float64, assets []Charger)
}

// AssetCharger implements the ChargeRegulator interface
type AssetCharger struct {
	carbonIntensityFetcher nationalgrid.CarbonIntensityFetcher
}

// NewAssetCharger creates an AssetCharger
func NewAssetCharger(ci nationalgrid.CarbonIntensityFetcher) *AssetCharger {
	return &AssetCharger{carbonIntensityFetcher: ci}
}

func (ac *AssetCharger) CheckCarbonIntensity(signal chan<- float64) {
	for {
		intensity, err := ac.carbonIntensityFetcher.FetchCurrentCarbonIntensity()
		if err != nil {
			// Don't want to interrupt the service
			// Fow now we just log and continue
			log.Println(err)
			ac.waitForInterval()
			continue
		}

		switch i := intensity.Intensity.Index; i {
		case "high", "very high":
			signal <- -1.0
		case "low", "very low":
			signal <- 1.0
		default:
			// Means the intensity index is moderate so we wait and then continue
			fmt.Printf("Carbon Intensity is %s. Let's wait.. \n", i)
			ac.waitForInterval()
			continue
		}

		ac.waitForInterval()
	}
}

func (ac *AssetCharger) UpdateAssets(signal <-chan float64, assets []Charger) {
	for {
		chargeSignal := <-signal
		for _, a := range assets {
			err := a.SetCharging(chargeSignal)
			if err != nil {
				// Don't want to interrupt the service
				// Fow now we just log and continue
				log.Println(err)
				continue
			}
		}
	}
}

func (ac *AssetCharger) waitForInterval() {
	seconds := os.Getenv("INTENSITY_CHECK_INTERVAL_SECONDS")
	if seconds == "" {
		log.Fatal("INTENSITY_CHECK_INTERVAL_SECONDS must be set")
	}

	s, err := strconv.Atoi(seconds)
	if err != nil {
		fmt.Println("INTENSITY_CHECK_INTERVAL_SECONDS must be a number")
		log.Fatal(err)
	}

	time.Sleep(time.Duration(s) * time.Second)
}
