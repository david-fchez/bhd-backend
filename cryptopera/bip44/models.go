package bip44

import (
	"github.com/gcash/bchutil"
	"github.com/gcash/bchutil/hdkeychain"
)

type ChangeType uint32

const (
	ExternalChangeType ChangeType = 0
	InternalChangeType ChangeType = 1
)


type Purpose uint32

const (
	BIP44Purpose Purpose = 44
)

type CoinType uint32

const (
	BchCoinType        CoinType = 145
	SlpBchCoinType     CoinType = 245
	BchTestnetCoinType CoinType = 145
)

type HDStartPath struct {
	PurposeIndex  uint32 `json:"purpose_index"`
	CoinTypeIndex uint32 `json:"coin_type"`
	AccountIndex  uint32 `json:"account_index"`
}

type HDEndPath struct {
	ChangeIndex  uint32 `json:"change_index"`
	AddressIndex uint32 `json:"address_index"`
}

type Address struct {
	HDStartPath HDStartPath `json:"hd_start_path"`
	HDEndPath   HDEndPath   `json:"hd_end_path"`
	Value       string      `json:"value"`
	PubAddress		*bchutil.AddressPubKeyHash
	PrivateKey     *hdkeychain.ExtendedKey
}


