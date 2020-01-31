package main

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
)

// ArtifactMetadata encodes data about offers for an artifact.
type ArtifactMetadata struct {
	Offers []ArtifactOffer `json:"offers" toml:"offers"`
}

// ArtifactOffer represents an offer relevant to an artifact.
type ArtifactOffer struct {
	API     string `json:"api" toml:"api"`
	OfferID string `json:"offerID" toml:"offerID"`
	Public  string `json:"public" toml:"public"`
}

// ParseArtifactMetadata validates and parses parsed JSON data as a ArtifactMetadata.
func ParseArtifactMetadata(unstructured interface{}) (*ArtifactMetadata, error) {
	if validV1ArtifactMetadata(unstructured) {
		return parseV1ArtifactMetadata(unstructured), nil
	}
	return nil, errors.New("unknown schema")
}

var v1ArtifactMetadataSchemaLoader *gojsonschema.SchemaLoader
var v1ArtifactMetadataSchema *gojsonschema.Schema

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

func init() {
	v1ArtifactMetadataSchemaLoader = gojsonschema.NewSchemaLoader()
	for _, schema := range schema1_0_0Pre {
		loader := gojsonschema.NewStringLoader(schema)
		v1ArtifactMetadataSchemaLoader.AddSchemas(loader)
	}
	artifactLoader := gojsonschema.NewStringLoader(artifact1_0_0Pre)
	v1ArtifactMetadataSchema, _ = v1ArtifactMetadataSchemaLoader.Compile(artifactLoader)
}

func validV1ArtifactMetadata(parsed interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ArtifactMetadataSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

func parseV1ArtifactMetadata(parsed interface{}) *ArtifactMetadata {
	object := parsed.(map[string]interface{})
	offers := object["offers"].([]interface{})
	var returned ArtifactMetadata
	for _, element := range offers {
		asMap := element.(map[string]interface{})
		public, ok := asMap["public"].(string)
		if ok == false {
			public = ""
		}
		returned.Offers = append(returned.Offers, ArtifactOffer{
			API:     asMap["api"].(string),
			OfferID: asMap["offerID"].(string),
			Public:  public,
		})
	}
	return &returned
}
