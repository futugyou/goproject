package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage/enums"
	"github.com/futugyou/alphavantage/functions"
)

func main() {
	// StockDataAPIs()
	// AlphaIntelligence()
	// Fundamentals()
	// ForeignExchangeRates()
	// DigitalCurrency()
	// Commodities()
	EconomicIndicators()
}

func EconomicIndicators() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewEconomicIndicatorsClient(apikey)
	p := alphavantage.EconomicIndicatorsParameter{
		Interval: enums.LWeekly,
		Function: functions.RealGDP,
	}

	result, err := s.GetEconomicIndicatorsDara(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result.Name, len(result.Data))
}

func Commodities() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewCommoditiesClient(apikey)
	// CrudeOilWti(s)
	// CrudeOilBrent(s)
	AllCommodities(s)
}

func AllCommodities(s *alphavantage.CommoditiesClient) {
	p := alphavantage.AllCommoditiesParameter{
		Interval: enums.LWeekly,
	}

	result, err := s.AllCommodities(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result.Name, len(result.Data))
}

func CrudeOilBrent(s *alphavantage.CommoditiesClient) {
	p := alphavantage.CrudeOilBrentParameter{
		Interval: enums.LWeekly,
	}

	result, err := s.CrudeOilBrent(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result.Name, len(result.Data))
}

func CrudeOilWti(s *alphavantage.CommoditiesClient) {
	p := alphavantage.CrudeOilWtiParameter{
		Interval: enums.LWeekly,
	}

	result, err := s.CrudeOilWti(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result.Name, len(result.Data))
}

func DigitalCurrency() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewDigitalCurrencyClient(apikey)
	// CryptoExchange(s)
	CryptoIntraday(s)
	// CurrencyDaily(s)
	// CurrencyWeekly(s)
	// CurrencyMonthly(s)
}

func CurrencyMonthly(s *alphavantage.DigitalCurrencyClient) {
	p := alphavantage.CurrencyMonthlyParameter{
		Symbol: "BTC",
		Market: "AFN",
	}

	result, err := s.CurrencyMonthly(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.Symbol, v.Market, v.MarketOpen, v.MarketHigh, v.MarketLow, v.MarketOpen,
			v.USDOpen, v.USDHigh, v.USDLow, v.USDOpen, v.Timestamp, v.USDmarketCap, v.Volume)
	}
}

func CurrencyWeekly(s *alphavantage.DigitalCurrencyClient) {
	p := alphavantage.CurrencyWeeklyParameter{
		Symbol: "BTC",
		Market: "AFN",
	}

	result, err := s.CurrencyWeekly(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.Symbol, v.Market, v.MarketOpen, v.MarketHigh, v.MarketLow, v.MarketOpen,
			v.USDOpen, v.USDHigh, v.USDLow, v.USDOpen, v.Timestamp, v.USDmarketCap, v.Volume)
	}
}

func CurrencyDaily(s *alphavantage.DigitalCurrencyClient) {
	p := alphavantage.CurrencyDailyParameter{
		Symbol: "BTC",
		Market: "AFN",
	}

	result, err := s.CurrencyDaily(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.Symbol, v.Market, v.MarketOpen, v.MarketHigh, v.MarketLow, v.MarketOpen,
			v.USDOpen, v.USDHigh, v.USDLow, v.USDOpen, v.Timestamp, v.USDmarketCap, v.Volume)
	}
}

func CryptoIntraday(s *alphavantage.DigitalCurrencyClient) {
	p := alphavantage.CryptoIntradayParameter{
		Symbol:   "ETH",
		Market:   "USD",
		Interval: enums.T15min,
	}

	result, err := s.CryptoIntraday(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.Symbol, v.Market, v.Close, v.High, v.Low, v.Open, v.Timestamp)
	}
}

func CryptoExchange(s *alphavantage.DigitalCurrencyClient) {
	p := alphavantage.CryptoExchangeParameter{
		FromCurrency: "BTC",
		ToCurrency:   "CNY",
	}

	result, err := s.CryptoExchange(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.DigitalExchangeRate.FromCurrencyCode, result.DigitalExchangeRate.ToCurrencyCode,
		result.DigitalExchangeRate.AskPrice, result.DigitalExchangeRate.BidPrice,
		result.DigitalExchangeRate.ExchangeRate, result.DigitalExchangeRate.LastRefreshed, result.DigitalExchangeRate.TimeZone)
}

