package octranspoapi

import (
	"golang.org/x/time/rate"
)

// ApiURLPrefix is the address at which the API is available.
const ApiURLPrefix = "https://api.octranspo1.com/v1.2/"

// Connection holds the Application ID and API key needed to make requests.
// It also has a rate limiter, used by the Connection's methods to limit calls on the API.
type Connection struct {
	Id      string
	Key     string
	Limiter *rate.Limiter
}

// NewConnection returns a new connection without a rate limit.
func NewConnection(id, key string) Connection {
	return Connection{
		Id:      id,
		Key:     key,
		Limiter: rate.NewLimiter(rate.Inf, 0),
	}
}

// NewConnectionWithRateLimit returns a new connection with a rate limit set.
func NewConnectionWithRateLimit(id, key string, perSecond float64, burst int) Connection {
	return Connection{
		Id:      id,
		Key:     key,
		Limiter: rate.NewLimiter(rate.Limit(perSecond), burst),
	}
}
