/*
 * MIT License

 * Copyright (c) 2022 V1def

 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:

 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.

 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"flag"
	"log"

	"github.com/v1def/subscribe-ethereum-transfer/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	clientAddr   string // Client ethereum address.
	contractAddr string // Ethereum contranct address.
)

// A function that is called before the main function.
func init() {
	// Add a new flag for the ethereum client address.
	flag.StringVar(&clientAddr, "client-addr", "", "Ethereum client address")
	// Add a new flag for the ethereum contract address.
	flag.StringVar(&contractAddr, "contract-addr", "", "Ethereum contract address")

	// Parsing all the flags that were specified at startup.
	flag.Parse()

	// Empty check.
	if clientAddr == "" || contractAddr == "" {
		log.Fatal("error client or contranct address is empty")
	}
}

// The main function that is called when running a this program.
func main() {
	// Create a new client connections.
	client, err := ethclient.Dial(clientAddr)
	if err != nil {
		log.Fatalf("error creating client: %s", err.Error())
	}
	defer client.Close()

	log.Print("Ethereum client connected!")

	// Connect to ethereum contract.
	tc, err := contract.NewContract(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("error connections to contract: %s", err.Error())
	}

	// Create transfers channel.
	transfers := make(chan *contract.ContractTransfer)

	// Subscribe to ethereum contract transfers event.
	sub, err := tc.WatchTransfer(nil, transfers, nil, nil)
	if err != nil {
		log.Fatalf("error subcribe to event: %s", err.Error())
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf(": %s", err.Error())
		case t := <-transfers:
			log.Println("===================================================")
			log.Printf("From:         %s", t.From)
			log.Printf("To:           %s", t.To)
			log.Printf("Value:        %s", t.Value)
			log.Printf("Block number: %d", t.Raw.BlockNumber)
			log.Printf("Block hash:   %s", t.Raw.BlockHash)
		}
	}
}
