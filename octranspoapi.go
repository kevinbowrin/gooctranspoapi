package octranspoapi

import (
	"golang.org/x/time/rate"
)

const ApiURLPrefix = "https://api.octranspo1.com/v1.2/"

type Connection struct {
	id      string
	key     string
	limiter *rate.Limiter
}

func NewConnection(id, key string, options ...func(*Connection) error) (*Connection, error) {
	c := &Connection{
		id:      id,
		key:     key,
		limiter: rate.NewLimiter(rate.Inf, 0),
	}
	for _, opt := range options {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func RateLimit(perSecond float64) func(*Connection) error {
	return func(c *Connection) error {
		c.limiter = rate.NewLimiter(rate.Limit(perSecond), 1)
		return nil
	}
}
