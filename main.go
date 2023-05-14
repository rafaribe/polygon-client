package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create an HTTP client to make requests to the Polygon RPC endpoint
	client := http.Client{
		Timeout: time.Second * 5,
	}
	// This could be an env var or an array of endpoints if we wish to expand the application further
	polygonRpcEndpoint := "https://polygon-rpc.com"
	// Start an infinite loop to periodically make requests to the RPC endpoint
	for {
		// Get the latest block number
		blockNumberReq := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_blockNumber",
			"id":      1,
		}

		blockNumberRespBytes, err := makeRPCRequest(&client, polygonRpcEndpoint, blockNumberReq)
		if err != nil {
			log.Printf("error getting block number: %v", err)
			continue
		}

		var blockNumberResp BlockNumberResponse
		if err := json.Unmarshal(blockNumberRespBytes, &blockNumberResp); err != nil {
			log.Printf("error unmarshalling block number response: %v", err)
			continue
		}

		// Get the latest block
		blockReq := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_getBlockByNumber",
			"params": []interface{}{
				blockNumberResp.Result,
				true,
			},
			"id": 2,
		}

		blockRespBytes, err := makeRPCRequest(&client, polygonRpcEndpoint, blockReq)
		if err != nil {
			log.Printf("error getting block: %v", err)
			continue
		}

		var blockResp BlockResponse
		if err := json.Unmarshal(blockRespBytes, &blockResp); err != nil {
			log.Printf("error unmarshalling block response: %v", err)
			continue
		}

		log.Printf("Latest block number: %s", blockNumberResp.Result)
		log.Printf("Latest block hash: %s", blockResp.Result)

		// Wait for 5 seconds before making the next request
		time.Sleep(time.Second * 5)
	}
}
