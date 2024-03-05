package msg

import (
	"bhd/utils"
	"bytes"
	"errors"
	"strconv"
)

// GetBlocksMsg request the sequence of blocks that occur after a specific block. If the specified block is on the
// server's most-work chain, the server responds with a set of up to 500 inv messages identifying the next
// blocks on that chain. If the specified block is not on the most-work chain, the server uses block information
// in the locator structure to determine the fork point and provides inv messages from that point.
type GetBlocksMsg struct {
	ProtocolVersion uint32
	Count           uint64
	Items           []Hash
	StopAtHash      Hash
}

func (m *GetBlocksMsg) GetCommandString() string {
	return CmdGetBlocks
}

// Pack constructs the binary content of the version message
func (m *GetBlocksMsg) Pack() []byte {
	var buf bytes.Buffer
	buf.Write(utils.UInt32ToByte(m.ProtocolVersion))
	buf.Write(utils.VarIntToByte(uint64(len(m.Items))))
	for _, item := range m.Items {
		buf.Write(item[:])
	}
	buf.Write(m.StopAtHash[:])
	return buf.Bytes()
}

func DecodeGetBlocksMsg(reader *bytes.Reader) (*GetBlocksMsg, error) {
	ver := &GetBlocksMsg{}
	ver.Count = utils.ReadVarInt(reader)
	for i := 0; i < int(ver.Count); i++ {
		item := NewHash()
		n, err := reader.Read(item)
		if err != nil {
			return nil, err
		}
		if n != 32 {
			return nil, errors.New("hash should be 32 bytes,read:" + strconv.Itoa(n))
		}
		ver.Items = append(ver.Items, item)
	}
	return ver, nil
}

func NewGetBlocksMsg() *GetBlocksMsg {
	var addr = &GetBlocksMsg{
		ProtocolVersion: ProtocolVersion,
		Count:           0,
		Items:           make([]Hash, 0),
		StopAtHash:      EmptyHash,
	}
	return addr
}

func (b *GetBlocksMsg) AddBlock(hash Hash) {
	b.Items = append(b.Items, hash)
}

type CompressedTargetFormat struct {
	Exponent    uint8
	Significand [3]byte
}

/*
the block structure and encode / decode functions
*/

type BlockHeader struct {
	BlockVersion  int32
	PrevBlockHash Hash
	MerkleRoot    Hash
	Timestamp     uint32
	HashTarget    CompressedTargetFormat
	Nonce         uint32
}

// Pack set includeTransactionCount when you are packing
// normal block with transaction, just for pure headers this is
// always 0
func (b *BlockHeader) Pack() []byte {
	var buf bytes.Buffer
	buf.Write(utils.UInt32ToByte(uint32(b.BlockVersion)))
	buf.Write(b.PrevBlockHash[:])
	buf.Write(b.MerkleRoot[:])
	buf.Write(utils.UInt32ToByte(b.Timestamp))
	buf.WriteByte(b.HashTarget.Exponent)
	buf.Write(b.HashTarget.Significand[:])
	buf.Write(utils.UInt32ToByte(b.Nonce))
	return buf.Bytes()
}

// Hash will return hash of the block
func (b *BlockHeader) Hash() Hash {
	var buf bytes.Buffer
	buf.Write(utils.UInt32ToByte(uint32(b.BlockVersion)))
	buf.Write(b.PrevBlockHash[:])
	buf.Write(b.MerkleRoot[:])
	buf.Write(utils.UInt32ToByte(b.Timestamp))
	buf.WriteByte(b.HashTarget.Exponent)
	buf.Write(b.HashTarget.Significand[:])
	buf.Write(utils.UInt32ToByte(b.Nonce))
	return DoubleHashB(buf.Bytes())
}

func DecodeBlockHeader(reader *bytes.Reader) (*BlockHeader, error) {
	var ver = BlockHeader{
		BlockVersion:  0,
		PrevBlockHash: NewHash(),
		MerkleRoot:    NewHash(),
		Timestamp:     0,
		HashTarget: CompressedTargetFormat{
			Exponent:    0,
			Significand: [3]byte{0x00, 0x00, 0x00},
		},
		Nonce: 0,
	}
	ver.BlockVersion = int32(utils.ReadUint32(reader))
	_, err := reader.Read(ver.PrevBlockHash[:])
	if err != nil {
		return nil, err
	}
	_, err = reader.Read(ver.MerkleRoot[:])
	if err != nil {
		return nil, err
	}
	ver.Timestamp = utils.ReadUint32(reader)
	ver.HashTarget.Exponent, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}
	_, err = reader.Read(ver.HashTarget.Significand[:])
	if err != nil {
		return nil, err
	}
	ver.Nonce = utils.ReadUint32(reader)
	return &ver, nil
}

type BlockMsg struct {
	BlockHeader
	TransactionCount uint64
	Transactions     []Tx
}

func (m *BlockMsg) GetCommandString() string {
	return CmdBlock
}

// Pack constructs the binary content of the version message
func (m *BlockMsg) Pack() []byte {
	var buf bytes.Buffer
	buf.Write(m.BlockHeader.Pack())
	buf.Write(utils.VarIntToByte(uint64(len(m.Transactions))))
	for _, tx := range m.Transactions {
		var txSerialized = tx.Pack()
		buf.Write(txSerialized)
	}
	return buf.Bytes()
}

func DecodeBlockMsg(reader *bytes.Reader) (*BlockMsg, error) {
	ver := &BlockMsg{
		TransactionCount: 0,
		Transactions:     make([]Tx, 0),
	}
	blockHeader, err := DecodeBlockHeader(reader)
	if err != nil {
		return nil, err
	}
	ver.BlockHeader = *blockHeader
	ver.TransactionCount = utils.ReadVarInt(reader)
	for i := 0; i < int(ver.TransactionCount); i++ {
		tx, err := DecodeTxMsg(reader)
		if err != nil {
			return nil, err
		}
		ver.Transactions = append(ver.Transactions, *tx)
	}
	return ver, nil
}

func NewBlockMsg() (*BlockMsg, error) {
	ver := &BlockMsg{
		BlockHeader: BlockHeader{
			BlockVersion:  0,
			PrevBlockHash: NewHash(),
			MerkleRoot:    NewHash(),
			Timestamp:     0,
			HashTarget: CompressedTargetFormat{
				Exponent:    0,
				Significand: [3]byte{0x00, 0x00, 0x00},
			},
			Nonce: 0,
		},
		TransactionCount: 0,
		Transactions:     make([]Tx, 0),
	}
	return ver, nil
}
