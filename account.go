package main

// Account contains information about a licensor account.
type Account interface {
	API() string
	LicensorID() string
	Token() string
}
