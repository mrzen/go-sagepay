Sagepay Go
==========

An API client library for Sagepay and Go.


Features
--------

* Overridable HTTP Client
* Context support for tracibility
* Minimal Dependencies
* Pluggable credential sources


Examples
--------

1. Get a merchant session key

````go
import "github.com/mrzen/sagepay"

func GetSessionKey(ctx context.Context) (*sagepay.SessionKey, error) {
    // Gets credentials from `SAGE_USERNAME` and `SAGE_PASSWORD` env
    sage := sagepay.New(ctx, sagepay.EnvironmentCredentialsProvider{})
    return sage.GetSessionKey(ctx)
}
````