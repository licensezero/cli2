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
	API       string
	OfferID   string
	OrderID   string
	Effective string
	Expires   string
	Price     Price
	Licensor  Licensor
	Licensee  Licensee
	Vendor    Vendor
}

// Price represents a monetary amount.
type Price struct {
	Currency string
	Amount   uint
}

// Licensor represents a party that offered licenses for sale.
type Licensor struct {
	EMail        string
	Jurisdiction string
	Name         string
	LicensorID   string
}

// Licensee represents a party that bought a license.
type Licensee struct {
	EMail        string
	Jurisdiction string
	Name         string
}

// Vendor represents a party that sold a license.
type Vendor struct {
	EMail        string
	Name         string
	Jurisdiction string
	Website      string
}

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) (receipts []Receipt, errors []error, err error) {
	directoryPath := path.Join(configPath, "receipts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return
		}
		return nil, nil, directoryReadError
	}
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "receipts", name)
		receipt, err := readReceipt(filePath)
		if err != nil {
			errors = append(errors, err)
		} else {
			receipts = append(receipts, *receipt)
		}
	}
	return
}

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
	return ParseReceipt(unstructured)
}

// ParseReceipt validates and parses parsed JSON data as a Receipt.
func ParseReceipt(unstructured interface{}) (*Receipt, error) {
	if validV1Receipt(unstructured) {
		return parseV1Receipt(unstructured), nil
	}
	return nil, errors.New("unknown schema")
}

var schema1_0_0Pre = []string{
	jurisdiction1_0_0Pre,
	key1_0_0Pre,
	price1_0_0Pre,
	signature1_0_0Pre,
	time1_0_0Pre,
	url1_0_0Pre,
}

var v1SchemaLoader *gojsonschema.SchemaLoader
var v1ReceiptSchema *gojsonschema.Schema

func init() {
	v1SchemaLoader = gojsonschema.NewSchemaLoader()
	for _, schema := range schema1_0_0Pre {
		loader := gojsonschema.NewStringLoader(schema)
		v1SchemaLoader.AddSchemas(loader)
	}
	receiptLoader := gojsonschema.NewStringLoader(receipt1_0_0Pre)
	v1ReceiptSchema, _ = v1SchemaLoader.Compile(receiptLoader)
}

func validV1Receipt(parsed interface{}) bool {
	dataLoader := gojsonschema.NewGoLoader(parsed)
	result, err := v1ReceiptSchema.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

func parseV1Receipt(parsed interface{}) *Receipt {
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
	var vendor Vendor
	vendorMap, ok := values["vendor"].(map[string]interface{})
	if ok == true {
		vendor = Vendor{
			EMail:        vendorMap["email"].(string),
			Jurisdiction: vendorMap["jurisdiction"].(string),
			Name:         vendorMap["name"].(string),
			Website:      vendorMap["website"].(string),
		}
	}
	// Parse optional price.
	var price Price
	priceMap, ok := values["price"].(map[string]interface{})
	if ok == true {
		price = Price{
			Currency: priceMap["currency"].(string),
			Amount:   uint(priceMap["amount"].(float64)),
		}
	}
	return &Receipt{
		API:       values["api"].(string),
		OfferID:   values["offerID"].(string),
		OrderID:   values["orderID"].(string),
		Effective: values["effective"].(string),
		Expires:   expires,
		Price:     price,
		Licensor: Licensor{
			EMail:        licensor["email"].(string),
			Jurisdiction: licensor["jurisdiction"].(string),
			Name:         licensor["name"].(string),
			LicensorID:   licensor["licensorID"].(string),
		},
		Licensee: Licensee{
			EMail:        licensee["email"].(string),
			Jurisdiction: licensee["jurisdiction"].(string),
			Name:         licensee["name"].(string),
		},
		Vendor: vendor,
	}
}
