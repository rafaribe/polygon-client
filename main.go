package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type BlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type BlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Number           string `json:"number"`
		Hash             string `json:"hash"`
		ParentHash       string `json:"parentHash"`
		Nonce            string `json:"nonce"`
		Sha3Uncles       string `json:"sha3Uncles"`
		LogsBloom        string `json:"logsBloom"`
		TransactionsRoot string `json:"transactionsRoot"`
		StateRoot        string `json:"stateRoot"`
		Miner            string `json:"miner"`
		Difficulty       string `json:"difficulty"`
		TotalDifficulty  string `json:"totalDifficulty"`
		ExtraData        string `json:"extraData"`
		Size             string `json:"size"`
		GasLimit         string `json:"gasLimit"`
		GasUsed          string `json:"gasUsed"`
		Timestamp        string `json:"timestamp"`
		Transactions     []struct {
			BlockHash        string `json:"blockHash"`
			BlockNumber      string `json:"blockNumber"`
			From             string `json:"from"`
			Gas              string `json:"gas"`
			GasPrice         string `json:"gasPrice"`
			Hash             string `json:"hash"`
			Input            string `json:"input"`
			Nonce            string `json:"nonce"`
			To               string `json:"to"`
			TransactionIndex string `json:"transactionIndex"`
			Value            string `json:"value"`
			V                string `json:"v"`
			R                string `json:"r"`
			S                string `json:"s"`
		} `json:"transactions"`
		Uncles []string `json:"uncles"`
	} `json:"result"`
}

func makeRPCRequest(client http.Client, reqBody map[string]interface{}) ([]byte, error) {
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://polygon-rpc.com", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = http.MaxBytesReader(nil, req.Body, 1048576)
	req.Body = nopCloser{bytes.NewReader(reqJSON)}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response body: %v", err)
	}

	return respBody, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func main() {
	// Create an HTTP client to make requests to the Polygon RPC endpoint
	client := http.Client{
		Timeout: time.Second * 5,
	}

	// Start an infinite loop to periodically make requests to the RPC endpoint
	for {
		// Get the latest block number
		blockNumberReq := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_blockNumber",
			"id":      1,
		}

		blockNumberRespBytes, err := makeRPCRequest(client, blockNumberReq)
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

		blockRespBytes, err := makeRPCRequest(client, blockReq)
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
		log.Printf("Latest block hash: %s", blockResp.Result.Hash)

		// Wait for 5 seconds before making the next request
		time.Sleep(time.Second * 5)
	}
}
