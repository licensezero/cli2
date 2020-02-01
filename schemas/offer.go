package schemas

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/cli2/abstract"
)

const offer1_0_0Pre = `{
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

func ParseOffer(unstructured interface{}) (*abstract.Offer, error) {
	if validV1Offer(unstructured) {
		return parseV1Offer(unstructured), nil
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

func parseV1Offer(unstructured interface{}) *abstract.Offer {
	object := unstructured.(map[string]interface{})
	pricing := object["pricing"].(map[string]interface{})
	single := pricing["single"].(map[string]interface{})
	return &abstract.Offer{
		URL:        object["url"].(string),
		LicensorID: object["licensorID"].(string),
		Pricing: abstract.Pricing{
			Single: abstract.Price{
				Currency: single["currency"].(string),
				Amount:   uint(single["currency"].(float64)),
			},
		},
	}
}
