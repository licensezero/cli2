package abstract

import (
	"errors"
	"github.com/mitchellh/mapstructure"
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
          "$ref": "offer.json"
        },
        "site": {
          "$ref": "offer.json"
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

func ParseOffer(unstructured interface{}) (Offer, error) {
	if validV1Offer(unstructured) {
		var v1 offer1_0_0Pre
		err := mapstructure.Decode(unstructured, &v1)
		if err != nil {
			return nil, err
		}
		return v1, nil
	}
	return nil, errors.New("unknown schema")
}

func validV1Offer(unstructured interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := v1OfferSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
