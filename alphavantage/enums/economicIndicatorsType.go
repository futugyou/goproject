package enums

type EconomicIndicatorsType interface {
	private()
	String() string
}

type economicIndicatorsType string

func (c economicIndicatorsType) private() {}
func (c economicIndicatorsType) String() string {
	return string(c)
}

const RealGDP economicIndicatorsType = "REAL_GDP"
const RealGDPPerCapita economicIndicatorsType = "REAL_GDP_PER_CAPITA"
const TreasuryYield economicIndicatorsType = "TREASURY_YIELD"
const FederalFunds economicIndicatorsType = "FEDERAL_FUNDS_RATE"
const CPI economicIndicatorsType = "CPI"
const Inflation economicIndicatorsType = "INFLATION"
const RetailSales economicIndicatorsType = "RETAIL_SALES"
const DurableGoods economicIndicatorsType = "DURABLES"
const UnemploymentRate economicIndicatorsType = "UNEMPLOYMENT"
const NonfarmPayroll economicIndicatorsType = "NONFARM_PAYROLL"
