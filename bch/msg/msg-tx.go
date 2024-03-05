package msg

import (
	"bhd/log"
	"bhd/utils"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchutil"
)

type Script []byte

func (m Script) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(m))
}

func (m Script) ToHexString() string {
	return hex.EncodeToString(m)
}

type TxInput struct {
	PreviousOutputHash Hash
	PreviousIndex      uint32
	UnlockingScript    Script
	AddressStr         string
	AddressRaw         []byte
	SequenceNumber     uint32
}

func (inp *TxInput) Size() int {
	return 32 + 4 + len(inp.UnlockingScript) + 4
}

type TxOutput struct {
	Value         uint64
	LockingScript Script
	AddressStr    string
	AddressRaw    []byte
	Token         *CashToken
}

func (o *TxOutput) Size() int {
	return 8 + len(o.LockingScript)
}

// Tx provides the contents of a transaction.
type Tx struct {
	TxHash   Hash
	Version  uint32
	Inputs   []TxInput
	Outputs  []TxOutput
	LockTime uint32
}

func (m *Tx) Size() int {
	size := 32 + 4
	for _, inp := range m.Inputs {
		size += inp.Size()
	}
	for _, out := range m.Outputs {
		size += out.Size()
	}
	size += 4
	return size
}

func (m *Tx) GetCommandString() string {
	return CmdTx
}

// SerializeSize will serialize the transaction
// and return the exact size of tx
func (m *Tx) SerializeSize(includeCashBack bool) uint64 {
	var txSize = 4                // transaction version
	txSize += 2                   // input count
	txSize += len(m.Inputs) * 100 // sig script
	txSize += len(m.Inputs) * 4   // prev point index
	txSize += len(m.Inputs) * 32  // prev point hash
	// outputs
	txSize += 2 // output count
	for _, out := range m.Outputs {
		txSize += 8 // value
		txSize += 2 // script len
		txSize += len(out.LockingScript)
	}
	txSize += 4 // lock time
	// sig script?
	txSize += len(m.Inputs) * 100
	if includeCashBack {
		txSize += 8  // value
		txSize += 2  // script len
		txSize += 32 // actual script
	}
	return uint64(txSize)
}

// Pack constructs the binary content of the version message
func (m *Tx) Pack() []byte {
	var buf bytes.Buffer
	// yeah, only supporting version 1
	buf.Write(utils.UInt32ToByte(m.Version))
	buf.Write(utils.VarIntToByte(uint64(len(m.Inputs))))
	for _, in := range m.Inputs {
		buf.Write(in.PreviousOutputHash)
		buf.Write(utils.UInt32ToByte(in.PreviousIndex))
		buf.Write(utils.VarIntToByte(uint64(len(in.UnlockingScript))))
		buf.Write(in.UnlockingScript)
		buf.Write(utils.UInt32ToByte(in.SequenceNumber))
	}
	buf.Write(utils.VarIntToByte(uint64(len(m.Outputs))))
	for _, out := range m.Outputs {
		buf.Write(utils.UInt64ToByte(out.Value))
		buf.Write(utils.VarIntToByte(uint64(len(out.LockingScript))))
		buf.Write(out.LockingScript)
	}
	buf.Write(utils.UInt32ToByte(m.LockTime))
	return buf.Bytes()
}

// GetHash  the hash of the tx, it has to serialize
// transaction to get its proper hash
func (m *Tx) GetHash() Hash {
	return DoubleHashB(m.Pack())
}

// decodeToken will parse the token structure
// as described in the https://github.com/cashtokens/cashtokens
func decodeToken(script []byte) (*CashToken, int, error) {
	var reader = bytes.NewReader(script)
	_, _ = reader.ReadByte()
	var tokenId = NewHash()
	_, err := reader.Read(tokenId)
	if err != nil {
		return nil, 0, err
	}
	tokenBitfield, err := reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}

	var ct = &CashToken{
		Category:      tokenId,
		TokenBitField: tokenBitfield,
	}
	var hasCommitmentLength = tokenBitfield&0x40 == 0x040
	var hasNft = tokenBitfield&0x20 == 0x20
	if hasNft {
		if tokenBitfield&0x00 == 0x00 {
			ct.TokenType = TokenImmutable
		}
		if tokenBitfield&0x01 == 0x01 {
			ct.TokenType = TokenMutable
		}
		if tokenBitfield&0x02 == 0x02 {
			ct.TokenType = TokenMinting
		}
	}
	var hasAmount = tokenBitfield&0x10 == 0x10
	if hasCommitmentLength {
		var commitmentLen = utils.ReadVarInt(reader)
		ct.Commitment = make([]byte, commitmentLen)
		_, err = reader.Read(ct.Commitment)
		if err != nil {
			return nil, 0, err
		}
	}
	if hasAmount {
		ct.Amount = utils.ReadVarInt(reader)
	}

	return ct, reader.Len(), nil
}

