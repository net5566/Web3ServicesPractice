package types

type Block struct {
	BlockNum   int64  `gorm:"primaryKey;autoIncrement:false"`
	BlockHash  string `gorm:"not null"`
	BlockTime  int64  `gorm:"not null"`
	ParentHash string `gorm:"not null"`
}
