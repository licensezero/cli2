package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

// Artifact encodes data about offers for an artifact.
type Artifact interface {
	Offers() []ArtifactOffer
}

// ArtifactOffer represents an offer relevant to an artifact.
type ArtifactOffer struct {
	API     string
	OfferID string
	Public  string
}

type artifact1_0_0Pre struct {
	OfferArray []artifactOffer1_0_0Pre `mapstructure:"offers"`
}

type artifactOffer1_0_0Pre struct {
	API     string
	OfferID string
	Public  string
}

func (a artifact1_0_0Pre) Offers() (offers []ArtifactOffer) {
	for _, item := range a.OfferArray {
		offers = append(offers, ArtifactOffer{
			API:     item.API,
			OfferID: item.OfferID,
			Public:  item.Public,
		})
	}
	return
}

const artifact1_0_0PreSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/artifact.json",
  "type": "object",
  "required": [
    "offers"
  ],
  "additionalProperties": false,
  "properties": {
    "offers": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [
          "offerID",
          "api"
        ],
        "additionalProperties": false,
        "properties": {
          "offerID": {
            "title": "UUIDv4 offer identifier",
            "type": "string",
            "format": "uuid"
          },
          "api": {
            "title": "licensing API",
            "type": "string",
            "format": "uri",
            "pattern": "^https://",
            "examples": [
              "https://api.licensezero.com"
            ]
          },
          "public": {
            "title": "public license identifier",
            "type": "string",
            "pattern": "^[A-Za-z0-9-.]+",
            "examples": [
              "Parity-7.0.0"
            ]
          }
        }
      }
    }
  }
}`

// ParseArtifact validates and parses parsed JSON data as a Artifact.
func ParseArtifact(unstructured interface{}) (a Artifact, err error) {
	if validV1Artifact(unstructured) {
		return parseV1Artifact(unstructured), nil
	}
	return a, errors.New("unknown schema")
}

var v1ArtifactSchema *gojsonschema.Schema = nil

func validV1Artifact(parsed interface{}) bool {
	if v1ArtifactSchema == nil {
		v1ArtifactSchema, _ = schemaLoader().Compile(
			gojsonschema.NewStringLoader(artifact1_0_0PreSchema),
		)
	}
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ArtifactSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

func parseV1Artifact(unstructured interface{}) (a artifact1_0_0Pre) {
	mapstructure.Decode(unstructured, &a)
	return
}
