package main

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
	"testing"
)

/*

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
	if receipt.API() != "https://api.licensezero.com" {
		t.Error("failed to parse API")
	}
	if receipt.OrderID() != "2c743a84-09ce-4549-9f0d-19d8f53462bb" {
		t.Error("failed to parse orderID")
	}
	if receipt.OfferID() != "9aab7058-599a-43db-9449-5fc0971ecbfa" {
		t.Error("failed to parse offerID")
	}
	if err := receipt.ValidateSignature(); err == nil {
		t.Error("validates invalid signature")
	}
}

func TestReceiptSignature(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	message := `{"form":"Test license form.","values":{"api":"https://api.licensezero.com","effective":"2018-11-13T20:20:39Z","licensee":{"email":"licensee@example.com","jurisdiction":"US-TX","name":"Joe"},"licensor":{"email":"licensor@example.com","jurisdiction":"US-CA","licensorID":"59e70a4d-ffee-4e9d-a526-7a9ff9161664","name":"Jane"},"offerID":"9aab7058-599a-43db-9449-5fc0971ecbfa","orderID":"2c743a84-09ce-4549-9f0d-19d8f53462bb"}}`
	signature := ed25519.Sign(privateKey, []byte(message))
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	combined := "{" +
		quote("key") + ":" + quote(publicKeyHex) + "," +
		quote("signature") + ":" + quote(signatureHex) + "," +
		quote("license") + ":" + message +
		"}"
	var unstructured interface{}
	err := json.Unmarshal([]byte(combined), &unstructured)
	if err != nil {
		t.Error(err)
	}
	receipt, err := ParseReceipt(unstructured)
	if err := receipt.ValidateSignature(); err != nil {
		t.Error("invalidates invalid signature")
	}
}
*/

func TestReceiptWithPriceAndVendorSignature(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	message := `{"form":"Test license form.","values":{"api":"https://api.licensezero.com","effective":"2018-11-13T20:20:39Z","licensee":{"email":"licensee@example.com","jurisdiction":"US-TX","name":"Joe"},"licensor":{"email":"licensor@example.com","jurisdiction":"US-CA","licensorID":"59e70a4d-ffee-4e9d-a526-7a9ff9161664","name":"Jane"},"offerID":"9aab7058-599a-43db-9449-5fc0971ecbfa","orderID":"2c743a84-09ce-4549-9f0d-19d8f53462bb","price":{"amount":1000,"currency":"USD"},"vendor":{"email":"vendor@example.com","jurisdiction":"US-CA","name":"Vendor","website":"https://example.com"}}}`
	signature := ed25519.Sign(privateKey, []byte(message))
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	combined := "{" +
		quote("key") + ":" + quote(publicKeyHex) + "," +
		quote("signature") + ":" + quote(signatureHex) + "," +
		quote("license") + ":" + message +
		"}"
	var unstructured interface{}
	err := json.Unmarshal([]byte(combined), &unstructured)
	if err != nil {
		t.Error(err)
	}
	receipt, err := ParseReceipt(unstructured)
	if err := receipt.ValidateSignature(); err != nil {
		t.Error("invalidates invalid signature")
	}
}

func quote(input string) string {
	return "\"" + input + "\""
}
