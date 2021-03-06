package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MessageResponse recived from the target url containing response status, code and content
type MessageResponse struct {
	Status  string
	Code    int
	Content []byte
}

// Sender is a http request wrapper
type Sender interface {
	SendPOST(message []byte, url string) (MessageResponse, error)
}

// DefaultSender is the wrapper around go http libary used for sending http requests; It implements the `Sender` interface
type DefaultSender struct {
	httpClient http.Client
}

// SendPOST default implementation
func (ds *DefaultSender) SendPOST(message []byte, url string) (retVal MessageResponse, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return
	}

	resp, err := ds.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	retVal = MessageResponse{
		Status:  resp.Status,
		Code:    resp.StatusCode,
		Content: body,
	}

	return
}
