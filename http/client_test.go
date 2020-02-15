package http

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatingNewHTTPClient(t *testing.T) {
	t.Run("successful creation of new http client, default sender", func(t *testing.T) {
		config := ClientConfig{
			URL:   "https://test.com",
			Limit: 10,
		}

		_, err := NewHTTPClient(config)
		assert.NoError(t, err)
	})

	t.Run("successful creation of new http client, custom sender", func(t *testing.T) {
		config := ClientConfig{
			URL:    "https://test.com",
			Limit:  10,
			Sender: &ValidSenderMock{},
		}

		_, err := NewHTTPClient(config)
		assert.NoError(t, err)
		assert.IsType(t, &ValidSenderMock{}, config.Sender)
	})

	t.Run("failed creating http client, limit = 0", func(t *testing.T) {
		config := ClientConfig{
			URL:   "https://test.com",
			Limit: 0,
		}

		_, err := NewHTTPClient(config)
		assert.Error(t, err)
	})
}

func TestSendMessage(t *testing.T) {
	t.Run("one message sent succesfully", func(t *testing.T) {
		config := ClientConfig{
			URL:    "https://test.com",
			Limit:  10,
			Sender: &ValidSenderMock{},
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		counter := 0
		client, _ := NewHTTPClient(config)
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			if err == nil {
				counter++
			}
		})

		wg.Wait()
		assert.Equal(t, 1, counter)
	})

	t.Run("three messages sent succesfully, with limit 1", func(t *testing.T) {
		config := ClientConfig{
			URL:    "https://test.com",
			Limit:  1,
			Sender: &ValidSenderMock{},
		}

		wg := sync.WaitGroup{}
		wg.Add(3)
		counter := 0
		client, _ := NewHTTPClient(config)
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			if err == nil {
				counter++
			}
		})
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			if err == nil {
				counter++
			}
		})
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			if err == nil {
				counter++
			}
		})

		wg.Wait()
		assert.Equal(t, 3, counter)
	})

	t.Run("failed sending one message", func(t *testing.T) {
		config := ClientConfig{
			URL:    "https://test.com",
			Limit:  10,
			Sender: &InvalidSenderMock{},
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		client, _ := NewHTTPClient(config)
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			assert.Error(t, err)
		})

		wg.Wait()
	})

	t.Run("one message sent, bad request recived", func(t *testing.T) {
		config := ClientConfig{
			URL:    "https://test.com",
			Limit:  10,
			Sender: &BadRequestSenderMock{},
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		client, _ := NewHTTPClient(config)
		client.SendMessage([]byte("test message"), func(content []byte, err error) {
			defer wg.Done()
			assert.Error(t, err)
		})

		wg.Wait()
	})
}

type ValidSenderMock struct{}

func (vs *ValidSenderMock) SendPOST(message []byte, url string) (retVal MessageResponse, err error) {
	return MessageResponse{Code: StatusOK}, nil
}

type InvalidSenderMock struct{}

func (ivs *InvalidSenderMock) SendPOST(message []byte, url string) (retVal MessageResponse, err error) {
	return MessageResponse{}, fmt.Errorf("test error")
}

type BadRequestSenderMock struct{}

func (ivs *BadRequestSenderMock) SendPOST(message []byte, url string) (retVal MessageResponse, err error) {
	return MessageResponse{Code: StatusBadRequest}, nil
}
