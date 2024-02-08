package enums

type SentimentTopics interface {
	private()
	String() string
}

type sentimentTopics string

func (c sentimentTopics) private() {}
func (c sentimentTopics) String() string {
	return string(c)
}

const Blockchain sentimentTopics = "blockchain"
const Earnings sentimentTopics = "earnings"
const IPO sentimentTopics = "ipo"
const MergersAcquisitions sentimentTopics = "mergers_and_acquisitions"
const FinancialMarkets sentimentTopics = "financial_markets"
const EconomyFiscalPolicy sentimentTopics = "economy_fiscal"
const EconomyMonetaryPolicy sentimentTopics = "economy_monetary"
const EconomyMacroOverall sentimentTopics = "economy_macro"
const EnergyTransportation sentimentTopics = "energy_transportation"
const Finance sentimentTopics = "finance"
const LifeSciences sentimentTopics = "life_sciences"
const Manufacturing sentimentTopics = "manufacturing"
const RealEstateConstruction sentimentTopics = "real_estate"
const RetailWholesale sentimentTopics = "retail_wholesale"
const Technology sentimentTopics = "technology"
