package main

// Licensor represents a party that offered licenses for sale.
type Licensor struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	LicensorID   string `json:"licensorID"`
	Name         string `json:"name"`
}
