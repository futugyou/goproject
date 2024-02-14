package news

import "time"

type NewsEntity struct {
	Id                    string    `bson:"_id"`
	Title                 string    `bson:"title"`
	URL                   string    `bson:"url"`
	TimePublished         time.Time `bson:"time_published"`
	Authors               []string  `bson:"authors"`
	Summary               string    `bson:"summary"`
	BannerImage           string    `bson:"banner_image"`
	Source                string    `bson:"source"`
	CategoryWithinSource  string    `bson:"category_within_source"`
	SourceDomain          string    `bson:"source_domain"`
	Topics                []string  `bson:"topics"`
	OverallSentimentScore float64   `bson:"overall_sentiment_score"`
	OverallSentimentLabel string    `bson:"overall_sentiment_label"`
	TickerSentiment       []string  `bson:"ticker_sentiment"`
}

func (NewsEntity) GetType() string {
	return "news"
}
