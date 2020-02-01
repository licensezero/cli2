package abstract

import (
	"github.com/xeipuuv/gojsonschema"
)

var v1ArtifactMetadataSchema *gojsonschema.Schema
var v1OfferSchema *gojsonschema.Schema
var v1ReceiptSchema *gojsonschema.Schema

func init() {
	subschemas := []string{
		jurisdiction1_0_0PreSchema,
		key1_0_0PreSchema,
		price1_0_0PreSchema,
		signature1_0_0PreSchema,
		time1_0_0PreSchema,
		url1_0_0PreSchema,
	}
	v1SubschemaLoader := gojsonschema.NewSchemaLoader()
	for _, schema := range subschemas {
		loader := gojsonschema.NewStringLoader(schema)
		v1SubschemaLoader.AddSchemas(loader)
	}
	// Artifact
	v1ArtifactMetadataSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(artifact1_0_0PreSchema),
	)
	// Offer
	v1OfferSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(offer1_0_0PreSchema),
	)
	// Receipt
	v1ReceiptSchema, _ = v1SubschemaLoader.Compile(
		gojsonschema.NewStringLoader(receipt1_0_0PreSchema),
	)
}