func ForeignExchangeRates() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewForeignExchangeRatesClient(apikey)
	// CurrencyExchange(s)
	// FxIntraday(s)
	// FxDaily(s)
	// FxWeekly(s)
	FxMonthly(s)
}

func FxMonthly(s *alphavantage.ForeignExchangeRatesClient) {
	p := alphavantage.FxMonthlyParameter{
		FromSymbol: "EUR",
		ToSymbol:   "USD",
	}

	result, err := s.FxMonthly(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.FromSymbol, v.ToSymbol, v.Close, v.High, v.Low, v.Open, v.Timestamp)
	}
}

func FxWeekly(s *alphavantage.ForeignExchangeRatesClient) {
	p := alphavantage.FxWeeklyParameter{
		FromSymbol: "EUR",
		ToSymbol:   "USD",
	}

	result, err := s.FxWeekly(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.FromSymbol, v.ToSymbol, v.Close, v.High, v.Low, v.Open, v.Timestamp)
	}
}

func FxDaily(s *alphavantage.ForeignExchangeRatesClient) {
	p := alphavantage.FxDailyParameter{
		FromSymbol: "EUR",
		ToSymbol:   "USD",
	}

	result, err := s.FxDaily(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.FromSymbol, v.ToSymbol, v.Close, v.High, v.Low, v.Open, v.Timestamp)
	}
}

func FxIntraday(s *alphavantage.ForeignExchangeRatesClient) {
	p := alphavantage.FxIntradayParameter{
		FromSymbol: "EUR",
		ToSymbol:   "USD",
		Interval:   enums.T5min,
	}

	result, err := s.FxIntraday(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result {
		fmt.Println(v.FromSymbol, v.ToSymbol, v.Close, v.High, v.Low, v.Open, v.Timestamp)
	}
}

func CurrencyExchange(s *alphavantage.ForeignExchangeRatesClient) {
	p := alphavantage.CurrencyExchangeParameter{
		FromCurrency: "USD",
		ToCurrency:   "JPY",
	}

	result, err := s.CurrencyExchange(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.ForeignExchangeRate.FromCurrencyCode, result.ForeignExchangeRate.ToCurrencyCode)
}

func Fundamentals() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewFundamentalsClient(apikey)
	// CompanyOverview(s)
	// IncomeStatement(s)
	// BalanceSheet(s)
	// CashFlow(s)
	// Earnings(s)
	// ListingStatus(s)
	// EarningsCalendar(s)
	IpoCalendar(s)
}

func CompanyOverview(s *alphavantage.FundamentalsClient) {
	p := alphavantage.CompanyOverviewParameter{
		Symbol: "IBM",
	}
	result, err := s.CompanyOverview(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Address, result.AssetType)
}

func IncomeStatement(s *alphavantage.FundamentalsClient) {
	p := alphavantage.IncomeStatementParameter{
		Symbol: "IBM",
	}
	result, err := s.IncomeStatement(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Symbol)
	for _, v := range result.AnnualReports {
		fmt.Println(v.CostOfRevenue)
	}
	for _, v := range result.QuarterlyReports {
		fmt.Println(v.ComprehensiveIncomeNetOfTax)
	}
}

func BalanceSheet(s *alphavantage.FundamentalsClient) {
	p := alphavantage.BalanceSheetParameter{
		Symbol: "IBM",
	}
	result, err := s.BalanceSheet(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Symbol)
	for _, v := range result.AnnualReports {
		fmt.Println(v.AccumulatedDepreciationAmortizationPPE)
	}
	for _, v := range result.QuarterlyReports {
		fmt.Println(v.CashAndCashEquivalentsAtCarryingValue)
	}
}

func CashFlow(s *alphavantage.FundamentalsClient) {
	p := alphavantage.CashFlowParameter{
		Symbol: "IBM",
	}
	result, err := s.CashFlow(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Symbol)
	for _, v := range result.AnnualReports {
		fmt.Println(v.CapitalExpenditures)
	}
	for _, v := range result.QuarterlyReports {
		fmt.Println(v.ChangeInCashAndCashEquivalents)
	}
}

