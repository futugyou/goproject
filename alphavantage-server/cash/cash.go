package cash

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncCashSheetData(symbol string) {
	log.Printf("%s cash sentiment data sync start. \n", symbol)
	// get Cash data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewFundamentalsClient(apikey)
	p := alphavantage.CashFlowParameter{
		Symbol: symbol,
	}
	s, err := client.CashFlow(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}

	// create Cash list
	data := make([]CashEntity, 0)
	for ii := 0; ii < len(s.AnnualReports); ii++ {
		e := createCashEntity(s.AnnualReports[ii], symbol, "Annual")
		data = append(data, e)
	}
	for ii := 0; ii < len(s.QuarterlyReports); ii++ {
		e := createCashEntity(s.QuarterlyReports[ii], symbol, "Quarterly")
		data = append(data, e)
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCashRepository(config)
	r, err := repository.InsertMany(context.Background(), data, CashFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("Cash sentiment data sync finish")
}

func createCashEntity(cash alphavantage.CashFlowReport, symbol string, dataType string) CashEntity {
	return CashEntity{
		Id:                                    symbol + "-" + dataType + "-" + cash.FiscalDateEnding,
		Symbol:                                symbol,
		DataType:                              dataType,
		FiscalDateEnding:                      cash.FiscalDateEnding,
		ReportedCurrency:                      cash.ReportedCurrency,
		OperatingCashflow:                     cash.OperatingCashflow,
		PaymentsForOperatingActivities:        cash.PaymentsForOperatingActivities,
		ProceedsFromOperatingActivities:       cash.ProceedsFromOperatingActivities,
		ChangeInOperatingLiabilities:          cash.ChangeInOperatingLiabilities,
		ChangeInOperatingAssets:               cash.ChangeInOperatingAssets,
		DepreciationDepletionAndAmortization:  cash.DepreciationDepletionAndAmortization,
		CapitalExpenditures:                   cash.CapitalExpenditures,
		ChangeInReceivables:                   cash.ChangeInReceivables,
		ChangeInInventory:                     cash.ChangeInInventory,
		ProfitLoss:                            cash.ProfitLoss,
		CashflowFromInvestment:                cash.CashflowFromInvestment,
		CashflowFromFinancing:                 cash.CashflowFromFinancing,
		ProceedsFromRepaymentsOfShortTermDebt: cash.ProceedsFromRepaymentsOfShortTermDebt,
		PaymentsForRepurchaseOfCommonStock:    cash.PaymentsForRepurchaseOfCommonStock,
		PaymentsForRepurchaseOfEquity:         cash.PaymentsForRepurchaseOfEquity,
		PaymentsForRepurchaseOfPreferredStock: cash.PaymentsForRepurchaseOfPreferredStock,
		DividendPayout:                        cash.DividendPayout,
		DividendPayoutCommonStock:             cash.DividendPayoutCommonStock,
		DividendPayoutPreferredStock:          cash.DividendPayoutPreferredStock,
		ProceedsFromIssuanceOfCommonStock:     cash.ProceedsFromIssuanceOfCommonStock,
		ProceedsFromIssuanceOfPreferredStock:  cash.ProceedsFromIssuanceOfPreferredStock,
		ProceedsFromRepurchaseOfEquity:        cash.ProceedsFromRepurchaseOfEquity,
		ProceedsFromSaleOfTreasuryStock:       cash.ProceedsFromSaleOfTreasuryStock,
		ChangeInCashAndCashEquivalents:        cash.ChangeInCashAndCashEquivalents,
		ChangeInExchangeRate:                  cash.ChangeInExchangeRate,
		NetIncome:                             cash.NetIncome,
		ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet: cash.ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet,
	}
}

func CashFilter(e CashEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "dataType", Value: e.DataType}, {Key: "fiscalDateEnding", Value: e.FiscalDateEnding}}
}
