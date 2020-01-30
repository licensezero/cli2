package main

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"path"
)

// Receipt represents a receipt for a license.
type Receipt struct {
	OfferID   string
	OrderID   string
	Effective string
	Expires   string
	Licensor  struct {
		EMail        string
		Jurisdiction string
		Name         string
		LicensorID   string
	}
	Licensee struct {
		EMail        string
		Jurisdiction string
		Name         string
	}
	Vendor struct {
		API          string
		EMail        string
		Name         string
		Jurisdiction string
		Website      string
	}
}

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) ([]Receipt, error) {
	directoryPath := path.Join(configPath, "receipts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []Receipt{}, nil
		}
		return nil, directoryReadError
	}
	var receipts []Receipt
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "receipts", name)
		receipt, err := readReceipt(filePath)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, *receipt)
	}
	return receipts, nil
}

//go:generate ./schema-to-string jurisdiction1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/jurisdiction.json
//go:generate ./schema-to-string key1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/key.json
//go:generate ./schema-to-string name1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/name.json
//go:generate ./schema-to-string price1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/price.json
//go:generate ./schema-to-string receipt1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/receipt.json
//go:generate ./schema-to-string signature1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/signature.json
//go:generate ./schema-to-string time1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/time.json
//go:generate ./schema-to-string url1_0_0Pre https://schemas.licensezero.com/1.0.0-pre/url.json

func readReceipt(filePath string) (*Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	err = json.Unmarshal(data, &unstructured)
	if err != nil {
		return nil, err
	}
	asMap, ok := unstructured.(map[string]interface{})
	if ok != true {
		return nil, errors.New("not an Object")
	}
	schemaID, ok := asMap["$schema"].(string)
	if ok != true {
		return nil, errors.New("$schema is not a string")
	}
	if schemaID == "https://schemas.licensezero.com/1.0.0-pre/receipt.json" {
		return parseV1Receipt(asMap)
	}
	return nil, errors.New("unknown schema: " + schemaID)
}

var schema1_0_0Pre = []string{
	jurisdiction1_0_0Pre,
	key1_0_0Pre,
	price1_0_0Pre,
	signature1_0_0Pre,
	time1_0_0Pre,
	url1_0_0Pre,
}

func parseV1Receipt(parsed map[string]interface{}) (*Receipt, error) {
	schemaLoader := gojsonschema.NewSchemaLoader()
	for _, schema := range schema1_0_0Pre {
		loader := gojsonschema.NewStringLoader(schema)
		schemaLoader.AddSchemas(loader)
	}
	receiptLoader := gojsonschema.NewStringLoader(receipt1_0_0Pre)
	schema, err := schemaLoader.Compile(receiptLoader)
	if err != nil {
		return nil, err
	}

	document := gojsonschema.NewGoLoader(parsed)
	result, err := schema.Validate(document)
	if err != nil {
		return nil, err
	}
	if !result.Valid() {
		return nil, errors.New("does not conform to schema")
	}

	var receipt Receipt
	license := parsed["license"].(map[string]interface{})
	values := license["values"].(map[string]interface{})
	receipt.OfferID = values["offerID"].(string)
	receipt.OrderID = values["orderID"].(string)
	receipt.Effective = values["effective"].(string)
	receipt.Expires = values["expires"].(string)

	licensor := values["licensor"].(map[string]string)
	receipt.Licensor.EMail = licensor["email"]
	receipt.Licensor.Jurisdiction = licensor["jurisdiction"]
	receipt.Licensor.Name = licensor["name"]
	receipt.Licensor.LicensorID = licensor["licensorID"]

	licensee := values["licensee"].(map[string]string)
	receipt.Licensee.EMail = licensee["email"]
	receipt.Licensee.Jurisdiction = licensee["jurisdiction"]
	receipt.Licensee.Name = licensee["name"]

	vendor := values["vendor"].(map[string]string)
	receipt.Vendor.API = vendor["api"]
	receipt.Vendor.EMail = vendor["email"]
	receipt.Vendor.Jurisdiction = vendor["jurisdiction"]
	receipt.Vendor.Name = vendor["name"]
	receipt.Vendor.Website = vendor["website"]

	return &receipt, nil
}
