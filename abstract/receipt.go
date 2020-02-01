package abstract

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
