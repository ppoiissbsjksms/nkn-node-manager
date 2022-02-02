package controllers

import (
	"encoding/hex"
	"github.com/nknorg/nkn-sdk-go"
	"io"
	"io/ioutil"
	"net/http"
	"nkn-node-manager/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateWalletInput struct {
	KeyStore   string `json:"keystore" binding:"required"`
	PassWord   string `json:"password" binding:"required"`
	IP         string `json:"ip"`
	Idle       bool   `json:"idle"`
	LastActive uint   `json:"lastActive"`
}

// FindWallets GET /wallets
// Find all wallet
func FindWallets(c *gin.Context) {
	var wallets []models.Wallet
	models.DB.Find(&wallets)

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

// GenerateID GET /generate/:address
// send 10 nkn to address
func GenerateID(c *gin.Context) {
	// Get model if exist

	walletBytes, err := ioutil.ReadFile("wallet.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	passBytes, err := ioutil.ReadFile("wallet.pswd")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	walletString := string(walletBytes)
	passString := strings.TrimSpace(string(passBytes))
	conf := &nkn.WalletConfig{
		Password:          passString,
		SeedRPCServerAddr: nkn.NewStringArray("http://seed.nkn.org:30003"),
	}
	w, err := nkn.WalletFromJSON(walletString, conf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	txnHash, err := w.Transfer(c.Param("address"), "10", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tx": txnHash})
}

// FindWallet GET /wallets/:address
// Find a wallet
func FindWallet(c *gin.Context) {
	// Get model if exist
	var wallet models.Wallet
	if err := models.DB.Where("address = ?", c.Param("address")).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

// FindIdleWallet GET /wallets/free
// Find a wallet
func FindIdleWallet(c *gin.Context) {
	// Get model if exist
	var wallet models.Wallet
	ts := time.Now().Unix()
	if err := models.DB.Where("idle = true AND last_active < ?", ts-3600).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"idle": wallet})
}

// DeleteWallet GET /remove/:address
// Remove a wallet
func DeleteWallet(c *gin.Context) {
	// Get model if exist
	var wallet models.Wallet
	if err := models.DB.Where("address = ?", c.Param("address")).Delete(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"delete": c.Param("address")})
}

// UploadWallet POST /wallet
// Create new wallet
func UploadWallet(c *gin.Context) {
	// Validate input
	var input CreateWalletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create wallet
	w, err := nkn.WalletFromJSON(input.KeyStore, &nkn.WalletConfig{Password: input.PassWord})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet := models.Wallet{KeyStore: input.KeyStore, PassWord: input.PassWord, IP: c.ClientIP(),
		Idle: false, Address: w.Address(), PublicKey: hex.EncodeToString(w.PubKey()), LastActive: time.Now().Unix()}
	models.DB.Create(&wallet)

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

// UploadWalletFile POST /walletform
// Create new wallet
func UploadWalletFile(c *gin.Context) {
	mf, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := mf.File["keystore"][0].Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	walletJson := new(strings.Builder)
	_, err = io.Copy(walletJson, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pass, err := mf.File["password"][0].Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := new(strings.Builder)
	_, err = io.Copy(password, pass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := strings.TrimSpace(password.String())
	w, err := nkn.WalletFromJSON(walletJson.String(), &nkn.WalletConfig{Password: p})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet := models.Wallet{KeyStore: walletJson.String(), PassWord: p, IP: c.ClientIP(),
		Idle: false, Address: w.Address(), PublicKey: hex.EncodeToString(w.PubKey()), LastActive: time.Now().Unix()}
	models.DB.Create(&wallet)

	c.JSON(http.StatusOK, gin.H{"data": wallet.Address})
}
