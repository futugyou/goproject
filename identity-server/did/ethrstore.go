package did

import (
	"context"
	"fmt"
)

// mock data
var Networks = map[string]NetworkConfiguration{
	"mainnet": {
		RpcUrl:      "https://mainnet.infura.io/v3/YOUR_KEY",
		ContractAdr: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b",
	},
	"sepolia": {
		RpcUrl:      "https://sepolia.infura.io/v3/YOUR_KEY",
		ContractAdr: "0x03d5003bf0e79c5f5223588f347eba39afbc3818",
	},
}

type NetworkConfigurationStore struct {
}

func (n *NetworkConfigurationStore) Get(ctx context.Context, name string) (*NetworkConfiguration, error) {
	config, ok := Networks[name]
	if !ok {
		return nil, fmt.Errorf("%s not support", name)
	}

	return &config, nil
}
