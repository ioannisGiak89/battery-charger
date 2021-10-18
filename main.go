package main

import (
	"ioannisGiak89/arenko/assets"
	"ioannisGiak89/arenko/factory"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appFactory := factory.NewAppFactory()
	regulator := appFactory.CreateAssetCharger()
	b1 := appFactory.CreateBattery("battery-a")
	b2 := appFactory.CreateBattery("battery-b")
	signal := make(chan float64)
	batteries := []assets.Charger{b1, b2}
	var wg sync.WaitGroup

	wg.Add(2)
	go regulator.CheckCarbonIntensity(signal)
	go regulator.UpdateAssets(signal, batteries)

	// Keep the app running
	wg.Wait()
}
