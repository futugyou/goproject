package did

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var _ IDidResolver = (*DidEthrResolver)(nil)

// DidEthrResolver resolves Ethereum DIDs
type DidEthrResolver struct {
	networkStore NetworkConfigurationStore
}

// GetMethod implements [IDidResolver].
func (r *DidEthrResolver) GetMethod() string {
	return "ethr"
}

func NewDidEthrResolver(store NetworkConfigurationStore) *DidEthrResolver {
	return &DidEthrResolver{networkStore: store}
}

// Resolve resolves a DID string to a DID Document
func (r *DidEthrResolver) Resolve(ctx context.Context, did string) (*DidDocument, error) {
	if strings.TrimSpace(did) == "" {
		return nil, errors.New("did cannot be empty")
	}

	// DID
	decID, err := DidEthrExtractor(did)
	if err != nil {
		return nil, err
	}
	did = decID.Format()

	// network
	network := decID.Network
	if network == "" {
		network = "mainnet"
	}
	netCfg, err := r.networkStore.Get(ctx, network)
	if err != nil {
		return nil, fmt.Errorf("failed to get network config: %w", err)
	}

	// ethclient
	client, err := ethclient.DialContext(ctx, netCfg.RpcUrl)
	if err != nil {
		return nil, err
	}

	// bind contracts
	registry, err := NewEthereumDIDRegistry(common.HexToAddress(netCfg.ContractAdr), client)
	if err != nil {
		return nil, err
	}

	var now uint64 = uint64(time.Now().Unix())

	if decID.Version != nil {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(*decID.Version)))
		if err != nil {
			return nil, err
		}

		now = block.Time()
	}

	// load events
	evts, err := r.readEvents(ctx, client, registry, decID, netCfg)
	if err != nil {
		return nil, err
	}

	builder := NewDidDocumentBuilderWithID(did)
	controllerVM := DidDocumentVerificationMethod{
		ID:                  fmt.Sprintf("%s#controller", did),
		Type:                "EcdsaSecp256k1RecoveryMethod2020",
		Controller:          did,
		BlockChainAccountId: BuildEthereumMainet(decID.Address).Format(),
	}
	builder.AddVerificationMethodSimple(controllerVM, authAndAssertionUsage(), nil)

	if decID.PublicKey != "" {
		payload, err := HexToByteArray(decID.PublicKey)
		if err != nil {
			return nil, err
		}
		asymmKey, err := NewGenericECDSAKey("ES256K")
		if err != nil {
			return nil, err
		}
		err = asymmKey.Import(payload, nil)
		if err != nil {
			return nil, err
		}

		isReference := true
		var encoding SignatureKeyEncodingTypes = HEX
		builder.AddVerificationMethod("EcdsaSecp256k1VerificationKey2019",
			asymmKey,
			did,
			authAndAssertionUsage(),
			&isReference,
			nil,
			&encoding,
			func(c *DidDocumentVerificationMethod) {
				c.ID = fmt.Sprintf("%s#controllerKey", did)
			})
	}

	r.consume(builder, decID, evts, now, decID.Version)
	return builder.Build(), nil
}

