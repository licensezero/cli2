package main

// Vendor represents a party that sold a license.
type Vendor struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	Website      string `json:"website"`
}
