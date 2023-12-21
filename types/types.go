package types

type Block struct {
	BlockNum   int64  `gorm:"primaryKey;autoIncrement:false"`
	BlockHash  string `gorm:"not null"`
	BlockTime  int64  `gorm:"not null"`
	ParentHash string `gorm:"not null"`
}

type Log struct {
	Data  string `bson:"data"`
	Index int    `bson:"index"`
}

type Transaction struct {
	Hash  string `bson:"tx_hash" json:"tx_hash"`
	From  string `bson:"from" json:"from"`
	To    string `bson:"to" json:"to"`
	Value string `bson:"value" json:"value"`
	Nonce int    `bson:"nonce" json:"nonce"`
	Data  string `bson:"data" json:"data"`
	Logs  []Log  `bson:"logs" json:"logs"`
}

type BlockTransactions struct {
	BlockNum     int      `bson:"block_num"`
	Transactions []string `bson:"transactions"`
}
