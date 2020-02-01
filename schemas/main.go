package schemas

import (
	"github.com/xeipuuv/gojsonschema"
)

var v1ArtifactMetadataSchema *gojsonschema.Schema
var v1OfferSchema *gojsonschema.Schema
var v1ReceiptSchema *gojsonschema.Schema

func init() {
	subschemas := []string{
		Jurisdiction1_0_0Pre,
		Key1_0_0Pre,
		Price1_0_0Pre,
		Signature1_0_0Pre,
		Time1_0_0Pre,
		URL1_0_0Pre,
	}
	v1SubschemaLoader := gojsonschema.NewSchemaLoader()
	for _, schema := range subschemas {
		loader := gojsonschema.NewStringLoader(schema)
		v1SubschemaLoader.AddSchemas(loader)
	}
	// Artifact
	v1ArtifactMetadataSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(artifact1_0_0Pre),
	)
	// Offer
	v1OfferSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(offer1_0_0Pre),
	)
	// Receipt
	v1ReceiptSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(receipt1_0_0Pre),
	)
}
