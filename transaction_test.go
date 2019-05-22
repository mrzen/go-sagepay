package sagepay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestClientCreateTransactionPayment(t *testing.T) {

	client := getTestClient()

	sk, err := client.GetSessionKey(context.TODO(), "sandboxEC")

	if err != nil {
		t.Error(err)
	}

	tok, err := createTestCardToken(sk.Key)

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	rand.Seed(time.Now().Unix())
	transactionReference := fmt.Sprintf("transaction-%08x", rand.Int63())

	tx := &TransactionRequest{
		Type:              TransactionTypePayment,
		PaymentMethod:     RequestPaymentMethod{},
		CustomerFirstName: "John",
		CustomerLastName:  "Smith",
		ApplyThreeDSecure: ThreeDSModeDefault,
		Reference:         transactionReference,
		Description:       "TEST PAYMENT",
		Amount:            1000,
		Currency:          "GBP",
		BillingAddress: BillingAddress{
			Line1:      "88",
			PostalCode: "412",
			City:       "Penistone",
			Country:    "GB",
		},
	}

	tx.PaymentMethod.Card.SessionKey = sk.Key
	tx.PaymentMethod.Card.Identifier = tok
	tx.PaymentMethod.Card.Save = false

	tr, err := client.CreateTransaction(context.TODO(), tx)

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log("Transaction ID: ", tr.ID)
	t.Log("Transaction Reference: ", transactionReference)

	t.Logf("%+v", tr)

}

// Function to make test card tokens
func createTestCardToken(msk string) (string, error) {

	card := map[string]interface{}{
		"cardDetails": map[string]string{
			"cardholderName": "JOHN SMITH",
			"cardNumber":     "4929000000006",
			"expiryDate":     "1224",
			"securityCode":   "123",
		},
	}

	body, _ := json.Marshal(card)

	fmt.Println(string(body))

	req, _ := http.NewRequest(
		http.MethodPost,
		TestHost+"/card-identifiers",
		bytes.NewReader(body),
	)

	req.Header.Set("Authorization", "Bearer "+msk)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sagepay.go TEST LIB")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {

		a, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return "", errors.New(string(a))
	}

	rp := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&rp); err != nil {
		return "", err
	}

	return rp["cardIdentifier"].(string), nil
}
