package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/ed25519"
)

// Receipt represents a receipt for a license.
type Receipt interface {
	API() string
	OfferID() string
	OrderID() string
	Effective() string
	Expires() string
	Price() Price
	Licensor() Licensor
	Licensee() Licensee
	Vendor() Vendor
	Form() string
	ValidateSignature() error
}

type receipt1_0_0Pre struct {
	Key       string
	Signature string
	License   struct {
		Values struct {
			API       string
			OfferID   string
			OrderID   string
			Effective string
			Price     Price
			Expires   string
			Licensor  Licensor
			Licensee  Licensee
			Vendor    Vendor
		}
		Form string
	}
}

func (r receipt1_0_0Pre) API() string {
	return r.License.Values.API
}

func (r receipt1_0_0Pre) OfferID() string {
	return r.License.Values.OfferID
}

func (r receipt1_0_0Pre) OrderID() string {
	return r.License.Values.OrderID
}

func (r receipt1_0_0Pre) Effective() string {
	return r.License.Values.Effective
}

func (r receipt1_0_0Pre) Expires() string {
	return r.License.Values.Expires
}

func (r receipt1_0_0Pre) Price() Price {
	return r.License.Values.Price
}

func (r receipt1_0_0Pre) Licensor() Licensor {
	return r.License.Values.Licensor
}

func (r receipt1_0_0Pre) Licensee() Licensee {
	return r.License.Values.Licensee
}

func (r receipt1_0_0Pre) Vendor() Vendor {
	return r.License.Values.Vendor
}

func (r receipt1_0_0Pre) Form() string {
	return r.License.Form
}

func (r receipt1_0_0Pre) ValidateSignature() error {
	serialized := serializeV1License(&r)
	return checkSignature(r.Key, r.Signature, serialized)
}

// Manually implement JSON serialization for receipt license objects
// to ensure that keys are serialized in sorted order and optional
// properties, like price, get correctly omitted.
func serializeV1License(r *receipt1_0_0Pre) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"form\":" + stringify(r.License.Form) + ",")
	values := r.License.Values

	// Values
	buffer.WriteString("\"values\":{")
	buffer.WriteString("\"api\":" + stringify(values.API) + ",")
	buffer.WriteString("\"effective\":" + stringify(values.Effective) + ",")
	if expires := values.Expires; expires != "" {
		buffer.WriteString("\"expires\":" + stringify(expires) + ",")
	}

	// Licensee
	licensee := values.Licensee
	buffer.WriteString("\"licensee\":{")
	buffer.WriteString("\"email\":" + stringify(licensee.EMail) + ",")
	buffer.WriteString("\"jurisdiction\":" + stringify(licensee.Jurisdiction) + ",")
	buffer.WriteString("\"name\":" + stringify(licensee.Name))
	buffer.WriteString("},")

	// Licensor
	licensor := values.Licensor
	buffer.WriteString("\"licensor\":{")
	buffer.WriteString("\"email\":" + stringify(licensor.EMail) + ",")
	buffer.WriteString("\"jurisdiction\":" + stringify(licensor.Jurisdiction) + ",")
	buffer.WriteString("\"licensorID\":" + stringify(licensor.LicensorID) + ",")
	buffer.WriteString("\"name\":" + stringify(licensor.Name))
	buffer.WriteString("},")

	// OfferID
	buffer.WriteString("\"offerID\":" + stringify(values.OfferID) + ",")

	// OrderID
	buffer.WriteString("\"orderID\":" + stringify(values.OrderID))

	// Price
	if price := values.Price; price.Currency != "" && price.Amount != 0 {
		buffer.WriteString(",")
		buffer.WriteString("\"price\":{")
		buffer.WriteString("\"amount\":" + stringify(price.Amount) + ",")
		buffer.WriteString("\"currency\":" + stringify(price.Currency))
		buffer.WriteString("}")
	}

	// Vendor
	if vendor := values.Vendor; vendor.Name != "" {
		buffer.WriteString(",")
		buffer.WriteString("\"vendor\":{")
		buffer.WriteString("\"email\":" + stringify(vendor.EMail) + ",")
		buffer.WriteString("\"jurisdiction\":" + stringify(vendor.Jurisdiction) + ",")
		buffer.WriteString("\"name\":" + stringify(vendor.Name) + ",")
		buffer.WriteString("\"website\":" + stringify(vendor.Website))
		buffer.WriteString("}")
	}

	buffer.WriteString("}") // values
	buffer.WriteString("}") // object

	return buffer.Bytes()
}

