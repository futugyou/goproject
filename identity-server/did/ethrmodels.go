package did

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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

func DidEthrExtractor(did string) (*DecentralizedIdentifierEthr, error) {
	decentralizedIdentifier, err := DidExtractor(did)
	if err != nil {
		return nil, err
	}

	if decentralizedIdentifier.Method != "ethr" {
		return nil, fmt.Errorf("method must be equals to ethr")
	}

	splittedIdentifier := strings.Split(decentralizedIdentifier.Identifier, ":")
	if len(splittedIdentifier) > 2 {
		return nil, fmt.Errorf("The did identifier cannot contains more than 2 parts")
	}
	network := ""
	address := splittedIdentifier[0]
	if len(splittedIdentifier) == 2 {
		network = splittedIdentifier[0]
		address = splittedIdentifier[1]
	}
	version := 0
	re := regexp.MustCompile(`(\w|\d)*\?versionId=(\d)*`)

	if re.MatchString(address) {
		splittedAdr := strings.Split(address, "?")
		lastPart := splittedAdr[len(splittedAdr)-1]
		versionIdStr := strings.Replace(lastPart, "versionId=", "", 1)
		version, err = strconv.Atoi(versionIdStr)
		if err != nil {
			return nil, err
		}

		address = splittedAdr[0]
		decentralizedIdentifier.Identifier = strings.Split(decentralizedIdentifier.Identifier, "?")[0]
	}

	identifierPayload, err := HexToByteArray(address)
	if err != nil {
		return nil, err
	}

	pk := ""
	if len(address) > 42 {
		publicKey, err := crypto.UnmarshalPubkey(identifierPayload)
		if err != nil {
			return nil, err
		}
		pk = strings.Replace(address, "0x", "", 0)
		address = crypto.PubkeyToAddress(*publicKey).String()
	}

	return &DecentralizedIdentifierEthr{
		Scheme:     decentralizedIdentifier.Scheme,
		Method:     decentralizedIdentifier.Method,
		Identifier: decentralizedIdentifier.Identifier,
		Fragment:   decentralizedIdentifier.Fragment,
		Address:    address,
		Network:    network,
		PublicKey:  pk,
		Version:    version,
	}, nil
}
