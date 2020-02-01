package main

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
)

// Offer represents an offer to sell licenses.
type Offer interface {
	URL() string
	LicensorID() string
	Pricing() Pricing
}

// Pricing represents a price list.
type Pricing struct {
	Single    Price
	Relicense Price
}

type offer1_0_0Pre struct {
	licensorID string  `mapstructure:"licensorID"`
	pricing    Pricing `mapstructure:"pricing"`
	url        string  `mapstructure:"url"`
}

func (o offer1_0_0Pre) URL() string {
	return o.url
}

func (o offer1_0_0Pre) LicensorID() string {
	return o.licensorID
}

func (o offer1_0_0Pre) Pricing() Pricing {
	return o.pricing
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
	return nil, errors.New("unknown schema")
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

func parseV1Offer(unstructured interface{}) (o offer1_0_0Pre) {
	asMap := unstructured.(map[string]interface{})
	o.url = asMap["url"].(string)
	o.licensorID = asMap["licensorID"].(string)
	pricing := asMap["pricing"].(map[string]interface{})
	single := pricing["single"].(map[string]interface{})
	o.pricing.Single.Currency = single["currency"].(string)
	o.pricing.Single.Amount = uint(single["amount"].(float64))
	return
}
