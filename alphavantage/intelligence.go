package alphavantage

// parameter for NEWS_SENTIMENT API
type SentimentParameter struct {
	// The stock/crypto/forex symbols of your choice.
	// For example: tickers=IBM will filter for articles that mention the IBM ticker; tickers=COIN,CRYPTO:BTC,FOREX:USD
	// will filter for articles that simultaneously mention Coinbase (COIN), Bitcoin (CRYPTO:BTC), and US Dollar (FOREX:USD) in their content.
	Tickers string `json:"tickers"`
	// The news topics of your choice. For example: topics=technology will filter for articles that write about the technology sector;
	// topics=technology,ipo will filter for articles that simultaneously cover technology and IPO in their content.
	Topics string `json:"topics"`
	// The time range of the news articles you are targeting, in YYYYMMDDTHHMM format.
	TimeFrom string `json:"time_from"`
	// The time range of the news articles you are targeting, in YYYYMMDDTHHMM format.
	TimeTo string `json:"time_to"`
	// By default, sort=LATEST and the API will return the latest articles first.
	// You can also set sort=EARLIEST or sort=RELEVANCE based on your use case.
	Sort string `json:"sort"`
	// By default, limit=50 and the API will return up to 50 matching results. You can also set limit=1000 to output up to 1000 results.
	Limit int `json:"limit"`
}

type NewsSentiment struct {
	Items                    string `json:"items"`
	SentimentScoreDefinition string `json:"sentiment_score_definition"`
	RelevanceScoreDefinition string `json:"relevance_score_definition"`
	Feed                     []Feed `json:"feed"`
}

type Feed struct {
	Title                 string            `json:"title"`
	URL                   string            `json:"url"`
	TimePublished         string            `json:"time_published"`
	Authors               []string          `json:"authors"`
	Summary               string            `json:"summary"`
	BannerImage           *string           `json:"banner_image"`
	Source                string            `json:"source"`
	CategoryWithinSource  string            `json:"category_within_source"`
	SourceDomain          string            `json:"source_domain"`
	Topics                []Topic           `json:"topics"`
	OverallSentimentScore float64           `json:"overall_sentiment_score"`
	OverallSentimentLabel string            `json:"overall_sentiment_label"`
	TickerSentiment       []TickerSentiment `json:"ticker_sentiment"`
}
type Topic struct {
	Topic          string `json:"topic"`
	RelevanceScore string `json:"relevance_score"`
}

type TickerSentiment struct {
	Ticker               string `json:"ticker"`
	RelevanceScore       string `json:"relevance_score"`
	TickerSentimentScore string `json:"ticker_sentiment_score"`
	TickerSentimentLabel string `json:"ticker_sentiment_label"`
}

type IntelligenceClient struct {
	httpClient *httpClient
	apikey     string
}

// The APIs in this section contain advanced market intelligence built with our decades of expertise in AI,
// machine learning, and quantitative finance.
// We hope these highly differentiated alternative datasets can help turbocharge your trading strategy,
// market research, and financial software application to the next level.
func NewIntelligenceClientClient(apikey string) *IntelligenceClient {
	return &IntelligenceClient{
		httpClient: newHttpClient(),
		apikey:     apikey,
	}
}
