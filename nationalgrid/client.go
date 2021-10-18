package nationalgrid

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// client defines the National Grid interface
type NGClient interface {
	// Get does a get request to an endpoint
	Get(endpoint string) (*Response, error)
}

// HTTPClient interface. This interface is implemented by http.Client and is used for mocking
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements the client interface
type Client struct {
	baseURL string
	client  HTTPClient
}

// Creates a new Client
func NewClient(httpClient HTTPClient, baseURL string) *Client {
	return &Client{
		client:  httpClient,
		baseURL: baseURL,
	}
}

func (cl *Client) Get(endpoint string) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(
		"%s%s",
		cl.baseURL,
		endpoint,
	), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := cl.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return &Response{
			StatusCode:   res.StatusCode,
			ResponseBody: resBody,
		}, nil
	}

	return nil, &HttpError{
		StatusCode: res.StatusCode,
		Message:    string(resBody),
	}
}
