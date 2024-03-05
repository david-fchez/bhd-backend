package bhdmodels

import (
	"bhd/bch/msg"
	"bhd/utils"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/wire"
)

const (
	BhdGetTransactionRequestType        = 5
	BhdGetTransactionResponseType       = 6
	BhdGetTransactionsRequestType       = 7
	BhdGetTransactionsResponseType      = 8
	BhdMemPoolFilterRequest             = 9
	BhdMemPoolFilterResponse            = 10
	BhdBroadcastTransactionRequestType  = 11
	BhdBroadcastTransactionResponseType = 12
	BhdGetUxtosRequestType              = 13
	BhdGetUxtosResponseType             = 14
	BhdMemPoolTransactionRequestType    = 15
	BhdGetBalanceRequestType            = 16
	BhdGetBalanceResponseType           = 17
	BhdSendCoinsRequestType             = 18
	BhdSendCoinsResponseType            = 19
	BhdRegisterBchAddressType           = 20
	BchAddressPrefix                    = "bitcoincash:"
	BchAddressPrefixForGeneration       = "bitcoincash"
)

/*
	the tx structures, transactions, uxtos, etc
*/

type TxIn struct {
	Sequence  uint32 `json:"sequence"`
	Value     int64  `json:"value"`
	PrevHash  string `json:"prevHash"`
	PrevIndex uint32 `json:"prevIndex"`
	PubScript string `json:"pubScript"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Value      int64  `json:"value"`
	Spent      bool   `json:"spent"`
	PkScript   string `json:"pkScript"`
	Address    string `json:"address"`
	IsCashback bool   `json:"-"`
	AddressRaw []byte `json:"-"`
}

type Tx struct {
	Hash       string   `json:"hash"`
	DateTime   int64    `json:"dateTime"`
	Size       int32    `json:"size"`
	Height     int32    `json:"height"`
	Index      int32    `json:"index"`
	Version    int32    `json:"version"`
	LockTime   uint32   `json:"lockTime"`
	Inputs     []*TxIn  `json:"inputs"`
	Outputs    []*TxOut `json:"outputs"`
	InputVal   int64    `json:"inputVal"`
	OutputVal  int64    `json:"outputVal"`
	CashBack   int64    `json:"cashBack"`
	NetworkFee int64    `json:"networkFee"`
}

// Unpack will parse the content of the
// binary array and return the instance
// of the new struct
func Unpack[T any](content []byte) (T, error) {
	// extract header
	var obj = new(T)
	err := json.Unmarshal(content, obj)
	if err != nil {
		return *obj, err
	}
	return *obj, nil
}

// ToJson will return json
// representation of transaction
func (tx *Tx) ToJson() string {
	content, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(content)
}

// TxToByteArray will convert Tx to byte array
func TxToByteArray(tx *Tx) []byte {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(tx)
	if err != nil {
		panic(err)
	}
	content := buff.Bytes()
	return content
}

// TxFromByteArray will convert byte array to tx
func TxFromByteArray(content []byte) *Tx {
	var reader = bytes.NewReader(content)
	dec := gob.NewDecoder(reader)
	var oldTx Tx
	dec.Decode(&oldTx)
	return &oldTx
}

var (
	PduStartSignature = []byte{255, 254, 253, 252, 251, 250}
)

// structToByte will pack the struct
// into json and append the header
func structToByte(v interface{}, requestType int16, pduId uint32) []byte {
	content, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	pdu := make([]byte, 6)
	pdu[0] = PduStartSignature[0]
	pdu[1] = PduStartSignature[1]
	pdu[2] = PduStartSignature[2]
	pdu[3] = PduStartSignature[3]
	pdu[4] = PduStartSignature[4]
	pdu[5] = PduStartSignature[5]
	pdu = append(pdu, utils.UInt32ToByte(pduId)...)
	pdu = append(pdu, utils.UInt32ToByte(uint32(len(content)))...)
	pdu = append(pdu, utils.Int16ToByte(requestType)...)
	pdu = append(pdu, content...)
	return pdu
}

// ToBCHDWireFormat will create wire.Tx struct instance
// as it was too much to move all of the signature code
// to this project
func ToBCHDWireFormat(tx *Tx) (*wire.MsgTx, error) {
	wireMsg := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     nil,
		TxOut:    nil,
		LockTime: 0,
	}
	// inputs with signature
	for _, txIn := range tx.Inputs {
		hash, err := chainhash.NewHashFromStr(txIn.PrevHash)
		if err != nil {
			return nil, err
		}
		sign, err := hex.DecodeString(txIn.Signature)
		if err != nil {
			return nil, err
		}
		wireMsg.TxIn = append(wireMsg.TxIn, &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *hash,
				Index: txIn.PrevIndex,
			},
			SignatureScript: sign,
			Sequence:        txIn.Sequence,
		})
	}
	// outputs
	for _, txOut := range tx.Outputs {
		payScript, err := hex.DecodeString(txOut.PkScript)
		if err != nil {
			return nil, err
		}
		wireMsg.TxOut = append(wireMsg.TxOut, &wire.TxOut{
			Value:    txOut.Value,
			PkScript: payScript,
		})
	}
	return wireMsg, nil
}

// ToWireFormat will convert from internal transaction
// model to wire ready format
func (tx *Tx) ToWireFormat() (*msg.Tx, error) {
	bchTx := &msg.Tx{
		Version:  msg.TransactionVersion,
		LockTime: 0,
	}
	// inputs with signature
	for _, txIn := range tx.Inputs {
		sign, err := hex.DecodeString(txIn.Signature)
		if err != nil {
			return nil, err
		}
		prevHash, err := msg.NewHashFromString(txIn.PrevHash)
		if err != nil {
			return nil, err
		}
		bchTx.Inputs = append(bchTx.Inputs, msg.TxInput{
			PreviousOutputHash: prevHash,
			PreviousIndex:      txIn.PrevIndex,
			UnlockingScript:    sign,
			SequenceNumber:     txIn.Sequence,
		})
	}
	// outputs
	for _, txOut := range tx.Outputs {
		payScript, err := hex.DecodeString(txOut.PkScript)
		if err != nil {
			return nil, err
		}
		bchTx.Outputs = append(bchTx.Outputs, msg.TxOutput{
			Value:         uint64(txOut.Value),
			LockingScript: payScript,
		})
	}
	return bchTx, nil
}

// MemPoolFilterRequest is a list of addresses
// for which the server will inform client
// that there is new transaction in the pool
type MemPoolFilterRequest struct {
	PlayerId string   `json:"playerId"`
	Address  []string `json:"address"`
	BasePdu
}

func (r *MemPoolFilterRequest) Pack() []byte {
	return structToByte(r, BhdMemPoolFilterRequest, r.PduId)
}

type MemPoolFilterResponse struct {
	Info string
	BasePdu
}

func (r *MemPoolFilterResponse) Pack() []byte {
	return structToByte(r, BhdMemPoolFilterResponse, r.PduId)
}

// TransactionInMemPoolRequest response to request
// to get all transactions
type TransactionInMemPoolRequest struct {
	Transaction *Tx `json:"transaction"`
	BasePdu
}

func (r *TransactionInMemPoolRequest) Pack() []byte {
	return structToByte(r, BhdMemPoolTransactionRequestType, r.PduId)
}

// Uxto is non spent transaction output
type Uxto struct {
	Hash     string `json:"hash"`
	Index    int32  `json:"index"`
	PkScript string `json:"pkScript"`
	Value    int64  `json:"value"`
	BasePdu
}

type ProvideUxtoRequest struct {
	PlayerId string   `json:"playerId"`
	Address  []string `json:"address"`
	Skip     int      `json:"skip"`
	PageSize int      `json:"pageSize"`
	BasePdu
}

func (r *ProvideUxtoRequest) Pack() []byte { return structToByte(r, BhdGetUxtosRequestType, r.PduId) }

type ProvideUxtoResponse struct {
	BalanceSat int64   `json:"balanceSat"`
	BalanceBch float64 `json:"balanceBch"`
	Inputs     []*Uxto `json:"inputs"`
	InputCount int     `json:"inputCount"`
	BasePdu
}

func (r *ProvideUxtoResponse) Pack() []byte { return structToByte(r, BhdGetUxtosResponseType, r.PduId) }

// BroadcastTransactionRequest is request to broadcast transaction
type BroadcastTransactionRequest struct {
	PlayerId          string `json:"playerId"`
	SignedTransaction *Tx    `json:"signedTransaction"`
	BasePdu
}

func (r *BroadcastTransactionRequest) Pack() []byte {
	return structToByte(r, BhdBroadcastTransactionRequestType, r.PduId)
}

// BroadcastTransactionResponse is response to ping request
type BroadcastTransactionResponse struct {
	Info string `json:"info"`
	BasePdu
}

func (r *BroadcastTransactionResponse) Pack() []byte {
	return structToByte(r, BhdBroadcastTransactionResponseType, r.PduId)
}

// ProvideTransactionRequest is request to provide transactions
// for the addresses in the list
type ProvideTransactionRequest struct {
	PlayerId string `json:"playerId"`
	Hash     string `json:"hash"`
	BasePdu
}

func (r *ProvideTransactionRequest) Pack() []byte {
	return structToByte(r, BhdGetTransactionRequestType, r.PduId)
}

// ProvideTransactionResponse response to request
// to get all transactions
type ProvideTransactionResponse struct {
	Transaction *Tx `json:"transaction"`
	BasePdu
}

func (r *ProvideTransactionResponse) Pack() []byte {
	return structToByte(r, BhdGetTransactionResponseType, r.PduId)
}

// ProvideTransactionsRequest is request to provide transactions
// for the addresses in the list
type ProvideTransactionsRequest struct {
	PlayerId string   `json:"playerId"`
	Address  []string `json:"address"`
	Skip     int      `json:"skip"`
	PageSize int      `json:"pageSize"`
	BasePdu
}

func (r *ProvideTransactionsRequest) Pack() []byte {
	return structToByte(r, BhdGetTransactionsRequestType, r.PduId)
}

// ProvideTransactionsResponse response to request
// to get all transactions
type ProvideTransactionsResponse struct {
	Transactions     []*Tx `json:"transactions"`
	TransactionCount int   `json:"transactionCount"`
	BasePdu
}

func (r *ProvideTransactionsResponse) Pack() []byte {
	return structToByte(r, BhdGetTransactionsResponseType, r.PduId)
}

// GetBalanceRequest is request to provide balance of the wallet
type GetBalanceRequest struct {
	PlayerId string   `json:"playerId"`
	Address  []string `json:"address"`
	BasePdu
}

func (r *GetBalanceRequest) Pack() []byte {
	return structToByte(r, BhdGetBalanceRequestType, r.PduId)
}

// GetBalanceResponse returns wallet balance
type GetBalanceResponse struct {
	BalanceSat int64   `json:"balanceSat"`
	BalanceBch float64 `json:"balanceBch"`
	BasePdu
}

func (r *GetBalanceResponse) Pack() []byte {
	return structToByte(r, BhdGetBalanceResponseType, r.PduId)
}

// SendCoinsRequest is the request to transfer money from one address
// to another
type SendCoinsRequest struct {
	PlayerId              string `json:"playerId"`
	OriginBchAddress      string `json:"originBchAddress"`
	DestinationBchAddress string `json:"destinationBchAddress"`
	AmountToTransfer      int64  `json:"amountToTransfer"`
	BasePdu
}

func (r *SendCoinsRequest) Pack() []byte {
	return structToByte(r, BhdSendCoinsRequestType, r.PduId)
}

// SendCoinsResponse is the response to SendCoinRequest and
// it will return the transaction that should be signed by the
// client and the call BroadcastTransactionRequest
type SendCoinsResponse struct {
	OriginBchAddress      string `json:"originBchAddress"`
	DestinationBchAddress string `json:"destinationBchAddress"`
	AmountToTransfer      int64  `json:"amountToTransfer"`
	Transaction           *Tx    `json:"transaction"`
	BasePdu
}

func (r *SendCoinsResponse) Pack() []byte {
	return structToByte(r, BhdSendCoinsResponseType, r.PduId)
}

// RegisterBchAddressRequest is request to register
// player address so that system can start collecting
// relevant transactions
type RegisterBchAddressRequest struct {
	BchAddress string `json:"bchAddress"`
	BasePdu
}

func (r *RegisterBchAddressRequest) Pack() []byte {
	return structToByte(r, BhdRegisterBchAddressType, r.PduId)
}
