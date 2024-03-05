package msg

import (
	"bhd/utils"
	"bytes"
)

// HeadersMsg provides a contiguous set of block headers.
// No more than 2000 block headers may be sent at one time. Block headers in this array MUST be sequential, ordered by height and without range gaps.
type HeadersMsg struct {
	Count uint64
	Items []*BlockHeader
}

func (m *HeadersMsg) GetCommandString() string {
	return CmdHeaders
}

// Pack constructs the binary content of the version message
func (m *HeadersMsg) Pack() []byte {
	var buf bytes.Buffer
	buf.Write(utils.VarIntToByte(uint64(len(m.Items))))
	for _, item := range m.Items {
		buf.Write(item.Pack())
		// the depreceated transaction count
		buf.WriteByte(0x00)
	}
	return buf.Bytes()
}

func DecodeHeadersMsg(reader *bytes.Reader) (*HeadersMsg, error) {
	ver := &HeadersMsg{
		Count: 0,
		Items: make([]*BlockHeader, 0),
	}
	ver.Count = utils.ReadVarInt(reader)
	for i := 0; i < int(ver.Count); i++ {
		item, err := DecodeBlockHeader(reader)
		if err != nil {
			return nil, err
		}
		// there's one depreceated transaction count
		// field that must be read
		_ = utils.ReadVarInt(reader)
		ver.Items = append(ver.Items, item)
	}
	return ver, nil
}

func NewHeadersMsg() *HeadersMsg {
	return &HeadersMsg{
		Items: make([]*BlockHeader, 0),
	}
}
