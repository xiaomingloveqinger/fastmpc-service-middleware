package types

import (
	"encoding/json"
	"errors"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	common2 "github.com/anyswap/fastmpc-service-middleware/chains/common"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"math/big"
	"strings"
)

type EVMChain struct {
	ChainType
}

var EvmChain = NewEVMChain()

func NewEVMChain() *EVMChain {
	return &EVMChain{
		EVM,
	}
}

type UnsignedEVMTx struct {
	From        string `json:"from"`
	To          string `json:"to"`
	ChainId     string `json:"chainId"`
	Value       string `json:"value"`
	Nonce       uint64 `json:"nonce"`
	Gas         uint64 `json:"gas"`
	GasPrice    int64  `json:"gasPrice"`
	Data        string `json:"data"`
	OriginValue string `json:"originValue"`
	Name        string `json:"name"`
}

func (*EVMChain) ValidateParam(txdata string) bool {
	unTx := &UnsignedEVMTx{}
	err := json.Unmarshal([]byte(txdata), unTx)
	if err != nil {
		return false
	}
	if !common.CheckEthereumAddress(unTx.From) {
		return false
	}
	if !common.CheckEthereumAddress(unTx.To) {
		return false
	}
	if common.IsBlank(unTx.ChainId) || common.IsBlank(unTx.Name) || common.IsBlank(unTx.Value) ||
		common.IsBlank(unTx.OriginValue) || unTx.Gas == 0 || unTx.GasPrice == 0 {
		return false
	}
	if len(unTx.Value) < 2 || len(unTx.ChainId) < 2 {
		return false
	}
	_, ok := new(big.Int).SetString(unTx.ChainId[2:], 16)
	if !ok {
		return false
	}
	_, ok = new(big.Int).SetString(unTx.Value[2:], 16)
	if !ok {
		return false
	}

	return true
}

func (e *EVMChain) GetUnsignedTransactionHash(data string) (string, error) {
	if !e.ValidateParam(data) {
		return "", errors.New("invalid unsigned tx data")
	}
	unTx := &UnsignedEVMTx{}
	_ = json.Unmarshal([]byte(data), unTx)
	chainId, _ := new(big.Int).SetString(unTx.ChainId[2:], 16)
	signer := NewEIP155Signer(chainId)
	amount, _ := new(big.Int).SetString(unTx.Value[2:], 16)
	var payload []byte
	if unTx.Data == "" {
		payload = []byte{}
	} else {
		payload = []byte(unTx.Data)
	}
	tx := NewTransaction(unTx.Nonce, common2.HexToAddress(unTx.To), amount, unTx.Gas, new(big.Int).SetInt64(unTx.GasPrice), payload)
	return signer.Hash(tx).String(), nil
}

func (e *EVMChain) ValidateUnsignedTransactionHash(src string, cmp string) bool {
	if !e.ValidateParam(src) {
		log.Error("ValidateUnsignedTransactionHash", "src", "not valid src")
		return false
	}
	hash, err := e.GetUnsignedTransactionHash(src)
	if err != nil {
		log.Error("ValidateUnsignedTransactionHash", "hash", err)
		return false
	}
	if !strings.EqualFold(hash, cmp) {
		return false
	}
	return true
}
