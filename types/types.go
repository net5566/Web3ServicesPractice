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
	Hash  string `bson:"tx_hash"`
	From  string `bson:"from"`
	To    string `bson:"to"`
	Value string `bson:"value"`
	Nonce int    `bson:"nonce"`
	Data  string `bson:"data"`
	Logs  []Log  `bson:"logs"`
}

type BlockTransactions struct {
	BlockNum     int      `bson:"block_num"`
	Transactions []string `bson:"transactions"`
}
