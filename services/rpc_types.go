package services

type RPCBlock struct {
	Number       string
	Hash         string
	ParentHash   string
	Timestamp    string
	Transactions []RPCTransaction
}

type RPCTransaction struct {
	Number string `json:"blockNumber"`
	Hash   string `json:"hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Value  string `json:"value"`
	Nonce  string `json:"nonce"`
	Data   string `json:"input"`
}

type RPCTransactionReceipt struct {
	Logs []RPCLog `json:"logs"`
}

type RPCLog struct {
	Data  string `json:"data"`
	Index string `json:"logIndex"`
}
