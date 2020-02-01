package abstract

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
	key       string `json:"key"`
	signature string `json:"signature"`
	license   struct {
		values struct {
			api       string   `json:"api"`
			offerID   string   `json:"offerID"`
			orderID   string   `json:"orderID"`
			effective string   `json:"effective"`
			price     Price    `json:"price"`
			expires   string   `json:"expires"`
			licensor  Licensor `json:"licensor"`
			licensee  Licensee `json:"licensee"`
			vendor    Vendor   `json:"vendor"`
		} `json:"values"`
		form string `json:"form"`
	} `json:"license"`
}

func (r receipt1_0_0Pre) API() string {
	return r.license.values.api
}

func (r receipt1_0_0Pre) OfferID() string {
	return r.license.values.offerID
}

func (r receipt1_0_0Pre) OrderID() string {
	return r.license.values.orderID
}

func (r receipt1_0_0Pre) Effective() string {
	return r.license.values.effective
}

func (r receipt1_0_0Pre) Expires() string {
	return r.license.values.expires
}

func (r receipt1_0_0Pre) Price() Price {
	return r.license.values.price
}

func (r receipt1_0_0Pre) Licensor() Licensor {
	return r.license.values.licensor
}

func (r receipt1_0_0Pre) Licensee() Licensee {
	return r.license.values.licensee
}

func (r receipt1_0_0Pre) Vendor() Vendor {
	return r.license.values.vendor
}

func (r receipt1_0_0Pre) Form() string {
	return r.license.form
}

func (r receipt1_0_0Pre) ValidateSignature() error {
	serialized, err := json.Marshal(r.license)
	if err != nil {
		return errors.New("could not serialize")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not compact JSON")
	}
	return checkSignature(r.key, r.signature, compacted.Bytes())
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
func ParseReceipt(unstructured interface{}) (Receipt, error) {
	if validV1Receipt(unstructured) {
		var v1 receipt1_0_0Pre
		err := mapstructure.Decode(unstructured, &v1)
		if err != nil {
			return nil, err
		}
		return v1, nil
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
