package enums

type SentimentSort interface {
	private()
	String() string
}

type sentimentSort string

func (c sentimentSort) private() {}
func (c sentimentSort) String() string {
	return string(c)
}

const LATEST sentimentSort = "LATEST"
const EARLIEST sentimentSort = "EARLIEST"
const RELEVANCE sentimentSort = "RELEVANCE"
