package constants

type LogAction string

const (
	ActionMiddlewareStart LogAction = "MiddlewareStart"
	ActionMiddlewareEnd   LogAction = "MiddlewareEnd"
	ActionMiddlewareError LogAction = "MiddlewareError"

	// HTTP Request Actions
	ActionRequestReceived  LogAction = "RequestReceived"
	ActionRequestCompleted LogAction = "RequestCompleted"
	ActionRequestError     LogAction = "RequestError"

	// Authentication Actions
	ActionAuthAttempt LogAction = "AuthAttempt"
	ActionAuthSuccess LogAction = "AuthSuccess"
	ActionAuthFailed  LogAction = "AuthFailed"

	ActionUserCreated LogAction = "UserCreated"
	ActionUserUpdated LogAction = "UserUpdated"

	// System Actions
	ActionDatabaseQuery   LogAction = "DatabaseQuery"
	ActionCacheOperation  LogAction = "CacheOperation"
	ActionExternalAPICall LogAction = "ExternalAPICall"
)
