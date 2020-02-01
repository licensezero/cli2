package main

// Price represents a price in a specific currency.
type Price struct {
	Currency string
	Amount   uint
}

const price1_0_0PreSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/price.json",
  "title": "price",
  "type": "object",
  "required": [
    "amount",
    "currency"
  ],
  "additionalProperties": false,
  "properties": {
    "amount": {
      "title": "purchase price in minor units of currency",
      "type": "integer",
      "minimum": 1,
      "examples": [
        0,
        100
      ]
    },
    "currency": {
      "title": "purchase price currency code",
      "$ref": "currency.json",
      "examples": [
        "USD"
      ]
    }
  }
}`
