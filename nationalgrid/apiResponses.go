package nationalgrid

import "fmt"

// IntensityResponse represents the intensity response from National Grid API
type IntensityResponse struct {
	Data []IntensityData
}

// IntensityData holds information about carbon intensity in the time frame between from and to
type IntensityData struct {
	From      string
	To        string
	Intensity Intensity
}

// Intensity holds information about the carbon intensity
type Intensity struct {
	Forecast int
	Actual   int
	Index    string
}

// Response is used to keep the status code and the response body from National Grid
type Response struct {
	StatusCode   int
	ResponseBody []byte
}

// HttpError is a custom error to hold information regarding the HTTTP errors
type HttpError struct {
	StatusCode int
	Message    string
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("Error code %d: %s", he.StatusCode, he.Message)
}
