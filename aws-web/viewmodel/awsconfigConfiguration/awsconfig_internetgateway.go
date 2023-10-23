package awsconfigConfiguration

type InternetGatewayConfiguration struct {
	Attachments       []Attachment `json:"attachments"`
	InternetGatewayID string       `json:"internetGatewayId"`
	OwnerID           string       `json:"ownerId"`
	Tags              []Tag        `json:"tags"`
}

type Attachment struct {
	State string `json:"state"`
	VpcID string `json:"vpcId"`
}
