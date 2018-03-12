package cli

import (
	. "github.com/linxinzhe/go-simple-coin/datastruct"
	"flag"
	"os"
	"fmt"
	"strconv"
)

type CLI struct {
	BC *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  addblock - add a blocks in the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	createBlockchainData := createBlockchainCmd.String("address", "", "user address")

	fromData := sendCmd.String("from", "", "from address")
	toData := sendCmd.String("to", "", "to address")
	amountData := sendCmd.Int("amount", 0, "to address")

	getBalanceData := getBalanceCmd.String("address", "", "user address")

	switch os.Args[1] {
	case "printchain":
		_ = printChainCmd.Parse(os.Args[2:])
	case "createblockchain":
		_ = createBlockchainCmd.Parse(os.Args[2:])
	case "getbalance":
		_ = getBalanceCmd.Parse(os.Args[2:])
	case "send":
		_ = sendCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainData == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainData)
	}

	if sendCmd.Parsed() {
		if *fromData == "" || *toData == "" || *amountData <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*fromData, *toData, *amountData)
	}
}

func (cli *CLI) printChain() {
	bci := cli.BC.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)

		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) createBlockchain(address string) {
	bc := NewBlockchain(address)
	defer bc.DB.Close()
	cli.BC = bc
	cli.printChain()
}

func (cli *CLI) send(from string, to string, amount int) {
	bc := NewBlockchain(from)
	cli.BC = bc
	defer bc.DB.Close()
	transaction := NewUTXOTransaction(from, to, amount, cli.BC)
	cli.BC.MineBlock([]*Transaction{transaction})
	fmt.Println("Success!")
}

func (cli *CLI) getBalance(address string) int {
	bc := NewBlockchain(address)
	cli.BC = bc
	txOutputs := cli.BC.FindUTXO(address)
	balance := 0
	for _, out := range txOutputs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
	return balance
}
