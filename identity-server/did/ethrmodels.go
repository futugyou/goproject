package did

import (
	"time"
)

type DecentralizedIdentifierEthr struct {
	Scheme     string
	Method     string
	Identifier string
	Fragment   string
	Network    string
	Address    string
	PublicKey  string
	Version    int
}

func (d DecentralizedIdentifierEthr) GetDidWithoutFragment() string {
	return d.Scheme + ":" + d.Method + ":" + d.Identifier
}

func (d DecentralizedIdentifierEthr) Format() string {
	result := d.Scheme + ":" + d.Method
	if len(d.Network) > 0 {
		result = result + ":" + d.Network
	}

	return result + ":" + d.Identifier
}

type ERC1056Event struct {
	Identity       string
	Value          string
	Owner          string
	Delegate       string
	DelegateType   string
	ValidTo        int64
	PreviousChange int64
	BlockNumber    int64
}

type ERC1056EventType uint32

const (
	DIDAttributeChanged ERC1056EventType = iota
	DIDOwnerChanged
	DIDDelegateChanged
)

type NetworkConfiguration struct {
	Name           string
	RpcUrl         string
	ContractAdr    string
	CreateDateTime time.Time
	UpdateDateTime time.Time
}

type CAIP10BlockChainAccount struct {
	Namespace      string
	Reference      string
	AccountAddress string
}

func BuildEthereumMainet(accountAddress string) *CAIP10BlockChainAccount {
	return &CAIP10BlockChainAccount{
		AccountAddress: accountAddress,
		Namespace:      "eip155",
		Reference:      "1",
	}
}

func (c CAIP10BlockChainAccount) Format() string {
	if len(c.Reference) > 0 {
		return c.Namespace + ":" + c.Reference + ":" + c.AccountAddress
	}

	return c.Namespace + ":" + c.AccountAddress
}
