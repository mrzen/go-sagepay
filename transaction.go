package sagepay

import (
	"context"
	"net/http"
)

// TransactionType represents a valid transaction type
type TransactionType string

// ThreeDSMode is an enum for valid modes for 3DS
type ThreeDSMode string

const (
	// ThreeDSModeDefault uses the account default setting for 3DS
	ThreeDSModeDefault ThreeDSMode = "UseMSPSetting"

	// ThreeDSModeForce forces the use of 3DS
	ThreeDSModeForce ThreeDSMode = "Force"

	// ThreeDSModeDisable forcfully disables the use of 3DS
	ThreeDSModeDisable ThreeDSMode = "Disable"

	// ThreeDSModeIgnoreRules enables 3DS but disables FP rules.
	ThreeDSModeIgnoreRules ThreeDSMode = "ForceIgnoringRules"

	// TransactionTypePayment instigates a one-time payment
	TransactionTypePayment TransactionType = "Payment"

	// TransactionTypeRepeat instigates a repeat payment
	TransactionTypeRepeat TransactionType = "Repeat"

	// TransactionTypeRefund instigates a refund
	TransactionTypeRefund TransactionType = "Refund"
)

// TransactionRequest represents the request data for creating a transaction
type TransactionRequest struct {
	Type          TransactionType      `json:"transactionType"`
	PaymentMethod RequestPaymentMethod `json:"paymentMethod"`
	Amount        int64                `json:"amount"`
	Currency      string               `json:"currency"`
	Description   string               `json:"description"`
	Reference     string               `json:"vendorTxCode"`

	EntryMethod       string      `json:"entryMethod,omitempty"`
	ApplyThreeDSecure ThreeDSMode `json:"apply3dSecure"`

	CustomerFirstName string `json:"customerFirstName"`
	CustomerLastName  string `json:"customerLastName"`
	CustomerEmail     string `json:"customerEmail,omitempty"`
	CustomerPhone     string `json:"customerPhone,omitempty"`

	BillingAddress `json:"billingAddress"`
}

// RequestPaymentMethod represents a payment method
type RequestPaymentMethod struct {
	Card struct {
		SessionKey string `json:"merchantSessionKey"`
		Identifier string `json:"cardIdentifier"`
		Reusable   bool   `json:"reusable"`
		Save       bool   `json:"save"`
	} `json:"card"`
}

// BillingAddress is the billing address
type BillingAddress struct {
	Line1      string `json:"address1"`
	Line2      string `json:"address2,omitempty"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type CreditCard struct {
	PANSuffix string `json:"lastFourDigits"`
	Type      string `json:"cardType"`
}

// TransactionResponse is the response data for creating a transaction
type TransactionResponse struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusDetail"`

	ID   string          `json:"transactionId"`
	Type TransactionType `json:"transactionType"`

	ResponseCode string `json:"bankResponseCode"`
	AuthCode     string `json:"bankAuthCode"`

	Status   string `json:"status"`
	Currency string `json:"currency"`
	Method   struct {
		Card CreditCard `json:"card"`
	} `json:"paymentMethod"`

	ThreeDSecure *struct {
		Status string `json:"status"`
	} `json:"3DSecure"`

	Amount struct {
		Total         int64 `json:"totalAmount"`
		Sale          int64 `json:"saleAmount"`
		ServiceCharge int64 `json:"surchargeAmount"`
	}

	// 3DS Data
	AcsURL string `json:"acsUrl,omitempty"`
	PAReq  string `json:"paReq,omitempty"`
}

// CreateTransaction will create a transaction with the given TransactionRequest
func (c Client) CreateTransaction(ctx context.Context, transaction *TransactionRequest) (*TransactionResponse, error) {

	path := "/transactions"

	res := TransactionResponse{}

	if err := c.JSON(ctx, http.MethodPost, path, transaction, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetTransaction gets a single transaction by ID
func (c *Client) GetTransaction(ctx context.Context, transactionID string) (*TransactionResponse, error) {
	path := "/transactions/" + transactionID

	res := TransactionResponse{}

	if err := c.JSON(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
