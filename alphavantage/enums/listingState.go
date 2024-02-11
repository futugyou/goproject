package enums

type ListingState interface {
	privateListingState()
	String() string
}

type listingState string

func (c listingState) privateListingState() {}
func (c listingState) String() string {
	return string(c)
}

const Active listingState = "active"
const Delisted listingState = "delisted"
