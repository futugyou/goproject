package commodities

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage/enums"
)

func CreateCommoditiesIndex(ctx context.Context) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	if err := repository.CreateIndex(ctx); err != nil {
		log.Println(err)
	}
}

func SyncDailyCommoditiesData(ctx context.Context) {
	log.Println("commodities daily data sync start")
	// get commodities data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewCommoditiesClient(apikey)

	// get daily commodities data
	data := wti(client)
	data = append(data, brent(client)...)
	data = append(data, gas(client)...)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("commodities daily data sync finish")
}

func SyncMonthlyCommoditiesData(ctx context.Context) {
	log.Println("commodities monthly data sync start")
	// get commodities data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewCommoditiesClient(apikey)

	// get monthly commodities data
	data := copper(client)
	data = append(data, copper(client)...)
	data = append(data, aluminum(client)...)
	data = append(data, wheat(client)...)
	data = append(data, corn(client)...)
	data = append(data, cotton(client)...)
	data = append(data, sugar(client)...)
	data = append(data, coffee(client)...)
	data = append(data, all(client)...)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("commodities monthly data sync finish")
}

func convertData(name string, interval string, unit string, datatype string, d []alphavantage.Datum) []CommoditiesEntity {
	data := make([]CommoditiesEntity, 0)
	for _, v := range d {
		data = append(data, CommoditiesEntity{
			Id:       datatype + v.Date,
			Name:     name,
			Interval: interval,
			Unit:     unit,
			Date:     v.Date,
			Value:    v.Value,
			DataType: datatype,
		})
	}
	return data
}

func wti(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get wti from alphavantage")
	p := alphavantage.CrudeOilWtiParameter{
		Interval: enums.CommoditiesDaily,
	}
	s, err := client.CrudeOilWti(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "wti", s.Data)
}

func brent(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get brent from alphavantage")
	p := alphavantage.CrudeOilBrentParameter{
		Interval: enums.CommoditiesDaily,
	}
	s, err := client.CrudeOilBrent(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "brent", s.Data)
}

func gas(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get gas from alphavantage")
	p := alphavantage.NaturalGasParameter{
		Interval: enums.CommoditiesDaily,
	}
	s, err := client.NaturalGas(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "gas", s.Data)
}

func copper(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get copper from alphavantage")
	p := alphavantage.CopperParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Copper(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "copper", s.Data)
}

func aluminum(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get aluminum from alphavantage")
	p := alphavantage.AluminumParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Aluminum(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "aluminum", s.Data)
}

func wheat(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get wheat from alphavantage")
	p := alphavantage.WheatParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Wheat(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "wheat", s.Data)
}

func corn(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get corn from alphavantage")
	p := alphavantage.CornParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Corn(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "corn", s.Data)
}

func cotton(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get cotton from alphavantage")
	p := alphavantage.CottonParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Cotton(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "cotton", s.Data)
}

func sugar(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get sugar from alphavantage")
	p := alphavantage.SugarParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Sugar(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "sugar", s.Data)
}

func coffee(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get coffee from alphavantage")
	p := alphavantage.CoffeeParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.Coffee(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "coffee", s.Data)
}

func all(client *alphavantage.CommoditiesClient) []CommoditiesEntity {
	log.Println("get all from alphavantage")
	p := alphavantage.AllCommoditiesParameter{
		Interval: enums.CommoditiesMonthly2,
	}
	s, err := client.AllCommodities(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}
	return convertData(s.Name, s.Interval, s.Unit, "all", s.Data)
}

func SyncDailyEconomicData(ctx context.Context) {
	log.Println("economic daily data sync start")
	// get economic data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewEconomicIndicatorsClient(apikey)

	// get daily economic data
	data := treasury(client)
	data = append(data, interest(client)...)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("economic daily data sync finish")
}

func SyncMonthlyEconomicData(ctx context.Context) {
	log.Println("economic monthly data sync start")
	// get economic data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewEconomicIndicatorsClient(apikey)

	// get monthly economic data
	data := cpi(client)
	data = append(data, retail(client)...)
	data = append(data, durable(client)...)
	data = append(data, unemployment(client)...)
	data = append(data, payroll(client)...)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("economic monthly data sync finish")
}

func SyncQuarterlyEconomicData(ctx context.Context) {
	log.Println("economic quarterly data sync start")
	// get economic data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewEconomicIndicatorsClient(apikey)

	// get quarterly economic data
	data := realgdp(client)
	data = append(data, realgdpcapita(client)...)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("economic quarterly data sync finish")
}

func SyncAnnualEconomicData(ctx context.Context) {
	log.Println("economic annual data sync start")
	// get economic data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewEconomicIndicatorsClient(apikey)

	// get annual economic data
	data := inflation(client)

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewCommoditiesRepository(config)
	r, err := repository.InsertMany(ctx, data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("economic annual data sync finish")
}

func realgdp(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get real gdp from alphavantage")
	p := alphavantage.RealGdpParameter{
		Interval: enums.EconomicGdpQuarterly,
	}
	s, err := client.RealGdp(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "realgdp", s.Data)
}

func realgdpcapita(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get real-gdp-per-capita from alphavantage")
	s, err := client.RealGdpPerCapita()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "realgdpcapita", s.Data)
}

func treasury(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get treasury from alphavantage")
	p := alphavantage.TreasuryYieldParameter{
		Interval: enums.EconomicTreasuryDaily,
		Maturity: enums.M5year,
	}
	s, err := client.TreasuryYield(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "treasury", s.Data)
}

func interest(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get interest from alphavantage")
	p := alphavantage.InterestRateParameter{
		Interval: enums.EconomicFundsDaily,
	}
	s, err := client.InterestRate(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "interest", s.Data)
}

func cpi(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get cpi from alphavantage")
	p := alphavantage.CPIParameter{
		Interval: enums.EconomicCPIMonthly,
	}
	s, err := client.CPI(p)
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "cpi", s.Data)
}

func inflation(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get inflation from alphavantage")
	s, err := client.Inflation()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "inflation", s.Data)
}

func retail(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get retail-sales from alphavantage")
	s, err := client.RetailSales()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "retail", s.Data)
}

func durable(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get durable-goods from alphavantage")
	s, err := client.DurableGoods()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "durable", s.Data)
}

func unemployment(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get unemployment from alphavantage")
	s, err := client.UnemploymentRate()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "unemployment", s.Data)
}

func payroll(client *alphavantage.EconomicIndicatorsClient) []CommoditiesEntity {
	log.Println("get nonfarm-payroll from alphavantage")
	s, err := client.NonfarmPayroll()
	if err != nil || s == nil {
		log.Println(err)
		return []CommoditiesEntity{}
	}

	return convertData(s.Name, s.Interval, s.Unit, "payroll", s.Data)
}

func CommoditiesFilter(e CommoditiesEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "type", Value: e.DataType}, {Key: "date", Value: e.Date}, {Key: "interval", Value: e.Interval}}
}
