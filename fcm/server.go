package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/practice/bank_fcm/model"
)

var (
	// retryableErrors whether the error is a retryable
	retryableErrors = map[string]bool{
		"Unavailable":         true,
		"InternalServerError": true,
	}
)

//Server stores client with api key to firebase
type Server struct {
	APIKey     string
	HTTPClient *http.Client
}

//NewServer creates a new client
func NewServer(apiKey string) *Server {
	return NewServerWithClient(apiKey, &http.Client{})
}

// NewServerWithClient creates a new client
func NewServerWithClient(apiKey string, httpClient *http.Client) *Server {
	return &Server{
		APIKey:     apiKey,
		HTTPClient: httpClient,
	}
}

//AuthorizationToken provide authorization token
func (s *Server) AuthorizationToken() string {
	return fmt.Sprintf("key=%v", s.APIKey)
}

// Send message to FCM
func (s *Server) Send(message model.Message) (model.Response, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return model.Response{}, err
	}

	req, err := http.NewRequest(model.MethodPOST, model.FCMServerURL, bytes.NewBuffer(data))
	if err != nil {
		return model.Response{}, err
	}

	req.Header.Set("Authorization", s.AuthorizationToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return model.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return model.Response{}, fmt.Errorf("%d status code", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Response{}, err
	}

	response := model.Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	response.StatusCode = resp.StatusCode
	response.RetryAfter = resp.Header.Get(model.HeaderRetryAfter)
	if err := s.Failed(&response); err != nil {
		return response, err
	}
	response.Ok = true

	return response, nil
}

// Failed method indicates if the server couldn't process
// the request in time.
func (s *Server) Failed(response *model.Response) error {
	for _, response := range response.Results {
		if retryableErrors[response.Error] {
			return fmt.Errorf("Failed %s", response.Error)
		}
	}

	return nil
}
