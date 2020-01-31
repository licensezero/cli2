package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestReadReceipts(t *testing.T) {
	WithTestDir(t, func(directory string) {
		receipts := path.Join(directory, "receipts")
		err := os.MkdirAll(receipts, 0700)
		if err != nil {
			t.Error(err)
		}

		withVendor := path.Join(receipts, "withVendor.json")
		err = ioutil.WriteFile(
			withVendor,
			[]byte(`
{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "expires": "2019-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      },
      "price": {
        "amount": 1000,
        "currency": "USD"
      },
      "vendor": {
        "email": "support@artlessdevices.com",
        "name": "Artless Devices LLC",
        "jurisdiction": "US-CA",
        "website": "https://artlessdevices.com"
      }
    },
    "form": "Test license form."
  }
}
			`),
			0700,
		)
		if err != nil {
			t.Error(err)
		}

		withoutVendor := path.Join(receipts, "withoutVendor.json")
		err = ioutil.WriteFile(
			withoutVendor,
			[]byte(`
{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      }
    },
    "form": "Test license form."
  }
}
			`),
			0700,
		)
		if err != nil {
			t.Error(err)
		}

		invalid := path.Join(receipts, "invalid.json")
		err = ioutil.WriteFile(invalid, []byte(`{}`), 0700)
		if err != nil {
			t.Error(err)
		}

		results, receiptErrors, readError := ReadReceipts(directory)
		if readError != nil {
			t.Error("read error")
		}

		if len(results) != 2 {
			t.Error("did not find receipt")
		}

		first := results[0]
		if first.API != "https://api.licensezero.com" {
			t.Error("failed to parse API")
		}
		if first.OrderID != "2c743a84-09ce-4549-9f0d-19d8f53462bb" {
			t.Error("failed to parse orderID")
		}
		if first.OfferID != "9aab7058-599a-43db-9449-5fc0971ecbfa" {
			t.Error("failed to parse orderID")
		}
		if first.Effective != "2018-11-13T20:20:39Z" {
			t.Error("failed to parse effective date")
		}
		if first.Expires != "2019-11-13T20:20:39Z" {
			t.Error("added expiration date")
		}
		if first.Price.Amount != 1000 {
			t.Error("failed to parse price amount")
		}
		if first.Price.Currency != "USD" {
			t.Error("failed to parse price currency")
		}

		second := results[1]
		if second.Vendor.Name != "" {
			t.Error("failed to parse missing vendor")
		}

		if len(receiptErrors) != 1 {
			t.Error("missing invalid error")
		}
	})
}

func TestParseReceipt(t *testing.T) {
	bytes := []byte(`{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      }
    },
    "form": "Test license form."
  }
}`)
	var unstructured interface{}
	err := json.Unmarshal(bytes, &unstructured)
	if err != nil {
		t.Error(err)
	}
	receipt, err := ParseReceipt(unstructured)
	if receipt.API != "https://api.licensezero.com" {
		t.Error("failed to parse API")
	}
	if receipt.OrderID != "2c743a84-09ce-4549-9f0d-19d8f53462bb" {
		t.Error("failed to parse orderID")
	}
	if receipt.OfferID != "9aab7058-599a-43db-9449-5fc0971ecbfa" {
		t.Error("failed to parse orderID")
	}
}

func WithTestDir(t *testing.T, script func(string)) {
	directory, err := ioutil.TempDir("/tmp", "licensezero-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(directory)
	script(directory)
}
