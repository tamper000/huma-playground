package cache

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-go"
)

const ttl = time.Minute * 5

type Cache struct {
	client valkey.Client
}

func New() (*Cache, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		return nil, err
	}

	return &Cache{client: client}, nil
}

func (c *Cache) CreateLink(ctx context.Context, url, ID string) error {
	cmd := c.client.B().Set().Key(ID).Value(url).Nx().Ex(ttl).Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return IDExists
		}

		return err
	}

	return nil
}

func (c *Cache) GetLink(ctx context.Context, ID string) (string, error) {
	cmd := c.client.B().Get().Key(ID).Build()

	resp := c.client.Do(ctx, cmd)
	err := resp.Error()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", IDNotExists
		}

		return "", err
	}

	return resp.ToString()
}
