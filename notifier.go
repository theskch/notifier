package notifier

// Callback function for handeling responses
type Callback func([]byte, error)

// Client interface used for notification
type Client interface {
	// SendMessage sends a message and recieves a response trough callback
	SendMessage(message []byte, call Callback)
}
