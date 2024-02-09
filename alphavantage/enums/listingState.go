package enums

type ListingState interface {
	private()
	String() string
}

type listingState string

func (c listingState) private() {}
func (c listingState) String() string {
	return string(c)
}

const Active listingState = "active"
const Delisted listingState = "delisted"
