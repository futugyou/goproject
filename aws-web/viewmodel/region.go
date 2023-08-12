package viewmodel

type AwsRegion struct {
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
	Region   string `json:"region"`
}
