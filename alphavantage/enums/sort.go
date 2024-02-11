package enums

type SentimentSort interface {
	privateSentimentSort()
	String() string
}

type sentimentSort string

func (c sentimentSort) privateSentimentSort() {}
func (c sentimentSort) String() string {
	return string(c)
}

const LATEST sentimentSort = "LATEST"
const EARLIEST sentimentSort = "EARLIEST"
const RELEVANCE sentimentSort = "RELEVANCE"
