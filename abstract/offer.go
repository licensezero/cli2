package abstract

// Offer represents an offer to sell licenses.
type Offer struct {
	API        string
	OfferID    string
	URL        string
	LicensorID string
	Pricing    Pricing
}

type Pricing struct {
	Single Price
}
