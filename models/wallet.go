package models

type Wallet struct {
	ID        uint   `json:"id"`
	KeyStore  string `json:"keystore"`
	PassWord  string `json:"password"`
	IP        string `json:"ip"`
	Active    bool   `json:"active"`
	Address   string `json:"address" gorm:"primary_key"`
	PublicKey string `json:"publicKey"`
}