func (r *DidEthrResolver) readEvents(ctx context.Context, client *ethclient.Client, registry *EthereumDIDRegistry, decID *DecentralizedIdentifierEthr, netCfg *NetworkConfiguration) ([]ERC1056Event, error) {
	var result []ERC1056Event

	blockNumber := big.NewInt(0)
	if decID.Version != nil {
		blockNumber = big.NewInt(int64(*decID.Version))
	}

	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{common.HexToAddress(netCfg.ContractAdr)},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	for len(logs) > 0 {
		evt := r.extractEvent(logs, registry)
		if evt == nil {
			break
		}
		evt.BlockNumber = blockNumber
		result = append(result, *evt)

		if evt.PreviousChange.Cmp(blockNumber) < 0 {
			blockNumber = evt.PreviousChange
			query.FromBlock = blockNumber
			query.ToBlock = blockNumber
			logs, err = client.FilterLogs(ctx, query)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return result, nil
}

func (r *DidEthrResolver) extractEvent(logs []types.Log, registry *EthereumDIDRegistry) *ERC1056Event {
	for _, lg := range logs {
		// 1️⃣ try AttributeChanged
		itAttr, err := registry.FilterDIDAttributeChanged(&bind.FilterOpts{
			Start:   lg.BlockNumber,
			End:     &lg.BlockNumber,
			Context: context.TODO(),
		}, nil)
		if err == nil && itAttr.Next() {
			evt := itAttr.Event
			payload, _ := hex.DecodeString(strings.TrimPrefix(string(evt.Value), "0x"))
			return &ERC1056Event{
				Identity:       decodeString32(evt.Name),
				Value:          string(payload),
				PreviousChange: evt.PreviousChange,
				ValidTo:        evt.ValidTo,
				Type:           DIDAttributeChanged,
			}
		}

		// 2️⃣ try OwnerChanged
		itOwner, err := registry.FilterDIDOwnerChanged(&bind.FilterOpts{
			Start:   lg.BlockNumber,
			End:     &lg.BlockNumber,
			Context: context.TODO(),
		}, nil)
		if err == nil && itOwner.Next() {
			evt := itOwner.Event
			return &ERC1056Event{
				Identity:       evt.Identity.String(),
				Owner:          evt.Owner.String(),
				PreviousChange: evt.PreviousChange,
				Type:           DIDOwnerChanged,
			}
		}

		// 3️⃣ try DelegateChanged
		itDel, err := registry.FilterDIDDelegateChanged(&bind.FilterOpts{
			Start:   lg.BlockNumber,
			End:     &lg.BlockNumber,
			Context: context.TODO(),
		}, nil)
		if err == nil && itDel.Next() {
			evt := itDel.Event
			return &ERC1056Event{
				Identity:       evt.Identity.String(),
				Delegate:       evt.Delegate.String(),
				DelegateType:   decodeString32(evt.DelegateType),
				PreviousChange: evt.PreviousChange,
				Type:           DIDDelegateChanged,
			}
		}
	}
	return nil
}

func (r *DidEthrResolver) consume(builder *DidDocumentBuilder, decID *DecentralizedIdentifierEthr, events []ERC1056Event, now uint64, blockNumber *int) {
	for i := len(events)/2 - 1; i >= 0; i-- {
		opp := len(events) - 1 - i
		events[i], events[opp] = events[opp], events[i]
	}
	did := decID.Format()
	controller := decID.Address
	authRefs := make(map[string]string)
	sigRefs := make(map[string]string)
	keyAgreementRefs := make(map[string]string)
	services := make(map[string]DidDocumentService)
	verificationMethods := make(map[string]DidDocumentVerificationMethod)

	regex := regexp.MustCompile(`^did/(pub|svc)/(\w+)(/(\w+))?(/(\w+))?$`)
	serviceCount := 0
	delegateCount := 0

	for _, evt := range events {
		if blockNumber != nil && evt.BlockNumber.Int64() > int64(*blockNumber) {
			continue
		}

		if evt.ValidTo != nil && evt.ValidTo.Uint64() < now {
			continue
		}

		eventName := evt.Type.Name()
		eventIndex := ""
		switch evt.Type {
		case DIDDelegateChanged:
			eventIndex = eventName + "-" + evt.DelegateType + "-" + evt.Delegate
		case DIDAttributeChanged:
			eventIndex = eventName + "-" + evt.Identity + "-" + evt.Value
		default:
			eventIndex = eventName
		}

		splitted := strings.Split(evt.Identity, "/")

		if evt.ValidTo.Uint64() >= now {
			switch evt.Type {
			case DIDDelegateChanged:
				delegateCount++
				verificationMethodId := fmt.Sprintf("%s#delegate-%d", did, delegateCount)
				switch evt.DelegateType {
				case "sigAuth":
					authRefs[eventIndex] = verificationMethodId
					sigRefs[eventIndex] = verificationMethodId
				case "veriKey":
					verificationMethod := DidDocumentVerificationMethod{
						ID:                  verificationMethodId,
						Type:                "EcdsaSecp256k1RecoveryMethod2020",
						Controller:          did,
						BlockChainAccountId: BuildEthereumMainet(decID.Address).Format(),
					}
					verificationMethods[eventIndex] = verificationMethod
					sigRefs[eventIndex] = verificationMethod.ID
				}
			case DIDAttributeChanged:
				if !regex.MatchString(evt.Identity) {
					continue
				}

				section := splitted[1]
				switch section {
				case "svc":
					serviceCount++
					svcID := fmt.Sprintf("%s#service-%d", decID.Format(), serviceCount)
					t, _ := HexToByteArray(evt.Value)
					services[svcID] = DidDocumentService{
						ID:              svcID,
						Type:            splitted[2],
						ServiceEndpoint: string(t),
					}
				case "pub":
					delegateCount++
					keyAlgorithm := splitted[2]
					keyPurpose := splitted[3]
					encoding := splitted[4]
					resolvedType := resolveVerificationMethodType(keyAlgorithm)
					if len(resolvedType) == 0 {
						continue
					}
					verificationMethodId := fmt.Sprintf("%s#delegate-%d", did, delegateCount)
					verificationMethod := DidDocumentVerificationMethod{
						ID:         verificationMethodId,
						Type:       "EcdsaSecp256k1RecoveryMethod2020",
						Controller: controller,
					}
					verificationMethods[eventIndex] = verificationMethod
					switch keyPurpose {
					case "sigAuth":
						authRefs[eventIndex] = verificationMethod.ID
						sigRefs[eventIndex] = verificationMethod.ID
					case "enc":
						keyAgreementRefs[eventIndex] = verificationMethod.ID
					default:
						sigRefs[eventIndex] = verificationMethod.ID
					}

					payload, _ := HexToByteArray(evt.Value)

					switch encoding {
					case "hex":
						verificationMethod.PublicKeyHex = ToHex(payload, false)
					case "base64":
						verificationMethod.PublicKeyBase64 = base64.URLEncoding.EncodeToString(payload)
					case "base58":
						verificationMethod.PublicKeyBase58 = base58.Encode(payload)
					}
				}
			}
		} else if evt.Type == DIDOwnerChanged {
			controller = evt.Owner
		} else {
			section := splitted[1]
			if section == "svc" {
				serviceCount++
			}
			if section == "pub" {
				delegateCount++
			}
			delete(authRefs, eventIndex)
			delete(sigRefs, eventIndex)
			delete(services, eventIndex)
			delete(verificationMethods, eventIndex)
		}

	}

	for _, svc := range services {
		builder.AddDidDocumentService(svc)
	}
	for _, vm := range verificationMethods {
		usage := determineUsage(vm.ID, authRefs, keyAgreementRefs)
		builder.AddVerificationMethodSimple(vm, usage, nil)
	}
}

func authAndAssertionUsage() VerificationMethodUsage {
	return Authentication | AssertionMethod
}

func determineUsage(vmID string, authRefs, keyAgreementRefs map[string]string) VerificationMethodUsage {
	usage := AssertionMethod
	if containsValue(keyAgreementRefs, vmID) {
		usage = KeyAgreement
	} else if containsValue(authRefs, vmID) {
		usage |= Authentication
	}
	return usage
}

func resolveVerificationMethodType(keyAlg string) string {
	switch keyAlg {
	case "Secp256k1":
		return "EcdsaSecp256k1VerificationKey2019"
	case "Ed25519":
		return "Ed25519VerificationKey2018"
	case "X25519":
		return "X25519KeyAgreementKey2019"
	}

	return ""
}
