package bip44

import (
	"encoding/hex"
	"github.com/gcash/bchutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

type ExtendedKey struct {
	HdKey *hdkeychain.ExtendedKey
}


func NewKeyFromSeed(seed string, net Network) (*ExtendedKey, error) {
	pk, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}
	return NewKeyFromSeedBytes(pk, net)
}


func NewKeyFromSeedHex(seed string, net Network) (*ExtendedKey, error) {
	pk, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}
	return NewKeyFromSeedBytes(pk, net)
}

func NewKeyFromMnemonic(mnemonic, password string, net Network) (*ExtendedKey, error) {
	seed := bip39.NewSeed(mnemonic, password)
	return NewKeyFromSeedBytes(seed, MAINNET)
}

func NewKeyFromSeedBytes(seed []byte, net Network) (*ExtendedKey, error) {
	n, err := networkToChainConfig(net)
	if err != nil {
		return nil, err
	}

	xKey, err := hdkeychain.NewMaster(seed, n)

	if err != nil {
		return nil, err
	}

	return &ExtendedKey{
		HdKey: xKey,
	}, nil
}

func (e *ExtendedKey) BIP44AccountKey(coinType CoinType, accIndex uint32, includePrivateKey bool) (*AccountKey, error) {

	return e.baseDeriveAccount(BIP44Purpose, coinType, accIndex, includePrivateKey)
}

func (e *ExtendedKey) baseDeriveAccount(purpose Purpose, coinType CoinType, accIndex uint32, includePrivateKey bool) (*AccountKey, error) {

	var purposeIndex = uint32(purpose)
	var coinTypeIndex = uint32(coinType)

	if e.HdKey.IsPrivate() {
		purposeIndex = hdkeychain.HardenedKeyStart + purposeIndex
		coinTypeIndex = hdkeychain.HardenedKeyStart + coinTypeIndex
		accIndex = hdkeychain.HardenedKeyStart + accIndex
	}

	purposeK, err := e.HdKey.Child(purposeIndex)
	if err != nil {
		return nil, err
	}

	cTypeK, err := purposeK.Child(coinTypeIndex)
	if err != nil {
		return nil, err
	}

	accK, err := cTypeK.Child(accIndex)
	if err != nil {
		return nil, err
	}

	hdStartPath := HDStartPath{
		PurposeIndex:  purposeIndex,
		CoinTypeIndex: coinTypeIndex,
		AccountIndex:  accIndex,
	}

	if includePrivateKey {
		return &AccountKey{
			ExKey: accK,
			startPath:   hdStartPath,
		}, nil
	}

	pub, err := accK.Neuter()
	if err != nil {
		return nil, err
	}

	return &AccountKey{
		ExKey: pub,
		startPath:   hdStartPath,
	}, nil
}