package internal

import (
	"context"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmanilo/eth-parser/parser"
)

type ethParser struct {
	hasBlockToParse   chan bool
	blockToParse      int64
	parsedBlock       atomic.Int64
	latestBlockNumber atomic.Int64

	observersLock    sync.RWMutex
	observers        map[string]bool
	transactionsLock sync.RWMutex
	transactions     map[string][]parser.Transaction
}

func NewParser(ctx context.Context) parser.Parser {
	parser := &ethParser{
		hasBlockToParse: make(chan bool, 10),
		observers:       make(map[string]bool),
		transactions:    make(map[string][]parser.Transaction),
	}

	go parser.run(ctx)

	return parser
}

func (p *ethParser) run(ctx context.Context) {
	blockTicker := time.NewTicker(12 * time.Second)

	p.updateLatestBlockNumber()
	p.blockToParse = p.latestBlockNumber.Load()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutdown parser")
			return

		case <-blockTicker.C:
			p.updateLatestBlockNumber()

		case <-p.hasBlockToParse:
			if p.blockToParse <= p.latestBlockNumber.Load() {
				p.parseBlockTransactions()
				p.blockToParse++
			}
		}
	}
}

func (p *ethParser) updateLatestBlockNumber() {
	blockNumber := getCurrentBlockNumber()
	p.latestBlockNumber.Store(hexToInt(blockNumber))

	p.hasBlockToParse <- true
}

func (p *ethParser) parseBlockTransactions() {
	p.observersLock.RLock()
	hasObservers := len(p.observers) > 0

	p.observersLock.RUnlock()

	if !hasObservers {
		return
	}

	p.observersLock.RLock()
	snapshot := makeSnapshot(p.observers)
	p.observersLock.RUnlock()

	block := getBlock(intToHex(p.blockToParse))
	for _, tx := range block.Transactions {
		if snapshot[tx.From] {
			p.addTransaction(tx.From, parser.Transaction{
				Hash:        tx.Hash,
				BlockNumber: p.blockToParse,
				Type:        parser.OutboundTx,
			})
		}

		if snapshot[tx.To] {
			p.addTransaction(tx.To, parser.Transaction{
				Hash:        tx.Hash,
				BlockNumber: p.blockToParse,
				Type:        parser.InboundTx,
			})
		}
	}

	p.parsedBlock.Store(p.blockToParse)
}

func (p *ethParser) addTransaction(address string, tx parser.Transaction) {
	p.transactionsLock.Lock()
	defer p.transactionsLock.Unlock()

	p.transactions[address] = append(p.transactions[address], tx)
}

func makeSnapshot(input map[string]bool) map[string]bool {
	result := make(map[string]bool, len(input))
	for k, v := range input {
		result[k] = v
	}

	return result
}

func (p *ethParser) GetCurrentBlock() int {
	blockNumber := getCurrentBlockNumber()
	return int(hexToInt(blockNumber))
}

func (p *ethParser) Subscribe(address string) bool {
	addr := strings.ToLower(address)
	if !isValidHex(addr) {
		return false
	}

	p.observersLock.Lock()
	p.observers[addr] = true
	p.observersLock.Unlock()

	log.Println("added observer -", addr)

	return true
}

func (p *ethParser) GetTransactions(address string) []parser.Transaction {
	var transactions []parser.Transaction

	addr := strings.ToLower(address)

	p.transactionsLock.RLock()
	defer p.transactionsLock.RUnlock()

	for _, tx := range p.transactions[addr] {
		transactions = append(transactions, tx)
	}

	return transactions
}
