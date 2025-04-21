package constant

const (
	TrackIDHeader = "track_id"
	Undefined     = "undefined"
)

type ctxRequestIDKey int

const (
	RequestIDKey ctxRequestIDKey = 0
)

// Common variable
const (
	ContextTimeout = 3 // seconds
	WarnTime       = 3000
)
