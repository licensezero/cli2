package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

// Offer represents an offer to sell licenses.
type Offer struct {
	URL        string  `mapstructure:"url"`
	LicensorID string  `mapstructure:"licensorID"`
	Pricing    Pricing `mapstructure:"pricing"`
}

// Pricing represents a price list.
type Pricing struct {
	Single    Price `mapstructure:"single"`
	Relicense Price `mapstructure:"relicense"`
}

const offer1_0_0PreSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/offer.json",
  "type": "object",
  "required": [
    "url",
    "licensorID",
    "pricing"
  ],
  "additionalProperties": true,
  "properties": {
    "url": {
      "type": "string",
      "format": "uri"
    },
    "licensorID": {
      "type": "string",
      "format": "uuid"
    },
    "pricing": {
      "type": "object",
      "properties": {
        "single": {
          "$ref": "price.json"
        },
        "site": {
          "$ref": "price.json"
        }
      },
      "patternProperties": {
        "^d+$": {
          "$ref": "price.json"
        }
      }
    }
  }
}`

// ParseOffer parses instructed offer data.
func ParseOffer(unstructured interface{}) (Offer, error) {
	if validV1Offer(unstructured) {
		return parseV1Offer(unstructured), nil
	}
	return Offer{}, errors.New("unknown schema")
}

var v1OfferSchema *gojsonschema.Schema = nil

func validV1Offer(unstructured interface{}) bool {
	if v1OfferSchema == nil {
		schema, err := schemaLoader().Compile(
			gojsonschema.NewStringLoader(offer1_0_0PreSchema),
		)
		if err != nil {
			panic(err)
		}
		v1OfferSchema = schema
	}
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := v1OfferSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

func parseV1Offer(unstructured interface{}) (o Offer) {
	mapstructure.Decode(unstructured, &o)
	return
}
