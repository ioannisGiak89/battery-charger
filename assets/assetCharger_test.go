package assets_test

import (
	"ioannisGiak89/arenko/assets"
	"ioannisGiak89/arenko/nationalgrid"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type mockFetcher struct {
	MockFetchCurrentCarbonIntensity func() (*nationalgrid.IntensityData, error)
}

func (mf *mockFetcher) FetchCurrentCarbonIntensity() (*nationalgrid.IntensityData, error) {
	return mf.MockFetchCurrentCarbonIntensity()
}

func TestAssetCharger_CheckCarbonIntensity(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)

	t.Run("should write -1.0 to channel if intensity is high", func(t *testing.T) {
		assetCharger := assets.NewAssetCharger(
			&mockFetcher{
				MockFetchCurrentCarbonIntensity: func() (*nationalgrid.IntensityData, error) {
					return &nationalgrid.IntensityData{
						From: "2018-01-20T12:00Z",
						To:   "2018-01-20T12:30Z",
						Intensity: nationalgrid.Intensity{
							Forecast: 266,
							Actual:   100,
							Index:    "high",
						},
					}, nil
				},
			},
		)
		signal := make(chan float64, 1)
		go assetCharger.CheckCarbonIntensity(signal)
		assert.Equal(t, -1.0, <-signal)
	})

	t.Run("should write -1.0 to channel if intensity is very high", func(t *testing.T) {
		assetCharger := assets.NewAssetCharger(
			&mockFetcher{
				MockFetchCurrentCarbonIntensity: func() (*nationalgrid.IntensityData, error) {
					return &nationalgrid.IntensityData{
						From: "2018-01-20T12:00Z",
						To:   "2018-01-20T12:30Z",
						Intensity: nationalgrid.Intensity{
							Forecast: 266,
							Actual:   100,
							Index:    "very high",
						},
					}, nil
				},
			},
		)
		signal := make(chan float64, 1)
		go assetCharger.CheckCarbonIntensity(signal)
		assert.Equal(t, -1.0, <-signal)
	})

	t.Run("should write 1.0 to channel if intensity is low", func(t *testing.T) {
		assetCharger := assets.NewAssetCharger(
			&mockFetcher{
				MockFetchCurrentCarbonIntensity: func() (*nationalgrid.IntensityData, error) {
					return &nationalgrid.IntensityData{
						From: "2018-01-20T12:00Z",
						To:   "2018-01-20T12:30Z",
						Intensity: nationalgrid.Intensity{
							Forecast: 266,
							Actual:   100,
							Index:    "low",
						},
					}, nil
				},
			},
		)
		signal := make(chan float64, 1)
		go assetCharger.CheckCarbonIntensity(signal)
		assert.Equal(t, 1.0, <-signal)
	})

	t.Run("should write 1.0 to channel if intensity is very low", func(t *testing.T) {
		assetCharger := assets.NewAssetCharger(
			&mockFetcher{
				MockFetchCurrentCarbonIntensity: func() (*nationalgrid.IntensityData, error) {
					return &nationalgrid.IntensityData{
						From: "2018-01-20T12:00Z",
						To:   "2018-01-20T12:30Z",
						Intensity: nationalgrid.Intensity{
							Forecast: 266,
							Actual:   100,
							Index:    "very low",
						},
					}, nil
				},
			},
		)
		signal := make(chan float64, 1)
		go assetCharger.CheckCarbonIntensity(signal)
		assert.Equal(t, 1.0, <-signal)
	})
}
