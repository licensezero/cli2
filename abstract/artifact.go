package abstract

// ArtifactMetadata encodes data about offers for an artifact.
type ArtifactMetadata struct {
	Offers []ArtifactOffer `json:"offers" toml:"offers"`
}

// ArtifactOffer represents an offer relevant to an artifact.
type ArtifactOffer struct {
	API     string `json:"api" toml:"api"`
	OfferID string `json:"offerID" toml:"offerID"`
	Public  string `json:"public" toml:"public"`
}
