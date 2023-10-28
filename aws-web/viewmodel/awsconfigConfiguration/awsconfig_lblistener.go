package awsconfigConfiguration

type LoadBalancerListenerConfiguration struct {
	SSLPolicy       string          `json:"SslPolicy"`
	LoadBalancerArn string          `json:"LoadBalancerArn"`
	DefaultActions  []DefaultAction `json:"DefaultActions"`
	Port            int64           `json:"Port"`
	Certificates    []Certificate   `json:"Certificates"`
	Protocol        string          `json:"Protocol"`
	ListenerArn     string          `json:"ListenerArn"`
	AlpnPolicy      []interface{}   `json:"AlpnPolicy"`
}

type Certificate struct {
	CertificateArn string `json:"CertificateArn"`
}

type DefaultAction struct {
	Type                string              `json:"Type"`
	Order               int64               `json:"Order"`
	FixedResponseConfig FixedResponseConfig `json:"FixedResponseConfig"`
	TargetGroupArn      string              `json:"TargetGroupArn"`
	ForwardConfig       ForwardActionConfig `json:"ForwardConfig"`
}

type FixedResponseConfig struct {
	MessageBody string `json:"MessageBody"`
	StatusCode  string `json:"StatusCode"`
	ContentType string `json:"ContentType"`
}

type ForwardActionConfig struct {
	TargetGroups []TargetGroupTuple `json:"TargetGroups"`
}

type TargetGroupTuple struct {
	TargetGroupArn string `json:"TargetGroupArn"`
}