func stringify(input interface{}) (result string) {
	data, _ := json.Marshal(input)
	return string(data)
}

func checkSignature(publicKey string, signature string, json []byte) error {
	signatureBytes := make([]byte, hex.DecodedLen(len(signature)))
	_, err := hex.Decode(signatureBytes, []byte(signature))
	if err != nil {
		return errors.New("invalid signature")
	}
	publicKeyBytes := make([]byte, hex.DecodedLen(len(publicKey)))
	_, err = hex.Decode(publicKeyBytes, []byte(publicKey))
	if err != nil {
		return errors.New("invalid public key")
	}
	signatureValid := ed25519.Verify(
		publicKeyBytes,
		json,
		signatureBytes,
	)
	if !signatureValid {
		return errors.New("invalid signature")
	}
	return nil
}

const receipt1_0_0PreSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/receipt.json",
  "title": "license receipt",
  "comment": "A receipt represents confirmation of the sale of a software license.",
  "type": "object",
  "required": [
    "key",
    "license",
    "signature"
  ],
  "additionalProperties": false,
  "properties": {
    "key": {
      "title": "public signing key of the license vendor",
      "$ref": "key.json"
    },
    "license": {
      "title": "license manifest",
      "type": "object",
      "required": [
        "form",
        "values"
      ],
      "properties": {
        "form": {
          "title": "license form",
          "type": "string",
          "minLength": 1
        },
        "values": {
          "type": "object",
          "required": [
            "api",
            "effective",
            "licensee",
            "licensor",
            "offerID",
            "orderID"
          ],
          "additionalProperties": false,
          "properties": {
            "api": {
              "title": "license API",
              "$ref": "url.json"
            },
            "effective": {
              "title": "effective date",
              "$ref": "time.json"
            },
            "expires": {
              "title": "expiration date of the license",
              "$ref": "time.json"
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
            "price": {
              "title": "purchase price",
              "$ref": "price.json"
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
                "licensorID",
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
                "licensorID": {
                  "title": "licensor identifier",
                  "type": "string",
                  "format": "uuid"
                },
                "name": {
                  "$ref": "name.json",
                  "examples": [
                    "Joe Licensor"
                  ]
                }
              }
            },
            "vendor": {
              "title": "licesne vendor",
              "comment": "information on the party that sold the license, such as an agent or reseller, if the licensor did not sell the license themself",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "name",
                "website"
              ],
              "additionalProperties": false,
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
                  "example": [
                    "Artless Devices LLC"
                  ]
                },
                "website": {
                  "$ref": "url.json"
                }
              }
            }
          }
        }
      }
    },
    "signature": {
      "title": "signature of the license vendor",
      "$ref": "signature.json"
    }
  }
}`

// ParseReceipt validates and parses parsed JSON data as a Receipt.
func ParseReceipt(unstructured interface{}) (Receipt, error) {
	if validV1Receipt(unstructured) {
		return parseV1Receipt(unstructured), nil
	}
	return nil, errors.New("unknown schema")
}

var v1ReceiptSchema *gojsonschema.Schema = nil

func validV1Receipt(parsed interface{}) bool {
	if v1ReceiptSchema == nil {
		schema, err := schemaLoader().Compile(
			gojsonschema.NewStringLoader(receipt1_0_0PreSchema),
		)
		if err != nil {
			panic(err)
		}
		v1ReceiptSchema = schema
	}
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ReceiptSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	if result.Valid() {
		return true
	}
	return false
}

func parseV1Receipt(unstructured interface{}) (r receipt1_0_0Pre) {
	mapstructure.Decode(unstructured, &r)
	return
}
