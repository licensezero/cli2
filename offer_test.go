package main

import (
	"encoding/json"
	"testing"
)

func TestParseOffer(t *testing.T) {
	bytes := []byte(`{
	"url": "http://example.com",
	"licensorID": "d56ee0a6-4ed3-4793-9485-6135644c158f",
	"pricing": {
		"single": {
			"currency": "USD",
			"amount": 1000
		}
	}
}`)
	var unstructured interface{}
	err := json.Unmarshal(bytes, &unstructured)
	if err != nil {
		t.Error(err)
	}
	offer, err := ParseOffer(unstructured)
	if err != nil {
		t.Error(err)
	}
	if offer.URL() != "http://example.com" {
		t.Error("failed to parse URL")
	}
	if offer.LicensorID() != "d56ee0a6-4ed3-4793-9485-6135644c158f" {
		t.Error("failed to parse licensorID")
	}
}
