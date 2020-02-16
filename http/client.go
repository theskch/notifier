// Package http provides implementation of the notifier client that uses HTTP to sent messages
package http

import (
	"fmt"
	"net/http"

	"github.com/theskch/notifier"
)

// ClientConfig used to configure http client notifier
type ClientConfig struct {
	URL   string
	Limit int
	//if left null, default sender will be used
	Sender Sender
}

// NewHTTPClient returnes a http imepelemntation of notifier client.
// url is the address to which the message is sent.
// limit is the number of simultaneous messages that can be sent. Setting the limit to 1 will send sync messages. Limit must be greater than 0.
func NewHTTPClient(config ClientConfig) (notifier.Client, error) {
	if config.Limit <= 0 {
		return nil, fmt.Errorf("invalid client configuration, limit must be greater than 0")
	}

	retVal := client{
		url: config.URL,
		sem: make(chan int, config.Limit),
	}

	if config.Sender != nil {
		retVal.sender = config.Sender
	} else {
		retVal.sender = &DefaultSender{
			httpClient: http.Client{},
		}
	}

	return &retVal, nil
}

// client an http impelemntation of NotifierClient
type client struct {
	url    string
	sem    chan int // channel used for semaphore implementation
	sender Sender
}

// SendMessage http client implementation
func (c *client) SendMessage(message []byte, call notifier.Callback) {
	go func() {
		c.sem <- 1 // blocks the message brodcast if the limit has been reached
		c.sendPOSTRequest(message, call)
	}()
}

func (c *client) sendPOSTRequest(message []byte, call notifier.Callback) {
	defer func() {
		<-c.sem
	}()

	var err error
	var retVal []byte

	response, err := c.sender.SendPOST(message, c.url)
	if err != nil {
		call([]byte{}, err)
		return
	}

	if err == nil {
		// If response code is one of these 3, set the content for the call backe; else set the error and ingore the content
		if response.Code == StatusOK || response.Code == StatusCreated || response.Code == StatusAccepted {
			retVal = response.Content
		} else {
			err = fmt.Errorf("failed sending notification: response code: %d  response status: %s", response.Code, response.Status)
		}
	}

	if call != nil {
		call(retVal, err)
	}
}
