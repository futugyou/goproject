package commodities

type CommoditiesEntity struct {
	Id       string `bson:"_id"`
	Name     string `bson:"name"`
	DataType string `bson:"type"`
	Interval string `bson:"interval"`
	Unit     string `bson:"unit"`
	Date     string `bson:"date"`
	Value    string `bson:"value"`
}

func (CommoditiesEntity) GetType() string {
	return "commoditiess"
}
