# Notifier
A simple notifier used for sending messages to the service. Currently only HTTP notifier is implemented. 

# Usage
Notifier should be used trough `Client` interface in `notifier` package.

## HTTP notifier
`HTTPNotifier` is the implementation of the `Client` interfaces that uses HTTP POST requests to send messages. Number of simultaneous requests is limited by configuration. Response from the server is send back to the caller using the `callback` function.