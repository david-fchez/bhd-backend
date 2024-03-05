/*
		This is the definition, implementation of the exported API from backend to DART frontend.
	    Functions are receiving and returning simple types only. Few rules about naming:
		- All wallet functions start with W
		- All app functions start with A
		- All tdb functions start with D
		- all sales functions start with S
*/
package app

import (
	"C"
	"bhd/cryptopera"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"github.com/skip2/go-qrcode"
	"bunnyhedger.com/shared/bhdmodels"
	"strconv"
)

// ToCString must be in the same package so this is simple util function
func ToCString(s string) *C.char {
	return C.CString(s)
}

func FromCString(s *C.char) string {
	return C.GoString(s)
}

// InitializeWallet will create new wallet from a mnemonic string
// if you try to create multiple account from same seed it will be ignored
func InitializeWallet(mnemonic string) *ApiReturnStruct {
	var rValue = &ApiReturnStruct{}
	err := cryptopera.NewWalletFromBip39Seed(mnemonic)
	if err != nil {
		rValue.ErrorID = 6
		rValue.ErrorDescription = err.Error()
	}
	rValue.Content = ""
	return rValue
}

// GetWalletMnemonic returns wallet mnemonic to the frontend client
func GetWalletMnemonic() *ApiReturnStruct {
	var rValue = &ApiReturnStruct{}
	if cryptopera.Service == nil {
		rValue.ErrorID = 2
		rValue.ErrorDescription = "Wallet not initialized"
		return rValue
	}
	rValue.Content = cryptopera.Service.Mnemonic
	return rValue
}

// GetPublicBchAddress returns wallet mnemonic to the frontend client
func GetPublicBchAddress() *ApiReturnStruct {
	var rValue = &ApiReturnStruct{}
	if cryptopera.Service == nil {
		rValue.ErrorID = 2
		rValue.ErrorDescription = "Wallet not initialized"
		return rValue
	}
	rValue.Content = cryptopera.Service.PubAddress
	return rValue
}

// SignTransaction call this when one needs to have transaction signed
func SignTransaction(txStr string) *ApiReturnStruct {
	var rValue = &ApiReturnStruct{} // check if wallet exists
	if cryptopera.Service == nil {
		rValue.ErrorID = 2
		rValue.ErrorDescription = "Wallet not initialized"
		return rValue
	}
	var tx = &bhdmodels.Tx{}
	err := json.Unmarshal([]byte(txStr), tx)
	if err != nil {
		rValue.ErrorID = 3
		rValue.ErrorDescription = "Cannot deserialize tx request due to:" + err.Error()
		return rValue
	}

	err = cryptopera.Service.SignTransaction(cryptopera.Service.Mnemonic, tx)
	if err != nil {
		rValue.ErrorID = 4
		rValue.ErrorDescription = "Cannot sign tx due to:" + err.Error()
		return rValue
	}

	txSignedStr, err := json.Marshal(tx)
	if err != nil {
		rValue.ErrorID = 5
		rValue.ErrorDescription = "Cannot serialize signed tx due to:" + err.Error()
		return rValue
	}

	rValue.Content = string(txSignedStr)
	return rValue
}

// GetBchAddressQrCode returns png encoded hex string
func GetBchAddressQrCode(qrCode string, size string) *ApiReturnStruct {
	var png []byte
	var rValue = &ApiReturnStruct{}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		rValue.ErrorID = 5
		rValue.ErrorDescription = "Bad qr code size parameter:" + err.Error()
		return rValue
	}
	// generate qr code
	png, err = qrcode.Encode(qrCode, qrcode.Medium, sizeInt)
	if err != nil {
		rValue.ErrorID = 5
		rValue.ErrorDescription = "Cannot serialize qr code due to:" + err.Error()
		return rValue
	}
	rValue.Content = hex.EncodeToString(png)
	return rValue
}
