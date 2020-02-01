package abstract

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

// ArtifactMetadata encodes data about offers for an artifact.
type ArtifactMetadata interface {
	Offers() []ArtifactOffer
}

// ArtifactOffer represents an offer relevant to an artifact.
type ArtifactOffer struct {
	API     string
	OfferID string
	Public  string
}

type artifact1_0_0Pre struct {
	offers []struct {
		api     string `json:"api"`
		offerID string `json:"offerID"`
		public  string `json:"public"`
	}
}

func (a artifact1_0_0Pre) Offers() (offers []ArtifactOffer) {
	for _, item := range a.offers {
		offers = append(offers, ArtifactOffer{
			API:     item.api,
			OfferID: item.offerID,
			Public:  item.public,
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

// ParseArtifactMetadata validates and parses parsed JSON data as a ArtifactMetadata.
func ParseArtifactMetadata(unstructured interface{}) (ArtifactMetadata, error) {
	if validV1ArtifactMetadata(unstructured) {
		var v1 artifact1_0_0Pre
		err := mapstructure.Decode(unstructured, &v1)
		if err != nil {
			return nil, err
		}
		return v1, nil
	}
	return nil, errors.New("unknown schema")
}

func validV1ArtifactMetadata(parsed interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ArtifactMetadataSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