func DecodeTxMsg(reader *bytes.Reader) (*Tx, error) {
	ver := &Tx{
		Version:  0,
		Inputs:   make([]TxInput, 0),
		Outputs:  make([]TxOutput, 0),
		LockTime: 0,
	}
	ver.Version = utils.ReadUint32(reader)
	// read input count
	inputCount := int(utils.ReadVarInt(reader))
	// read inputs
	for i := 0; i < inputCount; i++ {
		var in = TxInput{
			PreviousOutputHash: NewHash(),
			PreviousIndex:      0,
			UnlockingScript:    nil,
			SequenceNumber:     0,
		}
		_, err := reader.Read(in.PreviousOutputHash)
		if err != nil {
			return nil, err
		}
		in.PreviousIndex = utils.ReadUint32(reader)
		var scriptLen = utils.ReadVarInt(reader)
		in.UnlockingScript = make([]byte, scriptLen)
		_, err = reader.Read(in.UnlockingScript)
		if err != nil {
			return nil, err
		}
		in.SequenceNumber = utils.ReadUint32(reader)
		// parse address
		dataElements, err := txscript.ExtractDataElements(in.UnlockingScript)
		if len(dataElements) >= 2 {
			var pubAddress = dataElements[1]
			addr, err := bchutil.NewAddressPubKeyHash(bchutil.Hash160(pubAddress), &chaincfg.MainNetParams)
			if err != nil {
				log.Error("Failed to get address from pubkey hash")
			} else {
				in.AddressRaw = addr.ScriptAddress()
				in.AddressStr = addr.String()
			}
		}
		ver.Inputs = append(ver.Inputs, in)
	}
	// read outputs
	var outputLen = int(utils.ReadVarInt(reader))
	for i := 0; i < outputLen; i++ {
		var out = TxOutput{
			Value:         0,
			LockingScript: nil,
		}
		out.Value = utils.ReadUint64(reader)
		var scriptLen = utils.ReadVarInt(reader)
		out.LockingScript = make([]byte, scriptLen)
		_, err := reader.Read(out.LockingScript)
		if err != nil {
			return nil, err
		}
		// if this is token output then
		// decode the token structure
		// and update the locking script
		if len(out.LockingScript) > 0 && out.LockingScript[0] == CashTokenPrefix {
			cashToken, remainingBytes, err := decodeToken(out.LockingScript)
			if err != nil {
				log.Error("Failed to parse token structure due to", err.Error())
			}
			out.Token = cashToken
			// what remains of the script
			// is regular address
			var addrScript = out.LockingScript[len(out.LockingScript)-remainingBytes:]
			_, addresses, _, err := txscript.ExtractPkScriptAddrs(addrScript, &chaincfg.MainNetParams)
			if err == nil && len(addresses) > 0 {
				out.AddressStr = addresses[0].EncodeAddress()
				out.AddressRaw = addresses[0].ScriptAddress()
			}
		} else {
			// decode address normaly
			_, addresses, _, err := txscript.ExtractPkScriptAddrs(out.LockingScript, &chaincfg.MainNetParams)
			if err != nil {
				log.Error("Failed to extract address from tx output")
			}
			if err == nil && len(addresses) > 0 {
				out.AddressStr = addresses[0].EncodeAddress()
				out.AddressRaw = addresses[0].ScriptAddress()
			}
		}
		ver.Outputs = append(ver.Outputs, out)
	}
	ver.LockTime = utils.ReadUint32(reader)
	ver.TxHash = ver.GetHash()
	return ver, nil
}

func NewTxMsg() *Tx {
	return &Tx{
		TxHash:   EmptyHash,
		Version:  0x01,
		Inputs:   make([]TxInput, 0),
		Outputs:  make([]TxOutput, 0),
		LockTime: 0,
	}
}

// ToJson converts tx to json string
func (m *Tx) ToJson() string {
	content, err := json.MarshalIndent(m, "   ", "")
	if err != nil {
		return err.Error()
	}
	return string(content)
}
