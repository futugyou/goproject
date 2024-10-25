package income

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncIncomeStatementData(ctx context.Context, symbol string) {
	log.Printf("%s income sentiment data sync start. \n", symbol)
	// get income data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewFundamentalsClient(apikey)
	p := alphavantage.IncomeStatementParameter{
		Symbol: symbol,
	}
	s, err := client.IncomeStatement(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}

	// create income list
	data := make([]IncomeEntity, 0)
	for ii := 0; ii < len(s.AnnualReports); ii++ {
		e := createIncomeEntity(s.AnnualReports[ii], symbol, "Annual")
		data = append(data, e)
	}
	for ii := 0; ii < len(s.QuarterlyReports); ii++ {
		e := createIncomeEntity(s.QuarterlyReports[ii], symbol, "Quarterly")
		data = append(data, e)
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewIncomeRepository(config)
	r, err := repository.InsertMany(ctx, data, incomeFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("income sentiment data sync finish")
}

func createIncomeEntity(income alphavantage.IncomeStatementReport, symbol string, dataType string) IncomeEntity {
	return IncomeEntity{
		Id:                                symbol + "-" + dataType + "-" + income.FiscalDateEnding,
		Symbol:                            symbol,
		DataType:                          dataType,
		FiscalDateEnding:                  income.FiscalDateEnding,
		ReportedCurrency:                  income.ReportedCurrency,
		GrossProfit:                       income.GrossProfit,
		TotalRevenue:                      income.TotalRevenue,
		CostOfRevenue:                     income.CostOfRevenue,
		CostofGoodsAndServicesSold:        income.CostofGoodsAndServicesSold,
		OperatingIncome:                   income.OperatingIncome,
		SellingGeneralAndAdministrative:   income.SellingGeneralAndAdministrative,
		ResearchAndDevelopment:            income.ResearchAndDevelopment,
		OperatingExpenses:                 income.OperatingExpenses,
		InvestmentIncomeNet:               income.InvestmentIncomeNet,
		NetInterestIncome:                 income.NetInterestIncome,
		InterestIncome:                    income.InterestIncome,
		InterestExpense:                   income.InterestExpense,
		NonInterestIncome:                 income.NonInterestIncome,
		OtherNonOperatingIncome:           income.OtherNonOperatingIncome,
		Depreciation:                      income.Depreciation,
		DepreciationAndAmortization:       income.DepreciationAndAmortization,
		IncomeBeforeTax:                   income.IncomeBeforeTax,
		IncomeTaxExpense:                  income.IncomeTaxExpense,
		InterestAndDebtExpense:            income.InterestAndDebtExpense,
		NetIncomeFromContinuingOperations: income.NetIncomeFromContinuingOperations,
		ComprehensiveIncomeNetOfTax:       income.ComprehensiveIncomeNetOfTax,
		Ebit:                              income.Ebit,
		Ebitda:                            income.Ebitda,
		NetIncome:                         income.NetIncome,
	}
}

func incomeFilter(e IncomeEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "dataType", Value: e.DataType}, {Key: "fiscalDateEnding", Value: e.FiscalDateEnding}}
}
