package awsconfigConfiguration

type ServiceDiscoveryConfiguration struct {
	ID          string    `json:"Id"`
	NamespaceID string    `json:"NamespaceId"`
	Arn         string    `json:"Arn"`
	Name        string    `json:"Name"`
	Type        string    `json:"Type"`
	Description string    `json:"Description"`
	DNSConfig   DNSConfig `json:"DnsConfig"`
}

type DNSConfig struct {
	NamespaceID   string      `json:"NamespaceId"`
	RoutingPolicy string      `json:"RoutingPolicy"`
	DNSRecords    []DNSRecord `json:"DnsRecords"`
}
