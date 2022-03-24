/*
 * Copyright © 2021-2022 V1def

 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.

 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.

 * You should have received a copy of the GNU Lesser General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
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
