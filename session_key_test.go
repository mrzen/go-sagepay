package sagepay

import (
	"context"
	"testing"
)

func TestClientGetSessionKey(t *testing.T) {

	client := getTestClient()

	sk, err := client.GetSessionKey(context.TODO(), "sandboxEC")

	if err != nil {
		t.Error(err)
	}

	t.Log(sk)

}
