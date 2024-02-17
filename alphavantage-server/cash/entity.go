package cash

type CashEntity struct {
	Id                                                        string `bson:"_id"`
	Symbol                                                    string `bson:"symbol"`
	DataType                                                  string `bson:"dataType"`
	FiscalDateEnding                                          string `bson:"fiscalDateEnding"`
	ReportedCurrency                                          string `bson:"reportedCurrency"`
	OperatingCashflow                                         string `bson:"operatingCashflow"`
	PaymentsForOperatingActivities                            string `bson:"paymentsForOperatingActivities"`
	ProceedsFromOperatingActivities                           string `bson:"proceedsFromOperatingActivities"`
	ChangeInOperatingLiabilities                              string `bson:"changeInOperatingLiabilities"`
	ChangeInOperatingAssets                                   string `bson:"changeInOperatingAssets"`
	DepreciationDepletionAndAmortization                      string `bson:"depreciationDepletionAndAmortization"`
	CapitalExpenditures                                       string `bson:"capitalExpenditures"`
	ChangeInReceivables                                       string `bson:"changeInReceivables"`
	ChangeInInventory                                         string `bson:"changeInInventory"`
	ProfitLoss                                                string `bson:"profitLoss"`
	CashflowFromInvestment                                    string `bson:"cashflowFromInvestment"`
	CashflowFromFinancing                                     string `bson:"cashflowFromFinancing"`
	ProceedsFromRepaymentsOfShortTermDebt                     string `bson:"proceedsFromRepaymentsOfShortTermDebt"`
	PaymentsForRepurchaseOfCommonStock                        string `bson:"paymentsForRepurchaseOfCommonStock"`
	PaymentsForRepurchaseOfEquity                             string `bson:"paymentsForRepurchaseOfEquity"`
	PaymentsForRepurchaseOfPreferredStock                     string `bson:"paymentsForRepurchaseOfPreferredStock"`
	DividendPayout                                            string `bson:"dividendPayout"`
	DividendPayoutCommonStock                                 string `bson:"dividendPayoutCommonStock"`
	DividendPayoutPreferredStock                              string `bson:"dividendPayoutPreferredStock"`
	ProceedsFromIssuanceOfCommonStock                         string `bson:"proceedsFromIssuanceOfCommonStock"`
	ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `bson:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
	ProceedsFromIssuanceOfPreferredStock                      string `bson:"proceedsFromIssuanceOfPreferredStock"`
	ProceedsFromRepurchaseOfEquity                            string `bson:"proceedsFromRepurchaseOfEquity"`
	ProceedsFromSaleOfTreasuryStock                           string `bson:"proceedsFromSaleOfTreasuryStock"`
	ChangeInCashAndCashEquivalents                            string `bson:"changeInCashAndCashEquivalents"`
	ChangeInExchangeRate                                      string `bson:"changeInExchangeRate"`
	NetIncome                                                 string `bson:"netIncome"`
}

func (CashEntity) GetType() string {
	return "cashs"
}
