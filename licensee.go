package main

// Licensee represents a party that bought a license.
type Licensee struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}
