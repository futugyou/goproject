package balance

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncBalanceSheetData(symbol string) {
	log.Printf("%s balance sentiment data sync start. \n", symbol)
	// get Balance data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewFundamentalsClient(apikey)
	p := alphavantage.BalanceSheetParameter{
		Symbol: symbol,
	}
	s, err := client.BalanceSheet(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}

	// create Balance list
	data := make([]BalanceEntity, 0)
	for ii := 0; ii < len(s.AnnualReports); ii++ {
		e := createBalanceEntity(s.AnnualReports[ii], symbol, "Annual")
		data = append(data, e)
	}
	for ii := 0; ii < len(s.QuarterlyReports); ii++ {
		e := createBalanceEntity(s.QuarterlyReports[ii], symbol, "Quarterly")
		data = append(data, e)
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewBalanceRepository(config)
	r, err := repository.InsertMany(context.Background(), data, BalanceFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("balance data sync finish")
}

func createBalanceEntity(Balance alphavantage.BalanceSheetReport, symbol string, dataType string) BalanceEntity {
	return BalanceEntity{
		Id:                                     symbol + "-" + dataType + "-" + Balance.FiscalDateEnding,
		Symbol:                                 symbol,
		DataType:                               dataType,
		FiscalDateEnding:                       Balance.FiscalDateEnding,
		ReportedCurrency:                       Balance.ReportedCurrency,
		TotalAssets:                            Balance.TotalAssets,
		TotalCurrentAssets:                     Balance.TotalCurrentAssets,
		CashAndCashEquivalentsAtCarryingValue:  Balance.CashAndCashEquivalentsAtCarryingValue,
		CashAndShortTermInvestments:            Balance.CashAndShortTermInvestments,
		Inventory:                              Balance.Inventory,
		CurrentNetReceivables:                  Balance.CurrentNetReceivables,
		TotalNonCurrentAssets:                  Balance.TotalNonCurrentAssets,
		PropertyPlantEquipment:                 Balance.PropertyPlantEquipment,
		AccumulatedDepreciationAmortizationPPE: Balance.AccumulatedDepreciationAmortizationPPE,
		IntangibleAssets:                       Balance.IntangibleAssets,
		IntangibleAssetsExcludingGoodwill:      Balance.IntangibleAssetsExcludingGoodwill,
		Goodwill:                               Balance.Goodwill,
		Investments:                            Balance.Investments,
		LongTermInvestments:                    Balance.LongTermInvestments,
		ShortTermInvestments:                   Balance.ShortTermInvestments,
		OtherCurrentAssets:                     Balance.OtherCurrentAssets,
		OtherNonCurrentAssets:                  Balance.OtherNonCurrentAssets,
		TotalLiabilities:                       Balance.TotalLiabilities,
		TotalCurrentLiabilities:                Balance.TotalCurrentLiabilities,
		CurrentAccountsPayable:                 Balance.CurrentAccountsPayable,
		DeferredRevenue:                        Balance.DeferredRevenue,
		CurrentDebt:                            Balance.CurrentDebt,
		ShortTermDebt:                          Balance.ShortTermDebt,
		TotalNonCurrentLiabilities:             Balance.TotalNonCurrentLiabilities,
		CapitalLeaseObligations:                Balance.CapitalLeaseObligations,
		LongTermDebt:                           Balance.LongTermDebt,
		CurrentLongTermDebt:                    Balance.CurrentLongTermDebt,
		LongTermDebtNoncurrent:                 Balance.LongTermDebtNoncurrent,
		ShortLongTermDebtTotal:                 Balance.ShortLongTermDebtTotal,
		OtherCurrentLiabilities:                Balance.OtherCurrentLiabilities,
		OtherNonCurrentLiabilities:             Balance.OtherNonCurrentLiabilities,
		TotalShareholderEquity:                 Balance.TotalShareholderEquity,
		TreasuryStock:                          Balance.TreasuryStock,
		RetainedEarnings:                       Balance.RetainedEarnings,
		CommonStock:                            Balance.CommonStock,
		CommonStockSharesOutstanding:           Balance.CommonStockSharesOutstanding,
	}
}

func BalanceFilter(e BalanceEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "dataType", Value: e.DataType}, {Key: "fiscalDateEnding", Value: e.FiscalDateEnding}}
}
