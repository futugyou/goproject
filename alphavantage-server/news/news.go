package news

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncNewsSentimentData(symbol string) {
	log.Printf("%s news sentiment data sync start. \n", symbol)
	// get news data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewIntelligenceClient(apikey)
	t := time.Now()
	y, m, d := t.Date()
	p := alphavantage.SentimentParameter{
		Tickers:  symbol,
		TimeFrom: time.Date(y, m, d-7, 0, 0, 0, 1, t.Location()).Format("20060102T1504"),
		TimeTo:   time.Date(y, m, d, 0, 0, 0, -1, t.Location()).Format("20060102T1504"),
	}
	s, err := client.NewsSentiment(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}
	log.Printf("%s feed count is %d \n", symbol, len(s.Feed))
	// create news list
	data := make([]NewsEntity, 0)
	for ii := 0; ii < len(s.Feed); ii++ {
		timePublished, err := time.Parse("20060102T150405", s.Feed[ii].TimePublished)
		if err != nil {
			log.Println(err)
			continue
		}
		e := NewsEntity{
			Id:                    s.Feed[ii].Title + s.Feed[ii].TimePublished,
			Title:                 s.Feed[ii].Title,
			URL:                   s.Feed[ii].URL,
			TimePublished:         timePublished,
			Authors:               s.Feed[ii].Authors,
			Summary:               s.Feed[ii].Summary,
			BannerImage:           s.Feed[ii].BannerImage,
			Source:                s.Feed[ii].Source,
			CategoryWithinSource:  s.Feed[ii].CategoryWithinSource,
			SourceDomain:          s.Feed[ii].SourceDomain,
			Topics:                getTopics(s.Feed[ii].Topics),
			OverallSentimentScore: s.Feed[ii].OverallSentimentScore,
			OverallSentimentLabel: s.Feed[ii].OverallSentimentLabel,
			TickerSentiment:       getTickerSentiment(s.Feed[ii].TickerSentiment),
		}
		data = append(data, e)
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewNewsRepository(config)
	r, err := repository.InsertMany(context.Background(), data, NewsFilter)
	if err != nil {
		log.Println(err)
		return
	}
	
	r.String() 
	log.Println("news sentiment data sync finish")
}

func getTickerSentiment(tickerSentiment []alphavantage.TickerSentiment) []string {
	r := make([]string, 0)
	for _, v := range tickerSentiment {
		r = append(r, v.Ticker)
	}
	return r
}

func getTopics(topic []alphavantage.Topic) []string {
	r := make([]string, 0)
	for _, v := range topic {
		r = append(r, v.Topic)
	}
	return r
}

func NewsFilter(e NewsEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "title", Value: e.Title}}
}
