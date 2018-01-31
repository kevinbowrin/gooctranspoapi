package gooctranspoapi

import (
	"net/url"
)

const ApiURLPrefix = "https://api.octranspo1.com/v1.2/"

type Connection struct {
	id        string
	key       string
	rateLimit int
}

func RateLimit(perSecond int) func(*Connection) error {
	return func(c *Connection) error {
		c.rateLimit = perSecond
		return nil
	}
}

func Setup(id, key string, options ...func(*Connection) error) (*Connection, error) {
	c := &Connection{
		id:  id,
		key: key,
	}
	for _, opt := range options {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Connection) setupQuery() url.Values {
	query := url.Values{}
	query.Set("appID", c.id)
	query.Set("apiKey", c.key)
	query.Set("format", "json")
	return query
}
