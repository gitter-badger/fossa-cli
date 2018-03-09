package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("")

// API provides a default implementation of APIService.
type API struct {
	client      http.Client
	server      *url.URL
	apiKey      string
	initialized bool
}

// NewAPI constructs an API instance.
func NewAPI(server, APIKey string) API {
	api := API{
		client: http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		},
	}
	api.Initialize(server, APIKey)
	return api
}

// Initialize sets up an API instance.
func (api *API) Initialize(server, APIKey string) error {
	u, err := url.Parse(server)
	if err != nil {
		return err
	}
	api.server = u
	api.apiKey = APIKey
	api.initialized = true
	return nil
}

// MustInitialize crashes and logs a fatal exception if `Initialize` fails.
func (api *API) MustInitialize(server, APIKey string) {
	err := api.Initialize(server, APIKey)
	if err != nil {
		logger.Fatalf("Could not initialize API: %s", err.Error())
	}
}

// Get is a convenience method for MakeAPIRequest.
func (api *API) Get(URL string, body []byte) (res string, statusCode int, err error) {
	return api.stringAPIRequest(api.server, http.MethodGet, URL, api.apiKey, body)
}

// Post is a convenience method for MakeAPIRequest.
func (api *API) Post(URL string, body []byte) (res string, statusCode int, err error) {
	return api.stringAPIRequest(api.server, http.MethodPost, URL, api.apiKey, body)
}

// GetJSON is a convenience method for MakeAPIRequest.
func (api *API) GetJSON(URL string, body []byte, v interface{}) (statusCode int, err error) {
	return api.jsonAPIRequest(api.server, http.MethodGet, URL, api.apiKey, body, v)
}

// PostJSON is a convenience method for MakeAPIRequest.
func (api *API) PostJSON(URL string, body []byte, v interface{}) (statusCode int, err error) {
	return api.jsonAPIRequest(api.server, http.MethodPost, URL, api.apiKey, body, v)
}

func (api *API) stringAPIRequest(server *url.URL, method, URL, APIKey string, body []byte) (string, int, error) {
	endpoint, err := server.Parse(URL)
	if err != nil {
		return "", 0, err
	}
	res, code, err := api.MakeAPIRequest(endpoint, method, APIKey, body)
	if err != nil {
		return "", code, err
	}
	return string(res), code, nil
}

func (api *API) jsonAPIRequest(server *url.URL, method, URL, APIKey string, body []byte, v interface{}) (int, error) {
	endpoint, err := server.Parse(URL)
	if err != nil {
		return 0, err
	}
	res, code, err := api.MakeAPIRequest(endpoint, method, APIKey, body)
	if err != nil {
		return code, err
	}
	jsonErr := json.Unmarshal(res, v)
	if jsonErr != nil {
		return code, fmt.Errorf("could not unmarshal JSON API response: %s", jsonErr.Error())
	}
	return code, nil
}

func isTimeout(err error) bool {
	switch e := err.(type) {
	case net.Error:
		return e.Timeout()
	case *url.Error:
		return e.Err == io.EOF
	}
	return false
}

// MakeAPIRequest runs and logs a request backed by an `http.Client`.
func (api *API) MakeAPIRequest(endpoint *url.URL, method, APIKey string, body []byte) (res []byte, statusCode int, err error) {
	logger.Debugf("Making API request: %#v %#v %#v %#v %#v", *api, *endpoint, method, APIKey, body)

	// Check preconditions.
	if !api.initialized {
		return nil, 0, errors.New("API is not initialized")
	}

	// Construct request.
	req, err := http.NewRequest(method, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, 0, fmt.Errorf("could not construct API HTTP request: %s", err.Error())
	}
	req.Close = true
	req.Header.Set("Authorization", "token "+APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request.
	response, err := api.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, 0, ErrTimeout(fmt.Errorf("API request timed out: %s", err.Error()))
		}
		return nil, 0, fmt.Errorf("could not send API HTTP request: %s", err.Error())
	}
	defer response.Body.Close()

	// Read request.
	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("could not read API HTTP response: %s", err.Error())
	}

	logger.Debugf("Got API response: %#v", string(body))
	return res, response.StatusCode, nil
}
