package msg

import (
	"bhd/utils"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

/*
	Bitcoin Cash (BCH) is a cryptocurrency that was created and launched to bring decentralization back to cryptocurrency.
	It is the result of a 2017 Bitcoin "hard fork," which occurs when an existing blockchain splits into two.
	Bitcoin Cash allows a greater number of transactions in a single block than Bitcoin, which should lower fees and transaction times.
	Bitcoin Cash is designed to be used as a cheap payment system, much in the way Bitcoin was designed to be originally.
	Transactions fees are generally less than $.01, and transaction confirmation times are significantly less than Bitcoin's, generally within seconds.
*/

import (
	"errors"
)

var (
	// MagicValue - the network identifier is used to separate blockchains and test networks.
	// This reduces unnecessary load on peers, allowing them to rapidly ban nodes rather then forcing the peer to
	// do a blockchain analysis before banning or disconnecting. For Bitcoin Cash main net, the net magic field is
	// always 0xE3E1F3E8 (the ASCII string, "cash", with each byte's highest bit set).
	// Any message received that does not begin with the net magic is invalid./
	MagicValue                   = []byte{0xE3, 0xE1, 0xF3, 0xE8}
	EmptyPayloadChecksum         = []byte{0x5D, 0xF6, 0xE0, 0xE2}
	SupportedServices            = utils.ByteToUInt64([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	DefaultHandShakeTimeout      = 3 * time.Second
	DefaultQueryReadTimeout      = 10 * time.Second
	ServerListenPort             = 8333
	MinConnectedPeersInPool      = 24 // MinConnectedPeersInPool is the number of always connected peers in the PeerService
	CashTokenPrefix         byte = 0xef
	BlockFolder                  = "/blocks"
	CacheBlocksLocally           = true
)

/*
The Bitcoin Cash Peer-to-Peer (P2P) Network protocol is a binary protocol used by Full Nodes and SPV Nodes, transmitted over TCP/IP.
Individual nodes on the Bitcoin Cash network connect and create a mesh network where each node is indirectly connected to many others via just a couple of hops.
In the original Satoshi implementation of the P2P protocol the design of INV and getdata have been used for propagating transaction data
using the rules of the gossip protocol values: forwarding validated transactions to a few peer-nodes who send it to others until the
entire network has the transaction. This emergent behavior of the P2P layer allows fast propagation without undue strain on any individual node.

The P2P protocol is designed around messages. Each message is separate and self-contained. Nodes should be tolerant of message-types they do
not understand. It is best to simply ignore those.

Generally speaking, each message is an event that the node can choose to respond to. Events can be notifications of new data (transactions/blocks/etc),
requests for such data to be sent, or the sending of the data itself. In some specific cases a message can indicate the rejection of another message,
though this is optional and should not be relied upon.

These design decisions were made with consideration to communication with untrusted/uncooperative partners.
Developer Notes: A common message strategy is to wait for any message that provides the required data (with a timeout), and then separately
issue the request in a retry loop to multiple peers.

*/

const (
	ProtocolVersion        uint32 = 70015
	AgentString            string = "/cryptopera 1.0.0/"
	ListeningPort                 = 8333
	MaxCommandStringLength        = 12
)

// The command string is a fixed-length 12 byte ASCII string. Commands may not be longer than 12 bytes.
// Commands that are shorter than 12 bytes are right-padded with null bytes (0x00). The command string is used
// to determine the type of message being transmitted. Messages with an unrecognized command string are ignored by most
// implementations but may result in a ban by implementations that diverge from the Satoshi-client defacto standard.
const (
	CmdVersion      = "version"
	CmdXVersion     = "xversion"
	CmdVerAck       = "verack"
	CmdAvaHello     = "avahello"
	CmdXVerAck      = "xverack"
	CmdGetAddr      = "getaddr"
	CmdAddr         = "addr"
	CmdAddrV2       = "addrv2"
	CmdGetBlocks    = "getblocks"
	CmdInv          = "inv"
	CmdGetData      = "getdata"
	CmdNotFound     = "notfound"
	CmdBlock        = "block"
	CmdTx           = "tx"
	CmdGetHeaders   = "getheaders"
	CmdHeaders      = "headers"
	CmdPing         = "ping"
	CmdPong         = "pong"
	CmdMemPool      = "mempool"
	CmdFilterAdd    = "filteradd"
	CmdFilterClear  = "filterclear"
	CmdFilterLoad   = "filterload"
	CmdMerkleBlock  = "merkleblock"
	CmdReject       = "reject"
	CmdSendHeaders  = "sendheaders"
	CmdFeeFilter    = "feefilter"
	CmdGetCFilters  = "getcfilters"
	CmdGetCFHeaders = "getcfheaders"
	CmdGetCFCheckpt = "getcfcheckpt"
	CmdGetCFMempool = "getcfmempool"
	CmdCFilter      = "cfilter"
	CmdCFHeaders    = "cfheaders"
	CmdCFCheckpt    = "cfcheckpt"
	CmdSendCmpct    = "sendcmpct"
	CmdCmpctBlock   = "cmpctblock"
	CmdGetBlockTxns = "getblocktxn"
	CmdBlockTxns    = "blocktxn"
	CmdDsProof      = "dsproof-beta"
	CmdSendAddrv2   = "sendaddrv2"
	CmdProtoConf    = "protoconf"
)

const (
	InvTypeTransaction   = 0x0001
	InvTypeBlock         = 0x0002
	InvTypeFilteredBlock = 0x0003
	InvTypeCompactBlock  = 0x0004
	InvTypeXThinBlock    = 0x0005
	InvTypeGrapheneBlock = 0x0006
	InvTypeDblSpendProof = 0x94A0
)

const (
	BloomFilterUpdateNone          = 0
	BloomFilterUpdateAll           = 1
	BloomFilterUpdateP2PPubKeyOnly = 2
)

const (
	TransactionVersion = 0x00000001
)

const (
	RejectMalformed       = 0x01
	RejectInvalid         = 0x10
	RejectObsolete        = 0x11
	RejectDuplicate       = 0x12
	RejectNonstandard     = 0x40
	RejectDust            = 0x41
	RejectInsufficientFee = 0x42
	RejectCheckpoint      = 0x43
)

// AcceptedNodesAgents - BCH nodes, what to tell you, its fucked
// will go only with these bitcoin cash nodes
// bchd is fucked, bitcoin sv...don't even know
// what those are
var AcceptedNodesAgents = []string{
	"Bitcoin Cash Node:26",
	"Bitcoin Cash Node:27",
	"Bitcoin Cash Node:28",
	"Bitcoin Cash Node:29",
	"Bitcoin Cash Node:30",
	"Bitcoin Cash Node:31",
	"Bitcoin Cash Node:32",
	"Bitcoin Cash Node:33",
	"Bitcoin Cash Node:34",
}

var (
	EmptyHash = Hash(NewHash())
)

type Hash []byte

func NewHash() []byte {
	return make([]byte, 32)
}

func (m Hash) Reverse() Hash {
	var hsh = NewHash()
	j := 0
	for i := len(m) - 1; i >= 0; i-- {
		hsh[j] = m[i]
		j++
	}
	return hsh
}

// ZeroPrefixLen returns number of zeros at start of the hash
func (m Hash) ZeroPrefixLen() int {
	return 64 - len(strings.TrimLeft(m.ToString(), "0"))
}

func NewHashFromString(strHash string) (Hash, error) {
	var hash = NewHash()
	bytes, err := hex.DecodeString(strHash)
	if err != nil {
		return hash, err
	}
	if len(bytes) != 32 {
		return hash, errors.New("the string length is not 64")
	}
	copy(hash, bytes)
	return Hash(hash).Reverse(), nil
}

func (m Hash) MarshalJSON() ([]byte, error) {
	// note that hash is reversed (little endian)
	return json.Marshal(hex.EncodeToString(m.Reverse()))
}

func (m Hash) ToString() string {
	// note that hash is reversed (little endian)
	return hex.EncodeToString(m.Reverse())
}

// IsEmpty returns true if the hash is all zero
func (m Hash) IsEmpty() bool {
	var sum int = 0
	for _, el := range m {
		sum += int(el)
	}
	return sum == 0
}

type ProtocolPdu interface {
	GetCommandString() string
	Pack() []byte
}

// Pack returns binary representation of message
func Pack(pdu ProtocolPdu) []byte {
	// get the payload
	var buf bytes.Buffer
	buf.Write(MagicValue)
	var command [MaxCommandStringLength]byte
	cmd := pdu.GetCommandString()
	if len(cmd) > MaxCommandStringLength {
		panic(cmd + " is too long for command string")
	}
	copy(command[:], []byte(cmd))
	buf.Write(command[:])
	// size of the payload
	payload := pdu.Pack()
	buf.Write(utils.UInt32ToByte(uint32(len(payload))))
	// checksum of payload, first 4 digits
	if len(payload) > 0 {
		buf.Write(DoubleHashB(payload)[0:4])
	} else {
		buf.Write(EmptyPayloadChecksum)
	}
	buf.Write(payload)
	return buf.Bytes()
}

// DoubleHashB calculates hash(hash(b)) and returns the resulting bytes.
func DoubleHashB(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}
