package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Amount float32
	Name   string
	Des    string
	Time   time.Time
}
type Data struct {
	TotalAmount float32
	Records     []Record
}

var theData Data
var cacheFileExist = checkFileExist(getCurrentPath() + "cache.json")

func main() {
	reader := bufio.NewReader(os.Stdin)
	mainWindow(reader)
}
func printAllData() {
	for i := 0; i < len(theData.Records); i++ {
		fmt.Println(theData.Records[i].Time.Format("[0601021504]"), theData.Records[i].Name, theData.Records[i].Amount, theData.Records[i].Des)
	}
}
func loadCache() {
	cacheByte, _ := ioutil.ReadFile(getCurrentPath() + "cache.json")
	json.Unmarshal(cacheByte, &theData)
}
func saveCache() {
	cacheByte, _ := json.Marshal(theData)
	ioutil.WriteFile(getCurrentPath()+"cache.json", cacheByte, 0666)
}
func mainWindow(reader *bufio.Reader) {
	if cacheFileExist {
		loadCache()
	}
	dailyGain(100)
	fmt.Printf("你已累计拥有￥%.2f！\n", theData.TotalAmount)
	for {
		fmt.Println("1. 记录支出")
		fmt.Println("2. 记录收入")
		fmt.Println("0. 显示全部资金")
		fmt.Print(">")
		text, _ := reader.ReadString('\n')
		text = strings.Split(text, "\r\n")[0]
		switch text {
		case "1":
			logOutcome(reader)
		case "2":
			logIncome(reader)
		case "0":
			fmt.Println("总资金：", theData.TotalAmount)
		default:
			printAllData()
		}
		calcAllAmount()
		saveCache()
	}
}
func dailyGain(amount float32) {
	var newRecord Record = Record{Amount: amount, Name: "日增", Des: "", Time: time.Now()}
	theData.Records = append(theData.Records, newRecord)
	fmt.Printf("恭喜你今天获得￥%.2f！\n", amount)
}
func calcAllAmount() {
	var sum float32 = 0
	for i := 0; i < len(theData.Records); i++ {
		sum += theData.Records[i].Amount
	}
	theData.TotalAmount = sum
}
func logOutcome(reader *bufio.Reader) {
	var newRecord Record
	fmt.Print("支出的数额：￥")
	text, _ := reader.ReadString('\n')
	text = strings.Split(text, "\r\n")[0]
	num, _ := strconv.ParseFloat(text, 32)
	amount := float32(num)
	fmt.Print("支出的条目名称：")
	name, _ := reader.ReadString('\n')
	name = strings.Split(name, "\r\n")[0]
	fmt.Print("支出的条目描述：")
	des, _ := reader.ReadString('\n')
	des = strings.Split(des, "\r\n")[0]
	newRecord.Amount = -amount
	newRecord.Name = name
	newRecord.Des = des
	newRecord.Time = time.Now()
	theData.Records = append(theData.Records, newRecord)
}
func logIncome(reader *bufio.Reader) {
	var newRecord Record
	fmt.Print("收入的数额：￥")
	text, _ := reader.ReadString('\n')
	text = strings.Split(text, "\r\n")[0]
	num, _ := strconv.ParseFloat(text, 32)
	amount := float32(num)
	fmt.Print("收入的条目名称：")
	name, _ := reader.ReadString('\n')
	name = strings.Split(name, "\r\n")[0]
	fmt.Print("收入的条目描述：")
	des, _ := reader.ReadString('\n')
	des = strings.Split(des, "\r\n")[0]
	newRecord.Amount = amount
	newRecord.Name = name
	newRecord.Des = des
	newRecord.Time = time.Now()
	theData.Records = append(theData.Records, newRecord)
}
func index(reader *bufio.Reader) {
	fmt.Println("你今天有支出吗？(Y/n)")
	for {
		fmt.Print(">")
		text, _ := reader.ReadString('\n')
		text = strings.Split(text, "\r\n")[0]
		if judgeYN(text, "y") {
			fmt.Print("支出的数额：￥")
			text, _ := reader.ReadString('\n')
			text = strings.Split(text, "\r\n")[0]
			num, _ := strconv.ParseFloat(text, 32)
			amount := float32(num)
			fmt.Println("支出：￥", amount)
		}
	}
}
func judgeYN(text string, def string) bool {
	switch def {
	case "y":
		if text == "N" || text == "n" {
			return false
		}
		return true
	case "n":
		if text == "Y" || text == "y" {
			return true
		}
		return false
	default:
		if text == "Y" || text == "y" {
			return true
		} else if text == "N" || text == "n" {
			return false
		}
		fmt.Println("[WARN]judgeYN-default-skip")
		return false
	}
}
func getCurrentPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	return string(path[0 : i+1])
}
func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
