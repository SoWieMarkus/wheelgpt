package identity

func (c *Client) Validate(token string) (bool, error) {
	return true, nil
}

func (c *Client) Token() (string, error) {
	return "", nil
}
