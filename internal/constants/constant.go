// Package constants provides all the constant values that are used in this service
package constants

const (
	ParseErr               = "Error parsing request"
	IncompleteData         = "Error :: Incomplete data"
	PublisherRegistration  = "Successfully Registered as publisher to the channel"
	SubscriberRegistration = "Successfully Registered as subscriber to the channel"
	InvalidUUID            = "Error :: Invalid Universal UID"
	PublisherNotFound      = "Error :: No publisher found with the specified name for specified channel"
	NotifiedSub            = "Notified All Subscriber"
	NoRoute                = "Error :: No Route Found"
	ConnectedServer        = "Successfully connected to the server"
	StringUnmarshalError   = "cannot unmarshal string into Go struct field"
	ValidatorFail          = "failed on the 'required' tag"
)
