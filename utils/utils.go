package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"nkn-node-manager/models"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const IP_FILE = "ip.txt"

func Check_nkn_nodes() {
	var succ, fail uint32 = 0, 0
	ipList, err := GetIpList()
	if err != nil {
		fmt.Println(err)
		return
	}
	msgChan := make(chan string)
	for i, ip := range ipList {
		if i%100 == 0 && i > 1 {
			time.Sleep(5 * time.Second)
		}
		go func(ip string) {
			msg, ok := GetNodeState(ip)
			if ok {
				atomic.AddUint32(&succ, 1)
			} else {
				atomic.AddUint32(&fail, 1)
			}
			msgChan <- msg
		}(ip)
	}

	for {
		select {
		case m := <-msgChan:
			fmt.Println(m)
		default:
		}
		if succ+fail == uint32(len(ipList)) {
			break
		}
	}

	fmt.Printf("total: %d, succ: %d, failed: %d", succ+fail, succ, fail)
}

func GetIpList() ([]string, error) {
	var ipList []string
	inFile, err := os.Open(IP_FILE)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		ip := scanner.Text()
		if checkIPAddress(ip) {
			ipList = append(ipList, ip)
		}
		//fmt.Println(scanner.Text()) // the line
	}
	return ipList, nil
}

func GetIpListFromDB() ([]string, error) {
	var ipList []string
	inFile, err := os.Open(IP_FILE)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		ip := scanner.Text()
		if checkIPAddress(ip) {
			ipList = append(ipList, ip)
		}
		//fmt.Println(scanner.Text()) // the line
	}
	return ipList, nil
}

func GetNodeState(ip string) (string, bool) {

	jsonBody := []byte(`{"jsonrpc": "2.0","method": "getnodestate","id": "1","params": {}}`)
	body := bytes.NewBuffer(jsonBody)

	// Create client
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("http://%s:30003", ip)

	// Create request
	req, err := http.NewRequest("POST", url, body)

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp := new(http.Response)
	for i := 3; i > 0; i-- {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return fmt.Sprintf("ip: %-15s | err message: %s", ip, err), false
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		fmt.Println(err)
	}

	status, ok := result["result"]
	if ok {
		syncState := status.(map[string]interface{})["syncState"]
		height := status.(map[string]interface{})["height"]
		h := strconv.FormatFloat(height.(float64), 'f', -1, 64)
		return fmt.Sprintf("ip: %-15s | height: %s | sync state: %s", ip, h, syncState), true
	} else {
		e, ok := result["error"]
		if !ok {
			return fmt.Sprintf("ip: %-15s | err message: %s", ip, "can not access"), true
		}
		errorMsg := e.(map[string]interface{})["message"]
		return fmt.Sprintf("ip: %-15s | err message: %s", ip, errorMsg), true
	}
}

func checkIPAddress(ip string) bool {
	if len(ip) == 0 {
		return false
	}
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}

func CheckOffline() {
	for {
		var wallets []models.Wallet
		models.DB.Find(&wallets)
		for _, w := range wallets {
			var ww models.Wallet
			error, active := GetNodeState(w.IP)
			fmt.Println(error)
			if !active {
				if err := models.DB.Where("ip = ?", w.IP).First(&ww).Update("idle", true).Error; err != nil {
					fmt.Println("not active err:", err)
				}
			} else {
				ts := time.Now().Unix()
				if err := models.DB.Where("ip = ?", w.IP).First(&ww).Update("lastUpdate", ts).Error; err != nil {
					//fmt.Println("active err:", err)
				}
				if err := models.DB.Where("ip = ?", w.IP).First(&ww).Update("idle", false).Error; err != nil {
					//fmt.Println("active err:", err)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
