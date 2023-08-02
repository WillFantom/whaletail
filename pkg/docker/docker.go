package docker

import "context"

func (c *Client) Online() bool {
	client, err := c.local()
	if err != nil {
		return false
	}
	_, err = client.Info(context.Background())
	return (err == nil)
}
