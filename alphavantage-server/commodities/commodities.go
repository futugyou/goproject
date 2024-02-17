package commodities

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage/enums"
)

func SyncAllCommoditiesData() {
	log.Println("commodities data sync start")
	// get commodities data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewCommoditiesClient(apikey)

	// get all commodities data
	data := wti(client)
	data = append(data, brent(client)...)
	data = append(data, gas(client)...)
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
	r, err := repository.InsertMany(context.Background(), data, CommoditiesFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("commodities data sync finish")
}

func convertData(name string, interval string, unit string, d []alphavantage.Datum) []CommoditiesEntity {
	data := make([]CommoditiesEntity, 0)
	for _, v := range d {
		data = append(data, CommoditiesEntity{
			Id:       name + v.Date,
			Name:     name,
			Interval: interval,
			Unit:     unit,
			Date:     v.Date,
			Value:    v.Value,
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

	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
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
	return convertData(s.Name, s.Interval, s.Unit, s.Data)
}

func CommoditiesFilter(e CommoditiesEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "name", Value: e.Name}, {Key: "date", Value: e.Date}}
}
