package alphavantage

import (
	"strconv"
	"strings"

	"github.com/futugyou/alphavantage/enums"
)

// parameter for NEWS_SENTIMENT API
type SentimentParameter struct {
	// The stock/crypto/forex symbols of your choice.
	// For example: tickers=IBM will filter for articles that mention the IBM ticker; tickers=COIN,CRYPTO:BTC,FOREX:USD
	// will filter for articles that simultaneously mention Coinbase (COIN), Bitcoin (CRYPTO:BTC), and US Dollar (FOREX:USD) in their content.
	Tickers string `json:"tickers"`
	// The news topics of your choice. For example: topics=technology will filter for articles that write about the technology sector;
	// topics=technology,ipo will filter for articles that simultaneously cover technology and IPO in their content.
	Topics enums.SentimentTopics `json:"topics"`
	// The time range of the news articles you are targeting, in YYYYMMDDTHHMM format.
	TimeFrom string `json:"time_from"`
	// The time range of the news articles you are targeting, in YYYYMMDDTHHMM format.
	TimeTo string `json:"time_to"`
	// By default, sort=LATEST and the API will return the latest articles first.
	// You can also set sort=EARLIEST or sort=RELEVANCE based on your use case.
	Sort enums.SentimentSort `json:"sort"`
	// By default, limit=50 and the API will return up to 50 matching results. You can also set limit=1000 to output up to 1000 results.
	Limit int `json:"limit"`
}

func (p SentimentParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "NEWS_SENTIMENT"
	if len(strings.TrimSpace(p.Tickers)) > 0 {
		dic["tickers"] = strings.TrimSpace(p.Tickers)
	}
	if p.Topics != nil {
		dic["topics"] = p.Topics.String()
	}
	if len(strings.TrimSpace(p.TimeFrom)) > 0 {
		dic["time_from"] = strings.TrimSpace(p.TimeFrom)
	}
	if len(strings.TrimSpace(p.TimeTo)) > 0 {
		dic["time_to"] = strings.TrimSpace(p.TimeTo)
	}
	if p.Sort != nil {
		dic["sort"] = p.Sort.String()
	}
	if p.Limit > 50 && p.Limit < 1000 {
		dic["limit"] = strconv.Itoa(p.Limit)
	}
	return dic, nil
}

type NewsSentiment struct {
	Items                    string `json:"items,omitempty"`
	SentimentScoreDefinition string `json:"sentiment_score_definition,omitempty"`
	RelevanceScoreDefinition string `json:"relevance_score_definition,omitempty"`
	Feed                     []Feed `json:"feed,omitempty"`
	Information              string `json:"Information"`
	ErrorMessage             string `json:"Error Message"`
}

type Feed struct {
	Title                 string            `json:"title,omitempty"`
	URL                   string            `json:"url,omitempty"`
	TimePublished         string            `json:"time_published,omitempty"`
	Authors               []string          `json:"authors,omitempty"`
	Summary               string            `json:"summary,omitempty"`
	BannerImage           string            `json:"banner_image,omitempty"`
	Source                string            `json:"source,omitempty"`
	CategoryWithinSource  string            `json:"category_within_source,omitempty"`
	SourceDomain          string            `json:"source_domain,omitempty"`
	Topics                []Topic           `json:"topics,omitempty"`
	OverallSentimentScore float64           `json:"overall_sentiment_score,omitempty"`
	OverallSentimentLabel string            `json:"overall_sentiment_label,omitempty"`
	TickerSentiment       []TickerSentiment `json:"ticker_sentiment,omitempty"`
}
type Topic struct {
	Topic          string `json:"topic,omitempty"`
	RelevanceScore string `json:"relevance_score,omitempty"`
}

type TickerSentiment struct {
	Ticker               string `json:"ticker,omitempty"`
	RelevanceScore       string `json:"relevance_score,omitempty"`
	TickerSentimentScore string `json:"ticker_sentiment_score,omitempty"`
	TickerSentimentLabel string `json:"ticker_sentiment_label,omitempty"`
}

type TopGainersLosers struct {
	Metadata           string               `json:"metadata"`
	LastUpdated        string               `json:"last_updated"`
	TopGainers         []MostActivelyTraded `json:"top_gainers"`
	TopLosers          []MostActivelyTraded `json:"top_losers"`
	MostActivelyTraded []MostActivelyTraded `json:"most_actively_traded"`
	Information        string               `json:"Information"`
	ErrorMessage       string               `json:"Error Message"`
}

type MostActivelyTraded struct {
	Ticker           string `json:"ticker"`
	Price            string `json:"price"`
	ChangeAmount     string `json:"change_amount"`
	ChangePercentage string `json:"change_percentage"`
	Volume           string `json:"volume"`
}

type IntelligenceClient struct {
	innerClient
}

// The APIs in this section contain advanced market intelligence built with our decades of expertise in AI,
// machine learning, and quantitative finance.
// We hope these highly differentiated alternative datasets can help turbocharge your trading strategy,
// market research, and financial software application to the next level.
// TODO: https://alphavantageapi.co/timeseries/ return Forbidden
func NewIntelligenceClient(apikey string) *IntelligenceClient {
	return &IntelligenceClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

func (t *IntelligenceClient) NewsSentiment(p SentimentParameter) (*NewsSentiment, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &NewsSentiment{}

	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *IntelligenceClient) TopGainersLosers() (*TopGainersLosers, error) {
	dic := make(map[string]string)
	dic["function"] = "TOP_GAINERS_LOSERS"
	path := t.createQuerytUrl(dic)
	result := &TopGainersLosers{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
