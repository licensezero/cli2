package schemas

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/cli2/abstract"
)

const receipt1_0_0Pre = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/receipt.json",
  "title": "license receipt",
  "comment": "A receipt represents confirmation of the sale of a software license.",
  "type": "object",
  "required": [
    "key",
    "signature",
    "license"
  ],
  "additionalProperties": false,
  "properties": {
    "key": {
      "title": "public signing key of the license vendor",
      "$ref": "key.json"
    },
    "signature": {
      "title": "signature of the license vendor",
      "$ref": "signature.json"
    },
    "license": {
      "title": "license manifest",
      "type": "object",
      "required": [
        "values",
        "form"
      ],
      "properties": {
        "values": {
          "type": "object",
          "required": [
            "api",
            "offerID",
            "orderID",
            "effective",
            "licensor",
            "licensee"
          ],
          "additionalProperties": false,
          "properties": {
            "api": {
              "title": "license API",
              "$ref": "url.json"
            },
            "offerID": {
              "title": "offer identifier",
              "type": "string",
              "format": "uuid"
            },
            "orderID": {
              "title": "order identifier",
              "type": "string",
              "format": "uuid"
            },
            "effective": {
              "title": "effective date",
              "$ref": "time.json"
            },
            "price": {
              "title": "purchase price",
              "$ref": "price.json"
            },
            "expires": {
              "title": "expiration date of the license",
              "$ref": "time.json"
            },
            "licensee": {
              "title": "licensee",
              "comment": "The licensee is the one receiving the license.",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "name"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "jurisdiction": {
                  "$ref": "jurisdiction.json"
                },
                "name": {
                  "$ref": "name.json",
                  "examples": [
                    "Joe Licensee"
                  ]
                }
              }
            },
            "licensor": {
              "title": "licensor",
              "comment": "The licensor is the one giving the license.",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "name",
                "licensorID"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "jurisdiction": {
                  "$ref": "jurisdiction.json"
                },
                "name": {
                  "$ref": "name.json",
                  "examples": [
                    "Joe Licensor"
                  ]
                },
                "licensorID": {
                  "title": "licensor identifier",
                  "type": "string",
                  "format": "uuid"
                }
              }
            },
            "vendor": {
              "title": "licesne vendor",
              "comment": "information on the party that sold the license, such as an agent or reseller, if the licensor did not sell the license themself",
              "type": "object",
              "required": [
                "email",
                "name",
                "jurisdiction",
                "website"
              ],
              "additionalProperties": false,
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "name": {
                  "$ref": "name.json",
                  "example": [
                    "Artless Devices LLC"
                  ]
                },
                "jurisdiction": {
                  "$ref": "jurisdiction.json"
                },
                "website": {
                  "$ref": "url.json"
                }
              }
            }
          }
        },
        "form": {
          "title": "license form",
          "type": "string",
          "minLength": 1
        }
      }
    }
  }
}`

// ParseReceipt validates and parses parsed JSON data as a Receipt.
func ParseReceipt(unstructured interface{}) (*abstract.Receipt, error) {
	if validV1Receipt(unstructured) {
		return parseV1Receipt(unstructured), nil
	}
	return nil, errors.New("unknown schema")
}

func validV1Receipt(parsed interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ReceiptSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

func parseV1Receipt(parsed interface{}) *abstract.Receipt {
	object := parsed.(map[string]interface{})
	license := object["license"].(map[string]interface{})
	values := license["values"].(map[string]interface{})
	licensor := values["licensor"].(map[string]interface{})
	licensee := values["licensee"].(map[string]interface{})
	// Parse optional expiration date.
	expires, ok := values["expires"].(string)
	if ok == false {
		expires = ""
	}
	// Parse optional vendor information.
	var vendor abstract.Vendor
	vendorMap, ok := values["vendor"].(map[string]interface{})
	if ok == true {
		vendor = abstract.Vendor{
			EMail:        vendorMap["email"].(string),
			Jurisdiction: vendorMap["jurisdiction"].(string),
			Name:         vendorMap["name"].(string),
			Website:      vendorMap["website"].(string),
		}
	}
	// Parse optional price.
	var price abstract.Price
	priceMap, ok := values["price"].(map[string]interface{})
	if ok == true {
		price = abstract.Price{
			Currency: priceMap["currency"].(string),
			Amount:   uint(priceMap["amount"].(float64)),
		}
	}
	return &abstract.Receipt{
		API:       values["api"].(string),
		OfferID:   values["offerID"].(string),
		OrderID:   values["orderID"].(string),
		Effective: values["effective"].(string),
		Expires:   expires,
		Price:     price,
		Licensor: abstract.Licensor{
			EMail:        licensor["email"].(string),
			Jurisdiction: licensor["jurisdiction"].(string),
			Name:         licensor["name"].(string),
			LicensorID:   licensor["licensorID"].(string),
		},
		Licensee: abstract.Licensee{
			EMail:        licensee["email"].(string),
			Jurisdiction: licensee["jurisdiction"].(string),
			Name:         licensee["name"].(string),
		},
		Vendor: vendor,
	}
}
