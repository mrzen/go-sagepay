package sagepay

import (
	"context"
	"net/http"
)

// ThreeDSResult represents the result of a 3DS Challenge
type ThreeDSResult struct {
	Status string `json:"status"`
}

// RespondThreeDS responds to a 3DS challenge
func (c *Client) RespondThreeDS(ctx context.Context, transactionID, payload string) (*ThreeDSResult, error) {

	path := "/transactions/" + transactionID + "/3d-secure"

	body := map[string]string{
		"paRes": payload,
	}

	res := ThreeDSResult{}

	if err := c.JSON(ctx, http.MethodPost, path, body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
