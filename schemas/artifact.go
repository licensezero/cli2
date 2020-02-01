package schemas

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/cli2/abstract"
)

const artifact1_0_0Pre = `{
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
func ParseArtifactMetadata(unstructured interface{}) (*abstract.ArtifactMetadata, error) {
	if validV1ArtifactMetadata(unstructured) {
		return parseV1ArtifactMetadata(unstructured), nil
	}
	return nil, errors.New("unknown schema")
}

func parseV1ArtifactMetadata(parsed interface{}) *abstract.ArtifactMetadata {
	object := parsed.(map[string]interface{})
	offers := object["offers"].([]interface{})
	var returned abstract.ArtifactMetadata
	for _, element := range offers {
		asMap := element.(map[string]interface{})
		public, ok := asMap["public"].(string)
		if ok == false {
			public = ""
		}
		returned.Offers = append(returned.Offers, abstract.ArtifactOffer{
			API:     asMap["api"].(string),
			OfferID: asMap["offerID"].(string),
			Public:  public,
		})
	}
	return &returned
}
func validV1ArtifactMetadata(parsed interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ArtifactMetadataSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
