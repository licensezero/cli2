package main

import (
	"encoding/json"
	"testing"
)

func TestParseReceipt(t *testing.T) {
	bytes := []byte(`{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      }
    },
    "form": "Test license form."
  }
}`)
	var unstructured interface{}
	err := json.Unmarshal(bytes, &unstructured)
	if err != nil {
		t.Error(err)
	}
	receipt, err := ParseReceipt(unstructured)
	if receipt.API() != "https://api.licensezero.com" {
		t.Error("failed to parse API")
	}
	if receipt.OrderID() != "2c743a84-09ce-4549-9f0d-19d8f53462bb" {
		t.Error("failed to parse orderID")
	}
	if receipt.OfferID() != "9aab7058-599a-43db-9449-5fc0971ecbfa" {
		t.Error("failed to parse offerID")
	}
}
