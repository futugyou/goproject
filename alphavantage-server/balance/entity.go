package balance

type BalanceEntity struct {
	Id                                     string `bson:"_id"`
	Symbol                                 string `bson:"symbol"`
	DataType                               string `bson:"dataType"`
	FiscalDateEnding                       string `bson:"fiscalDateEnding"`
	ReportedCurrency                       string `bson:"reportedCurrency"`
	TotalAssets                            string `bson:"totalAssets"`
	TotalCurrentAssets                     string `bson:"totalCurrentAssets"`
	CashAndCashEquivalentsAtCarryingValue  string `bson:"cashAndCashEquivalentsAtCarryingValue"`
	CashAndShortTermInvestments            string `bson:"cashAndShortTermInvestments"`
	Inventory                              string `bson:"inventory"`
	CurrentNetReceivables                  string `bson:"currentNetReceivables"`
	TotalNonCurrentAssets                  string `bson:"totalNonCurrentAssets"`
	PropertyPlantEquipment                 string `bson:"propertyPlantEquipment"`
	AccumulatedDepreciationAmortizationPPE string `bson:"accumulatedDepreciationAmortizationPPE"`
	IntangibleAssets                       string `bson:"intangibleAssets"`
	IntangibleAssetsExcludingGoodwill      string `bson:"intangibleAssetsExcludingGoodwill"`
	Goodwill                               string `bson:"goodwill"`
	Investments                            string `bson:"investments"`
	LongTermInvestments                    string `bson:"longTermInvestments"`
	ShortTermInvestments                   string `bson:"shortTermInvestments"`
	OtherCurrentAssets                     string `bson:"otherCurrentAssets"`
	OtherNonCurrentAssets                  string `bson:"otherNonCurrentAssets"`
	TotalLiabilities                       string `bson:"totalLiabilities"`
	TotalCurrentLiabilities                string `bson:"totalCurrentLiabilities"`
	CurrentAccountsPayable                 string `bson:"currentAccountsPayable"`
	DeferredRevenue                        string `bson:"deferredRevenue"`
	CurrentDebt                            string `bson:"currentDebt"`
	ShortTermDebt                          string `bson:"shortTermDebt"`
	TotalNonCurrentLiabilities             string `bson:"totalNonCurrentLiabilities"`
	CapitalLeaseObligations                string `bson:"capitalLeaseObligations"`
	LongTermDebt                           string `bson:"longTermDebt"`
	CurrentLongTermDebt                    string `bson:"currentLongTermDebt"`
	LongTermDebtNoncurrent                 string `bson:"longTermDebtNoncurrent"`
	ShortLongTermDebtTotal                 string `bson:"shortLongTermDebtTotal"`
	OtherCurrentLiabilities                string `bson:"otherCurrentLiabilities"`
	OtherNonCurrentLiabilities             string `bson:"otherNonCurrentLiabilities"`
	TotalShareholderEquity                 string `bson:"totalShareholderEquity"`
	TreasuryStock                          string `bson:"treasuryStock"`
	RetainedEarnings                       string `bson:"retainedEarnings"`
	CommonStock                            string `bson:"commonStock"`
	CommonStockSharesOutstanding           string `bson:"commonStockSharesOutstanding"`
}

func (BalanceEntity) GetType() string {
	return "Balances"
}
