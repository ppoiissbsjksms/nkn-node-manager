package models

type Wallet struct {
	ID         uint   `json:"id"`
	KeyStore   string `json:"keystore"`
	PassWord   string `json:"password"`
	IP         string `json:"ip"`
	Address    string `json:"address" gorm:"primary_key"`
	PublicKey  string `json:"publicKey"`
	Idle       bool   `json:"idle"`
	LastActive int64  `json:"lastActive"`
}
