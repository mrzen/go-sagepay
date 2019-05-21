package sagepay

import (
	"context"
	"net/http"
	"time"
)

// SessionKey represents a session key
type SessionKey struct {
	Key    string    `json:"merchantSessionKey"`
	Expiry time.Time `json:"expiry"`
}

// GetSessionKey gets a new merchant session Key
func (c Client) GetSessionKey(ctx context.Context, vendorName string) (*SessionKey, error) {
	path := "/merchant-session-keys"

	body := map[string]string{
		"vendorName": vendorName,
	}

	res := SessionKey{}

	if err := c.JSON(ctx, http.MethodPost, path, body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
