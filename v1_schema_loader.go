package main

import (
	"github.com/xeipuuv/gojsonschema"
)

func schemaLoader() *gojsonschema.SchemaLoader {
	subschemas := []string{
		jurisdiction1_0_0PreSchema,
		key1_0_0PreSchema,
		price1_0_0PreSchema,
		signature1_0_0PreSchema,
		time1_0_0PreSchema,
		url1_0_0PreSchema,
	}
	loader := gojsonschema.NewSchemaLoader()
	for _, schema := range subschemas {
		loader.AddSchemas(gojsonschema.NewStringLoader(schema))
	}
	return loader
}
