package parser

type TxType string

const (
	InboundTx  TxType = "inbound"
	OutboundTx TxType = "outbound"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type Transaction struct {
	Hash        string `json:"hash"`
	BlockNumber int64  `json:"blockNumber"`
	Type        TxType `json:"type"`
}
