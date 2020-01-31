package main

import (
	"encoding/json"
	"testing"
)

func TestParseArtifactMetadata(t *testing.T) {
	bytes := []byte(`{
  "offers": [
    {
      "api": "https://api.licensezero.com",
      "offerID": "36fce1e2-5e96-41fc-8776-4e632b546d96",
			"public": "Parity-7.0.0"
    }
  ]
}`)
	var unstructured interface{}
	err := json.Unmarshal(bytes, &unstructured)
	if err != nil {
		t.Error(err)
	}
	metadata, err := ParseArtifactMetadata(unstructured)
	if err != nil {
		t.Error(err)
	}
	if len(metadata.Offers) != 1 {
		t.Error("failed to parse one offer")
	}
	first := metadata.Offers[0]
	if first.OfferID != "36fce1e2-5e96-41fc-8776-4e632b546d96" {
		t.Error("failed to parse offerID")
	}
	if first.API != "https://api.licensezero.com" {
		t.Error("failed to parse API")
	}
	if first.Public != "Parity-7.0.0" {
		t.Error("failed to parse public license identifier")
	}
}
