// Package notifier contains the interface for notification client. All future  notifiers should implement `Client` interface.
// `Client` inteerfaces is intendend to be used in the users code.
package notifier

// Callback function for handeling responses
type Callback func([]byte, error)

// Client interface used for sending notifications
type Client interface {
	// SendMessage sends a message and recieves a response trough callback
	SendMessage(message []byte, call Callback)
}
