package msg

const (
	TokenImmutable = 0x00
	TokenMutable   = 0x01
	TokenMinting   = 0x02
)

// CashToken (s) are ideologically similar to BEP-20 tokens found on BNB Chain or ERC-20 tokens found on Ethereum,
// in that they allow anybody to deploy tokens that represent practically any type of asset.
//
// These tokens are distinct from the native Bitcoin Cash gas unit (BCH) but can still be transferred on
// the blockchain via transactions. CashToken transactions are considered just as secure as non-token transactions, with no use of additional indexing software.
//
// The new token format enables a variety of business applications to be developed on the Bitcoin Cash blockchain,
// including identity tokens and decentralized exchanges.
// As of block #792773, the Bitcoin Cash blockchain is now capable of supporting CashTokens. As per data from 3XPL,
// more than 25,000 CashToken NFTs have been created since the upgrade, in addition to over 1,100 fungible tokens (FTs).
// Most of these represent NFT art collections.
type CashToken struct {
	Category      Hash
	TokenBitField uint8
	TokenType     int
	Commitment    []byte
	Amount        uint64
}
