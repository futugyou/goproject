### did

[ethr-did-registry](https://github.com/uport-project/ethr-did-registry/)

```sh
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
npm install -g solc
solcjs --abi contracts/EthereumDIDRegistry.sol -o build/
abigen --abi build/contracts_EthereumDIDRegistry_sol_EthereumDIDRegistry.abi --pkg did --type EthereumDIDRegistry --out registry.go
```
