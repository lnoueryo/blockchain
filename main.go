package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER = "THE BLOCKCHAIN"
	MINING_REWARD = 1.0
)

type Block struct {
	Nonce int `json:"nonce"`
	PreviousHash [32]byte `json:"previous_hash"`
	Timestamp int64 `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.Nonce = nonce
	b.PreviousHash = previousHash
	b.Timestamp = time.Now().UnixNano()
	b.Transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("Nonce			%d\n", b.Nonce)
	fmt.Printf("PreviousHash		%x\n", b.PreviousHash)
	fmt.Printf("Timestamp		%d\n", b.Timestamp)
	for _, t := range b.Transactions {
		t.Print()
	}
	// fmt.Printf("Transactions		%s\n", b.Transactions)
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block)MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Timestamp int64 `json: "timestamp"`
		Nonce int `json: "nonce"`
		PreviousHash [32]byte `json: "previous_hash"`
		Transaction []*Transaction `json: "transactions"`
	}{
		Timestamp: b.Timestamp,
		Nonce: b.Nonce,
		PreviousHash: b.PreviousHash,
		Transaction: b.Transactions,
	})
}



type BlockChain struct {
	TransactionPool		[]*Transaction
	Chain				[]*Block
	BlockChainAddress	string
}

func NewBlockChain(blockChainAddress string) *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.BlockChainAddress = blockChainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *BlockChain) CreateBlock(nounce int, previousHash [32]byte) *Block {
	b := NewBlock(nounce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*Transaction{}
	return b
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.TransactionPool = append(bc.TransactionPool, t)
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions, NewTransaction(
			t.SenderBlockChainAddress,
			t.RecipientBlockChainAddress,
			t.Value,
		))
	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		nonce,
		previousHash,
		0,
		transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.BlockChainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			value := t.Value
			if blockchainAddress == t.RecipientBlockChainAddress {
				totalAmount += value
			}

			if  blockchainAddress == t.SenderBlockChainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Transaction struct {
	SenderBlockChainAddress		string
	RecipientBlockChainAddress	string
	Value						float32
}

func NewTransaction (sender string, recipient string, value float32) *Transaction{
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address %s\n", t.SenderBlockChainAddress)
	fmt.Printf(" recipient_blockchain_address %s\n", t.RecipientBlockChainAddress)
	fmt.Printf(" value %.1f\n", t.Value)
}


func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	// block := &Block{Nonce: 1}
	// fmt.Printf("%x\n", block.Hash())
	myBlockChainAddress := "my_blockchain_address"
	blockChain := NewBlockChain(myBlockChainAddress)
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("C", "D", 2.0)
	blockChain.AddTransaction("A", "B", 9.0)
	blockChain.Mining()
	blockChain.Print()

	fmt.Printf("Miner %.1f\n", blockChain.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount("A"))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount("B"))
}