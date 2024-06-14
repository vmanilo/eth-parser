package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ethAPI = "https://cloudflare-eth.com"

type ethBlockNumberResult struct {
	Result string `json:"result"`
}

func getCurrentBlockNumber() string {
	resp, _ := http.Post(ethAPI, "application/json",
		bytes.NewBufferString(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`),
	)

	data, _ := io.ReadAll(resp.Body)

	var result ethBlockNumberResult
	_ = json.Unmarshal(data, &result)

	return result.Result
}

func getBlock(blockNumber string) *ethBlock {
	resp, _ := http.Post(ethAPI, "application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":1}`, blockNumber)),
	)

	data, _ := io.ReadAll(resp.Body)

	result := new(ethBlockResult)
	_ = json.Unmarshal(data, result)

	return result.Result
}

type ethBlockResult struct {
	Result *ethBlock `json:"result"`
}

type ethBlock struct {
	Hash         string                 `json:"hash"`
	Number       string                 `json:"number"`
	Transactions []*ethBlockTransaction `json:"transactions"`
}

type ethBlockTransaction struct {
	Hash string `json:"hash"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
