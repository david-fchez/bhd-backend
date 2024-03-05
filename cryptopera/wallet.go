package cryptopera

import (
	"bhd/bhdmodels"
	"bhd/cryptopera/bip44"
	"encoding/hex"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/txscript"
	hd "github.com/gcash/bchutil/hdkeychain"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
)

type KeyGenInfo struct {
	Key    *hd.ExtendedKey
	KeyIdx int32
}

var (
	Service *Wallet
)

type Wallet struct {
	NetParams  *chaincfg.Params
	PubAddress string
	Key        *bip44.ExtendedKey
	Mnemonic   string
}

// GetCryptoNetworkParams returns target blockchain network
func GetCryptoNetworkParams() *chaincfg.Params {
	return &chaincfg.MainNetParams
}

// NewWalletFromBip39Seed generates private key
// from the mnemonic string
func NewWalletFromBip39Seed(mnemonic string) error {
	var networkParams *chaincfg.Params
	var bip44NetType bip44.Network
	networkParams = GetCryptoNetworkParams()
	bip44NetType = bip44.MAINNET
	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(256)
		if err != nil {
			return errors.Wrap(err, "")
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	xKey, err := bip44.NewKeyFromMnemonic(mnemonic, "", bip44NetType)
	if err != nil {
		panic(err)
	}

	// convert to hex
	wl := &Wallet{
		NetParams: networkParams,
	}

	addr, err := wl.GenerateWalletAddress(xKey)
	wl.PubAddress = addr
	wl.Mnemonic = mnemonic
	Service = wl
	return nil
}

// GetAddress returns the address to use
func (w *Wallet) GetAddress() string {
	return w.PubAddress
}

// GenerateWalletAddress will create number (cnt) of addresses
// that will be used by the wallet as destination address
func (w *Wallet) GenerateWalletAddress(key *bip44.ExtendedKey) (string, error) {
	accountKey, err := key.BIP44AccountKey(bip44.BchCoinType, 0, true)
	if err != nil {
		panic(err)
	}
	netType := bip44.MAINNET
	externalAddress, err := accountKey.DeriveP2PKAddress(bip44.ExternalChangeType, uint32(255), netType)
	if err != nil {
		panic(err)
	}
	addr, err := externalAddress.PrivateKey.Address(w.NetParams)
	if err != nil {
		panic(err)
	}
	var bchPrefix string = bhdmodels.BchAddressPrefixForGeneration
	bchAddress, err := GetBech32Address(bchPrefix, addr.Hash160()[:])
	if err != nil {
		return "", err
	}
	return bchAddress, nil
}

// SignTransaction will sign the transaction, every input
// with proper key
func (w *Wallet) SignTransaction(mnemonic string, tx *bhdmodels.Tx) error {
	// sign the transactions by signing each input point
	_msg, err := bhdmodels.ToBCHDWireFormat(tx)
	if err != nil {
		return err
	}

	xKey, err := bip44.NewKeyFromMnemonic(mnemonic, "", bip44.MAINNET)
	if err != nil {
		return err
	}

	ak, err := xKey.BIP44AccountKey(bip44.BchCoinType, 0, true)
	addrRet, err := ak.DeriveP2PKAddress(bip44.ExternalChangeType, uint32(255), bip44.MAINNET)

	ecpPKey, err := addrRet.PrivateKey.ECPrivKey()

	var allSigned = true
	for i, el := range tx.Inputs {
		pubScript, err := hex.DecodeString(el.PubScript)
		if len(pubScript) == 0 {
			return errors.New("the input has empty pub script")
		}
		if err != nil {
			return err
		}
		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(pubScript, &chaincfg.MainNetParams)
		switch scriptClass {
		case txscript.NonStandardTy:
			allSigned = false
			continue
		case txscript.MultiSigTy:
			allSigned = false
			continue
		case txscript.NullDataTy:
			allSigned = false
			continue
		case txscript.ScriptHashTy:
			allSigned = false
			continue
		}
		sigScript, err := txscript.SignatureScript(_msg, i, el.Value, pubScript, txscript.SigHashAll, ecpPKey, true)
		if err != nil {
			return err
		}
		el.Signature = hex.EncodeToString(sigScript)
	}
	if !allSigned {
		return errors.New("not all inputs have been signed, some are not supported by this function")
	}

	err = w.ValidateTx(tx)
	if err != nil {
		return err
	}
	return nil
}

// ValidateTx with internal engine, useful for debugging
// not so much for production
func (w *Wallet) ValidateTx(tx *bhdmodels.Tx) error {
	msg, err := bhdmodels.ToBCHDWireFormat(tx)
	if err != nil {
		return err
	}

	flags := txscript.StandardVerifyFlags
	var isOkay = true
	for i, e := range tx.Inputs {
		pubScript, err := hex.DecodeString(e.PubScript)
		if err != nil {
			return err
		}
		vm, err := txscript.NewEngine(pubScript, msg, i, flags, nil, nil, nil, e.Value)
		if err != nil {
			return err
		}
		err = vm.Execute()
		if err != nil {
			isOkay = false
		}
	}
	if !isOkay {
		return errors.New("the tx" + tx.Hash + " failed to validate, no exact error given")
	}
	return nil
}
