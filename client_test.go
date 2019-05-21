package sagepay

import (
	"context"
	"net/http"
	"testing"
)

var demoCredentials = StaticCredentials(
	"dq9w6WkkdD2y8k3t4olqu8H6a0vtt3IY7VEsGhAtacbCZ2b5Ud",
	"hno3JTEwDHy7hJckU4WuxfeTrjD0N92pIaituQBw5Mtj7RG3V8zOdHCSPKwJ02wAV",
)

func getTestClient() *Client {
	c := New(context.TODO(), demoCredentials)

	c.SetTestMode(true)
	//c.DebugWriter = os.Stdout

	return c
}

func TestClientDo(t *testing.T) {
	c := getTestClient()

	req, _ := http.NewRequest(http.MethodGet, TestHost, nil)

	_, err := c.Do(context.TODO(), req)

	if err != nil {
		t.Error(err)
	}

	//res.Write(os.Stdout)
}