func Earnings(s *alphavantage.FundamentalsClient) {
	p := alphavantage.EarningsParameter{
		Symbol: "IBM",
	}
	result, err := s.Earnings(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Symbol)
	for _, v := range result.AnnualEarnings {
		fmt.Println(v.FiscalDateEnding, v.ReportedEPS)
	}
	for _, v := range result.QuarterlyEarnings {
		fmt.Println(v.FiscalDateEnding, v.ReportedEPS)
	}
}

func ListingStatus(s *alphavantage.FundamentalsClient) {
	p := alphavantage.ListingStatusParameter{}
	result, err := s.ListingStatus(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Symbol, v.DelistingDate, v.IpoDate)
	}
}

func EarningsCalendar(s *alphavantage.FundamentalsClient) {
	p := alphavantage.EarningsCalendarParameter{}
	result, err := s.EarningsCalendar(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Symbol, v.Name, v.ReportDate, v.FiscalDateEnding, v.Estimate, v.Currency)
	}
}

func IpoCalendar(s *alphavantage.FundamentalsClient) {
	result, err := s.IpoCalendar()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Symbol, v.Name, v.Exchange, v.IpoDate, v.PriceRangeHigh, v.PriceRangeLow, v.Currency)
	}
}

func AlphaIntelligence() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	s := alphavantage.NewIntelligenceClient(apikey)
	// NewsSentiment(s)
	TopGainersLosers(s)
}

func TopGainersLosers(s *alphavantage.IntelligenceClient) {
	result, err := s.TopGainersLosers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.LastUpdated, result.Metadata)
	for _, vv := range result.MostActivelyTraded {
		fmt.Println(vv.Volume, vv.ChangeAmount, vv.ChangePercentage, vv.Price, vv.Ticker)
	}
}

func NewsSentiment(s *alphavantage.IntelligenceClient) {
	p := alphavantage.SentimentParameter{
		Tickers: "IBM",
	}
	result, err := s.NewsSentiment(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result.Items, result.RelevanceScoreDefinition, result.SentimentScoreDefinition)
	for _, vv := range result.Feed {
		fmt.Println(vv.BannerImage, vv.CategoryWithinSource, vv.OverallSentimentLabel, vv.OverallSentimentScore, vv.Source, vv.SourceDomain)
	}
}

func StockDataAPIs() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")

	dic := make(map[string]string)
	dic["month"] = "2024-01"

	s := alphavantage.NewTimeSeriesClient(apikey)

	// TimeSeries(s, dic)
	// TimeSeriesAdjusted(s, dic)
	// GlobalQuote(s, dic)
	// SymbolSearch(s, dic)
	MarketStatus(s, dic)
}

func TimeSeries(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	p := alphavantage.TimeSeriesIntradayParameter{
		Symbol:     "IBM",
		Interval:   enums.T15min,
		Dictionary: dic,
	}
	result, err := s.TimeSeriesIntraday(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Symbol, v.Time, v.Open, v.High, v.Low, v.Close, v.Volume)
	}
}

func TimeSeriesAdjusted(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	pp := alphavantage.TimeSeriesMonthlyAdjustedParameter{
		Symbol:     "IBM",
		Dictionary: dic,
	}

	result1, err := s.TimeSeriesMonthlyAdjusted(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Time, v.Open, v.High, v.Low, v.Close, v.Volume, v.AdjustedClose, v.DividendAmount, v.SplitCoefficient)
	}
}

func GlobalQuote(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	pp := alphavantage.GlobalQuoteParameter{
		Symbol: "IBM",
	}

	result1, err := s.GlobalQuote(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Open, v.High, v.Low, v.Price, v.Volume, v.LatestDay, v.PreviousClose, v.Change, v.ChangePercent)
	}
}

func SymbolSearch(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	pp := alphavantage.SymbolSearchParameter{
		Keywords: "IBM",
	}

	result1, err := s.SymbolSearch(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Currency, v.MarketClose, v.MarketOpen, v.MatchScore, v.Name, v.Region, v.Timezone, v.Type)
	}
}

func MarketStatus(s *alphavantage.TimeSeriesClient, dic map[string]string) {

	v, err := s.MarketStatus()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(v.Endpoint)
	for _, vv := range v.Markets {
		fmt.Println(vv.CurrentStatus, vv.LocalClose, vv.LocalOpen, vv.MarketType, vv.Notes, vv.Notes, vv.PrimaryExchanges, vv.Region)
	}
}
